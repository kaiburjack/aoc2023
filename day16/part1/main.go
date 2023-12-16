package main

import (
	"bufio"
	"os"
)

type beamHead struct {
	x, y, dx, dy int
}

func bitmask(v int) uint8 {
	if v == 0 {
		return 0
	}
	return uint8((v + 3) >> 1)
}
func simulate(grid [][]uint8, x, y, dx, dy int, touched [][]uint8) int {
	beamHeads := make([]*beamHead, 0)
	beamHeads = append(beamHeads, &beamHead{x, y, dx, dy})
	numTouched := 0
	for changed := true; changed; {
		changed = false
		for i := 0; i < len(beamHeads); i++ {
			bh := beamHeads[i]
			bh.x += bh.dx
			bh.y += bh.dy
			touchIdx := bitmask(bh.dx) | bitmask(bh.dy)<<2
			if bh.x < 0 || bh.x >= len(grid[0]) || bh.y < 0 || bh.y >= len(grid) || touched[bh.y][bh.x]&touchIdx == touchIdx {
				beamHeads = append(beamHeads[:i], beamHeads[i+1:]...)
				i--
				continue
			}
			if touched[bh.y][bh.x] == 0 {
				numTouched++
			}
			touched[bh.y][bh.x] |= touchIdx
			changed = true
			switch grid[bh.y][bh.x] {
			case '/':
				bh.dx, bh.dy = -bh.dy, -bh.dx
			case '\\':
				bh.dx, bh.dy = bh.dy, bh.dx
			case '|':
				if bh.dx != 0 {
					bh.dx, bh.dy = 0, bh.dx
					beamHeads = append(beamHeads, &beamHead{bh.x, bh.y, 0, -bh.dy})
				}
			case '-':
				if bh.dy != 0 {
					bh.dx, bh.dy = bh.dy, 0
					beamHeads = append(beamHeads, &beamHead{bh.x, bh.y, -bh.dx, 0})
				}
			}
		}
	}
	return numTouched
}

func main() {
	file, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(file)
	grid := make([][]uint8, 0)
	touched := make([][]uint8, 0)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		row := make([]uint8, 0)
		trow := make([]uint8, 0)
		for _, c := range line {
			row = append(row, uint8(c))
			trow = append(trow, 0)
		}
		grid = append(grid, row)
		touched = append(touched, trow)
	}

	println(simulate(grid, -1, 0, 1, 0, touched))
}
