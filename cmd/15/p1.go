package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"regexp"
	"strconv"
)

type point struct {
	x int
	y int
}

type sensor struct {
	pos           point
	nearestBeacon point
	distance      int
}

type mapData struct {
	sensors []sensor
	minP    point
	maxP    point
}

func (md *mapData) getCh(x, y int) rune {
	p := point{x, y}
	for _, sn := range md.sensors {
		if sn.pos == p {
			return 'S'
		} else if sn.nearestBeacon == p {
			return 'B'
		}
	}
	for _, sn := range md.sensors {
		if manhattanDistance(p, sn.pos) <= sn.distance {
			return '#'
		}
	}
	return '.'
}

func (md *mapData) adjustBounds(sn sensor) {
	minX := sn.nearestBeacon.x
	if sn.pos.x < minX {
		minX = sn.pos.x
	}
	minX -= sn.distance
	if minX < md.minP.x {
		md.minP.x = minX
	}

	minY := sn.nearestBeacon.y
	if sn.pos.y < minY {
		minY = sn.pos.y
	}
	minY -= sn.distance
	if minY < md.minP.y {
		md.minP.y = minY
	}

	maxX := sn.nearestBeacon.x
	if sn.pos.x > maxX {
		maxX = sn.pos.x
	}
	maxX += sn.distance
	if maxX > md.maxP.x {
		md.maxP.x = maxX
	}

	maxY := sn.nearestBeacon.y
	if sn.pos.y > maxY {
		maxY = sn.pos.y
	}
	maxY += sn.distance
	if maxY > md.maxP.y {
		md.maxP.y = maxY
	}

}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func manhattanDistance(p1, p2 point) int {
	return abs(p1.x-p2.x) + abs(p1.y-p2.y)
}

func run() error {
	//const (
	//	filename    = "input15_test.txt"
	//	checkRow    = 10
	//	searchBound = 20
	//)
	const (
		filename    = "input15.txt"
		checkRow    = 2000000
		searchBound = 4000000
	)

	md := mapData{}
	r := regexp.MustCompile(`^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)$`)
	err := utils.OpenAndReadLines(filename, func(s string) error {
		m := r.FindStringSubmatch(s)
		sn := sensor{
			pos:           point{mustAtoi(m[1]), mustAtoi(m[2])},
			nearestBeacon: point{mustAtoi(m[3]), mustAtoi(m[4])},
		}
		sn.distance = manhattanDistance(sn.pos, sn.nearestBeacon)
		md.sensors = append(md.sensors, sn)
		return nil
	})
	if err != nil {
		return err
	}
	md.minP = point{md.sensors[0].pos.x, md.sensors[0].pos.y}
	md.maxP = point{md.sensors[0].pos.x, md.sensors[0].pos.y}
	for _, sn := range md.sensors {
		md.adjustBounds(sn)
	}

	count := 0
	for x := md.minP.x; x <= md.maxP.x; x++ {
		ch := md.getCh(x, checkRow)
		if ch == '#' {
			count++
		}
	}
	fmt.Printf("Part 1: %d\n", count)

	var found *point
	for y := 0; y <= searchBound; y++ {
		for x := 0; x <= searchBound; x++ {
			p := point{x, y}
			var skip bool
			for _, sn := range md.sensors {
				if sn.pos == p {
					skip = true
					break
				} else if sn.nearestBeacon == p {
					skip = true
					break
				}
				dy := abs(y - sn.pos.y)
				if dy+abs(x-sn.pos.x) <= sn.distance {
					// increment x beyond the blocked region
					x = sn.pos.x + sn.distance - dy
					skip = true
					break
				}
			}
			if skip {
				continue
			}
			found = &point{x, y}
			break
		}
		if found != nil {
			break
		}
	}
	fmt.Printf("Part 2: %d\n", found.x*4000000+found.y)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
