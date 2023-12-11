package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"math"
	"os"
)

type Input struct {
	Times     []float64 `parser:"'Time' ':' @Int+"`
	Distances []float64 `parser:"'Distance' ':' @Int+"`
}

// The puzzle is actually a mathematical solution to a quadratic equation
// which relates time spent to charge the boat to the distance it can then travel
// together with the maximum time allowed.
//
// So, solve the following inequality for 't' as integer:
// -t²+maxTime*t - minDistance > 0
func solveQuadratic(maxTime, minDistance float64) (int64, int64) {
	// t = (-b ± √(b² - 4ac))/2a
	var r0 = (maxTime - math.Sqrt(maxTime*maxTime-4.0*minDistance)) * 0.5
	var r1 = (maxTime + math.Sqrt(maxTime*maxTime-4.0*minDistance)) * 0.5
	return int64(math.Floor(r0 + 1)), int64(math.Ceil(r1 - 1))
}

func main() {
	fileName := "input.txt"
	file, _ := os.Open(fileName)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	var product int64 = 1
	for i := 0; i < len(input.Times); i++ {
		minTime, maxTime := solveQuadratic(input.Times[i], input.Distances[i])
		product *= maxTime - minTime + 1
	}
	println(product)
}
