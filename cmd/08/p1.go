package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"strconv"
)

type treeMap struct {
	data [][]int
}

func (tm *treeMap) height() int {
	return len(tm.data)
}

func (tm *treeMap) width() int {
	return len(tm.data[0])
}

func (tm *treeMap) treeID(r, c int) int {
	return r*tm.width() + c
}

func (tm *treeMap) treePos(id int) (int, int) {
	r := id / tm.width()
	c := id % tm.width()
	return r, c
}

func (tm *treeMap) getVisibleTrees(startRow, startCol, stepRow, stepCol int) []int {
	r, c := startRow, startCol
	maxTree := -1
	var results []int
	for r >= 0 && r < tm.height() && c >= 0 && c < tm.width() {
		curTree := tm.data[r][c]
		if curTree > maxTree {
			results = append(results, tm.treeID(r, c))
			maxTree = curTree
		}
		r += stepRow
		c += stepCol
	}
	return results
}

func (tm *treeMap) getScenicScore(r, c int) int {
	score := 1
	houseHeight := tm.data[r][c]
	dr, dc := []int{0, 1, 0, -1}, []int{1, 0, -1, 0}
	for i := 0; i < len(dr); i++ {
		cr, cc := r, c
		sightDistance := 0
		for {
			cr += dr[i]
			cc += dc[i]
			if cr < 0 || cc < 0 || cr >= tm.height() || cc >= tm.width() {
				break
			}
			if tm.data[cr][cc] >= houseHeight {
				sightDistance += 1
				break
			}
			sightDistance += 1
		}
		score *= sightDistance
	}
	return score
}

func run() error {
	var tm treeMap
	err := utils.OpenAndReadLines("input8.txt", func(s string) error {
		var treeRow []int
		for i := range s {
			h, err := strconv.Atoi(s[i : i+1])
			if err != nil {
				return err
			}
			treeRow = append(treeRow, h)
		}
		tm.data = append(tm.data, treeRow)
		return nil
	})
	if err != nil {
		return err
	}
	visibleTrees := make(map[int]struct{})
	for r := 0; r < tm.height(); r++ {
		trees := tm.getVisibleTrees(r, 0, 0, 1)
		for _, t := range trees {
			visibleTrees[t] = struct{}{}
		}
		trees = tm.getVisibleTrees(r, tm.width()-1, 0, -1)
		for _, t := range trees {
			visibleTrees[t] = struct{}{}
		}
	}
	for c := 0; c < tm.width(); c++ {
		trees := tm.getVisibleTrees(0, c, 1, 0)
		for _, t := range trees {
			visibleTrees[t] = struct{}{}
		}
		trees = tm.getVisibleTrees(tm.height()-1, c, -1, 0)
		for _, t := range trees {
			visibleTrees[t] = struct{}{}
		}
	}
	fmt.Printf("Part 1: %d\n", len(visibleTrees))
	bestScore := 0
	for r := 0; r < tm.height(); r++ {
		for c := 0; c < tm.width(); c++ {
			score := tm.getScenicScore(r, c)
			if score > bestScore {
				bestScore = score
			}
		}
	}
	fmt.Printf("Part 2: %d\n", bestScore)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
