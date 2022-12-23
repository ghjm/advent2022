package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"math"
	"os"
)

type point struct {
	x, y int
}

type elf struct {
	d        *data
	loc      point
	proposal *point
}

type direction struct {
	name  byte
	dir   point
	check []point
}

var directions = []direction{
	{'N', point{0, -1}, []point{{-1, -1}, {0, -1}, {1, -1}}},
	{'S', point{0, 1}, []point{{-1, 1}, {0, 1}, {1, 1}}},
	{'W', point{-1, 0}, []point{{-1, -1}, {-1, 0}, {-1, 1}}},
	{'E', point{1, 0}, []point{{1, -1}, {1, 0}, {1, 1}}},
}

func (e *elf) anyNeighbors() bool {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dy == 0 && dx == 0 {
				continue
			}
			_, ok := e.d.elves[point{e.loc.x + dx, e.loc.y + dy}]
			if ok {
				return true
			}
		}
	}
	return false
}

func (e *elf) makeProposal() {
	if !e.anyNeighbors() {
		e.proposal = nil
		return
	}
	for i := 0; i < 4; i++ {
		dir := directions[(e.d.dirStart+i)%4]
		elfFound := false
		for _, p := range dir.check {
			_, ok := e.d.elves[point{e.loc.x + p.x, e.loc.y + p.y}]
			if ok {
				elfFound = true
				break
			}
		}
		if !elfFound {
			p := point{e.loc.x + dir.dir.x, e.loc.y + dir.dir.y}
			e.proposal = &p
			return
		}
	}
	e.proposal = nil
}

type data struct {
	elves    map[point]*elf
	dirStart int
}

func (d *data) runStep() bool {
	for _, e := range d.elves {
		e.makeProposal()
	}
	props := make(map[point]*elf)
	for _, e := range d.elves {
		if e.proposal != nil {
			_, ok := props[*e.proposal]
			if ok {
				props[*e.proposal] = nil
			} else {
				props[*e.proposal] = e
			}
		}
	}
	anyMoved := false
	for _, e := range props {
		if e == nil || e.proposal == nil {
			continue
		}
		delete(d.elves, e.loc)
		e.loc = *e.proposal
		d.elves[*e.proposal] = e
		anyMoved = true
	}
	d.dirStart = (d.dirStart + 1) % 4
	return anyMoved
}

func (d *data) emptyTiles() int {
	minP, maxP := d.getLimits()
	return ((maxP.x - minP.x + 1) * (maxP.y - minP.y + 1)) - len(d.elves)
}

func (d *data) getLimits() (point, point) {
	minP := point{math.MaxInt, math.MaxInt}
	maxP := point{-math.MaxInt, -math.MaxInt}
	for _, e := range d.elves {
		if e.loc.x < minP.x {
			minP.x = e.loc.x
		}
		if e.loc.x > maxP.x {
			maxP.x = e.loc.x
		}
		if e.loc.y < minP.y {
			minP.y = e.loc.y
		}
		if e.loc.y > maxP.y {
			maxP.y = e.loc.y
		}
	}
	return minP, maxP
}

func (d *data) printBoard() {
	minP, maxP := d.getLimits()
	for y := minP.y; y <= maxP.y; y++ {
		for x := minP.x; x <= maxP.x; x++ {
			_, ok := d.elves[point{x, y}]
			if ok {
				fmt.Printf("#")
			} else if x == 0 || y == 0 {
				fmt.Printf("+")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}

func run() error {
	d := &data{
		elves: make(map[point]*elf),
	}
	y := 0
	err := utils.OpenAndReadLines("input23.txt", func(s string) error {
		for x, c := range s {
			if c == '#' {
				d.elves[point{x, y}] = &elf{d: d, loc: point{x, y}}
			}
		}
		y++
		return nil
	})
	if err != nil {
		return err
	}
	for i := 0; i < 10; i++ {
		//fmt.Printf("--- Round %d ---- Initial Direction %c -------------\n", i+1, directions[d.dirStart].name)
		d.runStep()
		//d.printBoard()
	}
	fmt.Printf("Part 1: %d\n", d.emptyTiles())
	move := 10
	for {
		move++
		anyMoved := d.runStep()
		if !anyMoved {
			fmt.Printf("Part 2: %d\n", move)
			break
		}
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
