package main

import (
	"bufio"
	"os"
)

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	// remember the position of the next obstacle per column
	// the "edge" of the map at the top counts as an obstacle (with value 0)
	var nextObstacle [100]uint8 // <- just allocate a large enough array
	// remember the number of rocks per column
	var numRocks [100]uint8 // <- just allocate a large enough array
	var total uint
	for y := uint8(1); fileScanner.Scan(); y++ {
		for i, c := range fileScanner.Text() {
			// add 1 for each rock that we've already seen in this column,
			// because with every new row that we discover in the input,
			// the rocks that we already saw are one more row away from the bottom.
			total += uint(numRocks[i])

			switch c {
			case '#':
				// we have an obstacle here, remember it so that
				// we know how far a rock can travel north once we
				// see the next rock in this column
				nextObstacle[i] = y
			case 'O':
				// we saw a rock, move it north and increment the
				// total by the number of positions moved
				numRocks[i]++
				total += uint(y - nextObstacle[i])
				// increment the index of the next obstacle by one
				// because the new rock "stacked" on top of it
				nextObstacle[i]++
			}
		}
	}

	// print the result
	println(total)
}
