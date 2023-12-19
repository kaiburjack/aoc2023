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
	X int `parser:"'{' 'x' '=' @Int ','"`
	M int `parser:"'m' '=' @Int ','"`
	A int `parser:"'a' '=' @Int ','"`
	S int `parser:"'s' '=' @Int '}'"`
}
type Input struct {
	Workflows []Workflow `parser:"@@+"`
	Parts     []Part     `parser:"@@+"`
}

func workflowByName(ws []Workflow, name string) Workflow {
	return ws[slices.IndexFunc(ws, func(w Workflow) bool {
		return w.Name == name
	})]
}

var cat2idx = map[string]int{"x": 0, "m": 1, "a": 2, "s": 3}

func combinations(min, max []int) int64 {
	return int64((max[0]-min[0])+1) * int64((max[1]-min[1])+1) * int64((max[2]-min[2])+1) * int64((max[3]-min[3])+1)
}

func validCombinations(min, max []int, w Workflow, ws []Workflow) int64 {
	sum := int64(0)
	for _, r := range w.Rules {
		if r.Op == "" {
			if r.Cat == "A" {
				sum += combinations(min, max)
			} else if r.Cat != "R" {
				sum += validCombinations(min, max, workflowByName(ws, r.Cat), ws)
			}
			return sum
		}
		ci := cat2idx[r.Cat]
		newMin, newMax := min[ci], max[ci]
		nextMin, nextMax := min[ci], max[ci]
		if r.Op == "<" {
			newMax = r.Num - 1
			nextMin = r.Num
		} else if r.Op == ">" {
			newMin = r.Num + 1
			nextMax = r.Num
		}
		mins := []int{min[0], min[1], min[2], min[3]}
		maxs := []int{max[0], max[1], max[2], max[3]}
		mins[ci], maxs[ci] = newMin, newMax
		if r.Dest == "A" {
			sum += combinations(mins, maxs)
		} else if r.Dest != "R" {
			sum += validCombinations(mins, maxs, workflowByName(ws, r.Dest), ws)
		}
		min[ci], max[ci] = nextMin, nextMax
	}
	return sum
}

func main() {
	fileName := "input.txt"
	file, _ := os.Open(fileName)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	in := workflowByName(input.Workflows, "in")
	mins := []int{1, 1, 1, 1}
	maxs := []int{4000, 4000, 4000, 4000}
	println(validCombinations(mins, maxs, in, input.Workflows))
}
