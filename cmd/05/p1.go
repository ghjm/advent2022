package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"regexp"
	"strconv"
)

func run() error {
	var initMap []string
	var moves []string
	inInit := true
	err := utils.OpenAndReadLines("input5.txt", func(s string) error {
		if s == "" {
			inInit = false
		} else if inInit {
			initMap = append(initMap, s)
		} else {
			moves = append(moves, s)
		}
		return nil
	})
	if err != nil {
		return err
	}
	legend := initMap[len(initMap)-1]
	initState := make(map[rune]string)
	for li, lc := range legend {
		if lc != ' ' {
			var stack string
			for level := len(initMap) - 2; level >= 0; level-- {
				row := initMap[level]
				if len(row) > li && row[li] != ' ' {
					stack = fmt.Sprintf("%s%c", stack, row[li])
				}
			}
			initState[lc] = stack
		}
	}
	state9000 := make(map[rune]string)
	state9001 := make(map[rune]string)
	for k, v := range initState {
		state9000[k] = v
		state9001[k] = v
	}
	re := regexp.MustCompile(`^move (\d+) from (\d+) to (\d+)$`)
	for _, move := range moves {
		m := re.FindStringSubmatch(move)
		if m == nil {
			return fmt.Errorf("line did not match regex: %s", move)
		}
		count, err := strconv.Atoi(m[1])
		if err != nil {
			return err
		}
		moveFrom, moveTo := rune(m[2][0]), rune(m[3][0])
		for i := 0; i < count; i++ {
			c := state9000[moveFrom][len(state9000[moveFrom])-1]
			state9000[moveFrom] = state9000[moveFrom][0 : len(state9000[moveFrom])-1]
			state9000[moveTo] = fmt.Sprintf("%s%c", state9000[moveTo], c)
		}
		stack := state9001[moveFrom][len(state9001[moveFrom])-count : len(state9001[moveFrom])]
		state9001[moveFrom] = state9001[moveFrom][0 : len(state9001[moveFrom])-count]
		state9001[moveTo] = fmt.Sprintf("%s%s", state9001[moveTo], stack)
	}
	fmt.Printf("Part 1: ")
	for _, lc := range legend {
		if lc != ' ' {
			fmt.Printf("%c", state9000[lc][len(state9000[lc])-1])
		}
	}
	fmt.Printf("\n")
	fmt.Printf("Part 2: ")
	for _, lc := range legend {
		if lc != ' ' {
			fmt.Printf("%c", state9001[lc][len(state9001[lc])-1])
		}
	}
	fmt.Printf("\n")

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
