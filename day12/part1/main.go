package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func matchesCount(s string, numDefects []int) int64 {
	if len(s) == 0 && len(numDefects) == 0 {
		return 1
	} else if len(s) == 0 {
		return 0
	}
	switch s[0] {
	case '.':
		return matchesCount(s[1:], numDefects)
	case '?':
		if len(numDefects) == 0 {
			return matchesCount(s[1:], numDefects)
		}
		return mustBeDefect(s, numDefects, numDefects[0]) + matchesCount(s[1:], numDefects)
	case '#':
		if len(numDefects) == 0 {
			return 0
		}
		return mustBeDefect(s, numDefects, numDefects[0])
	}
	return 0
}

func mustBeDefect(s string, numDefects []int, num int) int64 {
	if len(s) < num {
		return 0
	}
	for i := 0; i < num; i++ {
		if s[i] != '#' && s[i] != '?' {
			return 0
		}
	}
	if len(s) == num {
		return matchesCount(s[num:], numDefects[1:])
	} else if s[num] == '.' || s[num] == '?' {
		return matchesCount(s[num+1:], numDefects[1:])
	}
	return 0
}

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var matchesCountSum int64
	for fileScanner.Scan() {
		splitBySpace := strings.Split(fileScanner.Text(), " ")
		var damagedCount []int
		for _, number := range strings.Split(splitBySpace[1], ",") {
			n, _ := strconv.Atoi(strings.TrimSpace(number))
			damagedCount = append(damagedCount, n)
		}
		matchesCountSum += matchesCount(splitBySpace[0], damagedCount)
	}
	println(matchesCountSum)
}
