package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"os/exec"
	"regexp"
)

type expression struct {
	monkey1 string
	op      string
	monkey2 string
}

func eval(n1, n2 int, op string) int {
	var res int
	switch op {
	case "+":
		res = n1 + n2
	case "-":
		res = n1 - n2
	case "*":
		res = n1 * n2
	case "/":
		res = n1 / n2
	default:
		panic("unknown op")
	}
	return res
}

type data struct {
	knownNumbers map[string]int
	expressions  map[string]expression
	caresAbout   map[string][]string
}

func (d *data) findValues() {
	toCheck := []string{"root"}
	for len(toCheck) > 0 {
		monkey := toCheck[0]
		toCheck = toCheck[1:]
		_, ok := d.knownNumbers[monkey]
		if ok {
			continue
		}
		expr := d.expressions[monkey]
		n1, ok1 := d.knownNumbers[expr.monkey1]
		n2, ok2 := d.knownNumbers[expr.monkey2]
		if ok1 && ok2 {
			d.knownNumbers[monkey] = eval(n1, n2, expr.op)
			toCheck = append(toCheck, d.caresAbout[monkey]...)
			d.caresAbout[monkey] = nil
		} else {
			if !ok1 {
				toCheck = append(toCheck, expr.monkey1)
				d.caresAbout[expr.monkey1] = append(d.caresAbout[expr.monkey1], monkey)
			}
			if !ok2 {
				toCheck = append(toCheck, expr.monkey2)
				d.caresAbout[expr.monkey2] = append(d.caresAbout[expr.monkey2], monkey)
			}
		}
	}
}

func (d *data) symbolic(name string) any {
	if name == "humn" {
		return name
	}
	e, ok := d.expressions[name]
	if ok {
		m1 := d.symbolic(e.monkey1)
		m2 := d.symbolic(e.monkey2)
		v1, ok1 := m1.(int)
		v2, ok2 := m2.(int)
		if ok1 && ok2 {
			return eval(v1, v2, e.op)
		}
		if ok1 {
			m1 = fmt.Sprintf("%d", v1)
		}
		if ok2 {
			m2 = fmt.Sprintf("%d", v2)
		}
		return fmt.Sprintf("(%s %s %s)", m1, e.op, m2)
	}
	v, ok := d.knownNumbers[name]
	if ok {
		return v
	}
	panic(fmt.Sprintf("bad symbol: %s", name))
}

func (d *data) symbolicStr(name string) string {
	r := d.symbolic(name)
	switch v := r.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	}
	panic("unknown type")
}

func run() error {
	d := data{
		knownNumbers: make(map[string]int),
		expressions:  make(map[string]expression),
		caresAbout:   make(map[string][]string),
	}
	reNum := regexp.MustCompile(`^(\S+): (\d+)$`)
	reExpr := regexp.MustCompile(`^(\S+): (\S+) (.) (\S+)$`)
	err := utils.OpenAndReadLines("input21.txt", func(s string) error {
		m := reNum.FindStringSubmatch(s)
		if m != nil {
			d.knownNumbers[m[1]] = utils.MustAtoi(m[2])
			return nil
		}
		m = reExpr.FindStringSubmatch(s)
		if m != nil {
			d.expressions[m[1]] = expression{
				monkey1: m[2],
				op:      m[3],
				monkey2: m[4],
			}
			return nil
		}
		return fmt.Errorf("line did not match")
	})
	if err != nil {
		return err
	}
	d2 := data{
		knownNumbers: make(map[string]int),
		expressions:  make(map[string]expression),
		caresAbout:   make(map[string][]string),
	}
	for k, v := range d.knownNumbers {
		d2.knownNumbers[k] = v
	}
	for k, v := range d.expressions {
		d2.expressions[k] = v
	}
	d.findValues()
	fmt.Printf("Part 1: %d\n", d.knownNumbers["root"])

	e := d2.expressions["root"]
	lhs := d2.symbolicStr(e.monkey1)
	rhs := d2.symbolicStr(e.monkey2)

	f, err := os.CreateTemp("", "day21_sage_*")
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	fmt.Printf("Part 2:\n")
	utils.MustWriteString(f, `#!/usr/bin/env sage`+"\n")
	utils.MustWriteString(f, `import sys`+"\n")
	utils.MustWriteString(f, `from sage.all import *`+"\n")
	utils.MustWriteString(f, `var('humn')`+"\n")

	var varNames []string
	for k := range d2.knownNumbers {
		varNames = append(varNames, k)
	}
	for k := range d2.expressions {
		varNames = append(varNames, k)
	}
	utils.MustWriteString(f, fmt.Sprintf("print(solve(%s==%s,humn))\n", lhs, rhs))
	_ = f.Close()
	cmd := exec.Command("/usr/bin/sage", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
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
