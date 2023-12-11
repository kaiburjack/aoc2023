package main

import (
	"bufio"
	"io"
	"os"
)

type board struct {
	lines [][]byte
}

func (thiz board) isSymbolAround(x, y int) bool {
	for dy := -1; dy <= 1; dy++ {
		if y+dy < 0 || y+dy >= len(thiz.lines) {
			continue
		}
		for dx := -1; dx <= 1; dx++ {
			if x+dx < 0 || x+dx >= len(thiz.lines[y+dy]) {
				continue
			}
			b := thiz.lines[y+dy][x+dx]
			if b != '.' && (b < '0' || b > '9') {
				return true
			}
		}
	}
	return false
}

func (thiz board) sumOfNumbers() int64 {
	var sum int64
	for y := 0; y < len(thiz.lines); y++ {
		line := thiz.lines[y]
		for x := 0; x < len(line); x++ {
			b := line[x]
			if b >= '0' && b <= '9' {
				var num int
				var truePartNumber bool
				for ; x < len(line); x++ {
					b2 := line[x]
					if b2 >= '0' && b2 <= '9' {
						truePartNumber = truePartNumber || thiz.isSymbolAround(x, y)
						num = num*10 + int(b2-'0')
					} else {
						break
					}
				}
				if truePartNumber {
					sum += int64(num)
				}
			}
		}
	}
	return sum
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewReader(file)
	var brd board
	var currentLine []byte
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		switch b {
		case '\n':
			brd.lines = append(brd.lines, currentLine)
			currentLine = []byte{}
		default:
			currentLine = append(currentLine, b)
		}
	}
	println(brd.sumOfNumbers())
}
