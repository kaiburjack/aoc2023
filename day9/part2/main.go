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
	differences := make([]int64, len(numbers)-1)
	for i := 0; i < len(numbers)-1; i++ {
		differences[i] = numbers[i+1] - numbers[i]
	}
	return numbers[0] - extrapolationsSum(differences)
}

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var sumOfSums int64
	for fileScanner.Scan() {
		words := strings.Split(fileScanner.Text(), " ")
		numbers := make([]int64, len(words))
		for i, word := range words {
			numbers[i], _ = strconv.ParseInt(word, 10, 64)
		}
		sumOfSums += extrapolationsSum(numbers)
	}
	println(sumOfSums)
}
