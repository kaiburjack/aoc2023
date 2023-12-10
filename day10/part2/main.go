package main

import (
	"bufio"
	"os"
)

type pipeKind struct {
	kind       string
	connectors int
}

type pipe struct {
	kind       pipeKind
	neighbours []*pipe
	x, y       int
	isLoop     bool
	g          int
	isOuter    bool
}

var pipeKinds map[string]pipeKind
var directions = [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func init() {
	pipeKinds = make(map[string]pipeKind)
	pipeKinds["F"] = pipeKind{"F", 1 | 2}
	pipeKinds["|"] = pipeKind{"|", 2 | 8}
	pipeKinds["-"] = pipeKind{"-", 1 | 4}
	pipeKinds["L"] = pipeKind{"L", 1 | 8}
	pipeKinds["7"] = pipeKind{"7", 2 | 4}
	pipeKinds["J"] = pipeKind{"J", 4 | 8}
}

func hasPipe(pipes [][]*pipe, x, y int) bool {
	return x >= 0 && y >= 0 && y < len(pipes) && x < len(pipes[y])
}

// pipeKindBySurroundings determines the kind of a pipe based on its
// neighbors and their respective pipe kinds
// We use this for the "S" pipe
func pipeKindBySurroundings(pipes [][]*pipe, pip *pipe) pipeKind {
	for i, o := range directions {
		if hasPipe(pipes, pip.x+o[0], pip.y+o[1]) &&
			pipes[pip.y+o[1]][pip.x+o[0]].kind.connectors&(1<<((i+2)&3)) != 0 {
			pipes[pip.y][pip.x].kind.connectors |= 1 << i
		}
	}
	for _, v := range pipeKinds {
		if v.connectors == pipes[pip.y][pip.x].kind.connectors {
			return v
		}
	}
	panic("no pipe kind found")
}

// connect connects a pipe to its neighbours based on its kind
// which determines the sides at which the pipe can connect to
// its neighbours
func connect(rows [][]*pipe, pip *pipe) {
	for i, offset := range directions {
		if !hasPipe(rows, pip.x+offset[0], pip.y+offset[1]) {
			continue
		}
		neighbour := rows[pip.y+offset[1]][pip.x+offset[0]]
		if pip.kind.connectors&(1<<i) != 0 && neighbour.kind.connectors&(1<<((i+2)&3)) != 0 {
			pip.neighbours = append(pip.neighbours, neighbour)
		}
	}
}

func depthFirstSearchForExit1(pipes [][]*pipe, dir, x, y, g int, constraints int, tested *[]*pipe) bool {
	if !hasPipe(pipes, x, y) {
		// if we reached the edge of the grid, we found an exit
		return true
	}
	pip := pipes[y][x]
	if pip.isOuter {
		// this pipe was already marked as being an outer part
		return true
	} else if pip.g == g {
		// this pipe was already visited in this generation
		return false
	}
	pip.g = g // <- mark this pipe as visited in this generation
	var leftOk, rightOk, upOk, downOk bool
	if !pip.isLoop {
		// if this pipe is not part of the loop,
		// treat it as if it was an empty tile
		// and go in all directions without restrictions
		*tested = append(*tested, pip)
		leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^0, tested)
		rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^0, tested)
		upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
		downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
	} else {
		// test for every case (pipe kind and direction)
		// and put constraints on the recursive calls if
		// we need to "hug" a pipe to squeeze through
		// parallel, adjacent pipes.
		if pip.kind.kind == "-" {
			if dir == 1 {
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, constraints, tested)
				if constraints&4 != 0 {
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
				}
				if constraints&8 != 0 {
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
				}
			} else if dir == 3 {
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, constraints, tested)
				if constraints&4 != 0 {
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
				}
				if constraints&8 != 0 {
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
				}
			} else if dir == 2 {
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^8, tested)
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^8, tested)
			} else if dir == 4 {
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^4, tested)
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^4, tested)
			}
		} else if pip.kind.kind == "|" {
			if dir == 1 {
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^2, tested)
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^2, tested)
			} else if dir == 3 {
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^1, tested)
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^1, tested)
			} else if dir == 2 {
				if constraints&1 != 0 {
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^2, tested)
				}
				if constraints&2 != 0 {
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^1, tested)
				}
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, constraints, tested)
			} else if dir == 4 {
				if constraints&1 != 0 {
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^2, tested)
				}
				if constraints&2 != 0 {
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^1, tested)
				}
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, constraints, tested)
			}
		} else if pip.kind.kind == "F" {
			if dir == 1 {
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^8, tested)
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^2, tested)
			} else if dir == 3 {
				if constraints&4 != 0 {
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^0, tested)
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^2, tested)
				}
				if constraints&8 != 0 {
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^1, tested)
				}
			} else if dir == 2 {
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^0, tested)
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^8, tested)
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^2, tested)
			} else if dir == 4 {
				if constraints&1 != 0 {
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^0, tested)
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^8, tested)
				}
				if constraints&2 != 0 {
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^4, tested)
				}
			}
		} else if pip.kind.kind == "J" {
			if dir == 1 {
				if constraints&4 != 0 {
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^2, tested)
				}
				if constraints&8 != 0 {
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^0, tested)
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^1, tested)
				}
			} else if dir == 3 {
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^4, tested)
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^1, tested)
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
			} else if dir == 2 {
				if constraints&1 != 0 {
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^8, tested)
				}
				if constraints&2 != 0 {
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^0, tested)
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^4, tested)
				}
			} else if dir == 4 {
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^4, tested)
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^0, tested)
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^1, tested)
			}
		} else if pip.kind.kind == "L" {
			if dir == 1 {
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^4, tested)
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^2, tested)
			} else if dir == 3 {
				if constraints&4 != 0 {
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^1, tested)
				}
				if constraints&8 != 0 {
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^0, tested)
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^2, tested)
				}
			} else if dir == 2 {
				if constraints&1 != 0 {
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^0, tested)
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^0, tested)
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^4, tested)
				}
				if constraints&2 != 0 {
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^8, tested)
				}
			} else if dir == 4 {
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^0, tested)
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^4, tested)
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^2, tested)
			}
		} else if pip.kind.kind == "7" {
			if dir == 1 {
				if constraints&4 != 0 {
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^0, tested)
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^1, tested)
				}
				if constraints&8 != 0 {
					downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^2, tested)
				}
			} else if dir == 3 {
				upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^1, tested)
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^8, tested)
			} else if dir == 2 {
				leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^8, tested)
				rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^0, tested)
				downOk = depthFirstSearchForExit1(pipes, 2, x, y+1, g, ^1, tested)
			} else if dir == 4 {
				if constraints&1 != 0 {
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^4, tested)
				}
				if constraints&2 != 0 {
					rightOk = depthFirstSearchForExit1(pipes, 1, x+1, y, g, ^0, tested)
					upOk = depthFirstSearchForExit1(pipes, 4, x, y-1, g, ^0, tested)
					leftOk = depthFirstSearchForExit1(pipes, 3, x-1, y, g, ^8, tested)
				}
			}
		}
	}
	return leftOk || rightOk || upOk || downOk
}

