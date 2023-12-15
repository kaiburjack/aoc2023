package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewReader(file)
	h := uint8(0)
	sum := 0
	for {
		b, err := r.ReadByte()
		if err == io.EOF || b == ',' {
			sum += int(h)
			h = 0
			if err == io.EOF {
				break
			}
		} else {
			h = (h + b) * 17
		}
	}
	println(sum)
}
