package main

import (
	"bufio"
	"github.com/alecthomas/participle/v2"
	"math"
	"os"
	"strconv"
	"strings"
)

type Input struct {
	Time     string `"Time" ":" @(Int+)`
	Distance string `"Distance" ":" @(Int+)`
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
	if r0 == float64(int64(r0)) {
		// corner case for when r0 is exactly an integer, in which case this integer is NOT a solution
		r0 += 1
	}
	if r1 == float64(int64(r1)) {
		// corner case for when r1 is exactly an integer, in which case this integer is NOT a solution
		r1 -= 1
	}
	return int64(math.Ceil(r0)), int64(math.Floor(r1))
}

func main() {
	fileName := "input.txt"
	file, _ := os.OpenFile(fileName, os.O_RDONLY, 0)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	timeAsFloat64, _ := strconv.ParseFloat(strings.ReplaceAll(input.Time, " ", ""), 64)
	distanceAsFloat64, _ := strconv.ParseFloat(strings.ReplaceAll(input.Distance, " ", ""), 64)
	var minTime, maxTime = solveQuadratic(timeAsFloat64, distanceAsFloat64)
	println(maxTime - minTime + 1)
}
