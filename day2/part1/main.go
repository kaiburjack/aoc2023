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
	sumOfIds := 0
	cubes := [3]int{12, 13, 14}
	for _, g := range input.Games {
		possible := true
		for _, d := range g.Draws {
			for _, c := range d.Colors {
				if c.Red > cubes[0] || c.Green > cubes[1] || c.Blue > cubes[2] {
					possible = false
				}
			}
		}
		if possible {
			sumOfIds += g.Id
		}
	}
	println(sumOfIds)
}
