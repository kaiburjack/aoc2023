package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// matchesCount returns the solution for the given input string
// and number of damaged springs by recursing over the input string
// and accumulating the number of valid combinations that we can
// generate.
func matchesCount(s string, numDamanged []int) int64 {
	if len(s) == 0 && len(numDamanged) == 0 {
		// if we have no more characters to process and no more
		// damaged springs to find, we return "1 combination"
		// for the processed input/path until here.
		return 1
	} else if len(s) == 0 {
		// if we have no more characters to process, but still
		// have damaged springs to find, we return "zero combinations"
		return 0
	}

	switch s[0] {
	case '.':
		// this is a normal non-damanged spring, which we can ignore
		// we continue to recurse over the remaining input.
		return matchesCount(s[1:], numDamanged)
	case '?':
		// this is a spring that may or may not be damaged.
		if len(numDamanged) == 0 {
			// if we have no more damaged springs to find, we
			// interpret this as a normal non-damanged spring (.)
			// and continue to recurse over the remaining input.
			return matchesCount(s[1:], numDamanged)
		}
		// otherwise, we need to consider two possible paths:
		// 1. this is a damaged spring, and we need to complete the sequence
		//    of 'num' damaged springs before we can continue to recurse.
		// 2. this is a normal non-damaged spring, and we can continue to
		//    recurse over the remaining input (like in the case '.' above).
		return mustBeDamaged(s, numDamanged, numDamanged[0]) + matchesCount(s[1:], numDamanged)
	case '#':
		// this is a damaged spring, and we need to complete the sequence
		// of necessary damaged springs before we can continue.
		if len(numDamanged) == 0 {
			// if we have no more damaged springs to find, we cannot generate
			// a possible combination anymore, and return "zero combinations"
			return 0
		}
		// otherwise, we need to complete the sequence of 'num' damaged springs.
		// because this is a common case here and above, it is in its own function.
		return mustBeDamaged(s, numDamanged, numDamanged[0])
	}
	return 0
}

// mustBeDamaged is called when we know that we must generate _exactly_
// num damaged springs to now recurse per character.
func mustBeDamaged(s string, numDamaged []int, num int) int64 {
	if len(s) < num {
		// when we have less characters than we need to generate
		// damaged springs, return "zero combinations"
		return 0
	}
	if strings.ContainsRune(s[:num], '.') {
		// if any of the needed characters cannot be made
		// into a damaged spring, return "zero combinations"
		return 0
	}
	if len(s) == num {
		// if we have exactly the number of characters we need
		// may be done now, but let that be detected by matchesCount().
		return matchesCount(s[num:], numDamaged[1:])
	} else if s[num] == '.' || s[num] == '?' {
		// if there are more characters, the next one after the damaged
		// springs must be a normal spring.
		return matchesCount(s[num+1:], numDamaged[1:])
	}
	// any other case, and we return "zero combinations"
	return 0
}

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var matchesCountSum int64
	for fileScanner.Scan() {
		// parse each line as an input string and a list of
		// numbers of damaged springs
		splitBySpace := strings.Split(fileScanner.Text(), " ")
		var damagedCount []int
		for _, number := range strings.Split(splitBySpace[1], ",") {
			n, _ := strconv.Atoi(strings.TrimSpace(number))
			damagedCount = append(damagedCount, n)
		}
		// accumulate the number of combinations for each line
		matchesCountSum += matchesCount(splitBySpace[0], damagedCount)
	}
	println(matchesCountSum)
}
