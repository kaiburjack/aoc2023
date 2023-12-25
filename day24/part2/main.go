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

func twoRayIntersect2D(aox, aoy, adx, ady, box, boy, bdx, bdy int64) (int64, int64, bool) {
	dx, dy, det := box-aox, boy-aoy, bdx*ady-bdy*adx
	if det == 0 {
		return 0, 0, false
	}
	u := (bdx*dy - bdy*dx) / det
	return aox + u*adx, aoy + u*ady, true
}
func toInt(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
func zigZag(i int) int              { return (i >> 1) ^ -(i & 1) }
func projectY(hs *hailstone) int64  { return hs.y }
func projectZ(hs *hailstone) int64  { return hs.z }
func projectDY(hs *hailstone) int64 { return hs.dy }
func projectDZ(hs *hailstone) int64 { return hs.dz }

func findValuesForCommonIntersection(hs []*hailstone, rdx int64, projectY, projectDY func(*hailstone) int64) (int64, int64, bool) {
	ma := 500
	for rdyi := 0; rdyi <= ma*2; rdyi++ {
		rdy := int64(zigZag(rdyi))
		if ix, iy, ok := twoRayIntersect2D(
			hs[0].x, projectY(hs[0]), hs[0].dx-rdx, projectDY(hs[0])-rdy,
			hs[1].x, projectY(hs[1]), hs[1].dx-rdx, projectDY(hs[1])-rdy); ok {
			different := false
			for i := 2; i < len(hs) && !different; i++ {
				if ix2, iy2, ok := twoRayIntersect2D(
					hs[0].x, projectY(hs[0]), hs[0].dx-rdx, projectDY(hs[0])-rdy,
					hs[i].x, projectY(hs[i]), hs[i].dx-rdx, projectDY(hs[i])-rdy); !ok || ix != ix2 || iy != iy2 {
					different = true
				}
			}
			if !different {
				return ix, iy, true
			}
		}
	}
	return 0, 0, false
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
	hs := make([]*hailstone, 0)
	for r.Scan() {
		splitted := strings.Split(r.Text(), " @ ")
		posSplitted, velSplitted := strings.Split(splitted[0], ", "), strings.Split(splitted[1], ", ")
		x, y, z := strings.TrimSpace(posSplitted[0]), strings.TrimSpace(posSplitted[1]), strings.TrimSpace(posSplitted[2])
		dx, dy, dz := strings.TrimSpace(velSplitted[0]), strings.TrimSpace(velSplitted[1]), strings.TrimSpace(velSplitted[2])
		hs = append(hs, &hailstone{
			toInt(x), toInt(y), toInt(z),
			toInt(dx), toInt(dy), toInt(dz),
		})
	}
	for rdxi := 0; ; rdxi++ {
		rdx := int64(zigZag(rdxi))
		if ix, iy, found := findValuesForCommonIntersection(hs, rdx, projectY, projectDY); found {
			if _, iz, found := findValuesForCommonIntersection(hs, rdx, projectZ, projectDZ); found {
				fmt.Println(ix + iy + iz)
				break
			}
		}
	}
}
