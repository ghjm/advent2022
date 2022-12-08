package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
)

func run() error {
	matches, err := utils.OpenAndReadRegex("input4.txt", `^(\d+)-(\d+),(\d+)-(\d+)$`, true)
	if err != nil {
		return err
	}
	part1Count := 0
	part2Count := 0
	for _, m := range matches {
		mi, err := utils.StringsToInts(m, 1, 2, 3, 4)
		if err != nil {
			return err
		}
		e1min, e1max, e2min, e2max := mi[0], mi[1], mi[2], mi[3]
		if e1min > e2min {
			tmpMin := e1min
			tmpMax := e1max
			e1min = e2min
			e1max = e2max
			e2min = tmpMin
			e2max = tmpMax
		}
		if (e2min >= e1min && e2max <= e1max) || (e1min >= e2min && e1max <= e2max) {
			part1Count += 1
		}
		if e2min <= e1max {
			part2Count += 1
		}
	}
	fmt.Printf("Part 1: %d\n", part1Count)
	fmt.Printf("Part 2: %d\n", part2Count)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
