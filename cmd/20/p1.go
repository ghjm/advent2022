package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
)

type data struct {
	numbers []int64
	nodes   map[int64]*node
	zeroIdx int64
}

type node struct {
	value int64
	next  *node
	prev  *node
}

func (d *data) printList() {
	np := d.nodes[0]
	for {
		fmt.Printf("%d", np.value)
		if np.next == d.nodes[0] {
			break
		} else {
			fmt.Printf(", ")
			np = np.next
		}
	}
	fmt.Printf("\n")
}

func (d *data) move(n int64) {
	nn := d.nodes[n]
	if nn.value == 0 {
		return
	}
	nn.prev.next = nn.next
	nn.next.prev = nn.prev
	count := nn.value
	count = count % int64(len(d.numbers)-1)
	switch {
	case nn.value > 0:
		nt := nn
		for i := int64(0); i < count; i++ {
			nt = nt.next
		}
		nn.prev = nt
		nn.next = nt.next
		nt.next.prev = nn
		nt.next = nn
	case nn.value < 0:
		nt := nn
		for i := int64(0); i > count; i-- {
			nt = nt.prev
		}
		nn.next = nt
		nn.prev = nt.prev
		nt.prev.next = nn
		nt.prev = nn
	}
}

func (d *data) init() {
	d.nodes = make(map[int64]*node)
	for i, ni := range d.numbers {
		d.nodes[int64(i)] = &node{
			value: ni,
		}
		if ni == 0 {
			d.zeroIdx = int64(i)
		}
	}
	for i := range d.numbers {
		if i == 0 {
			d.nodes[int64(i)].prev = d.nodes[int64(len(d.numbers)-1)]
		} else {
			d.nodes[int64(i)].prev = d.nodes[int64(i-1)]
		}
		if i == len(d.numbers)-1 {
			d.nodes[int64(i)].next = d.nodes[0]
		} else {
			d.nodes[int64(i)].next = d.nodes[int64(i+1)]
		}
	}
}

func (d *data) result() int64 {
	var res int64
	finalValues := make([]int64, 0, len(d.numbers))
	np := d.nodes[d.zeroIdx]
	for {
		finalValues = append(finalValues, np.value)
		np = np.next
		if np == d.nodes[d.zeroIdx] {
			break
		}
	}
	for i := 1000; i <= 3000; i += 1000 {
		idx := i % len(d.numbers)
		res += finalValues[idx]
	}
	return res
}

func run() error {
	var d data
	err := utils.OpenAndReadLines("input20.txt", func(s string) error {
		d.numbers = append(d.numbers, utils.MustAtoi64(s))
		return nil
	})
	if err != nil {
		return err
	}
	d.init()
	for i := range d.numbers {
		d.move(int64(i))
	}
	fmt.Printf("Part 1: %d\n", d.result())

	var d2 data
	for _, n := range d.numbers {
		d2.numbers = append(d2.numbers, n*811589153)
	}
	d2.init()
	for c := 0; c < 10; c++ {
		for i := range d.numbers {
			d2.move(int64(i))
		}
	}
	fmt.Printf("Part 2: %d\n", d2.result())

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
