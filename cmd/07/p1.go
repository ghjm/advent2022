package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"strconv"
	"strings"
)

type Dir struct {
	Parent *Dir
	Dirs   map[string]*Dir
	Files  map[string]int
	Size   int
}

func newDir(parent *Dir) *Dir {
	return &Dir{
		Parent: parent,
		Dirs:   make(map[string]*Dir),
		Files:  make(map[string]int),
	}
}

func getAllDirs(d *Dir) []*Dir {
	dirList := []*Dir{d}
	for _, sd := range d.Dirs {
		dirList = append(dirList, getAllDirs(sd)...)
	}
	return dirList
}

func updateSizes(d *Dir) {
	d.Size = 0
	for _, fs := range d.Files {
		d.Size += fs
	}
	for _, sd := range d.Dirs {
		updateSizes(sd)
		d.Size += sd.Size
	}
}

func run() error {
	root := newDir(nil)
	curDir := root
	inLs := false
	err := utils.OpenAndReadLines("input7.txt", func(s string) error {
		if s[0] == '$' {
			inLs = false
			cmd := strings.Split(s[2:], " ")
			if len(cmd) < 1 {
				return fmt.Errorf("no command")
			}
			if cmd[0] == "cd" {
				if len(cmd) < 2 {
					return fmt.Errorf("no param for cd")
				}
				if cmd[1] == "/" {
					curDir = root
				} else if cmd[1] == ".." {
					if curDir.Parent == nil {
						return fmt.Errorf("cd above root")
					}
					curDir = curDir.Parent
				} else {
					d, ok := curDir.Dirs[cmd[1]]
					if !ok {
						return fmt.Errorf("cd into nonexistent")
					}
					curDir = d
				}
			} else if cmd[0] == "ls" {
				inLs = true
			}
		} else {
			if inLs {
				item := strings.Split(s, " ")
				if len(item) != 2 {
					return fmt.Errorf("invalid item")
				}
				if item[0] == "dir" {
					curDir.Dirs[item[1]] = newDir(curDir)
				} else {
					size, err := strconv.Atoi(item[0])
					if err != nil {
						return fmt.Errorf("invalid size: %w", err)
					}
					curDir.Files[item[1]] = size
				}
			} else {
				return fmt.Errorf("syntax error")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	updateSizes(root)

	curFree := 70000000 - root.Size
	needed := 30000000 - curFree
	var bestSize int

	totalSizes := 0
	for _, d := range getAllDirs(root) {
		if d.Size <= 100000 {
			totalSizes += d.Size
		}
		if d.Size >= needed {
			if bestSize == 0 || d.Size < bestSize {
				bestSize = d.Size
			}
		}
	}
	fmt.Printf("Part 1: %d\n", totalSizes)
	fmt.Printf("Part 2: %d\n", bestSize)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
