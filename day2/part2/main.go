package main

import (
	"bufio"
	"io"
	"os"
)

type game struct {
	draws []draw
}
type draw struct {
	colors [3]int // <- always in the order: [red, green, blue]
}
type parser struct {
	r *bufio.Reader
}

func (p *parser) color() int {
	b, _ := p.r.ReadByte()
	var c int
	switch b {
	case 'r':
		_ = p.expect("ed")
		c = 0
	case 'g':
		_ = p.expect("reen")
		c = 1
	case 'b':
		_ = p.expect("lue")
		c = 2
	}
	return c
}

func (p *parser) game() (game, error) {
	var g game
	err := p.expect("Game ")
	if err == io.EOF {
		return game{}, err
	}
	_ = p.integer()
	_ = p.expect(":")
	for {
		p.ws()
		d, _ := p.draw()
		g.draws = append(g.draws, d)
		b, err := p.r.ReadByte()
		if err == io.EOF {
			break
		}
		if b == '\n' {
			break
		}
	}
	return g, nil
}

func (p *parser) draw() (draw, error) {
	var d draw
	for {
		p.ws()
		count := p.integer()
		p.ws()
		col := p.color()
		d.colors[col] += count
		b, err := p.r.ReadByte()
		if err == io.EOF {
			break
		}
		if b == ',' {
			continue
		} else if b == ';' {
			_ = p.r.UnreadByte()
			break
		} else if b == '\n' {
			_ = p.r.UnreadByte()
			break
		}
	}
	return d, nil
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

func main() {
	file, _ := os.OpenFile("input.txt", os.O_RDONLY, 0)
	r := bufio.NewReader(file)
	p := parser{r}
	var sumOrPowers int64 = 0
	for {
		g, err := p.game()
		if err == io.EOF {
			break
		}
		minimums := [3]int{0, 0, 0}
		for _, d := range g.draws {
			for i := 0; i < 3; i++ {
				if d.colors[i] > minimums[i] {
					minimums[i] = d.colors[i]
				}
			}
		}
		power := int64(minimums[0]) * int64(minimums[1]) * int64(minimums[2])
		sumOrPowers += power
	}
	println(sumOrPowers)
}
