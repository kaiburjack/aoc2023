package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"os"
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

func main() {
	fileName := "input.txt"
	file, _ := os.Open(fileName)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	m := make(map[string]Row)
	for _, row := range input.Rows {
		m[row.Origin] = row
	}
	currentPos := "AAA"
	var step uint64
	for ; currentPos != "ZZZ"; step++ {
		if input.Instructions[step%uint64(len(input.Instructions))] == 'L' {
			currentPos = m[currentPos].Left
		} else {
			currentPos = m[currentPos].Right
		}
	}
	println(step)
}
