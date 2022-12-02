package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"sort"
	"strconv"
)

func readCals() ([][]int, error) {
	elves := [][]int{{}}
	err := utils.OpenAndReadAll("input1.txt", func(s string) error {
		if s == "" {
			elves = append(elves, []int{})
		} else {
			cals, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			elves[len(elves)-1] = append(elves[len(elves)-1], cals)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return elves, nil
}

func main() {
	elves, err := readCals()
	if err != nil {
		fmt.Printf("Errror: %s", err)
		os.Exit(1)
	}
	var elfTotals []int
	for _, elf := range elves {
		elfSum := 0
		for _, cal := range elf {
			elfSum += cal
		}
		elfTotals = append(elfTotals, elfSum)
	}
	sort.Ints(elfTotals)
	fmt.Printf("Part 1: %d\n", elfTotals[len(elfTotals)-1])
	fmt.Printf("Part 2: %d\n", elfTotals[len(elfTotals)-1]+elfTotals[len(elfTotals)-2]+elfTotals[len(elfTotals)-3])
}
