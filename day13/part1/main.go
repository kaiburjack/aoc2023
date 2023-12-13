package main

import (
	"bufio"
	"bytes"
	"os"
)

func transpose(slice [][]byte) [][]byte {
	xl, yl := len(slice[0]), len(slice)
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
		found := true
		for dy := 0; dy < min(y+1, len(pattern)-y-1); dy++ {
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
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	var pat [][]byte
	var sum int64
	hasNext := true
	for hasNext {
		hasNext = fileScanner.Scan()
		if hasNext && len(fileScanner.Bytes()) != 0 {
			pat = append(pat, []byte(fileScanner.Text()))
			continue
		}
		col := -1
		row := findReflectionLine(pat)
		if row == -1 {
			col = findReflectionLine(transpose(pat))
		}
		pat = nil
		if row > -1 {
			sum += 100 * (int64(row) + 1)
		} else if col > -1 {
			sum += int64(col) + 1
		}
	}
	println(sum)
}
