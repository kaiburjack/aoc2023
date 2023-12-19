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

var cat2idx = map[string]int{"x": 0, "m": 1, "a": 2, "s": 3}

func combinations(min, max []int) int64 {
	return int64((max[0]-min[0])+1) * int64((max[1]-min[1])+1) * int64((max[2]-min[2])+1) * int64((max[3]-min[3])+1)
}
func validCombinations(min, max []int, w Workflow, n2w map[string]Workflow) int64 {
	var sum int64
	for _, r := range w.Rules {
		if r.Op == "" {
			if r.Cat == "A" {
				sum += combinations(min, max)
			} else if r.Cat != "R" {
				sum += validCombinations(min, max, n2w[r.Cat], n2w)
			}
			return sum
		}
		ci := cat2idx[r.Cat]
		newMin, newMax := min[ci], max[ci]
		nextMin, nextMax := min[ci], max[ci]
		if r.Op == "<" {
			newMax, nextMin = r.Num-1, r.Num
		} else if r.Op == ">" {
			newMin, nextMax = r.Num+1, r.Num
		}
		mins := []int{min[0], min[1], min[2], min[3]}
		maxs := []int{max[0], max[1], max[2], max[3]}
		mins[ci], maxs[ci] = newMin, newMax
		if r.Dest == "A" {
			sum += combinations(mins, maxs)
		} else if r.Dest != "R" {
			sum += validCombinations(mins, maxs, n2w[r.Dest], n2w)
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
	n2w := make(map[string]Workflow)
	for _, w := range input.Workflows {
		n2w[w.Name] = w
	}
	mins := []int{1, 1, 1, 1}
	maxs := []int{4000, 4000, 4000, 4000}
	println(validCombinations(mins, maxs, n2w["in"], n2w))
}
