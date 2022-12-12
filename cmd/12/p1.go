package main

import (
	"fmt"
	"github.com/beefsack/go-astar"
	"github.com/ghjm/advent2022/pkg/utils"
	"math"
	"os"
)

type heightMap struct {
	tiles [][]*tile
	start *tile
	end   *tile
}

type tile struct {
	hm     *heightMap
	row    int
	col    int
	height int
}

type direction struct {
	dx int
	dy int
}

var directions = []direction{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

func (t *tile) PathNeighbors() []astar.Pather {
	var results []astar.Pather
	for _, d := range directions {
		r := t.row + d.dy
		c := t.col + d.dx
		if (r == t.row && c == t.col) || (r < 0 || r >= len(t.hm.tiles)) || (c < 0 || c >= len(t.hm.tiles[0])) {
			continue
		}
		tt := t.hm.tiles[r][c]
		if tt.height <= t.height+1 {
			results = append(results, tt)
		}
	}
	return results
}

func (t *tile) PathNeighborCost(to astar.Pather) float64 {
	return 1.0
}

func (t *tile) PathEstimatedCost(to astar.Pather) float64 {
	return math.Abs(float64(t.row-to.(*tile).row)) + math.Abs(float64(t.col-to.(*tile).col))
}

func run() error {
	hm := heightMap{}
	rowNo := 0
	var possibleStarts []*tile
	err := utils.OpenAndReadLines("input12.txt", func(s string) error {
		var curRow []*tile
		for i, c := range s {
			heightCh := c
			newTile := &tile{
				hm:  &hm,
				row: rowNo,
				col: i,
			}
			switch c {
			case 'S':
				hm.start = newTile
				heightCh = 'a'
			case 'E':
				hm.end = newTile
				heightCh = 'z'
			}
			newTile.height = int(heightCh) - int('a')
			curRow = append(curRow, newTile)
			if newTile.height == 0 {
				possibleStarts = append(possibleStarts, newTile)
			}
		}
		hm.tiles = append(hm.tiles, curRow)
		rowNo++
		return nil
	})
	if err != nil {
		return err
	}
	_, distance, found := astar.Path(hm.start, hm.end)
	if !found {
		return fmt.Errorf("no solution")
	}
	fmt.Printf("Part 1: %d\n", int(distance))
	bestDistance := distance
	for _, ps := range possibleStarts {
		_, distance, found = astar.Path(ps, hm.end)
		if found && distance < bestDistance {
			bestDistance = distance
		}
	}
	fmt.Printf("Part 2: %d\n", int(bestDistance))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
