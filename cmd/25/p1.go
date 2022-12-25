package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
)

var snafuDigits = map[rune]int64{'0': 0, '1': 1, '2': 2, '-': -1, '=': -2}

type digitInfo struct {
	char  rune
	value int64
}

var digitsToSnafu = map[int64]digitInfo{0: {'0', 0}, 1: {'1', 1}, 2: {'2', 2}, 3: {'=', -2}, 4: {'-', -1}}

func snafuToInt(snafu string) int64 {
	var a int64
	var pow int64 = 1
	for i := len(snafu) - 1; i >= 0; i-- {
		d, ok := snafuDigits[rune(snafu[i])]
		if !ok {
			panic("invalid snafu digit")
		}
		a += d * pow
		pow *= 5
	}
	return a
}

func intToSnafu(v int64) string {
	var s string
	for v != 0 {
		di := digitsToSnafu[utils.Mod64(v, 5)]
		s = fmt.Sprintf("%c%s", di.char, s)
		v -= di.value
		v = v / 5
	}
	if s == "" {
		s = "0"
	}
	return s
}

func run() error {
	var sum int64
	err := utils.OpenAndReadLines("input25.txt", func(s string) error {
		sum += snafuToInt(s)
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Part 1: %s\n", intToSnafu(sum))
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
