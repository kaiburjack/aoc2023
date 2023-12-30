package main

import (
	"bufio"
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

func dfsAllDistinctPaths(verts graph, k int, seen []bool, walked []bool) int {
	numPaths := 0
	for _, v := range verts.edges[0] {
		if v == k {
			numPaths++
			continue
		}
		for i := range seen {
			seen[i] = false
		}
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

func findDistinctPath(verts graph, k int, seen []bool, walked []bool, qe []int, numPaths *int, q *[][]int) bool {
	for _, c := range verts.edges[qe[0]] {
		if c == k {
			for _, c2 := range qe[1:] {
				walked[c2] = true
			}
			*numPaths++
			return true
		} else if !seen[c] && slices.Index(qe[1:], c) == -1 && !walked[c] {
			newItem := make([]int, len(qe[1:])+2)
			newItem[0] = c
			copy(newItem[1:], qe[1:])
			newItem[len(qe[1:])+1] = c
			*q = append(*q, newItem)
			seen[c] = true
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
	seen, used := make([]bool, len(verts.s2i)), make([]bool, len(verts.s2i))
	for _, k := range verts.s2i {
		for i := range used {
			used[i] = false
		}
		if dfsAllDistinctPaths(verts, k, seen, used) == 3 {
			c0++
		} else {
			c1++
		}
	}
	println(c0 * c1)
}
