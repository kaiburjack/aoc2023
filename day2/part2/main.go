package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"os"
)

type Game struct {
	Id    int    `parser:"'Game' @Int ':'"`
	Draws []Draw `parser:"@@+"`
}

type Draw struct {
	Colors []Color `parser:"@@+ ';'?"`
}

type Color struct {
	Red   int `parser:"  @Int 'red' ','?"`
	Green int `parser:"| @Int 'green' ','?"`
	Blue  int `parser:"| @Int 'blue' ','?"`
}

type Input struct {
	Games []Game `parser:"@@+"`
}

func main() {
	fileName := "input.txt"
	file, _ := os.Open(fileName)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	var sumOfPowers int64
	for _, g := range input.Games {
		minimums := [3]int{0, 0, 0}
		for _, d := range g.Draws {
			for _, c := range d.Colors {
				minimums[0] = max(c.Red, minimums[0])
				minimums[1] = max(c.Green, minimums[1])
				minimums[2] = max(c.Blue, minimums[2])
			}
		}
		power := int64(minimums[0]) * int64(minimums[1]) * int64(minimums[2])
		sumOfPowers += power
	}
	println(sumOfPowers)
}
