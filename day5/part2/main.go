package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Range represents a linear interval mapping between two ranges.
type Range struct {
	dstStart uint64
	srcStart uint64
	length   uint64
}

type Interval struct {
	start  uint64
	length uint64
}

//
// ALL OF BELOW IS JUST PARSER INFRASTRUCTURE :D
//

type parser struct {
	r *bufio.Reader
}

func (p *parser) ws() {
	for {
		b, err := p.r.ReadByte()
		if err == io.EOF {
			return
		}
		if b != ' ' {
			_ = p.r.UnreadByte()
			break
		}
	}
}

func (p *parser) integer() uint64 {
	var n uint64
	for {
		b, _ := p.r.ReadByte()
		if b >= '0' && b <= '9' {
			n = n*10 + uint64(b-'0')
		} else {
			_ = p.r.UnreadByte()
			return n
		}
	}
}

func (p *parser) Range() Range {
	dstStart := p.integer()
	p.ws()
	srcStart := p.integer()
	p.ws()
	length := p.integer()
	return Range{dstStart, srcStart, length}
}

func (p *parser) RangeList() ([]Range, error) {
	err := p.until(':')
	if err == io.EOF {
		return nil, err //
	}
	_ = p.expect("\n")
	var ranges []Range
	for {
		r := p.Range()
		ranges = append(ranges, r)
		_ = p.expect("\n")
		b, err := p.r.ReadByte()
		if err == io.EOF {
			return ranges, nil
		}
		if b == '\n' {
			break
		} else {
			_ = p.r.UnreadByte()
		}
	}
	return ranges, nil
}

func (p *parser) seed() Interval {
	p.ws()
	start := p.integer()
	p.ws()
	length := p.integer()
	return Interval{start, length}
}

func (p *parser) Seeds() []Interval {
	var seeds []Interval
	_ = p.expect("seeds:")
	for {
		s := p.seed()
		seeds = append(seeds, s)
		b, _ := p.r.ReadByte()
		if b == '\n' {
			break
		}
	}
	return seeds
}

func (p *parser) expect(s string) error {
	for i := 0; i < len(s); i++ {
		b, _ := p.r.ReadByte()
		if b != s[i] {
			return fmt.Errorf("expected %s, got %s", s, string(b))
		}
	}
	return nil
}

func (p *parser) until(e byte) error {
	for {
		b, err := p.r.ReadByte()
		if err == io.EOF {
			return err //
		}
		if b == e {
			break
		}
	}
	return nil
}

//
// END OF PARSER INFRASTRUCTURE :D
//

// Below is the actual "implementation"

func Map(inputIntervals []Interval, rl []Range) []Interval {
	var resultIntervals []Interval               // <- final mapping results
	var remainingIntervals = make([]Interval, 0) // <- intervals that are not mapped yet
	for _, r := range rl {
		inputIntervals = append(inputIntervals, remainingIntervals...)
		remainingIntervals = remainingIntervals[:0]
		for _, interval := range inputIntervals {
			inputIntervals = inputIntervals[1:]
			// Consider all cases:
			if interval.start >= r.srcStart && interval.start+interval.length <= r.srcStart+r.length {
				// interval completely within range
				resultIntervals = append(resultIntervals, Interval{r.dstStart + interval.start - r.srcStart, interval.length})
			} else if interval.start >= r.srcStart && interval.start <= r.srcStart+r.length {
				// interval starting in range but ending outside
				if r.srcStart+r.length > interval.start {
					resultIntervals = append(resultIntervals, Interval{r.dstStart + interval.start - r.srcStart, r.srcStart + r.length - interval.start})
				}
				remainingIntervals = append(remainingIntervals, Interval{interval.start + r.length - (interval.start - r.srcStart), interval.length - (r.length - (interval.start - r.srcStart))})
			} else if interval.start+interval.length >= r.srcStart && interval.start+interval.length <= r.srcStart+r.length {
				// interval ending in range but starting outside
				if interval.start+interval.length > r.srcStart {
					resultIntervals = append(resultIntervals, Interval{r.dstStart, interval.start + interval.length - r.srcStart})
				}
				remainingIntervals = append(remainingIntervals, Interval{interval.start, r.srcStart - interval.start})
			} else if interval.start < r.srcStart && interval.start+interval.length > r.srcStart+r.length {
				// range completely within interval
				resultIntervals = append(resultIntervals, Interval{r.dstStart, r.length})
				remainingIntervals = append(remainingIntervals, Interval{interval.start, r.srcStart - interval.start})
				remainingIntervals = append(remainingIntervals, Interval{interval.start + r.length, interval.length - (r.srcStart + r.length - interval.start)})
			} else {
				// range and interval not overlapping
				remainingIntervals = append(remainingIntervals, interval)
			}
		}
	}
	if len(remainingIntervals) > 0 {
		resultIntervals = append(resultIntervals, remainingIntervals...)
	}
	return resultIntervals
}

func main() {
	file, _ := os.OpenFile("input.txt", os.O_RDONLY, 0)
	r := bufio.NewReader(file)
	p := parser{r}
	// parse list of seed intervals
	intervals := p.Seeds()

	// parse list of interval mappings
	for {
		rl, err := p.RangeList()
		if err == io.EOF {
			break
		}
		intervals = Map(intervals, rl)
	}

	// find lowest number
	var lowest = intervals[0].start
	for i := 1; i < len(intervals); i++ {
		lowest = min(lowest, intervals[i].start)
	}
	fmt.Printf("lowest: %d\n", lowest)
}
