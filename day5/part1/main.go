package main

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/participle/v2"
	"os"
)

type Range struct {
	DstStart uint64 `parser:"@Int"`
	SrcStart uint64 `parser:"@Int"`
	Length   uint64 `parser:"@Int"`
}
type RangeList struct {
	Ranges []Range `parser:"~':'+ ':' @@+"`
}
type Input struct {
	Seeds      []uint64    `parser:"'seeds' ':' @Int+"`
	RangeLists []RangeList `parser:"@@+"`
}

func Map(src uint64, rl []Range) uint64 {
	for i := 0; i < len(rl); i++ {
		if src >= rl[i].SrcStart && src < rl[i].SrcStart+rl[i].Length {
			return rl[i].DstStart + (src - rl[i].SrcStart)
		}
	}
	return src
}
func main() {
	fileName := "input.txt"
	file, _ := os.Open(fileName)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	numbers := input.Seeds
	for _, rl := range input.RangeLists {
		for i, s := range numbers {
			result := Map(s, rl.Ranges)
			numbers[i] = result
		}
	}
	var lowest = numbers[0]
	for i := 1; i < len(numbers); i++ {
		lowest = min(lowest, numbers[i])
	}
	fmt.Printf("lowest: %d\n", lowest)
}
