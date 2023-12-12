package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type cacheKey struct {
	s string
	n [32]uint8
}

var cache = make(map[cacheKey]int64)

func key(s string, numDamaged []int) cacheKey {
	var n [32]uint8
	for i, v := range numDamaged {
		n[i] = uint8(v)
	}
	return cacheKey{s, n}
}

func matchesCount(s string, numDamaged []int) int64 {
	if cached, ok := cache[key(s, numDamaged)]; ok {
		return cached
	}
	if len(s) == 0 && len(numDamaged) == 0 {
		return 1
	} else if len(s) == 0 {
		return 0
	}
	var count int64
	switch s[0] {
	case '.':
		count = matchesCount(s[1:], numDamaged)
	case '?':
		if len(numDamaged) == 0 {
			count = matchesCount(s[1:], numDamaged)
		} else {
			count = mustBeDamaged(s, numDamaged, numDamaged[0]) + matchesCount(s[1:], numDamaged)
		}
	case '#':
		if len(numDamaged) == 0 {
			count = 0
		} else {
			count = mustBeDamaged(s, numDamaged, numDamaged[0])
		}
	}
	cache[key(s, numDamaged)] = count
	return count
}

func mustBeDamaged(s string, numDamaged []int, num int) int64 {
	if len(s) < num {
		return 0
	}
	if strings.ContainsRune(s[:num], '.') {
		return 0
	}
	if len(s) == num {
		return matchesCount(s[num:], numDamaged[1:])
	} else if s[num] == '.' || s[num] == '?' {
		return matchesCount(s[num+1:], numDamaged[1:])
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
		repeatedSprings := strings.Repeat("?"+splitBySpace[0], 5)[1:]
		repeatedDamageCounts := strings.Repeat(","+splitBySpace[1], 5)[1:]
		var damagedCount []int
		for _, number := range strings.Split(repeatedDamageCounts, ",") {
			n, _ := strconv.Atoi(strings.TrimSpace(number))
			damagedCount = append(damagedCount, n)
		}
		matchesCountSum += matchesCount(repeatedSprings, damagedCount)
	}
	println(matchesCountSum)
}
