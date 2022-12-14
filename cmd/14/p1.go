package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

type mapData struct {
	walls map[point]struct{}
	sand  map[point]struct{}
	minP  *point
	maxP  *point
	floor *int
}

func (md *mapData) setWall(x, y int) {
	p := point{x, y}
	if md.walls == nil {
		md.walls = make(map[point]struct{})
	}
	md.walls[p] = struct{}{}
	if md.minP == nil {
		md.minP = &point{500, 0}
	}
	if x < md.minP.x {
		md.minP.x = x
	}
	if y < md.minP.y {
		md.minP.y = y
	}
	if md.maxP == nil {
		md.maxP = &point{500, 0}
	}
	if x > md.maxP.x {
		md.maxP.x = x
	}
	if y > md.maxP.y {
		md.maxP.y = y
	}
}

func (md *mapData) pointCh(x, y int) rune {
	if md.floor != nil && y >= *md.floor {
		return '#'
	}
	_, ok := md.walls[point{x, y}]
	if ok {
		return '#'
	}
	_, ok = md.sand[point{x, y}]
	if ok {
		return 'o'
	}
	return '.'
}

func (md *mapData) placeSand(part1 bool) bool {
	sP := point{500, 0}
	var good bool
	for {
		if md.pointCh(sP.x, sP.y+1) == '.' {
			sP = point{sP.x, sP.y + 1}
		} else if md.pointCh(sP.x-1, sP.y+1) == '.' {
			sP = point{sP.x - 1, sP.y + 1}
		} else if md.pointCh(sP.x+1, sP.y+1) == '.' {
			sP = point{sP.x + 1, sP.y + 1}
		} else {
			break
		}
		if part1 {
			good = sP.x >= md.minP.x && sP.x <= md.maxP.x && sP.y >= md.minP.y && sP.y <= md.maxP.y
		} else {
			good = sP.x != 500 || sP.y != 0
		}
		if !good {
			break
		}
	}
	if good {
		if md.sand == nil {
			md.sand = make(map[point]struct{})
		}
		md.sand[sP] = struct{}{}
	}
	return good
}

func toPoint(s string) (*point, error) {
	dims := strings.Split(s, ",")
	if len(dims) != 2 {
		return nil, fmt.Errorf("invalid point")
	}
	var err error
	var p point
	p.x, err = strconv.Atoi(dims[0])
	if err != nil {
		return nil, fmt.Errorf("invalid point")
	}
	p.y, err = strconv.Atoi(dims[1])
	if err != nil {
		return nil, fmt.Errorf("invalid point")
	}
	return &p, nil
}

func run() error {
	md := mapData{}
	err := utils.OpenAndReadLines("input14.txt", func(s string) error {
		paths := strings.Split(s, " -> ")
		for i := range paths[1:] {
			fromP, err := toPoint(paths[i])
			if err != nil {
				return err
			}
			toP, err := toPoint(paths[i+1])
			if err != nil {
				return err
			}
			if fromP.x == toP.x {
				var start, stop int
				if fromP.y < toP.y {
					start = fromP.y
					stop = toP.y
				} else {
					start = toP.y
					stop = fromP.y
				}
				for i := start; i <= stop; i++ {
					md.setWall(fromP.x, i)
				}
			} else if fromP.y == toP.y {
				var start, stop int
				if fromP.x < toP.x {
					start = fromP.x
					stop = toP.x
				} else {
					start = toP.x
					stop = fromP.x
				}
				for i := start; i <= stop; i++ {
					md.setWall(i, fromP.y)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	iter := 0
	for {
		good := md.placeSand(true)
		if !good {
			break
		}
		iter++
	}
	fmt.Printf("Part 1: %d\n", iter)
	md2 := mapData{
		walls: md.walls,
		sand:  nil,
		minP:  md.minP,
		maxP:  md.maxP,
	}
	floor := md.maxP.y + 2
	md2.floor = &floor
	iter = 0
	for {
		iter++
		good := md2.placeSand(false)
		if !good {
			break
		}
	}
	fmt.Printf("Part 2: %d\n", iter)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
