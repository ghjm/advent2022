package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
)

type point struct {
	x, y int
}

var directions = []point{
	{1, 0},
	{0, 1},
	{-1, 0},
	{0, -1},
}

type cubeFaceInfo struct {
	face point
	rot  byte
}

var cubeFaces = map[cubeFaceInfo]cubeFaceInfo{
	{point{1, 0}, 'R'}: {point{2, 0}, 'N'},
	{point{1, 0}, 'U'}: {point{0, 3}, 'R'},
	{point{1, 0}, 'L'}: {point{0, 2}, 'F'},
	{point{1, 0}, 'D'}: {point{1, 1}, 'N'},

	{point{2, 0}, 'R'}: {point{1, 2}, 'F'},
	{point{2, 0}, 'U'}: {point{0, 3}, 'N'},
	{point{2, 0}, 'L'}: {point{1, 0}, 'N'},
	{point{2, 0}, 'D'}: {point{1, 1}, 'R'},

	{point{1, 1}, 'R'}: {point{0, 2}, 'L'},
	{point{1, 1}, 'U'}: {point{1, 0}, 'N'},
	{point{1, 1}, 'L'}: {point{0, 2}, 'L'},
	{point{1, 1}, 'D'}: {point{1, 2}, 'N'},

	{point{0, 2}, 'R'}: {point{1, 2}, 'N'},
	{point{0, 2}, 'U'}: {point{1, 1}, 'R'},
	{point{0, 2}, 'L'}: {point{1, 0}, 'F'},
	{point{0, 2}, 'D'}: {point{0, 3}, 'N'},

	{point{1, 2}, 'R'}: {point{2, 0}, 'F'},
	{point{1, 2}, 'U'}: {point{1, 1}, 'N'},
	{point{1, 2}, 'L'}: {point{0, 2}, 'N'},
	{point{1, 2}, 'D'}: {point{0, 3}, 'R'},

	{point{0, 3}, 'R'}: {point{1, 2}, 'L'},
	{point{0, 3}, 'U'}: {point{0, 2}, 'N'},
	{point{0, 3}, 'L'}: {point{1, 0}, 'L'},
	{point{0, 3}, 'D'}: {point{2, 0}, 'N'},
}

type data struct {
	board []string
	steps string
}

func (d *data) valueAt(p point) byte {
	if p.y < 0 || p.y >= len(d.board) {
		return ' '
	}
	if p.x < 0 || p.x >= len(d.board[p.y]) {
		return ' '
	}
	return d.board[p.y][p.x]
}

func (p point) cubeFace() point {
	return point{p.x / 50, p.y / 50}
}

func (d *data) move(loc point, facing int, cube bool) (point, int) {
	dir := directions[facing]
	if cube {
		oldFace := loc.cubeFace()
		naiveNewLoc := point{loc.x + dir.x, loc.y + dir.y}
		naiveNewFace := naiveNewLoc.cubeFace()
		if oldFace != naiveNewFace {
			newInfo, ok := cubeFaces[cubeFaceInfo{oldFace, "RDLU"[facing]}]
			if !ok {
				panic("invalid cube face move")
			}
			switch newInfo.rot {
			case 'N':
				switch facing {
				case 0: // R
					return point{newInfo.face.x*50 + 49, newInfo.face.y*50 + loc.y%50}, facing
				case 1: // D
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y * 50}, facing
				case 2: // L
					return point{newInfo.face.x * 50, newInfo.face.y*50 + loc.y%50}, facing
				case 3: // U
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y*50 + 49}, facing
				}
			case 'F':
				switch facing {
				case 0: // R
					return point{newInfo.face.x * 50, newInfo.face.y*50 + loc.y%50}, (facing + 2) % 4
				case 1: // D
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y*50 + 49}, (facing + 2) % 4
				case 2: // L
					return point{newInfo.face.x*50 + 49, newInfo.face.y*50 + loc.y%50}, (facing + 2) % 4
				case 3: // U
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y * 50}, (facing + 2) % 4
				}
			case 'R':
				switch facing {
				case 0: // R
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y*50 + 49}, (facing + 1) % 4
				case 1: // D
					return point{newInfo.face.x*50 + 49, newInfo.face.y*50 + loc.y%50}, (facing + 1) % 4
				case 2: // L
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y * 50}, (facing + 1) % 4
				case 3: // U
					return point{newInfo.face.x * 50, newInfo.face.y*50 + loc.y%50}, (facing + 1) % 4
				}
			case 'L':
				switch facing {
				case 0: // R
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y * 50}, (facing + 3) % 4
				case 1: // D
					return point{newInfo.face.x * 50, newInfo.face.y*50 + loc.y%50}, (facing + 3) % 4
				case 2: // L
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y*50 + 49}, (facing + 3) % 4
				case 3: // U
					return point{newInfo.face.x*50 + loc.x%50, newInfo.face.y * 50}, (facing + 3) % 4
				}
			}
		} else {
			return naiveNewLoc, facing
		}
	} else {
		res := point{loc.x, loc.y}
		for {
			res.x += dir.x
			res.y += dir.y
			if dir.y != 0 {
				switch {
				case res.y >= len(d.board):
					res.y = 0
				case res.y < 0:
					res.y = len(d.board) - 1
				}
			}
			if dir.x != 0 {
				switch {
				case res.x >= len(d.board[res.y]):
					res.x = 0
				case res.x < 0:
					res.x = len(d.board[res.y]) - 1
				}
			}
			if d.valueAt(res) != ' ' {
				return res, 0
			}
		}
	}
	panic("no move")
}

