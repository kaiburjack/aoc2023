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

func findReflectionLine(originalIndex int, pattern [][]byte) int {
	for y := 0; y < len(pattern)-1; y++ {
		found := true
		for dy := 0; dy < min(y+1, len(pattern)-y-1); dy++ {
			if !bytes.Equal(pattern[y-dy], pattern[y+dy+1]) {
				found = false
				break
			}
		}
		if found && originalIndex != y {
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
		tpat := transpose(pat)
		row := findReflectionLine(-1, pat)
		col := -1
		if row == -1 {
			col = findReflectionLine(-1, tpat)
		}
		sum += int64(findDifferentReflectionLine(pat, tpat, row, col))
		pat = nil
	}
	println(sum)
}

func findDifferentReflectionLine(pat [][]byte, tpat [][]byte, originalRow int, originalCol int) int {
	for y := 0; y < len(pat); y++ {
		for x := 0; x < len(pat[y]); x++ {
			oldVal := pat[y][x]
			if oldVal == '#' {
				pat[y][x], tpat[x][y] = '.', '.'
			} else {
				pat[y][x], tpat[x][y] = '#', '#'
			}
			if v := (findReflectionLine(originalRow, pat) + 1) * 100; v != 0 {
				return v
			} else if v = findReflectionLine(originalCol, tpat) + 1; v != 0 {
				return v
			}
			pat[y][x], tpat[x][y] = oldVal, oldVal
		}
	}
	return (originalRow+1)*100 + originalCol + 1
}
