package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"os"
	"slices"
)

type Workflow struct {
	Name  string `parser:"@Ident"`
	Rules []Rule `parser:"'{' @@ (',' @@)* '}'"`
}
type Rule struct {
	Cat  string `parser:"@Ident"`
	Op   string `parser:"@('<'|'>')?"`
	Num  int    `parser:"@Int?"`
	Dest string `parser:"(':' @Ident)?"`
}
type Part struct {
	X       int `parser:"'{' 'x' '=' @Int ','"`
	M       int `parser:"'m' '=' @Int ','"`
	A       int `parser:"'a' '=' @Int ','"`
	S       int `parser:"'s' '=' @Int '}'"`
	cat2val map[string]int
}
type Input struct {
	Workflows []Workflow `parser:"@@+"`
	Parts     []Part     `parser:"@@+"`
}

var compareFuncs = map[string]func(int, int) bool{
	"<": func(a int, b int) bool {
		return a < b
	},
	">": func(a int, b int) bool {
		return a > b
	},
}

func isAccepted(p Part, w Workflow, ws []Workflow) bool {
	for {
		for _, r := range w.Rules {
			dest := r.Dest
			if r.Op == "" {
				dest = r.Cat
			} else if !compareFuncs[r.Op](p.cat2val[r.Cat], r.Num) {
				continue
			}
			if dest == "A" {
				return true
			} else if dest == "R" {
				return false
			} else {
				w = ws[slices.IndexFunc(ws, func(w Workflow) bool {
					return w.Name == dest
				})]
				break
			}
		}
	}
}

func main() {
	fileName := "input.txt"
	file, _ := os.Open(fileName)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	in := input.Workflows[slices.IndexFunc(input.Workflows, func(w Workflow) bool {
		return w.Name == "in"
	})]
	ratings := int64(0)
	for _, p := range input.Parts {
		p.cat2val = map[string]int{"x": p.X, "m": p.M, "a": p.A, "s": p.S}
		if isAccepted(p, in, input.Workflows) {
			ratings += int64(p.X + p.M + p.A + p.S)
		}
	}
	println(ratings)
}
