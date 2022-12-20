package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"regexp"
	"strings"
)

type ingredient struct {
	resource string
	qty int
}

type blueprint struct {
	id int
	robots map[string][]ingredient
	maxes map[string]int
}

type data struct {
	blueprints []blueprint
}

type simState struct {
	blueprint      *blueprint
	parent         *simState
	tick           int
	timeLimit      int
	resources      map[string]int
	robots         map[string]int
}

func initState(b *blueprint, timeLimit int) *simState {
	state := &simState{
		blueprint: b,
		timeLimit: timeLimit,
		resources: make(map[string]int),
		robots:    make(map[string]int),
	}
	state.robots["ore"] = 1
	return state
}

func (state *simState) clone() *simState {
	newState := &simState{
		blueprint: state.blueprint,
		parent:    state,
		tick:      state.tick,
		timeLimit: state.timeLimit,
		resources: make(map[string]int),
		robots:    make(map[string]int),
	}
	for k, v := range state.resources {
		newState.resources[k] = v
	}
	for k, v := range state.robots {
		newState.robots[k] = v
	}
	return newState
} 

func (state *simState) step(conRes string) *simState {
	newState := state.clone()
	newState.tick++
	if conRes != "" {
		ingr := newState.checkAndGetIngredients(conRes)
		if ingr == nil {
			panic("un-constructable robot")
		}
		for _, ci := range ingr {
			newState.resources[ci.resource] -= ci.qty
		}
	}
	for k, v := range newState.robots {
		e, _ := newState.resources[k]
		newState.resources[k] = e + v
	}
	if conRes != "" {
		e, _ := newState.robots[conRes]
		newState.robots[conRes] = e+1
	}
	return newState
}

func (state *simState) checkAndGetIngredients(r string) []ingredient {
	ingr, ok := state.blueprint.robots[r]
	if !ok {
		panic("unknown robot type")
	}
	for _, ci := range ingr {
		e, ok := state.resources[ci.resource]
		if !ok || e < ci.qty {
			return nil
		}
	}
	return ingr
}

var reBlueprint = regexp.MustCompile(`^Blueprint (\d+): (.*)$`)
var reRecipe = regexp.MustCompile(`^Each (.*) robot costs (.*)$`)
var reIngredient = regexp.MustCompile(`^(\d+) (.*)$`)

func parseBlueprint(s string) blueprint {
	m := reBlueprint.FindStringSubmatch(s)
	if m == nil {
		panic("invalid blueprint")
	}
	b := blueprint{
		id:     utils.MustAtoi(m[1]),
		robots: make(map[string][]ingredient),
		maxes:  make(map[string]int),
	}
	for _, rec := range strings.Split(m[2], ".") {
		rec = strings.TrimSpace(rec)
		if rec == "" {
			continue
		}
		m = reRecipe.FindStringSubmatch(rec)
		if m == nil {
			panic("invalid recipe")
		}
		var ingr []ingredient
		robotResource := m[1]
		for _, iStr := range strings.Split(m[2], " and ") {
			iStr = strings.TrimSpace(iStr)
			m = reIngredient.FindStringSubmatch(iStr)
			if m == nil {
				panic("invalid ingredient")
			}
			ingr = append(ingr, ingredient{
				resource: m[2],
				qty:      utils.MustAtoi(m[1]),
			})
		}
		b.robots[robotResource] = ingr
	}
	for _, v := range b.robots {
		for _, ing := range v {
			maxQ, ok := b.maxes[ing.resource]
			if !ok || ing.qty > maxQ {
				b.maxes[ing.resource] = ing.qty
			}
		}
	}
	return b
}

func bestFrom(s *simState) *simState {
	if s.tick >= s.timeLimit {
		return s
	}
	best := bestFrom(s.step(""))
	for r := range s.blueprint.robots {
		rMax, ok := s.blueprint.maxes[r]
		if !ok || (s.robots[r] < rMax && s.resources[r] <= rMax+1) {
			ingr := s.checkAndGetIngredients(r)
			valid := ingr != nil
			if valid {
				found := false
				for _, ing := range ingr {
					if s.resources[ing.resource] <= ing.qty+s.robots[ing.resource] {
						found = true
						break
					}
				}
				if !found {
					valid = false
				}
			}
			if valid {
				rv := bestFrom(s.step(r))
				if rv.resources["geode"] > best.resources["geode"] {
					best = rv
				}
			}
		}
	}
	return best
}

func printHistory(s *simState) {
	var steps []*simState
	sp := s
	for sp != nil {
		steps = append(steps, sp)
		sp = sp.parent
	}
	for i := len(steps)-1; i>=0; i-- {
		sp = steps[i]
		if sp.tick == 0 {
			continue
		}
		fmt.Printf("== Minute %d ==\n", sp.tick)
		pr := i+1
		if pr < len(steps) {
			spr := steps[pr]
			for _, res := range []string{"ore", "clay", "obsidian", "geode"} {
				if sp.robots[res] != spr.robots[res] {
					fmt.Printf("Ordered new %s-collecting robot for %v\n", res, sp.blueprint.robots[res])
				}
			}
			for _, res := range []string{"ore", "clay", "obsidian", "geode"} {
				if spr.robots[res] > 0 && sp.resources[res] > 0 {
					fmt.Printf("%d %s-collecting robot collects %d %s; you now have %d\n", spr.robots[res], res, spr.robots[res], res, sp.resources[res])
				}
			}
			for _, res := range []string{"ore", "clay", "obsidian", "geode"} {
				if sp.robots[res] != steps[pr].robots[res] {
					fmt.Printf("New %s-collecting robot is ready; you now have %d\n", res, sp.robots[res])
				}
			}
		}
		fmt.Printf("\n")
	}
}

func run() error {
	d := data{}
	err := utils.OpenAndReadLines("input19.txt", func(s string) error {
		d.blueprints = append(d.blueprints, parseBlueprint(s))
		return nil
	})
	if err != nil {
		return err
	}

	part1 := 0
	for i, b := range d.blueprints {
		g := bestFrom(initState(&b, 24))
		part1 += (i+1) * g.resources["geode"]
	}
	fmt.Printf("Part 1: %d\n", part1)

	part2 := 1
	for _, b := range d.blueprints[:3] {
		g := bestFrom(initState(&b, 32))
		part2 *= g.resources["geode"]
	}
	fmt.Printf("Part 2: %d\n", part2)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
