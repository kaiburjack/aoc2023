package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func extrapolationsSum(numbers []int64) int64 {
	if len(numbers) <= 1 {
		return 0
	}
	for i := 0; i < len(numbers)-1; i++ {
		numbers[i] = numbers[i+1] - numbers[i]
	}
	return numbers[len(numbers)-1] + extrapolationsSum(numbers[:len(numbers)-1])
}

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	var sumOfSums int64
	numbers := make([]int64, 0)
	for fileScanner.Scan() {
		numbers = numbers[:0]
		words := strings.Split(fileScanner.Text(), " ")
		for _, word := range words {
			n, _ := strconv.ParseInt(word, 10, 64)
			numbers = append(numbers, n)
		}
		sumOfSums += extrapolationsSum(numbers)
	}
	println(sumOfSums)
}
