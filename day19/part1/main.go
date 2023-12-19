package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"os"
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
	X int `parser:"'{' 'x' '=' @Int"`
	M int `parser:"',' 'm' '=' @Int"`
	A int `parser:"',' 'a' '=' @Int"`
	S int `parser:"',' 's' '=' @Int '}'"`
}
type Input struct {
	Workflows []Workflow `parser:"@@+"`
	Parts     []Part     `parser:"@@+"`
}

var compareFuncs = map[string]func(int, int) bool{
	"<": func(a int, b int) bool { return a < b },
	">": func(a int, b int) bool { return a > b },
}

func isAccepted(c2v map[string]int, w Workflow, n2w map[string]Workflow) bool {
	for {
		for _, r := range w.Rules {
			dest := r.Dest
			if r.Op == "" {
				dest = r.Cat
			} else if !compareFuncs[r.Op](c2v[r.Cat], r.Num) {
				continue
			}
			if dest == "A" {
				return true
			} else if dest == "R" {
				return false
			} else {
				w = n2w[dest]
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
	n2w := make(map[string]Workflow)
	for _, w := range input.Workflows {
		n2w[w.Name] = w
	}
	var ratings int64
	for _, p := range input.Parts {
		c2v := map[string]int{"x": p.X, "m": p.M, "a": p.A, "s": p.S}
		if isAccepted(c2v, n2w["in"], n2w) {
			ratings += int64(p.X + p.M + p.A + p.S)
		}
	}
	println(ratings)
}
