package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"strconv"
	"strings"
)

func run() error {
	x := 1
	values := []int{1}
	err := utils.OpenAndReadLines("input10.txt", func(s string) error {
		if strings.HasPrefix(s, "noop") {
			values = append(values, x)
		} else if strings.HasPrefix(s, "addx") {
			v, err := strconv.Atoi(s[5:])
			if err != nil {
				return err
			}
			nextX := x + v
			values = append(values, []int{x, nextX}...)
			x = nextX
		} else {
			return fmt.Errorf("unknown op")
		}
		return nil
	})
	if err != nil {
		return err
	}
	strengthSum := 0
	for cycle := 1; cycle < len(values); cycle++ {
		if cycle%40 == 20 {
			strengthSum += values[cycle-1] * cycle
		}
	}
	fmt.Printf("Part 1: %d\n", strengthSum)
	if len(values) < 240 {
		return fmt.Errorf("not enough values")
	}
	var pixels [240]bool
	for i := range pixels {
		col := i % 40
		center := values[i]
		pixels[i] = col >= center-1 && col <= center+1
	}
	fmt.Printf("Part 2:\n")
	for row := 0; row < 6; row++ {
		for col := 0; col < 40; col++ {
			if pixels[row*40+col] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
