package main

import (
	"bufio"
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Hand struct {
	cards []byte
	bid   uint64
	type_ int
}

func ofAKind(n int) func([]byte, int) bool {
	return func(hand []byte, numJs int) bool {
		counts := distinctCountsDescending(hand)
		if len(counts) == 0 {
			return numJs >= n
		}
		return counts[0]+numJs >= n
	}
}

func detectTwoPairs(hand []byte, numJs int) bool {
	counts := distinctCountsDescending(hand)
	return counts[0] >= 2 && counts[1]+numJs >= 2
}

func distinctCountsDescending(hand []byte) []int {
	different := make(map[byte]int)
	for i := 0; i < len(hand); i++ {
		different[hand[i]]++
	}
	counts := make([]int, 0, len(different))
	for _, v := range different {
		counts = append(counts, v)
	}
	slices.Sort(counts)
	slices.Reverse(counts)
	return counts
}

func detectFullHouse(cardsWithoutJs []byte, numJs int) bool {
	counts := distinctCountsDescending(cardsWithoutJs)
	if len(counts) == 1 {
		return counts[0]+numJs >= 5
	} else if len(counts) == 2 {
		return counts[0] == 3 && counts[1]+numJs >= 2 || counts[0] == 2 && counts[1]+numJs >= 3
	}
	return false
}

var tests = []func([]byte, int) bool{
	ofAKind(5),
	ofAKind(4),
	detectFullHouse,
	ofAKind(3),
	detectTwoPairs,
	ofAKind(2),
}

func determineTypeOfHand(cardsWithoutJs []byte, numJs int) int {
	for i, v := range tests {
		if v(cardsWithoutJs, numJs) {
			return len(tests) - i
		}
	}
	return 0
}

func firstDifference(a, b []byte) int {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return -1
}

var ORDER = []byte("AKQT98765432J")

func main() {
	readFile, _ := os.Open("input.txt")
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	hands := make([]Hand, 0)
	for fileScanner.Scan() {
		handAndBid := strings.Split(fileScanner.Text(), " ")
		cards := handAndBid[0]
		bid, _ := strconv.Atoi(handAndBid[1])
		cardsWithoutJs := []byte(strings.ReplaceAll(cards, "J", ""))
		slices.Sort(cardsWithoutJs)
		type_ := determineTypeOfHand(cardsWithoutJs, strings.Count(cards, "J"))
		h := Hand{cards: []byte(cards), bid: uint64(bid), type_: type_}
		hands = append(hands, h)
	}

	slices.SortFunc(hands, func(a, b Hand) int {
		if a.type_ < b.type_ {
			return -1
		} else if a.type_ > b.type_ {
			return 1
		} else {
			i := firstDifference(a.cards, b.cards)
			k := bytes.IndexByte(ORDER, a.cards[i])
			j := bytes.IndexByte(ORDER, b.cards[i])
			if k < j {
				return 1
			} else if k > j {
				return -1
			} else {
				return 0
			}
		}
	})

	var rank uint64 = 1
	var sum uint64 = 0
	for _, hand := range hands {
		sum += rank * hand.bid
		rank++
	}
	println(sum)
}
