package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"github.com/zyedidia/generic/set"
	"os"
)

func priority(c rune) (int, error) {
	if c >= 'a' && c <= 'z' {
		return int(c) - int('a') + 1, nil
	} else if c >= 'A' && c <= 'Z' {
		return int(c) - int('A') + 27, nil
	} else {
		return 0, fmt.Errorf("unknown char")
	}
}

func part1(elves []string) (int, error) {
	prioSum := 0
	for _, elf := range elves {
		s1 := set.NewMapset([]rune(elf[0 : len(elf)/2])...)
		s2 := set.NewMapset([]rune(elf[len(elf)/2 : len(elf)])...)
		si := s1.Intersection(s2)
		if si.Size() != 1 {
			return 0, fmt.Errorf("wrong number of matches")
		}
		p, err := priority(si.Keys()[0])
		if err != nil {
			return 0, err
		}
		prioSum += p
	}
	return prioSum, nil
}

func part2(elves []string) (int, error) {
	var elfGroups [][3]string
	for i := 0; i < len(elves); i += 3 {
		elfGroups = append(elfGroups, [3]string{elves[i], elves[i+1], elves[i+2]})
	}
	prioSum := 0
	for _, eg := range elfGroups {
		s1 := set.NewMapset([]rune(eg[0])...)
		s2 := set.NewMapset([]rune(eg[1])...)
		s3 := set.NewMapset([]rune(eg[2])...)
		si := s1.Intersection(s2).Intersection(s3)
		if si.Size() != 1 {
			return 0, fmt.Errorf("wrong number of matches")
		}
		p, err := priority(si.Keys()[0])
		if err != nil {
			return 0, err
		}
		prioSum += p
	}
	return prioSum, nil
}

func main() {
	var elves []string
	err := utils.OpenAndReadAll("input3.txt", func(s string) error {
		elves = append(elves, s)
		return nil
	})
	p1, err := part1(elves)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Part 1: %d\n", p1)
	p2, err := part2(elves)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Part 1: %d\n", p2)
}
