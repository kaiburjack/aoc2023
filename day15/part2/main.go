package main

import (
	"bufio"
	"io"
	"os"
	"slices"
)

type lens struct {
	focalLength int
	label       string
}

func removeLens(boxes [][]lens, h uint8, label string) {
	boxes[h] = slices.DeleteFunc(boxes[h], func(l lens) bool {
		return l.label == label
	})
}

func putLens(boxes [][]lens, h uint8, v lens) {
	if i := slices.IndexFunc(boxes[h], func(e lens) bool {
		return v.label == e.label
	}); i != -1 {
		boxes[h][i] = v
	} else {
		boxes[h] = append(boxes[h], v)
	}
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewReader(file)
	var h uint8
	var state, focalLength int
	var boxes [256][]lens
	var label string

	for {
		b, err := r.ReadByte()
		if b == ',' || err == io.EOF {
			if state == 1 {
				putLens(boxes[:], h, lens{focalLength, label})
			} else if state == 2 {
				removeLens(boxes[:], h, label)
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
		} else if state == 0 {
			h = (h + b) * 17
			label += string(b)
		} else if state == 1 {
			focalLength = focalLength*10 + int(b-'0')
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
