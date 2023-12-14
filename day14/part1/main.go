package main

import (
	"bufio"
	"os"
)

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	var nextObstacle [1024]int
	var numRocks [1024]int
	var total int
	for y := 1; fileScanner.Scan(); y++ {
		for i, c := range fileScanner.Text() {
			total += numRocks[i]
			switch c {
			case '#':
				nextObstacle[i] = y
			case 'O':
				numRocks[i]++
				total += y - nextObstacle[i]
				nextObstacle[i]++
			}
		}
	}
	println(total)
}
