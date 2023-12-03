package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	file, _ := os.OpenFile("input.txt", os.O_RDONLY, 0)
	r := bufio.NewReader(file)
	var sum, state, first, second int
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if b == '\n' && state == 1 {
			sum += first*10 + second
			state = 0
		} else if b >= '0' && b <= '9' {
			if state == 0 {
				first = int(b - '0')
				second = first
				state = 1
			} else {
				second = int(b - '0')
				state = 1
			}
		}
	}
	println(sum)
}
