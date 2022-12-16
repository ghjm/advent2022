package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type room struct {
	paths []string
}

type state struct {
	loc             string
	timeRemaining   int
	availableValves int
}

type mapData struct {
	rooms       map[string]room
	valves      map[string]int
	alphaValves []string
	distances   utils.Map2D[string, string, int]
	memoValues  map[state]int
}

func (md *mapData) generateAlphaValves() {
	md.alphaValves = nil
	for k := range md.valves {
		md.alphaValves = append(md.alphaValves, k)
	}
	sort.Strings(md.alphaValves)
}

func (md *mapData) generateDistances() {
	for rn, rv := range md.rooms {
		md.distances.Set(rn, rn, 0)
		for _, p := range rv.paths {
			md.distances.Set(rn, p, 1)
		}
	}
	for k := range md.rooms {
		for i := range md.rooms {
			for j := range md.rooms {
				distIJ := md.distances.GetOrDefault(i, j, math.MaxInt/2)
				distIK := md.distances.GetOrDefault(i, k, math.MaxInt/2)
				distKJ := md.distances.GetOrDefault(k, j, math.MaxInt/2)
				if distIJ > distIK+distKJ {
					md.distances.Set(i, j, distIK+distKJ)
				}
			}
		}
	}
}

func (md *mapData) allValvesBitmap() int {
	return 1<<len(md.valves) - 1
}

func (md *mapData) valueOf(loc string, timeRemaining int, availableValves int) int {
	if timeRemaining <= 0 || availableValves == 0 {
		return 0
	}
	uniqID := state{loc, timeRemaining, availableValves}
	bestV, ok := md.memoValues[uniqID]
	if ok {
		return bestV
	}
	bestV = 0
	for vi, vn := range md.alphaValves {
		vm := 1 << vi
		if availableValves&vm == 0 {
			continue
		}
		nav := availableValves & ^vm
		nt := timeRemaining - md.distances.MustGet(loc, vn) - 1
		if nt < 0 {
			continue
		}
		v := md.valves[vn]*nt + md.valueOf(vn, nt, nav)
		if v > bestV {
			bestV = v
		}
	}
	md.memoValues[uniqID] = bestV
	return bestV
}

func run() error {
	md := mapData{
		rooms:      make(map[string]room),
		valves:     make(map[string]int),
		memoValues: make(map[state]int),
	}
	r := regexp.MustCompile(`^Valve (.*) has flow rate=(\d+); tunnels? leads? to valves? (.*)$`)
	err := utils.OpenAndReadLines("input16.txt", func(s string) error {
		m := r.FindStringSubmatch(s)
		if m == nil {
			return fmt.Errorf("regexp didn't match")
		}
		rate, err := strconv.Atoi(m[2])
		if err != nil {
			return err
		}
		md.rooms[m[1]] = room{
			paths: strings.Split(m[3], ", "),
		}
		if rate > 0 {
			md.valves[m[1]] = rate
		}
		return nil
	})
	if err != nil {
		return err
	}
	md.generateAlphaValves()
	md.generateDistances()
	fmt.Printf("Part 1: %d\n", md.valueOf("AA", 30, md.allValvesBitmap()))
	bestV := 0
	for pv := 0; pv <= md.allValvesBitmap(); pv++ {
		playerValves := pv
		elephantValves := md.allValvesBitmap() & ^pv
		v := md.valueOf("AA", 26, playerValves) + md.valueOf("AA", 26, elephantValves)
		if v > bestV {
			bestV = v
		}
	}
	fmt.Printf("Part 2: %d\n", bestV)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