type sim struct {
	d        *data
	location point
	facing   int
	history  map[point]int
}

func (s *sim) step(cube bool) {
	newP, rot := s.d.move(s.location, s.facing, cube)
	c := s.d.valueAt(newP)
	if c == '#' {
		return
	}
	if c == '.' {
		s.facing = rot
		s.location = newP
		s.history[newP] = s.facing
		return
	}
}

func (s *sim) runSteps(cube bool) {
	instr := s.d.steps
	for len(instr) > 0 {
		if instr[0] >= '0' && instr[0] <= '9' {
			p := 0
			for p < len(instr) && instr[p] >= '0' && instr[p] <= '9' {
				p++
			}
			num := instr[:p]
			instr = instr[p:]
			n := utils.MustAtoi(num)
			for i := 0; i < n; i++ {
				s.step(cube)
			}
		} else {
			cmd := instr[0]
			instr = instr[1:]
			switch cmd {
			case 'R':
				s.facing++
				if s.facing >= len(directions) {
					s.facing = 0
				}
			case 'L':
				s.facing--
				if s.facing < 0 {
					s.facing = len(directions) - 1
				}
			default:
				panic("unknown cmd")
			}
			s.history[s.location] = s.facing
		}
	}
}

func (s *sim) printHistory() {
	for y := 0; y < len(s.d.board); y++ {
		for x := 0; x < len(s.d.board[y]); x++ {
			p := point{x, y}
			v := s.d.valueAt(p)
			h, ok := s.history[p]
			if ok && v != '.' {
				panic(fmt.Sprintf("illegal move in history at x=%d,y=%d", p.x, p.y))
			}
			if ok {
				fmt.Printf("%c", ">v<^"[h])
			} else {
				fmt.Printf("%c", s.d.valueAt(p))
			}
		}
		fmt.Printf("\n")
	}
}

func initSim(d *data) *sim {
	s := &sim{
		d:        d,
		location: point{0, 0},
		facing:   0,
		history:  make(map[point]int),
	}
	for d.valueAt(s.location) == ' ' {
		s.location.x++
	}
	return s
}

func run() error {
	boardDone := false
	d := &data{}
	err := utils.OpenAndReadLines("input22.txt", func(s string) error {
		if s == "" {
			boardDone = true
		} else if boardDone {
			d.steps = s
		} else {
			d.board = append(d.board, s)
		}
		return nil
	})
	if err != nil {
		return err
	}
	s := initSim(d)
	s.runSteps(false)
	fmt.Printf("Part 1: %d\n", 1000*(s.location.y+1)+4*(s.location.x+1)+s.facing)
	s = initSim(d)
	s.runSteps(true)
	fmt.Printf("Part 2: %d\n", 1000*(s.location.y+1)+4*(s.location.x+1)+s.facing)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
