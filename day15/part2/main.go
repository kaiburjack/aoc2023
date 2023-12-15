package main

import (
	"bufio"
	"io"
	"os"
)

type lens struct {
	focalLength int
	label       string
}

func remove(boxes [][]lens, h uint8, label string) {
	for i, l := range boxes[h] {
		if l.label == label {
			boxes[h] = append(boxes[h][:i], boxes[h][i+1:]...)
			break
		}
	}
}

func putLens(boxes [][]lens, h uint8, lens lens) {
	found := false
	for i, l := range boxes[h] {
		if l.label == lens.label {
			boxes[h][i] = lens
			found = true
			break
		}
	}
	if !found {
		boxes[h] = append(boxes[h], lens)
	}
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewReader(file)
	var h uint8
	var state, focalLength int
	var boxes = make([][]lens, 256)
	var label string

	for {
		b, err := r.ReadByte()
		if b == ',' || err == io.EOF {
			if state == 1 {
				putLens(boxes, h, lens{focalLength, label})
			} else if state == 2 {
				remove(boxes, h, label)
			}
			label = ""
			state, h, focalLength = 0, 0, 0
			if err == io.EOF {
				break
			}
		} else if b == '=' {
			state = 1
		} else if b == '-' {
			state = 2
		} else {
			if state == 0 {
				h = (h + b) * 17
				label += string(b)
			} else if state == 1 {
				focalLength = focalLength*10 + int(b-'0')
			}
		}
	}

	focusingPower := 0
	for bi, b := range boxes {
		for i, lens := range b {
			focusingPower += (bi + 1) * (i + 1) * lens.focalLength
		}
	}
	println(focusingPower)
}
