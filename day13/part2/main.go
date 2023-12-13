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
		sum += int64(findDifferentReflectionLine(pat, tpat,
			findReflectionLine(-1, pat),
			findReflectionLine(-1, tpat)))
		pat = nil
	}
	println(sum)
}

func findDifferentReflectionLine(pat [][]byte, tpat [][]byte, originalRow int, originalCol int) int {
	for y := 0; y < len(pat); y++ {
		for x := 0; x < len(pat[y]); x++ {
			pat[y][x], tpat[x][y] = 81-pat[y][x], 81-tpat[x][y]
			if r, c := findReflectionLine(originalRow, pat),
				findReflectionLine(originalCol, tpat); r != -1 || c != -1 {
				return (r+1)*100 + c + 1
			}
			pat[y][x], tpat[x][y] = 81-pat[y][x], 81-tpat[x][y]
		}
	}
	return (originalRow+1)*100 + originalCol + 1
}
