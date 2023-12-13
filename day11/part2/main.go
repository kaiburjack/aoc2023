package main

import (
	"bufio"
	"math"
	"os"
)

type coord struct {
	x, y int64
	ox   int64
}

func manhattanDistance(c1, c2 coord) int64 {
	// sadly, there is no math.Abs for int, so we have to do a bit of type conversion
	return int64(math.Abs(float64(c1.ox-c2.ox))) + int64(math.Abs(float64(c1.y-c2.y)))
}

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	galaxies := make([]*coord, 0)

	// parse map and do y-offsetting already
	var y int64
	var scale int64 = 1000000
	var colHadGalaxy [1024]bool // <- just have a large enough array
	for fileScanner.Scan() {
		line := fileScanner.Text()
		emptyLine := true
		for x, c := range line {
			if c == '#' {
				// found a galaxy
				galaxies = append(galaxies, &coord{int64(x), y, int64(x)})
				// column cannot be empty anymore
				colHadGalaxy[x] = true
				// row cannot be empty anymore
				emptyLine = false
			}
		}
		if emptyLine {
			// offset y for later rows
			// the - 1 here is all of the trick of this part!
			// because in the first part we _actually_ made every
			// empty row TWICE as large (by adding one empty line)
			// so the scale was 2 in the first part, and therefore
			// the _increment_ was 1
			y += scale - 1
		}
		y++
	}

	// find empty columns and offset items
	// with x coordinate greater than that of
	// the empty column
	for x, hadGalaxy := range colHadGalaxy {
		if !hadGalaxy {
			for _, c := range galaxies {
				if c.x > int64(x) {
					// the - 1 here is all of the trick of this part!
					// because in the first part we _actually_ made every
					// empty column TWICE as large (by adding one empty column)
					// so the scale was 2 in the first part, and therefore
					// the _increment_ was 1
					c.ox += scale - 1
				}
			}
		}
	}

	// compute all manhattan distances
	// between pairs of galaxies
	var sumOfShortestPaths int64
	for i, gi := range galaxies {
		for _, gj := range galaxies[:i] {
			sumOfShortestPaths += manhattanDistance(*gi, *gj)
		}
	}

	println(sumOfShortestPaths)
}
