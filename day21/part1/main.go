package main

import (
	"bufio"
	"os"
)

var cache = map[[3]int]struct{}{}

func walk(grid [][]uint8, x, y int, step int, maxStep int) int64 {
	if x < 0 || y < 0 || x >= len(grid[0]) || y >= len(grid) || grid[y][x] == '#' {
		return 0
	}
	key := [3]int{x, y, step}
	if _, ok := cache[key]; ok {
		return 0
	} else if step == maxStep {
		cache[key] = struct{}{}
		return 1
	}
	cache[key] = struct{}{}
	v := walk(grid, x-1, y, step+1, maxStep) +
		walk(grid, x+1, y, step+1, maxStep) +
		walk(grid, x, y-1, step+1, maxStep) +
		walk(grid, x, y+1, step+1, maxStep)
	return v
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
	var grid [][]uint8
	for r.Scan() {
		grid = append(grid, []byte(r.Text()))
	}
	const N = 64
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[0]); x++ {
			if grid[y][x] == 'S' {
				println(walk(grid, x, y, 0, N))
				return
			}
		}
	}
}
