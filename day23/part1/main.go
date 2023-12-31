package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
)

type node struct {
	x, y  int
	id    uint8
	edges []edge
}

type edge struct {
	to *node
	d  uint
}

var dirs = [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
var arrows = []byte{'v', '>', '^', '<'}

func buildContractedEdges(grid [][]byte, sx, sy, px, py, ex, ey int, c *node, seen map[[2]int]*node) {
	for d := uint(1); ; d++ {
		var nextPossibles [][]int
		ai := slices.Index(arrows, grid[sy][sx])
		for i := 0; i < 4; i++ {
			// find possible next positions
			nx, ny := sx+dirs[i][0], sy+dirs[i][1]
			if nx < 0 || nx >= len(grid[0]) || ny < 0 || ny >= len(grid) ||
				grid[ny][nx] == '#' ||
				ai != -1 && ai != i ||
				px == nx && py == ny {
				continue
			}
			nextPossibles = append(nextPossibles, []int{nx, ny})
		}
		px, py = sx, sy
		if sx == ex && sy == ey {
			// if we're at the end, add an edge to the end node
			seenNode, ok := seen[[2]int{sx, sy}]
			if ok {
				c.edges = append(c.edges, edge{seenNode, d})
				break
			}
			endNode := &node{sx, sy, uint8(len(seen)), nil}
			seen[[2]int{sx, sy}] = endNode
			c.edges = append(c.edges, edge{endNode, d})
			break
		}
		if len(nextPossibles) == 1 {
			// if there's only one possible next position, move there
			// and continue the loop
			sx, sy = nextPossibles[0][0], nextPossibles[0][1]
		} else {
			// if there are multiple possible next positions, we're at a
			// junction. add an edge to the junction node and recurse
			// on each possible next position
			if seenNode, ok := seen[[2]int{sx, sy}]; ok {
				c.edges = append(c.edges, edge{seenNode, d})
				break
			}
			currentNode := &node{sx, sy, uint8(len(seen)), nil}
			seen[[2]int{sx, sy}] = currentNode
			c.edges = append(c.edges, edge{currentNode, d})
			for _, n := range nextPossibles {
				buildContractedEdges(grid, n[0], n[1], px, py, ex, ey, currentNode, seen)
			}
			break
		}
	}
}

func writeGraphvizDotFile(seen map[[2]int]*node) {
	f, _ := os.Create("graph.dot")
	_, _ = f.WriteString("digraph {\n")
	for _, n := range seen {
		color := "white"
		if n.id == 0 {
			color = "green"
		} else if n.edges == nil {
			color = "darkgoldenrod1"
		}
		_, _ = f.WriteString(fmt.Sprintf("\t%d [label=\"%d,%d\" fillcolor=\"%s\" style=\"filled\"];\n", n.id, n.x, n.y, color))
		for _, e := range n.edges {
			_, _ = f.WriteString(fmt.Sprintf("\t%d -> %d [label=\"%d\"];\n", n.id, e.to.id, e.d))
		}
	}
	_, _ = f.WriteString("}\n")
	_ = f.Close()
}

func longestPathDfs(n *node, seen []bool) uint {
	if len(n.edges) == 0 {
		return 0
	}
	var m uint
	seen[n.id] = true
	for _, e := range n.edges {
		if !seen[e.to.id] {
			m = max(m, longestPathDfs(e.to, seen)+e.d)
		}
	}
	seen[n.id] = false
	return m
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
	var grid [][]byte
	sx, sy := 0, 0
	for r.Scan() {
		line := r.Text()
		grid = append(grid, []byte(line))
		if len(grid) == 1 {
			sx = bytes.IndexByte(grid[0], '.')
		}
	}
	ex, ey := bytes.IndexByte(grid[len(grid)-1], '.'), len(grid)-1
	start := &node{sx, sy, 0, nil}
	seen := map[[2]int]*node{[2]int{sx, sy}: start}
	buildContractedEdges(grid, sx, sy+1, sx, sy, ex, ey, start, seen)
	writeGraphvizDotFile(seen)
	fmt.Println(longestPathDfs(start, make([]bool, len(seen))))
}
