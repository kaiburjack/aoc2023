package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
	var x, y, lastX, lastY, trenchArea, polygonArea2 int64
	dirToIndex := map[string]int{"R": 0, "D": 1, "L": 2, "U": 3}
	dx, dy := []int64{1, 0, -1, 0}, []int64{0, 1, 0, -1}
	for r.Scan() {
		lastX, lastY = x, y
		s := strings.Split(r.Text(), " ")
		steps, _ := strconv.ParseInt(s[1], 10, 64)
		x += steps * dx[dirToIndex[s[0]]]
		y += steps * dy[dirToIndex[s[0]]]
		trenchArea += steps
		polygonArea2 += (lastY + y) * (lastX - x)
	}
	println((polygonArea2+trenchArea)/2 + 1)
}
