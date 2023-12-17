package main

import (
	"bufio"
	"container/heap"
	"math"
	"os"
)

// build a priority queue using the container/heap package
type priorityQueue []pathVertex

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
	*v = append(*v, x.(pathVertex))
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
	d, x, y, dx, dy int
}

// vertices that we mark as already visited. These are NOT just the vertices of the grid,
// so don't _just_ depend on the position in the grid, but also on the direction we came from.
type visitedVertex struct {
	x, y, dx, dy int
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)

	// build heat loss grid
	var grid [][]int
	for r.Scan() {
		var row []int
		for _, c := range r.Bytes() {
			row = append(row, int(c-'0'))
		}
		grid = append(grid, row)
	}

	var q priorityQueue
	heap.Push(&q, pathVertex{})
	var visited = make(map[visitedVertex]struct{})
	var c pathVertex
	const MIN, MAX = 4, 10
	// we need to visit some more nodes
	for len(q) > 0 {
		c = heap.Pop(&q).(pathVertex)
		// check if c vertex is the end and we reached it with the minimum amount of steps in the same direction
		if c.x == len(grid[0])-1 && c.y == len(grid)-1 && (math.Abs(float64(c.dx)) >= MIN || math.Abs(float64(c.dy)) >= MIN) {
			break
		}

		// check if we can go in any direction
		if c.x > 0 && (c.dx == 0 && c.dy == 0 || c.dx < 0 || math.Abs(float64(c.dy)) >= MIN) && c.dx > -MAX {
			v := visitedVertex{c.x - 1, c.y, c.dx - 1, 0}
			if _, ok := visited[v]; !ok {
				nextCost := c.d + grid[c.y][c.x-1]
				visited[v] = struct{}{}
				heap.Push(&q, pathVertex{nextCost, c.x - 1, c.y, c.dx - 1, 0})
			}
		}
		if c.x < len(grid[0])-1 && (c.dx == 0 && c.dy == 0 || c.dx > 0 || math.Abs(float64(c.dy)) >= MIN) && c.dx < MAX {
			v := visitedVertex{c.x + 1, c.y, c.dx + 1, 0}
			if _, ok := visited[v]; !ok {
				nextCost := c.d + grid[c.y][c.x+1]
				visited[v] = struct{}{}
				heap.Push(&q, pathVertex{nextCost, c.x + 1, c.y, c.dx + 1, 0})
			}
		}
		if c.y > 0 && (c.dx == 0 && c.dy == 0 || c.dy < 0 || math.Abs(float64(c.dx)) >= MIN) && c.dy > -MAX {
			v := visitedVertex{c.x, c.y - 1, 0, c.dy - 1}
			if _, ok := visited[v]; !ok {
				nextCost := c.d + grid[c.y-1][c.x]
				visited[v] = struct{}{}
				heap.Push(&q, pathVertex{nextCost, c.x, c.y - 1, 0, c.dy - 1})
			}
		}
		if c.y < len(grid)-1 && (c.dx == 0 && c.dy == 0 || c.dy > 0 || math.Abs(float64(c.dx)) >= MIN) && c.dy < MAX {
			v := visitedVertex{c.x, c.y + 1, 0, c.dy + 1}
			if _, ok := visited[v]; !ok {
				nextCost := c.d + grid[c.y+1][c.x]
				visited[v] = struct{}{}
				heap.Push(&q, pathVertex{nextCost, c.x, c.y + 1, 0, c.dy + 1})
			}
		}
	}
	println(c.d)
}
