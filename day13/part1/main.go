package main

import (
	"bufio"
	"bytes"
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
		for dy := 0; dy < length; dy++ {
			if !bytes.Equal(pattern[y-dy], pattern[y+dy+1]) {
				found = false
				break
			}
		}
		if found {
			return y
		}
	}
	return -1
}

func main() {
	readFile, _ := os.Open("example.txt")
	fileScanner := bufio.NewScanner(readFile)
	var currentPattern [][]byte
	var sum int64
	hasNext := true
	for hasNext {
		hasNext = fileScanner.Scan()
		if hasNext && len(fileScanner.Bytes()) != 0 {
			currentPattern = append(currentPattern, []byte(fileScanner.Text()))
			continue
		}
		row := -1
		col := -1
		row = findReflectionLine(currentPattern)
		if row == -1 {
			col = findReflectionLine(transpose(currentPattern))
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
