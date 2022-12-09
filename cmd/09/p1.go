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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func follow(h point, t *point) {
	if t.y == h.y {
		if h.x-t.x > 1 {
			t.x = h.x - 1
		} else if t.x-h.x > 1 {
			t.x = h.x + 1
		}
	} else if t.x == h.x {
		if h.y-t.y > 1 {
			t.y = h.y - 1
		} else if t.y-h.y > 1 {
			t.y = h.y + 1
		}
	} else if abs(t.x-h.x) > 1 || abs(t.y-h.y) > 1 {
		if h.x > t.x {
			t.x += 1
		} else {
			t.x -= 1
		}
		if h.y > t.y {
			t.y += 1
		} else {
			t.y -= 1
		}
	}
}

func run() error {
	h := point{0, 0}
	t1 := point{0, 0}
	var t2 [9]point
	for i := 0; i < 9; i++ {
		t2[i] = point{0, 0}
	}
	directions := map[string]point{
		"R": {1, 0},
		"L": {-1, 0},
		"U": {0, -1},
		"D": {0, 1},
	}
	visited1 := make(map[point]struct{})
	visited1[t1] = struct{}{}
	visited2 := make(map[point]struct{})
	visited2[t1] = struct{}{}
	err := utils.OpenAndReadLines("input9.txt", func(s string) error {
		sc := strings.Split(s, " ")
		count, err := strconv.Atoi(sc[1])
		if err != nil {
			return err
		}
		for i := 0; i < count; i++ {
			dir := directions[sc[0]]
			h.x += dir.x
			h.y += dir.y
			follow(h, &t1)
			follow(h, &t2[0])
			for i := 1; i < 9; i++ {
				follow(t2[i-1], &t2[i])
			}
			visited1[t1] = struct{}{}
			visited2[t2[8]] = struct{}{}
			minX, maxX, minY, maxY := h.x, h.x, h.y, h.y
			for j := 0; j < 9; j++ {
				if t2[j].x < minX {
					minX = t2[j].x
				}
				if t2[j].x > maxX {
					maxX = t2[j].x
				}
				if t2[j].y < minY {
					minY = t2[j].y
				}
				if t2[j].x > maxY {
					maxY = t2[j].y
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Part 1: %d\n", len(visited1))
	fmt.Printf("Part 2: %d\n", len(visited2))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
