package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
)

var scores = map[string]int{
	"A X": 1 + 3,
	"A Y": 2 + 6,
	"A Z": 3 + 0,
	"B X": 1 + 0,
	"B Y": 2 + 3,
	"B Z": 3 + 6,
	"C X": 1 + 6,
	"C Y": 2 + 0,
	"C Z": 3 + 3,
}

var encScores = map[string]string{
	"A X": "A Z", // lose vs rock by playing scissors
	"A Y": "A X", // tie vs rock by playing rock
	"A Z": "A Y", // win vs rock by playing paper
	"B X": "B X", // lose vs paper by playing rock
	"B Y": "B Y", // tie vs paper by playing paper
	"B Z": "B Z", // win vs paper by playing scissors
	"C X": "C Y", // lose vs scissors by playing paper
	"C Y": "C Z", // tie vs scissors by playing scissors
	"C Z": "C X", // win vs scissors by playing rock
}

func main() {
	totalScore1 := 0
	totalScore2 := 0
	utils.OpenAndReadAll("input2.txt", func(s string) error {
		score, ok := scores[s]
		if !ok {
			return fmt.Errorf("not found in scores: %s", s)
		}
		totalScore1 += score
		encScore, ok := encScores[s]
		if !ok {
			return fmt.Errorf("not found in encScores: %s", s)
		}
		score, ok = scores[encScore]
		if !ok {
			return fmt.Errorf("not found in scores (p2): %s", encScore)
		}
		totalScore2 += score
		return nil
	})
	fmt.Printf("Part 1: %d\n", totalScore1)
	fmt.Printf("Part 2: %d\n", totalScore2)
}
