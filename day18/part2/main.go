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
	dx, dy := []int64{1, 0, -1, 0}, []int64{0, 1, 0, -1}
	for r.Scan() {
		lastX, lastY = x, y
		hex := strings.Trim(strings.Split(r.Text(), " ")[2], "(#)")
		steps, _ := strconv.ParseInt(hex[:5], 16, 64)
		x += steps * dx[hex[5]-'0']
		y += steps * dy[hex[5]-'0']
		trenchArea += steps
		polygonArea2 += (lastY + y) * (lastX - x)
	}
	println((polygonArea2+trenchArea)/2 + 1)
}
