package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"os"
)

type Game struct {
	Id    int    `"Game" @Int ":"`
	Draws []Draw `@@+`
}

type Draw struct {
	Colors []Color `@@+ ";"?`
}

type Color struct {
	Red   int `@Int "red" ","?`
	Green int `| @Int "green" ","?`
	Blue  int `| @Int "blue" ","?`
}

type Input struct {
	Games []Game `@@+`
}

func main() {
	fileName := "input.txt"
	file, _ := os.OpenFile(fileName, os.O_RDONLY, 0)
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
