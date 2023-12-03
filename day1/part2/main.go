package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	file, _ := os.OpenFile("input.txt", os.O_RDONLY, 0)
	r := bufio.NewReader(file)
	var sum, state, first, second int
	var sequences [10][]int
	numberDetected := func(n int) {
		if state == 0 {
			first = n
			second = first
			state = 1
		} else {
			second = n
		}
	}
	numbers := [10]string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if b >= '0' && b <= '9' {
			numberDetected(int(b - '0'))
			continue
		} else if b == '\n' {
			sum += first*10 + second
			state = 0
			sequences = [10][]int{}
			continue
		}
		for i := 1; i < 10; i++ {
			if b == numbers[i][0] {
				sequences[i] = append(sequences[i], len(numbers[i]))
			}
			for k := 0; k < len(sequences[i]); k++ {
				s := sequences[i][k]
				if numbers[i][len(numbers[i])-s] == b {
					sequences[i][k]--
					if sequences[i][k] == 0 {
						numberDetected(i)
					} else {
						continue
					}
				}
				sequences[i] = append(sequences[i][:k], sequences[i][k+1:]...)
				k--
			}
		}
	}
	println(sum)
}
