package main

import (
	"bufio"
	"os"
)

type beamHead struct {
	x, y, dx, dy int
}

func simulate(grid []uint8, w, x, y, dx, dy int) int {
	touched := make([]uint8, len(grid)>>1)
	beamHeads := []*beamHead{{x, y, dx, dy}}
	numTouched := 0
	for changed := true; changed; {
		changed = false
		for i := 0; i < len(beamHeads); i++ {
			bh := beamHeads[i]
			bh.x, bh.y = bh.x+bh.dx, bh.y+bh.dy
			dirMask, idx, shift := uint8((bh.dx+3)%3|(bh.dy+3)%3<<2), (w*bh.y+bh.x)>>1, (bh.x+bh.y)&1<<2
			if bh.x < 0 || bh.x >= w || bh.y < 0 || bh.y >= len(grid)/w || (touched[idx]>>shift)&dirMask == dirMask {
				beamHeads = append(beamHeads[:i], beamHeads[i+1:]...)
				i--
				continue
			}
			if touched[idx]>>shift&0xF == 0 {
				numTouched++
			}
			touched[idx] |= dirMask << shift
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
	r := bufio.NewScanner(file)
	var grid []uint8
	h := 0
	for r.Scan() {
		grid = append(grid, r.Bytes()...)
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
