package main

import (
	"bufio"
	"container/heap"
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
	var current pathVertex
	const MAX = 3
	// we need to visit some more nodes
	for len(q) > 0 {
		current = heap.Pop(&q).(pathVertex)
		// check if we reached the end
		if current.x == len(grid[0])-1 && current.y == len(grid)-1 {
			break
		}

		// check if we can go in any direction
		if current.x > 0 && current.dx <= 0 && current.dx > -MAX {
			v := visitedVertex{current.x - 1, current.y, current.dx - 1, 0}
			if _, ok := visited[v]; !ok {
				nextCost := current.d + grid[current.y][current.x-1]
				visited[v] = struct{}{}
				heap.Push(&q, pathVertex{nextCost, current.x - 1, current.y, current.dx - 1, 0})
			}
		}
		if current.x < len(grid[0])-1 && current.dx >= 0 && current.dx < MAX {
			v := visitedVertex{current.x + 1, current.y, current.dx + 1, 0}
			if _, ok := visited[v]; !ok {
				nextCost := current.d + grid[current.y][current.x+1]
				visited[v] = struct{}{}
				heap.Push(&q, pathVertex{nextCost, current.x + 1, current.y, current.dx + 1, 0})
			}
		}
		if current.y > 0 && current.dy <= 0 && current.dy > -MAX {
			v := visitedVertex{current.x, current.y - 1, 0, current.dy - 1}
			if _, ok := visited[v]; !ok {
				nextCost := current.d + grid[current.y-1][current.x]
				visited[v] = struct{}{}
				heap.Push(&q, pathVertex{nextCost, current.x, current.y - 1, 0, current.dy - 1})
			}
		}
		if current.y < len(grid)-1 && current.dy >= 0 && current.dy < MAX {
			v := visitedVertex{current.x, current.y + 1, 0, current.dy + 1}
			if _, ok := visited[v]; !ok {
				nextCost := current.d + grid[current.y+1][current.x]
				visited[v] = struct{}{}
				heap.Push(&q, pathVertex{nextCost, current.x, current.y + 1, 0, current.dy + 1})
			}
		}
	}
	println(current.d)
}
