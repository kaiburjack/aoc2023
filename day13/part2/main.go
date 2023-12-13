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

func findReflectionLine(originalRow int, pattern [][]byte) int {
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
		if found && originalRow != y {
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
	hasNext := true
	for hasNext {
		hasNext = fileScanner.Scan()
		if hasNext && len(fileScanner.Bytes()) != 0 {
			currentPattern = append(currentPattern, fileScanner.Bytes())
			continue
		}
		transposedPattern := transpose(currentPattern)
		row := -1
		col := -1
		row = findReflectionLine(-1, currentPattern)
		if row == -1 {
			col = findReflectionLine(-1, transposedPattern)
		}
		row, col = findDifferentReflectionLine(currentPattern, transposedPattern, row, col)
		currentPattern = nil
		sum += 100*(int64(row)+1) + int64(col) + 1
	}
	println(sum)
}

func findDifferentReflectionLine(currentPattern [][]byte, transposedPattern [][]byte, originalRow int, originalCol int) (int, int) {
	for y := 0; y < len(currentPattern); y++ {
		for x := 0; x < len(currentPattern[y]); x++ {
			oldVal := currentPattern[y][x]
			if oldVal == '#' {
				currentPattern[y][x] = '.'
				transposedPattern[x][y] = '.'
			} else {
				currentPattern[y][x] = '#'
				transposedPattern[x][y] = '#'
			}
			row := findReflectionLine(originalRow, currentPattern)
			if row != -1 {
				return row, -1
			}
			col := findReflectionLine(originalCol, transposedPattern)
			if col != -1 {
				return -1, col
			}
			currentPattern[y][x] = oldVal
			transposedPattern[x][y] = oldVal
		}
	}
	return originalRow, originalCol
}
