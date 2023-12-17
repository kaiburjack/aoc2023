package main

import (
	"bufio"
	"container/heap"
	"os"
)

// build a priority queue using the container/heap package
type priorityQueue []*pathVertex

func (v *priorityQueue) Less(i, j int) bool {
	return (*v)[i].d < (*v)[j].d
}
func (v *priorityQueue) Len() int {
	return len(*v)
}
func (v *priorityQueue) Swap(i, j int) {
	(*v)[i], (*v)[j] = (*v)[j], (*v)[i]
}
func (v *priorityQueue) Push(x interface{}) {
	*v = append(*v, x.(*pathVertex))
}
func (v *priorityQueue) Pop() interface{} {
	old := *v
	n := len(old)
	x := old[n-1]
	*v = old[0 : n-1]
	return x
}

// vertices that we follow via the priority queue and which form the path
type pathVertex struct {
	d      uint
	x, y   uint8
	dx, dy int8
}

// vertices that we mark as already visited. These are NOT just the vertices of the grid,
// so don't _just_ depend on the position in the grid, but also on the direction we came from.
type visitedVertex struct {
	x, y   uint8
	dx, dy int8
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)

	// build heat loss grid
	var grid [][]uint8
	for r.Scan() {
		var row []uint8
		for _, c := range r.Bytes() {
			row = append(row, c-'0')
		}
		grid = append(grid, row)
	}
	w, h := uint8(len(grid[0])), uint8(len(grid))

	var q priorityQueue
	heap.Push(&q, &pathVertex{})
	var visited = make(map[visitedVertex]struct{})
	var c *pathVertex
	const MAX = 3
	// we need to visit some more nodes
	for len(q) > 0 {
		c = heap.Pop(&q).(*pathVertex)
		// check if we reached the end
		if c.x == w-1 && c.y == h-1 {
			break
		}

		// check if we can go in any direction
		if c.x > 0 && c.dx <= 0 && c.dx > -MAX {
			v := visitedVertex{c.x - 1, c.y, c.dx - 1, 0}
			if _, ok := visited[v]; !ok {
				visited[v] = struct{}{}
				heap.Push(&q, &pathVertex{c.d + uint(grid[c.y][c.x-1]), c.x - 1, c.y, c.dx - 1, 0})
			}
		}
		if c.x < w-1 && c.dx >= 0 && c.dx < MAX {
			v := visitedVertex{c.x + 1, c.y, c.dx + 1, 0}
			if _, ok := visited[v]; !ok {
				visited[v] = struct{}{}
				heap.Push(&q, &pathVertex{c.d + uint(grid[c.y][c.x+1]), c.x + 1, c.y, c.dx + 1, 0})
			}
		}
		if c.y > 0 && c.dy <= 0 && c.dy > -MAX {
			v := visitedVertex{c.x, c.y - 1, 0, c.dy - 1}
			if _, ok := visited[v]; !ok {
				visited[v] = struct{}{}
				heap.Push(&q, &pathVertex{c.d + uint(grid[c.y-1][c.x]), c.x, c.y - 1, 0, c.dy - 1})
			}
		}
		if c.y < h-1 && c.dy >= 0 && c.dy < MAX {
			v := visitedVertex{c.x, c.y + 1, 0, c.dy + 1}
			if _, ok := visited[v]; !ok {
				visited[v] = struct{}{}
				heap.Push(&q, &pathVertex{c.d + uint(grid[c.y+1][c.x]), c.x, c.y + 1, 0, c.dy + 1})
			}
		}
	}
	println(c.d)
}
