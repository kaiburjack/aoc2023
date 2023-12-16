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
func simulate(grid []uint8, w, x, y, dx, dy int) int {
	touched := make([]uint8, len(grid))
	beamHeads := []*beamHead{{x, y, dx, dy}}
	numTouched := 0
	for changed := true; changed; {
		changed = false
		for i := 0; i < len(beamHeads); i++ {
			bh := beamHeads[i]
			bh.x += bh.dx
			bh.y += bh.dy
			dirMask := bitmask(bh.dx) | bitmask(bh.dy)<<2
			if bh.x < 0 || bh.x >= w || bh.y < 0 || bh.y >= len(grid)/w || touched[w*bh.y+bh.x]&dirMask == dirMask {
				beamHeads = append(beamHeads[:i], beamHeads[i+1:]...)
				i--
				continue
			}
			if touched[w*bh.y+bh.x] == 0 {
				numTouched++
			}
			touched[w*bh.y+bh.x] |= dirMask
			changed = true
			switch grid[w*bh.y+bh.x] {
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
	var grid []uint8
	h := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		for _, c := range line {
			grid = append(grid, uint8(c))
		}
		h++
	}
	w := len(grid) / h
	maxEnergized := 0
	for y := 0; y < len(grid)/w; y++ {
		maxEnergized = max(max(maxEnergized, simulate(grid, w, -1, y, 1, 0)), simulate(grid, w, w, y, -1, 0))
	}
	for x := 0; x < w; x++ {
		maxEnergized = max(max(maxEnergized, simulate(grid, w, x, -1, 0, 1)), simulate(grid, w, x, w, 0, -1))
	}
	println(maxEnergized)
}
