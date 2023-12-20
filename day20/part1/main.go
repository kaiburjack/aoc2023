package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type wire struct {
	sender   *module
	receiver *module
	state    int // <- for conjunctions to remember input state
}
type module struct {
	name    string
	ty      uint8
	outputs []*wire
	inputs  []*wire
	state   int // <- needed for flip-flops to remember own state
}
type pulse struct {
	w     *wire
	state int
}

func bfs(b *module, sent *[2]int64) {
	var queue []pulse
	sent[0]++ // <- initial button press sent to broadcaster
	for _, output := range b.outputs {
		queue = append(queue, pulse{state: 0, w: output})
	}
	for len(queue) > 0 { // <- while we have pulses to send
		p := queue[0] // <- pop pulse from queue
		queue = queue[1:]
		p.w.state = p.state                         // <- set pulse state on wire
		sent[p.state]++                             // <- count sent pulses
		if p.w.receiver.ty == '%' && p.state == 0 { // <- flip-flop
			p.w.receiver.state = 1 - p.w.receiver.state // <- swap state
			for _, output := range p.w.receiver.outputs {
				queue = append(queue, pulse{state: p.w.receiver.state, w: output})
			}
		} else if p.w.receiver.ty == '&' { // <- conjunction
			allInputsHigh := slices.IndexFunc(p.w.receiver.inputs, func(rw *wire) bool {
				return rw.state == 0 // <- low found
			}) == -1 // <- no low found?
			for _, output := range p.w.receiver.outputs {
				if allInputsHigh {
					queue = append(queue, pulse{state: 0, w: output})
				} else {
					queue = append(queue, pulse{state: 1, w: output})
				}
			}
		}
	}
}

func main() {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
	m := make(map[string]*module)
	var broadcaster *module
	for r.Scan() {
		line := r.Text()
		splitted := strings.Split(line, "->")
		typeAndName := strings.TrimSpace(splitted[0])
		var ty uint8
		var name string
		var mod, omod *module
		var ok bool

		// ensure module is created
		if typeAndName == "broadcaster" {
			name = typeAndName
		} else {
			ty = typeAndName[0]
			name = typeAndName[1:]
		}
		if mod, ok = m[name]; !ok {
			mod = &module{name: name}
			m[name] = mod
		}
		mod.ty = ty
		if typeAndName == "broadcaster" {
			broadcaster = mod
		}

		// add output -> input wires
		for _, output := range strings.Split(splitted[1], ",") {
			name = strings.TrimSpace(output)
			if omod, ok = m[name]; !ok {
				omod = &module{name: name}
				m[name] = omod
			}
			w := &wire{receiver: omod, sender: mod}
			mod.outputs = append(mod.outputs, w)
			omod.inputs = append(omod.inputs, w)
		}
	}

	var sent [2]int64
	for i := 0; i < 1000; i++ {
		bfs(broadcaster, &sent)
	}
	fmt.Printf("lo: %d, hi: %d, product: %d\n", sent[0], sent[1], sent[0]*sent[1])
}
