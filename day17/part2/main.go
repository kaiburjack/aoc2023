package main

import (
	"bufio"
	"container/heap"
	"math"
	"os"
)

// build a priority queue
type priorityQueue []*vertex

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
	*v = append(*v, x.(*vertex))
}

func (v *priorityQueue) Pop() interface{} {
	old := *v
	n := len(old)
	x := old[n-1]
	*v = old[0 : n-1]
	return x
}

type vertex struct {
	d            int
	x, y, dx, dy int
	prev         *vertex
	heatLoss     int
	isPartOfPath bool
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)

	// build heat loss grid
	var grid [][]vertex
	for r.Scan() {
		var row []vertex
		for _, c := range r.Bytes() {
			row = append(row, vertex{0, 0, 0, 0, 0, nil, int(c - '0'), false})
		}
		grid = append(grid, row)
	}

	var t priorityQueue
	heap.Push(&t, &vertex{0, 0, 0, 0, 0, nil, 0, false})
	var visited = make(map[vertex]int)
	var current *vertex
	const MIN, MAX = 4, 10

	// we need to visit some more nodes
	for len(t) > 0 {
		current = heap.Pop(&t).(*vertex)
		// check if current vertex is the end and we reached it with the minimum amount of steps in the same direction
		if current.x == len(grid[0])-1 && current.y == len(grid)-1 && (math.Abs(float64(current.dx)) >= MIN || math.Abs(float64(current.dy)) >= MIN) {
			break
		}

		// check if we can go in any direction
		if current.x > 0 && (current.dx == 0 && current.dy == 0 || current.dx < 0 || math.Abs(float64(current.dy)) >= MIN) && current.dx > -MAX {
			s := vertex{0, current.x - 1, current.y, current.dx - 1, 0, nil, 0, false}
			nextCost := current.d + grid[current.y][current.x-1].heatLoss
			if v, ok := visited[s]; !ok || v > nextCost {
				visited[s] = nextCost
				heap.Push(&t, &vertex{nextCost, current.x - 1, current.y, current.dx - 1, 0, current, 0, false})
			}
		}
		if current.x < len(grid[0])-1 && (current.dx == 0 && current.dy == 0 || current.dx > 0 || math.Abs(float64(current.dy)) >= MIN) && current.dx < MAX {
			state := vertex{0, current.x + 1, current.y, current.dx + 1, 0, nil, 0, false}
			nextCost := current.d + grid[current.y][current.x+1].heatLoss
			if s, ok := visited[state]; !ok || s > nextCost {
				visited[state] = nextCost
				heap.Push(&t, &vertex{nextCost, current.x + 1, current.y, current.dx + 1, 0, current, 0, false})
			}
		}
		if current.y > 0 && (current.dx == 0 && current.dy == 0 || current.dy < 0 || math.Abs(float64(current.dx)) >= MIN) && current.dy > -MAX {
			state := vertex{0, current.x, current.y - 1, 0, current.dy - 1, nil, 0, false}
			nextCost := current.d + grid[current.y-1][current.x].heatLoss
			if s, ok := visited[state]; !ok || s > nextCost {
				visited[state] = nextCost
				heap.Push(&t, &vertex{nextCost, current.x, current.y - 1, 0, current.dy - 1, current, 0, false})
			}
		}
		if current.y < len(grid)-1 && (current.dx == 0 && current.dy == 0 || current.dy > 0 || math.Abs(float64(current.dx)) >= MIN) && current.dy < MAX {
			state := vertex{0, current.x, current.y + 1, 0, current.dy + 1, nil, 0, false}
			nextCost := current.d + grid[current.y+1][current.x].heatLoss
			if s, ok := visited[state]; !ok || s > nextCost {
				visited[state] = nextCost
				heap.Push(&t, &vertex{nextCost, current.x, current.y + 1, 0, current.dy + 1, current, 0, false})
			}
		}
	}

	// print board with path
	d := current.d
	for current.prev != nil {
		grid[current.y][current.x].isPartOfPath = true
		current = current.prev
	}
	grid[0][0].isPartOfPath = true
	for _, row := range grid {
		for _, v := range row {
			if v.isPartOfPath {
				print("X")
			} else {
				print(".")
			}
		}
		println()
	}

	println(d)
}
