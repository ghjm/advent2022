package main

import (
	"fmt"
	"github.com/alexanderbez/gopq"
	"github.com/ghjm/advent2022/pkg/utils"
	"math"
	"os"
)

type point struct {
	x, y int
}

const Impossible = math.MaxInt/2

const (
	dirRight = 0
	dirDown = 1
	dirLeft = 2
	dirUp = 3
)

var dirChars = map[rune]int{'>': dirRight, 'v': dirDown, '<': dirLeft, '^': dirUp}

var directions = []point{
	{1,0},
	{0,1},
	{-1,0},
	{0,-1},
	{0,0},
}

type blizzard struct {
	initPos point
	direction int
}

type state struct {
	d *data
	pos point
	step int
}

type data struct {
	dimX, dimY int
	startX, endX int
	walls map[point]struct{}
	blizzards []blizzard
	blizzStep map[int]map[point]struct{}
}

func (d *data) calcBlizzStep(step int) map[point]struct{} {
	if d.blizzStep == nil {
		d.blizzStep = make(map[int]map[point]struct{})
	}
	bs, ok := d.blizzStep[step]
	if ok {
		return bs
	}
	bs = make(map[point]struct{})
	for _, b := range d.blizzards {
		bd := directions[b.direction]
		bp := point{b.initPos.x-1 + (bd.x * step), b.initPos.y-1 + (bd.y * step)}
		bp.x = utils.Mod(bp.x, d.dimX-2) + 1
		bp.y = utils.Mod(bp.y, d.dimY-2) + 1
		bs[bp] = struct{}{}
	}
	d.blizzStep[step] = bs
	return bs
}

func (d *data) valueAt(pos point, step int) byte {
	if pos.x < 0 || pos.x >= d.dimX || pos.y < 0 || pos.y >= d.dimY {
		return '#'
	}
	_, ok := d.walls[pos]; if ok {
		return '#'
	}
	_, ok = d.calcBlizzStep(step)[pos]; if ok {
		return '@'
	}
	return '.'
}

func (s *state) heuristic() int {
	return s.step + utils.Abs(s.d.dimY - s.pos.y) + utils.Abs(s.d.dimX - s.pos.x)
}

func (s *state) Index() (i int) { return }

func (s *state) SetIndex(_ int) {}

func (s *state) Priority(other any) bool {
	if so, ok := other.(*state); ok {
		return s.heuristic() < so.heuristic()
	}
	return false
}

func (d *data) distanceToExit(startPos point, endPos point, startStep int) int {
	toExplore := queue.NewPriorityQueue()
	toExplore.Push(&state{d, startPos, startStep})
	explored := make(map[state]struct{})
	stepsTo := make(map[point]int)
	bestSoFar := math.MaxInt
	for toExplore.Size() > 0 {
		tx, err := toExplore.Pop()
		if err != nil {
			panic(err)
		}
		tes := tx.(*state)
		if tes.step > bestSoFar {
			continue
		}
		{
			_, ok := explored[*tes]
			if ok {
				continue
			}
		}
		if tes.pos == endPos {
			e, ok := stepsTo[tes.pos]
			if !ok || e > tes.step {
				stepsTo[tes.pos] = tes.step
			}
			if tes.step < bestSoFar {
				bestSoFar = tes.step
			}
			continue
		}
		explored[*tes] = struct{}{}
		for _, dir := range directions {
			np := point{tes.pos.x+dir.x, tes.pos.y+dir.y}
			e, ok := stepsTo[np]
			if (!ok || tes.step < e) && tes.heuristic() < bestSoFar {
				if d.valueAt(np, tes.step+1) == '.' {
					toExplore.Push(&state{d, np, tes.step + 1})
				}
			}
		}
	}
	return stepsTo[endPos]
}

func (d *data) printMap(step int) {
	for y := 0; y < d.dimY; y++ {
		for x := 0; x < d.dimX; x++ {
			fmt.Printf("%c", d.valueAt(point{x, y}, step))
		}
		fmt.Printf("\n")
	}
}

func run() error {
	d := data{
		walls:          make(map[point]struct{}),
	}
	y := 0
	err := utils.OpenAndReadLines("input24.txt", func(s string) error {
		if len(s) > d.dimX {
			d.dimX = len(s)
		}
		for x, ch := range s {
			switch ch {
			case '#':
				d.walls[point{x,y}] = struct{}{}
			case '>', '<', '^', 'v':
				d.blizzards = append(d.blizzards, blizzard{
					initPos:   point{x, y},
					direction: dirChars[ch],
				})
			}
		}
		y++
		return nil
	})
	if err != nil {
		return err
	}
	d.dimY = y
	for x := 0; x < d.dimX; x++ {
		if d.valueAt(point{x, 0}, 0) == '.' {
			d.startX = x
		}
		if d.valueAt(point{x, d.dimY-1}, 0) == '.' {
			d.endX = x
		}
	}
	startP := point{d.startX, 0}
	endP := point{d.endX, d.dimY-1}
	d1 := d.distanceToExit(startP, endP, 0)
	fmt.Printf("Part 1: %d\n", d1)
	d2 := d.distanceToExit(endP, startP, d1)
	d3 := d.distanceToExit(startP, endP, d2)
	fmt.Printf("Part 2: %d\n", d3)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
