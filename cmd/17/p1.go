package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"strings"
)

var shapes = [][][]rune{
	{
		{'#', '#', '#', '#'},
	},
	{
		{0, '#', 0},
		{'#', '#', '#'},
		{0, '#', 0},
	},
	{
		{'#', '#', '#'},
		{0, 0, '#'},
		{0, 0, '#'},
	},
	{
		{'#'},
		{'#'},
		{'#'},
		{'#'},
	},
	{
		{'#', '#'},
		{'#', '#'},
	},
}

type mapData struct {
	jets     string
	data     [][7]rune
	curShape int
	curJet   int
	shapeX   int
	shapeY   int
	stateIds []int64
	heights  []int
}

func (md *mapData) collisionDetect(newX, newY int) bool {
	if newY < 0 || newX < 0 || newX > 7-len(shapes[md.curShape][0]) {
		return true
	}
	for y := 0; y < len(shapes[md.curShape]); y++ {
		for x := 0; x < len(shapes[md.curShape][0]); x++ {
			if shapes[md.curShape][y][x] == '#' && md.charAtNoShape(newX+x, newY+y) != '.' {
				return true
			}
		}
	}
	return false
}

func (md *mapData) charAtNoShape(x, y int) rune {
	var ret rune
	if ret == 0 && y >= 0 && y < len(md.data) && x >= 0 && x < 7 {
		ret = md.data[y][x]
	}
	if ret == 0 {
		ret = '.'
	}
	return ret
}

func (md *mapData) charAt(x, y int) rune {
	if y >= md.shapeY && y < md.shapeY+len(shapes[md.curShape]) &&
		x >= md.shapeX && x < md.shapeX+len(shapes[md.curShape][0]) {
		shape := shapes[md.curShape]
		if shape[y-md.shapeY][x-md.shapeX] == '#' {
			return '@'
		}
	}
	return md.charAtNoShape(x, y)
}

func (md *mapData) printMap() {
	maxY := len(md.data) - 1
	sMax := md.shapeY + len(shapes[md.curShape]) - 1
	if sMax > maxY {
		maxY = sMax
	}
	for y := maxY; y >= 0; y-- {
		fmt.Printf("|")
		for x := 0; x < 7; x++ {
			fmt.Printf("%c", md.charAt(x, y))
		}
		fmt.Printf("|\n")
	}
	fmt.Printf("+-------+\n\n")
}

func (md *mapData) simStep() {
	md.shapeY = len(md.data) + 3
	md.shapeX = 2
	for {
		jet := md.jets[md.curJet]
		if jet == '<' {
			if !md.collisionDetect(md.shapeX-1, md.shapeY) {
				md.shapeX--
			}
		} else if jet == '>' {
			if !md.collisionDetect(md.shapeX+1, md.shapeY) {
				md.shapeX++
			}
		} else {
			panic("invalid jet value")
		}
		md.curJet = (md.curJet + 1) % len(md.jets)
		if md.collisionDetect(md.shapeX, md.shapeY-1) {
			break
		}
		md.shapeY--
	}
	var stateId int64
	stateId = 1
	var statePower int64
	statePower = 1
	for y := 0; y < len(shapes[md.curShape]); y++ {
		dr := md.shapeY + y
		for dr >= len(md.data) {
			md.data = append(md.data, [7]rune{})
		}
		for x := 0; x < len(shapes[md.curShape][0]); x++ {
			sc := shapes[md.curShape][y][x]
			if sc != 0 {
				md.data[md.shapeY+y][md.shapeX+x] = sc
			}
			var stateMul int64
			if md.data[md.shapeY+y][md.shapeX+x] != 0 {
				stateMul = 1
			}
			stateId += stateMul * statePower
			statePower = statePower << 1
		}
	}
	stateId *= int64(md.curJet+1) * 3
	stateId *= int64(md.curShape+1) * 5
	md.stateIds = append(md.stateIds, stateId)
	md.heights = append(md.heights, len(md.data))
	md.curShape = (md.curShape + 1) % len(shapes)
}

func (md *mapData) continueSim(steps int) {
	for i := 0; i < steps; i++ {
		md.simStep()
	}
}

func (md *mapData) simulate(steps int) {
	md.data = make([][7]rune, 0)
	md.curShape = 0
	md.continueSim(steps)
}

func (md *mapData) checkPeriodicity() (bool, int, int) {
	st := make(map[int64]int)
	for i, v := range md.stateIds {
		pi, ok := st[v]
		if ok {
			good := true
			for d := 0; d < i-pi; d++ {
				if i+d >= len(md.stateIds) {
					good = false
					break
				}
				if md.stateIds[pi+d] != md.stateIds[i+d] {
					good = false
					break
				}
			}
			if good {
				return true, pi, i
			}
		} else {
			st[v] = i
		}
	}
	return false, 0, 0
}

func run() error {
	var md mapData
	err := utils.OpenAndReadLines("input17.txt", func(s string) error {
		md.jets = strings.TrimSpace(s)
		return nil
	})
	if err != nil {
		return err
	}
	md.simulate(2022)
	fmt.Printf("Part 1: %d\n", len(md.data))
	foundP := false
	var p1, p2 int
	for !foundP {
		foundP, p1, p2 = md.checkPeriodicity()
		if !foundP {
			md.continueSim(1000)
		}
	}
	heightSoFar := int64(md.heights[p1])
	stepsRemaining := int64(1000000000000 - p1)
	cycles := stepsRemaining / int64(p2-p1)
	heightSoFar += cycles * int64(md.heights[p2]-md.heights[p1])
	stepsRemaining -= cycles * int64(p2-p1)
	heightSoFar += int64(md.heights[p1+int(stepsRemaining)-1] - md.heights[p1])
	fmt.Printf("Part 2: %d\n", heightSoFar)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