func depthFirstSearchForExit(pipes [][]*pipe) int {
	g := 0 // <- generation of current search
	outer := make(map[*pipe]bool)
	for y := 0; y < len(pipes); y++ {
		for x := 0; x < len(pipes[y]); x++ {
			tested := make([]*pipe, 0)
			if pipes[y][x].isLoop {
				// this pipe belongs to the loop
				// (from start to start), so we can skip it
				continue
			}
			if depthFirstSearchForExit1(pipes, -1, x, y, g, ^0, &tested) {
				// mark all pipes as being outer pipes
				for _, v := range tested {
					outer[v] = true
					v.isOuter = true
				}
			}
		}
	}
	return len(outer)
}

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	rows := make([][]*pipe, 0)
	start := &pipe{}
	totalCount := 0

	// parse input
	for fileScanner.Scan() {
		pipesInRow := fileScanner.Text()
		row := make([]*pipe, 0)
		for _, p := range pipesInRow {
			if p == '.' {
				// here, we handle "." as a pipe with kind ""
				row = append(row, &pipe{kind: pipeKind{"", 0}, x: len(row), y: len(rows), g: 10000000})
				continue
			}
			pip := &pipe{x: len(row), y: len(rows), kind: pipeKinds[string(p)], g: 10000000}
			row = append(row, pip)
			if p == 'S' {
				start = pip
			}
		}
		totalCount += len(row)
		rows = append(rows, row)
	}

	// determine pipe kind of "S"
	start.kind = pipeKindBySurroundings(rows, start)

	// connect all pipes to their neighbours
	for y := 0; y < len(rows); y++ {
		for x := 0; x < len(rows[y]); x++ {
			pip := rows[y][x]
			if pip == nil {
				continue
			}
			connect(rows, pip)
		}
	}

	// find length of the loop
	d := 0
	var cameFrom = start
	var next = start
	for {
		d++
		next.isLoop = true
		// find one of the two neighbors which does
		// not lead us back to where we came from
		for _, n := range next.neighbours {
			if n == cameFrom {
				continue
			}
			cameFrom = next
			next = n
			break
		}
		if next == start {
			// we are back at the start
			break
		}
	}
	outer := depthFirstSearchForExit(rows)
	println(totalCount - outer - d)
}
