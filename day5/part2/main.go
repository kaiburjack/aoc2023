package main

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/participle/v2"
	"os"
)

// Range represents a linear interval mapping between two intervals (of the same length).
type Range struct {
	DstStart uint64 `@Int`
	SrcStart uint64 `@Int`
	Length   uint64 `@Int`
}

// Interval represents a contiguous interval/sequence of numbers.
type Interval struct {
	Start  uint64 `@Int`
	Length uint64 `@Int`
}

type RangeList struct {
	Ranges []Range `~":"+ ":" @@+`
}

type Input struct {
	Seeds      []Interval  `"seeds" ":" @@+`
	RangeLists []RangeList `@@+`
}

func Map(inputIntervals []Interval, ranges []Range) []Interval {
	var resultIntervals []Interval               // <- final mapping results
	var remainingIntervals = make([]Interval, 0) // <- intervals that are not mapped yet
	for _, r := range ranges {
		inputIntervals = append(inputIntervals, remainingIntervals...)
		remainingIntervals = remainingIntervals[:0]
		for _, interval := range inputIntervals {
			inputIntervals = inputIntervals[1:]
			// Consider all cases:
			if interval.Start >= r.SrcStart && interval.Start+interval.Length <= r.SrcStart+r.Length {
				// interval completely within range
				resultIntervals = append(resultIntervals,
					Interval{r.DstStart + interval.Start - r.SrcStart, interval.Length})
			} else if interval.Start >= r.SrcStart && interval.Start <= r.SrcStart+r.Length {
				// interval.Starting in range but ending outside
				if r.SrcStart+r.Length > interval.Start {
					resultIntervals = append(resultIntervals,
						Interval{r.DstStart + interval.Start - r.SrcStart, r.SrcStart + r.Length - interval.Start})
				}
				remainingIntervals = append(remainingIntervals,
					Interval{r.Length + r.SrcStart, interval.Length - (r.Length - (interval.Start - r.SrcStart))})
			} else if interval.Start+interval.Length >= r.SrcStart && interval.Start+interval.Length <= r.SrcStart+r.Length {
				// interval ending in range but.Starting outside
				if interval.Start+interval.Length > r.SrcStart {
					resultIntervals = append(resultIntervals,
						Interval{r.DstStart, interval.Start + interval.Length - r.SrcStart})
				}
				remainingIntervals = append(remainingIntervals,
					Interval{interval.Start, r.SrcStart - interval.Start})
			} else if interval.Start < r.SrcStart && interval.Start+interval.Length > r.SrcStart+r.Length {
				// range completely within interval
				resultIntervals = append(resultIntervals, Interval{r.DstStart, r.Length})
				remainingIntervals = append(remainingIntervals,
					Interval{interval.Start, r.SrcStart - interval.Start},
					Interval{interval.Start + r.Length, interval.Length - (r.SrcStart + r.Length - interval.Start)})
			} else {
				// range and interval not overlapping
				remainingIntervals = append(remainingIntervals, interval)
			}
		}
	}
	// collect all intervals that couldn't be mapped (mapped by identity function)
	if len(remainingIntervals) > 0 {
		resultIntervals = append(resultIntervals, remainingIntervals...)
	}
	return resultIntervals
}

func main() {
	fileName := "input.txt"
	file, _ := os.OpenFile(fileName, os.O_RDONLY, 0)
	parser, _ := participle.Build[Input]()
	input, _ := parser.Parse(fileName, bufio.NewReader(file))
	intervals := input.Seeds
	for _, r := range input.RangeLists {
		intervals = Map(intervals, r.Ranges)
	}
	var lowest = intervals[0].Start
	for i := 1; i < len(intervals); i++ {
		lowest = min(lowest, intervals[i].Start)
	}
	fmt.Printf("lowest: %d\n", lowest)
}
