package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type hailstone struct {
	x, y, z    int64
	dx, dy, dz int64
}

func twoRayIntersectWithinBounds2D(aox, aoy, adx, ady, box, boy, bdx, bdy, mi, ma int64) bool {
	dx, dy, det := box-aox, boy-aoy, bdx*ady-bdy*adx
	if det == 0 {
		return false
	}
	u, v := (bdx*dy-bdy*dx)/det, (adx*dy-ady*dx)/det
	if u < 0 || v < 0 {
		return false
	}
	ix, iy := aox+u*adx, aoy+u*ady
	if ix < mi || ix > ma || iy < mi || iy > ma {
		return false
	}
	return true
}

func toInt(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
	hs := make([]hailstone, 0)
	for r.Scan() {
		splitted := strings.Split(r.Text(), " @ ")
		pos, vel := splitted[0], splitted[1]
		posSplitted, velSplitted := strings.Split(pos, ", "), strings.Split(vel, ", ")
		x, y, z := strings.TrimSpace(posSplitted[0]), strings.TrimSpace(posSplitted[1]), strings.TrimSpace(posSplitted[2])
		dx, dy, dz := strings.TrimSpace(velSplitted[0]), strings.TrimSpace(velSplitted[1]), strings.TrimSpace(velSplitted[2])
		hs = append(hs, hailstone{
			toInt(x), toInt(y), toInt(z),
			toInt(dx), toInt(dy), toInt(dz),
		})
	}

	mi, ma := int64(200000000000000), int64(400000000000000)
	count := 0
	for i := 0; i < len(hs)-1; i++ {
		for j := i + 1; j < len(hs); j++ {
			if twoRayIntersectWithinBounds2D(
				hs[i].x, hs[i].y, hs[i].dx, hs[i].dy,
				hs[j].x, hs[j].y, hs[j].dx, hs[j].dy,
				mi, ma) {
				count++
			}
		}
	}
	fmt.Println(count)
}
