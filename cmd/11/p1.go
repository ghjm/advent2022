package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type item struct {
	worryLevel int
}

type monkey struct {
	items       []*item
	operation   func(int) int
	testDivisor int
	trueMonkey  string
	falseMonkey string
	inspections int
}

func run(part1 bool) error {
	res := map[string]*regexp.Regexp{
		"monkey":    regexp.MustCompile(`Monkey (\d+):`),
		"items":     regexp.MustCompile(`Starting items: (.*)$`),
		"operation": regexp.MustCompile(`Operation: new = (.*)$`),
		"test":      regexp.MustCompile(`Test: divisible by (\d+)`),
		"truecond":  regexp.MustCompile(`If true: throw to monkey (\d+)`),
		"falsecond": regexp.MustCompile(`If false: throw to monkey (\d+)`),
	}
	state := "monkey"
	monkeys := make(map[string]*monkey)
	var curMonkey *monkey
	err := utils.OpenAndReadLines("input11.txt", func(s string) error {
		if s == "" {
			return nil
		}
		m := res[state].FindStringSubmatch(s)
		if m == nil {
			return fmt.Errorf("match error")
		}
		switch state {
		case "monkey":
			curMonkey = &monkey{}
			monkeys[m[1]] = curMonkey
			state = "items"
		case "items":
			itemStrings := strings.Split(strings.ReplaceAll(m[1], " ", ""), ",")
			for _, itemStr := range itemStrings {
				wl, err := strconv.Atoi(itemStr)
				if err != nil {
					return fmt.Errorf("int parse error")
				}
				ni := &item{
					worryLevel: wl,
				}
				curMonkey.items = append(curMonkey.items, ni)
			}
			state = "operation"
		case "operation":
			expression, err := govaluate.NewEvaluableExpression(m[1])
			if err != nil {
				return fmt.Errorf("expression error: %w", err)
			}
			curMonkey.operation = func(wl int) int {
				result, err := expression.Evaluate(map[string]any{
					"old": wl,
				})
				if err != nil {
					panic(fmt.Sprintf("evaluation error: %s", err))
				}
				return int(result.(float64))
			}
			state = "test"
		case "test":
			divisor, err := strconv.Atoi(m[1])
			if err != nil {
				return fmt.Errorf("integer error: %w", err)
			}
			curMonkey.testDivisor = divisor
			state = "truecond"
		case "truecond":
			curMonkey.trueMonkey = m[1]
			state = "falsecond"
		case "falsecond":
			curMonkey.falseMonkey = m[1]
			state = "monkey"
		}
		return nil
	})
	if err != nil {
		return err
	}
	var monkeyKeys []string
	gcd := 1
	for k, v := range monkeys {
		monkeyKeys = append(monkeyKeys, k)
		gcd *= v.testDivisor
	}
	sort.Strings(monkeyKeys)
	var roundCount int
	if part1 {
		roundCount = 20
	} else {
		roundCount = 10000
	}
	for i := 0; i < roundCount; i++ {
		for _, mk := range monkeyKeys {
			m := monkeys[mk]
			thisItems := m.items
			m.items = nil
			for _, it := range thisItems {
				m.inspections++
				oldWL := it.worryLevel
				it.worryLevel = m.operation(it.worryLevel)
				if it.worryLevel < oldWL {
					return fmt.Errorf("integer overflow")
				}
				if part1 {
					it.worryLevel /= 3
				}
				it.worryLevel = it.worryLevel % gcd
				var targetMonkey string
				if it.worryLevel%m.testDivisor == 0 {
					targetMonkey = m.trueMonkey
				} else {
					targetMonkey = m.falseMonkey
				}
				monkeys[targetMonkey].items = append(monkeys[targetMonkey].items, it)
			}
		}
	}
	var inspList []int
	for _, mk := range monkeyKeys {
		inspList = append(inspList, monkeys[mk].inspections)
	}
	sort.Ints(inspList)
	topMonkeys := inspList[len(inspList)-2:]
	if part1 {
		fmt.Printf("Part 1: %d\n", topMonkeys[0]*topMonkeys[1])
	} else {
		fmt.Printf("Part 2: %d\n", topMonkeys[0]*topMonkeys[1])
	}
	return nil
}

func main() {
	err := run(true)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	err = run(false)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
