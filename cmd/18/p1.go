package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"math"
	"os"
	"strings"
)

type point3D struct {
	x int
	y int
	z int
}

var faces = []point3D{
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
	{0, 0, 1},
	{0, 0, -1},
}

func run() error {
	cubes := make(map[point3D]struct{})
	minX := math.MaxInt
	maxX := -math.MaxInt
	minY := math.MaxInt
	maxY := -math.MaxInt
	minZ := math.MaxInt
	maxZ := -math.MaxInt
	err := utils.OpenAndReadLines("input18.txt", func(s string) error {
		sl := strings.Split(s, ",")
		x := utils.MustAtoi(sl[0])
		y := utils.MustAtoi(sl[1])
		z := utils.MustAtoi(sl[2])
		cubes[point3D{x, y, z}] = struct{}{}
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		if z < minZ {
			minZ = z
		}
		if z > maxZ {
			maxZ = z
		}
		return nil
	})
	if err != nil {
		return err
	}
	count := 0
	for z := minZ; z <= maxZ; z++ {
		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				p := point3D{x, y, z}
				_, ok := cubes[p]
				if !ok {
					continue
				}
				for _, f := range faces {
					fp := point3D{x + f.x, y + f.y, z + f.z}
					_, ok = cubes[fp]
					if !ok {
						count++
					}
				}
			}
		}
	}
	fmt.Printf("Part 1: %d\n", count)
	var toExplore []point3D
	airCubes := make(map[point3D]struct{})
	for x := minX; x <= maxX; x++ {
		for y := minX; x <= maxX; x++ {
			toExplore = append(toExplore, point3D{x, y, minZ})
			airCubes[point3D{x, y, minZ - 1}] = struct{}{}
			toExplore = append(toExplore, point3D{x, y, maxZ})
			airCubes[point3D{x, y, maxZ + 1}] = struct{}{}
		}
	}
	for x := minX; x <= maxX; x++ {
		for z := minZ; z <= maxZ; z++ {
			toExplore = append(toExplore, point3D{x, minY, z})
			airCubes[point3D{x, minY - 1, z}] = struct{}{}
			toExplore = append(toExplore, point3D{x, maxY, z})
			airCubes[point3D{x, maxY + 1, z}] = struct{}{}
		}
	}
	for y := minY; y <= maxY; y++ {
		for z := minZ; z <= maxZ; z++ {
			toExplore = append(toExplore, point3D{minX, y, z})
			airCubes[point3D{minX - 1, y, z}] = struct{}{}
			toExplore = append(toExplore, point3D{maxX, y, z})
			airCubes[point3D{maxX + 1, y, z}] = struct{}{}
		}
	}
	for len(toExplore) > 0 {
		p := toExplore[0]
		toExplore = toExplore[1:]
		_, ok := cubes[p]
		if ok {
			continue
		}
		_, ok = airCubes[p]
		if ok {
			continue
		}
		airCubes[p] = struct{}{}
		for _, f := range faces {
			fp := point3D{p.x + f.x, p.y + f.y, p.z + f.z}
			if fp.x < minX-1 || fp.x > maxX+1 || fp.y < minY-1 || fp.y > maxY+1 || fp.z < minZ-1 || fp.z > maxZ+1 {
				continue
			}
			_, ok = cubes[fp]
			if ok {
				continue
			}
			_, ok = airCubes[fp]
			if ok {
				continue
			}
			toExplore = append(toExplore, fp)
		}
	}
	count = 0
	for z := minZ; z <= maxZ; z++ {
		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				p := point3D{x, y, z}
				_, ok := cubes[p]
				if !ok {
					continue
				}
				for _, f := range faces {
					fp := point3D{x + f.x, y + f.y, z + f.z}
					_, ok = cubes[fp]
					if !ok {
						_, ok = airCubes[fp]
						if ok {
							count++
						}
					}
				}
			}
		}
	}
	fmt.Printf("Part 2: %d\n", count)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
