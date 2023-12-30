package main

import (
	"bufio"
	"github.com/bits-and-blooms/bitset"
	"os"
	"slices"
	"strings"
)

type graph struct {
	s2i   map[string]int
	edges map[int][]int
}

func (v *graph) addVertex(s string) int {
	if _, ok := v.s2i[s]; !ok {
		v.s2i[s] = len(v.s2i)
	}
	return v.s2i[s]
}

func (v *graph) addEdge(v0, v1 int) {
	v.edges[v0] = append(v.edges[v0], v1)
	v.edges[v1] = append(v.edges[v1], v0)
}

func dfsAllDistinctPaths(verts graph, k int, seen *bitset.BitSet, walked *bitset.BitSet) int {
	numPaths := 0
	for _, v := range verts.edges[0] {
		if v == k {
			numPaths++
			continue
		}
		seen.ClearAll()
		q := [][]int{{v, v}}
		for len(q) > 0 && numPaths < 4 {
			qe := q[0]
			q = q[1:]
			if findDistinctPath(verts, k, seen, walked, qe, &numPaths, &q) {
				break
			}
		}
	}
	return numPaths
}

func findDistinctPath(verts graph, k int, seen *bitset.BitSet, walked *bitset.BitSet, qe []int, numPaths *int, q *[][]int) bool {
	for _, c := range verts.edges[qe[0]] {
		if c == k {
			for _, c2 := range qe[1:] {
				walked.Set(uint(c2))
			}
			*numPaths++
			return true
		} else if !seen.Test(uint(c)) && slices.Index(qe[1:], c) == -1 && !walked.Test(uint(c)) {
			newItem := make([]int, len(qe[1:])+2)
			newItem[0] = c
			copy(newItem[1:], qe[1:])
			newItem[len(qe[1:])+1] = c
			*q = append(*q, newItem)
			seen.Set(uint(c))
		}
	}
	return false
}

func main() {
	f, _ := os.Open("input.txt")
	s := bufio.NewScanner(f)
	verts := graph{make(map[string]int), make(map[int][]int)}
	for s.Scan() {
		line := s.Text()
		splitted := strings.Split(line, ": ")
		li := verts.addVertex(splitted[0])
		for _, c := range strings.Split(splitted[1], " ") {
			ci := verts.addVertex(c)
			verts.addEdge(li, ci)
		}
	}

	c0, c1 := 0, 0
	vs, used := bitset.New(uint(len(verts.s2i))), bitset.New(uint(len(verts.s2i)))
	for _, k := range verts.s2i {
		if dfsAllDistinctPaths(verts, k, vs, used.ClearAll().Set(0)) == 3 {
			c0++
		} else {
			c1++
		}
	}
	println(c0 * c1)
}
