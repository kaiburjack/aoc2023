package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"os"
	"slices"
)

type Input struct {
	Cards []Card `parser:"@@+"`
}

type Card struct {
	WinningNumbers []int `parser:"'Card' Int ':' @Int+"`
	MyNumbers      []int `parser:"'|' @Int+"`
}

func points(c Card) int {
	points := 0
	for _, n := range c.WinningNumbers {
		if slices.Contains(c.MyNumbers, n) {
			points = max(points<<1, 1)
		}
	}
	return points
}

func main() {
	fileName := "input.txt"
	file, _ := os.OpenFile(fileName, os.O_RDONLY, 0)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	sumOfPoints := 0
	for _, c := range input.Cards {
		sumOfPoints += points(c)
	}
	println(sumOfPoints)
}
