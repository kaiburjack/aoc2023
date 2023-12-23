package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"slices"
)

type node struct {
	id    uint8
	edges []*edge
	x, y  int
}

type edge struct {
	to *node
	d  int
}

var dirs = [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
var arrows = []byte{'v', '>', '^', '<'}

func buildContractedEdges(grid [][]byte, sx, sy, px, py, ex, ey int, c *node, seen map[[2]int]*node) {
	d := 0
	for {
		d++
		v := grid[sy][sx]
		ai := slices.Index(arrows, v)
		var possibleNext [][]int
		for i := 0; i < 4; i++ {
			nx, ny := sx+dirs[i][0], sy+dirs[i][1]
			if nx < 0 || nx >= len(grid[0]) || ny < 0 || ny >= len(grid) ||
				grid[ny][nx] == '#' ||
				ai != -1 && ai != i ||
				px == nx && py == ny {
				continue
			}
			possibleNext = append(possibleNext, []int{nx, ny})
		}
		px, py = sx, sy
		if sx == ex && sy == ey {
			seenNode, ok := seen[[2]int{sx, sy}]
			if ok {
				c.edges = append(c.edges, &edge{seenNode, d})
				break
			}
			endNode := &node{uint8(len(seen)), nil, sx, sy}
			seen[[2]int{sx, sy}] = endNode
			c.edges = append(c.edges, &edge{endNode, d})
			break
		}
		if len(possibleNext) == 0 {
			break
		} else if len(possibleNext) == 1 {
			sx, sy = possibleNext[0][0], possibleNext[0][1]
		} else {
			seenNode, ok := seen[[2]int{sx, sy}]
			if ok {
				c.edges = append(c.edges, &edge{seenNode, d})
				break
			}
			currentNode := &node{uint8(len(seen)), nil, sx, sy}
			seen[[2]int{sx, sy}] = currentNode
			c.edges = append(c.edges, &edge{currentNode, d})
			for _, n := range possibleNext {
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

func longestPathDfs(n *node, seen []bool) float64 {
	if len(n.edges) == 0 {
		return 0
	}
	var m float64
	seen[n.id] = true
	for _, e := range n.edges {
		if !seen[e.to.id] {
			m = math.Max(m, longestPathDfs(e.to, seen)+float64(e.d))
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
	start := &node{0, nil, sx, sy}
	seen := map[[2]int]*node{[2]int{sx, sy}: start}
	buildContractedEdges(grid, sx, sy+1, sx, sy, ex, ey, start, seen)
	writeGraphvizDotFile(seen)
	fmt.Println(longestPathDfs(start, make([]bool, len(seen))))
}
