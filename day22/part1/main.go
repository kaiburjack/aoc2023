package main

import (
	"bufio"
	"os"
	"slices"
	"strconv"
	"strings"
)

type brick struct {
	minX, minY, minZ int
	lenX, lenY, lenZ int
}

func intersect(a, b brick) bool {
	return a.minX <= b.minX+b.lenX-1 && a.minX+a.lenX-1 >= b.minX &&
		a.minY <= b.minY+b.lenY-1 && a.minY+a.lenY-1 >= b.minY &&
		a.minZ <= b.minZ+b.lenZ-1 && a.minZ+a.lenZ-1 >= b.minZ
}

func checkAnyIntersection(bs []brick, without int, i int) bool {
	a := bs[i]
	tester := brick{
		minX: a.minX, minY: a.minY, minZ: a.minZ - 1,
		lenX: a.lenX, lenY: a.lenY, lenZ: a.lenZ,
	}
	for j := i - 1; j >= 0; j-- {
		if intersect(tester, bs[j]) && j != without {
			return true
		}
	}
	return false
}

func pack(bs []brick, without int, onlyCheck bool) bool {
	for freeZ, i := 1, 0; i < len(bs); i++ {
		if bs[i].minZ > freeZ+1 {
			bs[i].minZ = freeZ + 1
		}
		for bs[i].minZ > 1 && !checkAnyIntersection(bs, without, i) {
			if onlyCheck {
				return true
			}
			bs[i].minZ--
		}
		freeZ = max(freeZ, bs[i].minZ+bs[i].lenZ)
	}
	return false
}

func toInts(s []string) []int {
	ret := make([]int, 0, len(s))
	for _, v := range s {
		n, _ := strconv.Atoi(v)
		ret = append(ret, n)
	}
	return ret
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
	var bricks []brick
	for r.Scan() {
		splitted := strings.Split(r.Text(), "~")
		mind, maxd := toInts(strings.Split(splitted[0], ",")), toInts(strings.Split(splitted[1], ","))
		bricks = append(bricks, brick{
			minX: mind[0], minY: mind[1], minZ: mind[2],
			lenX: maxd[0] - mind[0] + 1, lenY: maxd[1] - mind[1] + 1, lenZ: maxd[2] - mind[2] + 1,
		})
	}
	slices.SortFunc(bricks, func(a, b brick) int {
		return a.minZ - b.minZ
	})
	pack(bricks, -1, false)
	canRemoveCount := 0
	for i := 0; i < len(bricks); i++ {
		if !pack(bricks, i, true) {
			canRemoveCount++
		}
	}
	println(canRemoveCount)
}
