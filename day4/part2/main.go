package main

import (
	"bufio"
	"io"
	"os"
)

type parser struct {
	r *bufio.Reader
}

func (p *parser) expect(s string) error {
	for i := 0; i < len(s); i++ {
		_, err := p.r.ReadByte()
		if err == io.EOF {
			return err
		}
	}
	return nil
}

func (p *parser) integer() int {
	var n int
	for {
		b, _ := p.r.ReadByte()
		if b >= '0' && b <= '9' {
			n = n*10 + int(b-'0')
		} else {
			_ = p.r.UnreadByte()
			return n
		}
	}
}

func (p *parser) ws() {
	for {
		b, _ := p.r.ReadByte()
		if b != ' ' {
			_ = p.r.UnreadByte()
			break
		}
	}
}

func (p *parser) card() (int, error) {
	winningNumbers := make(map[int]bool)
	numberOfWinningNumbers := 0
	err := p.expect("Card")
	if err == io.EOF {
		return 0, err
	}
	p.ws()
	_ = p.integer()
	_ = p.expect(":")
	var state int
	for {
		p.ws()
		b, _ := p.r.ReadByte()
		if b == '|' {
			state = 1
		} else if b == '\n' {
			break
		} else {
			_ = p.r.UnreadByte()
		}
		p.ws()
		n := p.integer()
		if state == 0 {
			winningNumbers[n] = true
		} else {
			if winningNumbers[n] {
				numberOfWinningNumbers++
			}
		}
	}
	return numberOfWinningNumbers, nil
}

func main() {
	file, _ := os.OpenFile("input.txt", os.O_RDONLY, 0)
	r := bufio.NewReader(file)
	p := parser{r}
	var wonCopies []int64
	var numWonCards int64 = 0
	for {
		numberOfWinningNumbers, err := p.card()
		if err == io.EOF {
			break
		}
		var count int64 = 1
		if len(wonCopies) > 0 {
			count = 1 + wonCopies[0]
		} else {
			wonCopies = append(wonCopies, 1)
		}
		numWonCards += count
		wonCopies = wonCopies[1:]
		for i := 0; i < numberOfWinningNumbers; i++ {
			if len(wonCopies) > i {
				wonCopies[i] += count
			} else {
				wonCopies = append(wonCopies, count)
			}
		}
	}
	println(numWonCards)
}
