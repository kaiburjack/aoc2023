package main

import (
	"bufio"
	"os"
)

type pipeKind struct {
	connectors int
}

type pipe struct {
	kind       pipeKind
	neighbours []*pipe
	x, y       int
}

var pipeKinds map[string]pipeKind
var directions = [4][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func init() {
	pipeKinds = make(map[string]pipeKind)
	pipeKinds["F"] = pipeKind{1 | 2}
	pipeKinds["|"] = pipeKind{2 | 8}
	pipeKinds["-"] = pipeKind{1 | 4}
	pipeKinds["L"] = pipeKind{1 | 8}
	pipeKinds["7"] = pipeKind{2 | 4}
	pipeKinds["J"] = pipeKind{4 | 8}
}

func hasPipe(pipes [][]*pipe, x, y int) bool {
	return x >= 0 && y >= 0 && y < len(pipes) && x < len(pipes[y]) && pipes[y][x] != nil
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

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	rows := make([][]*pipe, 0)
	start := &pipe{}

	// parse input
	for fileScanner.Scan() {
		pipesInRow := fileScanner.Text()
		row := make([]*pipe, 0)
		for _, p := range pipesInRow {
			if p == '.' {
				// here, we handle "." as nil/null pipe
				row = append(row, nil)
				continue
			}
			pip := &pipe{x: len(row), y: len(rows), kind: pipeKinds[string(p)]}
			row = append(row, pip)
			if p == 'S' {
				start = pip
			}
		}
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
	var last = start
	var next = start
	for {
		d++
		// find one of the two neighbors which does
		// not lead us back to where we came from
		n := next.neighbours[0]
		if n == last {
			n = next.neighbours[1]
		}
		last = next
		next = n
		if next == start {
			// we are back at the start
			break
		}
	}

	// result is half of the loop length
	println(d / 2)
}
