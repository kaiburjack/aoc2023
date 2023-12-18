package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func areaUsingShoelaceFormula(vertices [][]int64) int64 {
	var area int64
	for i := 0; i < len(vertices); i++ {
		v0, v1 := vertices[i], vertices[(i+1)%len(vertices)]
		area += (v0[1] + v1[1]) * (v0[0] - v1[0])
	}
	return area / 2
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
	var vertices [][]int64
	var x, y, area int64
	dirToIndex := map[string]int{"R": 0, "D": 1, "L": 2, "U": 3}
	dx, dy := []int64{1, 0, -1, 0}, []int64{0, 1, 0, -1}
	for r.Scan() {
		s := strings.Split(r.Text(), " ")
		steps, _ := strconv.ParseInt(s[1], 10, 64)
		x += steps * dx[dirToIndex[s[0]]]
		y += steps * dy[dirToIndex[s[0]]]
		area += steps
		vertices = append(vertices, []int64{x, y})
	}
	println(areaUsingShoelaceFormula(vertices) + area/2 + 1)
}
