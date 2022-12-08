package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
)

func run() error {
	message, err := utils.OpenAndReadAll("input6.txt")
	if err != nil {
		return err
	}
	for i := 3; i < len(message); i++ {
		m := make(map[byte]struct{})
		for j := 0; j < 4; j++ {
			m[message[i-j]] = struct{}{}
		}
		if len(m) == 4 {
			fmt.Printf("Part 1: %d\n", i+1)
			break
		}
	}
	for i := 13; i < len(message); i++ {
		m := make(map[byte]struct{})
		for j := 0; j < 14; j++ {
			m[message[i-j]] = struct{}{}
		}
		if len(m) == 14 {
			fmt.Printf("Part 2: %d\n", i+1)
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
