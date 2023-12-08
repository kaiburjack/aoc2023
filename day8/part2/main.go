package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"os"
	"strings"
)

type Row struct {
	Origin string `parser:"@(Int? Ident)"`
	Left   string `parser:"'=' '(' @(Int? Ident) ','"`
	Right  string `parser:"@(Int? Ident) ')'"`
}

type Input struct {
	Instructions string `parser:"@Ident"`
	Rows         []Row  `parser:"@@+"`
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b uint64, rest ...uint64) uint64 {
	result := a * b / gcd(a, b)
	for i := 0; i < len(rest); i++ {
		result = lcm(result, rest[i])
	}
	return result
}

func main() {
	fileName := "input.txt"
	file, _ := os.OpenFile(fileName, os.O_RDONLY, 0)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	m := make(map[string]Row)
	for _, row := range input.Rows {
		m[row.Origin] = row
	}
	steps := make([]uint64, 0)
	for pos := range m {
		if !strings.HasSuffix(pos, "A") {
			continue
		}
		var step uint64
		for ; !strings.HasSuffix(pos, "Z"); step++ {
			if input.Instructions[step%uint64(len(input.Instructions))] == 'L' {
				pos = m[pos].Left
			} else {
				pos = m[pos].Right
			}
		}
		steps = append(steps, step)
	}
	println(lcm(steps[0], steps[1], steps[2:]...))
}
