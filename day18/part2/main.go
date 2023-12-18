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
	var x, y, trenchArea int64
	dx, dy := []int64{1, 0, -1, 0}, []int64{0, 1, 0, -1}
	for r.Scan() {
		col := strings.Split(r.Text(), " ")[2]
		colWithoutParenthesis := col[2 : len(col)-1]
		steps, _ := strconv.ParseInt(colWithoutParenthesis[:5], 16, 64)
		x += steps * dx[(colWithoutParenthesis[5]-'0')]
		y += steps * dy[(colWithoutParenthesis[5]-'0')]
		trenchArea += steps
		vertices = append(vertices, []int64{x, y})
	}
	println(areaUsingShoelaceFormula(vertices) + trenchArea/2 + 1)
}
