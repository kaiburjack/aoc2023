package main

import (
	"bufio"
	"bytes"
	"os"
)

// rotateRight rotates the board 90 degrees to the right.
// This is so that we can reuse the "tilt north" function below,
// since after rotating the board one step to the right, the previous "north"
// is now the new "west" and we can reuse the "tilt north" function below.
// This function returns a copy of the board.
func rotateRight(board [][]byte) [][]byte {
	result := make([][]byte, len(board[0]))
	for i := range result {
		result[i] = make([]byte, len(board))
		for j := range result[i] {
			result[i][j] = board[len(board)-j-1][i]
		}
	}
	return result
}

// boardsEqual checks if two boards are equal.
func boardsEqual(a [][]byte, b [][]byte) bool {
	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

// copyBoard copies the board.
func copyBoard(board [][]byte) [][]byte {
	result := make([][]byte, len(board))
	for i := range result {
		result[i] = make([]byte, len(board[0]))
		copy(result[i], board[i])
	}
	return result
}

// tiltNorthMut tilts the board "north".
// This function mutates the board.
func tiltNorthMut(board [][]byte) [][]byte {
	var nextObstacle [100]int
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			c := board[y][x]
			switch c {
			case '#':
				nextObstacle[x] = y + 1
			case 'O':
				board[y][x] = '.'
				board[nextObstacle[x]][x] = 'O'
				nextObstacle[x]++
			}
		}
	}
	return board
}

// simulateOneRoundMut runs one round of tilt and rotate operations
// each for north, west, south and east.
// This function mutates the board.
func simulateOneRoundMut(board [][]byte) [][]byte {
	for i := 0; i < 4; i++ {
		board = rotateRight(tiltNorthMut(board))
	}
	return board
}

// Evaluate the "north load" of the board.
func eval(board [][]byte) int {
	var numRocks [100]int
	var total int
	for y := 0; y < len(board); y++ {
		for i, c := range board[y] {
			total += numRocks[i]
			if c == 'O' {
				numRocks[i]++
				total++
			}
		}
	}
	return total
}

func main() {
	// parse the input into a 2D slice
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	var board [][]byte
	for fileScanner.Scan() {
		board = append(board, []byte(fileScanner.Text()))
	}

	// remember seen boards to detect loops
	var seenBoards [][][]byte

	// simulate the board for (at most) N rounds
	// we expect to (at some point) reach a loop where
	// we see a board we've seen before
	const N = 1000000000
	loopFound := false
	for round := 0; round < N && !loopFound; round++ {
		// actually simulate one round (north, west, south, east)
		board = simulateOneRoundMut(board)
		// check if we saw this resulting board before
		for seenIndex, seenBoard := range seenBoards {
			if !boardsEqual(board, seenBoard) {
				continue
			}
			// compute the length of the loop (current round index - index of first seen)
			loopLen := round - seenIndex
			// compute the remaining rounds that we need to simulate.
			// Because we already simulated "round" rounds before we saw the loop,
			// we need to simulate (N - round) additional rounds.
			// But we can skip all loops of length 'loopLen', so we
			// need to compute the remainder of (N - round) / loopLen.
			remainingRounds := (N - round - 1) % loopLen
			// simulate the remaining rounds
			for j := 0; j < remainingRounds; j++ {
				board = simulateOneRoundMut(board)
			}
			loopFound = true
			break
		}
		// remember this board for later to check for a loop
		seenBoards = append(seenBoards, copyBoard(board))
	}

	// evaluate the "north load" of the board
	println(eval(board))
}
