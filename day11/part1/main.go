package main

import (
	"bufio"
	"math"
	"os"
)

type coord struct {
	x, y int
	ox   int
}

func manhattanDistance(c1, c2 coord) int {
	// sadly, there is no math.Abs for int, so we have to do a bit of type conversion
	return int(math.Abs(float64(c1.ox-c2.ox))) + int(math.Abs(float64(c1.y-c2.y)))
}

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	galaxies := make([]*coord, 0)

	// parse map and do y-offsetting already
	var colHadGalaxy [1024]bool // <- just have a large enough array
	for y := 0; fileScanner.Scan(); y++ {
		emptyRow := true
		for x, c := range fileScanner.Text() {
			if c == '#' {
				// found a galaxy
				galaxies = append(galaxies, &coord{x, y, x})
				// column cannot be empty anymore
				colHadGalaxy[x] = true
				// row cannot be empty anymore
				emptyRow = false
			}
		}
		if emptyRow {
			// offset y for later rows
			y++
		}
	}

	// fix empty columns and offset items
	// with x coordinate greater than that of
	// the empty column
	for x, hadGalaxy := range colHadGalaxy {
		if !hadGalaxy {
			for _, c := range galaxies {
				if c.x > x {
					c.ox++
				}
			}
		}
	}

	// compute all manhattan distances
	// between pairs of galaxies
	sumOfShortestPaths := 0
	for i, gi := range galaxies {
		for _, gj := range galaxies[:i] {
			sumOfShortestPaths += manhattanDistance(*gi, *gj)
		}
	}

	println(sumOfShortestPaths)
}
