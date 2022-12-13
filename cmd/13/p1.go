package main

import (
	"fmt"
	"github.com/ghjm/advent2022/pkg/utils"
	"os"
	"sort"
	"strconv"
)

type tokenType int

const (
	leftBracket tokenType = iota
	rightBracket
	number
	eof
)

type token struct {
	tt tokenType
	v  int
}

type tokens struct {
	data []token
}

func (t *tokens) consume() {
	if len(t.data) > 0 {
		t.data = t.data[1:]
	}
}

func (t *tokens) peek() *token {
	if len(t.data) == 0 {
		return &token{eof, 0}
	}
	return &t.data[0]
}

func (t *tokens) parse() (any, error) {
	tok := t.peek()
	switch tok.tt {
	case number:
		t.consume()
		return tok.v, nil
	case leftBracket:
		var l []any
		t.consume()
		for {
			tok = t.peek()
			if tok.tt == rightBracket {
				t.consume()
				return l, nil
			} else {
				v, err := t.parse()
				if err != nil {
					return nil, err
				}
				l = append(l, v)
			}
		}
	case eof:
		return nil, nil
	}
	return nil, fmt.Errorf("unknown token")
}

func tokenize(s string) (*tokens, error) {
	toks := &tokens{}
	i := 0
	for i < len(s) {
		switch c := s[i]; {
		case c == '[':
			toks.data = append(toks.data, token{leftBracket, 0})
			i++
		case c == ']':
			toks.data = append(toks.data, token{rightBracket, 0})
			i++
		case c == ',':
			i++
		case c >= '0' && c <= '9':
			var val []byte
			for c >= '0' && c <= '9' {
				val = append(val, c)
				i++
				c = s[i]
			}
			v, err := strconv.Atoi(string(val))
			if err != nil {
				return nil, fmt.Errorf("invalid integer")
			}
			toks.data = append(toks.data, token{number, v})
		default:
			return nil, fmt.Errorf("invalid character")
		}
	}
	return toks, nil
}

func fullParse(s string) (any, error) {
	toks, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	tree, err := toks.parse()
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// returns 1 = in the right order, -1 = in the wrong order, 0 = equal / not decided
func compare(left, right any) (int, error) {
	switch leftV := left.(type) {
	case int:
		switch rightV := right.(type) {
		case int:
			if leftV < rightV {
				return 1, nil
			} else if leftV > rightV {
				return -1, nil
			} else {
				return 0, nil
			}
		case []any:
			return compare([]any{left}, right)
		}
	case []any:
		switch rightV := right.(type) {
		case int:
			return compare(left, []any{right})
		case []any:
			i := 0
			for {
				ltl := i >= len(leftV)
				ltr := i >= len(rightV)
				if ltl && ltr {
					return 0, nil
				} else if ltl {
					return 1, nil
				} else if ltr {
					return -1, nil
				}
				c, err := compare(leftV[i], rightV[i])
				if err != nil {
					return 0, err
				}
				if c != 0 {
					return c, nil
				}
				i++
			}
		}
	}
	return 0, fmt.Errorf("unknown type")
}

func run() error {
	var l1 any
	idx := 1
	sum := 0
	var allPackets []any
	err := utils.OpenAndReadLines("input13.txt", func(s string) error {
		if s == "" {
			return nil
		}
		tree, err := fullParse(s)
		if err != nil {
			return err
		}
		allPackets = append(allPackets, tree)
		if l1 == nil {
			l1 = tree
		} else {
			c, err := compare(l1, tree)
			if err != nil {
				return err
			}
			l1 = nil
			if c == 1 {
				sum += idx
			}
			idx++
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Part 1: %d\n", sum)
	for _, s := range []string{"[[2]]", "[[6]]"} {
		tree, err := fullParse(s)
		if err != nil {
			return err
		}
		allPackets = append(allPackets, tree)
	}
	sort.Slice(allPackets, func(i, j int) bool {
		c, err := compare(allPackets[i], allPackets[j])
		if err != nil {
			panic(err)
		}
		return c == 1
	})
	var pos2, pos6 int
	for i, p := range allPackets {
		s := fmt.Sprintf("%v", p)
		if s == "[[2]]" {
			pos2 = i + 1
		}
		if s == "[[6]]" {
			pos6 = i + 1
		}
	}
	fmt.Printf("Part 2: %d\n", pos2*pos6)
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
