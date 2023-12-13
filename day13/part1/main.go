package main

import (
	"bufio"
	"os"
)

func transpose(slice [][]byte) [][]byte {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]byte, xl)
	for i := range result {
		result[i] = make([]byte, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

func findReflectionLine(pattern [][]byte) int {
	for y := 0; y < len(pattern)-1; y++ {
		length := y + 1
		if length > len(pattern)/2 {
			length = len(pattern) - y - 1
		}
		found := true
		for dy := 0; dy < length && found; dy++ {
			for x := 0; x < len(pattern[y]); x++ {
				if pattern[y-dy][x] != pattern[y+dy+1][x] {
					found = false
					break
				}
			}
		}
		if found {
			return y
		}
	}
	return -1
}

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var currentPattern [][]byte
	var sum int64
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line != "" {
			currentPattern = append(currentPattern, []byte(line))
			continue
		}
		var row, col int
		transposedPattern := transpose(currentPattern)
		row = findReflectionLine(currentPattern)
		if row == -1 {
			col = findReflectionLine(transposedPattern)
		}
		currentPattern = nil
		if row > -1 {
			sum += 100 * (int64(row) + 1)
		} else if col > -1 {
			sum += int64(col) + 1
		}
	}
	println(sum)
}
