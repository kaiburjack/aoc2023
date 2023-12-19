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
type Input struct {
	Workflows []Workflow `parser:"@@+"`
}

var c2i = map[string]int{"x": 0, "m": 1, "a": 2, "s": 3}

func validCombinations(min, max []int, w Workflow, n2w map[string]Workflow) int64 {
	var sum int64
	for _, r := range w.Rules {
		mins := []int{min[0], min[1], min[2], min[3]}
		maxs := []int{max[0], max[1], max[2], max[3]}
		if r.Op == "<" {
			maxs[c2i[r.Cat]], min[c2i[r.Cat]] = r.Num-1, r.Num
		} else if r.Op == ">" {
			mins[c2i[r.Cat]], max[c2i[r.Cat]] = r.Num+1, r.Num
		} else {
			r.Dest = r.Cat
		}
		if r.Dest == "A" {
			sum += int64((maxs[0]-mins[0])+1) * int64((maxs[1]-mins[1])+1) * int64((maxs[2]-mins[2])+1) * int64((maxs[3]-mins[3])+1)
		} else if r.Dest != "R" {
			sum += validCombinations(mins, maxs, n2w[r.Dest], n2w)
		}
	}
	return sum
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
	mins := []int{1, 1, 1, 1}
	maxs := []int{4000, 4000, 4000, 4000}
	println(validCombinations(mins, maxs, n2w["in"], n2w))
}
