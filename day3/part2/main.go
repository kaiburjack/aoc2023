package main

import (
	"bufio"
	"io"
	"os"
)

type coordinate struct {
	x, y int
}

type gear struct {
	touched map[coordinate]bool
	ratio   int64
}

type board struct {
	lines [][]byte
	gears map[coordinate]*gear
}

func (thiz *board) detectGearsAround(cx, cy, x, y int, touchedGears map[coordinate]*gear) {
	for dy := -1; dy <= 1; dy++ {
		if y+dy < 0 || y+dy >= len(thiz.lines) {
			continue
		}
		for dx := -1; dx <= 1; dx++ {
			if x+dx < 0 || x+dx >= len(thiz.lines[y+dy]) {
				continue
			}
			b := thiz.lines[y+dy][x+dx]
			if b == '*' {
				c := coordinate{x + dx, y + dy}
				g := thiz.gears[c]
				if g == nil {
					g = &gear{ratio: 1, touched: make(map[coordinate]bool)}
					thiz.gears[c] = g
				}
				g.touched[coordinate{cx, cy}] = true
				touchedGears[c] = g
			}
		}
	}
}

func (thiz *board) findGears() {
	thiz.gears = make(map[coordinate]*gear)
	for y := 0; y < len(thiz.lines); y++ {
		line := thiz.lines[y]
		for x := 0; x < len(line); x++ {
			b := line[x]
			if b >= '0' && b <= '9' {
				var num int
				var x2 int
				gearsTouched := make(map[coordinate]*gear)
				for x2 = x; x2 < len(line); x2++ {
					b2 := line[x2]
					if b2 >= '0' && b2 <= '9' {
						thiz.detectGearsAround(x, y, x2, y, gearsTouched)
						num = num*10 + int(b2-'0')
					} else {
						break
					}
				}
				x = x2
				for _, g := range gearsTouched {
					g.ratio = g.ratio * int64(num)
				}
			}
		}
	}
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
	brd.findGears()
	var sumRatios int64 = 0
	for _, g := range brd.gears {
		if len(g.touched) == 2 {
			sumRatios += g.ratio
		}
	}
	println(sumRatios)
}
