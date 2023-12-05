package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Range represents a linear interval mapping between two intervals (of the same length).
type Range struct {
	dstStart uint64
	srcStart uint64
	length   uint64
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

func (p *parser) seed() uint64 {
	p.ws()
	return p.integer()
}

func (p *parser) Seeds() []uint64 {
	var seeds []uint64
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

func Map(src uint64, rl []Range) uint64 {
	for i := 0; i < len(rl); i++ {
		if src >= rl[i].srcStart && src <= rl[i].srcStart+rl[i].length {
			return rl[i].dstStart + (src - rl[i].srcStart)
		}
	}
	return src
}

// The idea here is to start with a list of seeds, and then apply the
// interval mappings one after the other. The result is a list of numbers
// whose lowest number is the answer.
func main() {
	file, _ := os.OpenFile("input.txt", os.O_RDONLY, 0)
	r := bufio.NewReader(file)
	p := parser{r}
	// parse list of seeds
	numbers := p.Seeds()

	// parse list of interval mappings
	for {
		rl, err := p.RangeList()
		if err == io.EOF {
			break
		}
		for i, s := range numbers {
			result := Map(s, rl)
			numbers[i] = result
		}
	}

	// find lowest number
	var lowest = numbers[0]
	for i := 1; i < len(numbers); i++ {
		lowest = min(lowest, numbers[i])
	}
	fmt.Printf("lowest: %d\n", lowest)
}
