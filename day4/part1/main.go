package main

import (
	"bufio"
	"io"
	"os"
)

type parser struct {
	r *bufio.Reader
}

type card struct {
	winningNumbers map[int]bool
	points         int
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

func (p *parser) integer() (int, error) {
	var n int
	for {
		b, _ := p.r.ReadByte()
		if b >= '0' && b <= '9' {
			n = n*10 + int(b-'0')
		} else {
			_ = p.r.UnreadByte()
			return n, nil
		}
	}
}

func (p *parser) ws() error {
	for {
		b, _ := p.r.ReadByte()
		if b != ' ' {
			_ = p.r.UnreadByte()
			break
		}
	}
	return nil
}

func (p *parser) card() (card, error) {
	c := card{winningNumbers: make(map[int]bool)}
	err := p.expect("Card")
	if err == io.EOF {
		return card{}, err //
	}
	_ = p.ws()
	_, _ = p.integer()
	_ = p.expect(":")
	var state int
	for {
		_ = p.ws()
		b, _ := p.r.ReadByte()
		if b == '|' {
			state = 1
		} else if b == '\n' {
			break
		} else {
			_ = p.r.UnreadByte()
		}
		_ = p.ws()
		n, _ := p.integer()
		if state == 0 {
			c.winningNumbers[n] = true
		} else {
			if c.winningNumbers[n] {
				c.points = max(c.points<<1, 1)
			}
		}
	}
	return c, nil
}

func main() {
	file, _ := os.OpenFile("input.txt", os.O_RDONLY, 0)
	r := bufio.NewReader(file)
	p := parser{r}
	sumOfPoints := 0
	for {
		c, err := p.card()
		if err == io.EOF {
			break
		}
		sumOfPoints += c.points
	}
	println(sumOfPoints)
}
