package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"os"
	"slices"
)

type Input struct {
	Cards []Card `@@+`
}

type Card struct {
	WinningNumbers []int `"Card" Int ":" @Int+`
	MyNumbers      []int `"|" @Int+`
}

func numberOfWinningNumbers(c Card) int {
	num := 0
	for _, n := range c.WinningNumbers {
		if slices.Contains(c.MyNumbers, n) {
			num++
		}
	}
	return num
}

func main() {
	fileName := "input.txt"
	file, _ := os.OpenFile(fileName, os.O_RDONLY, 0)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	var wonCopies []int64
	var numWonCards int64 = 0
	for _, c := range input.Cards {
		numberOfWinningNumbers := numberOfWinningNumbers(c)
		var count int64 = 1
		if len(wonCopies) > 0 {
			count = 1 + wonCopies[0]
		} else {
			wonCopies = append(wonCopies, 1)
		}
		numWonCards += count
		wonCopies = wonCopies[1:]
		for i := 0; i < numberOfWinningNumbers; i++ {
			if len(wonCopies) > i {
				wonCopies[i] += count
			} else {
				wonCopies = append(wonCopies, count)
			}
		}
	}
	println(numWonCards)
}
