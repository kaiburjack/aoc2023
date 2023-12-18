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
		col := strings.Split(r.Text(), " ")[2]
		colWithoutParenthesis := col[2 : len(col)-1]
		steps, _ := strconv.ParseInt(colWithoutParenthesis[:5], 16, 64)
		x += steps * dx[(colWithoutParenthesis[5]-'0')]
		y += steps * dy[(colWithoutParenthesis[5]-'0')]
		trenchArea += steps
		polygonArea2 += (lastY + y) * (lastX - x)
	}
	println((polygonArea2+trenchArea)/2 + 1)
}
