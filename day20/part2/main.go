package main

import (
	"bufio"
	"os"
	"slices"
	"strings"
)

type wire struct {
	sender   *module
	receiver *module
	state    uint8 // <- for conjunctions to remember input state
}
type module struct {
	name    string
	ty      uint8
	outputs []*wire
	inputs  []*wire
	state   uint8 // <- needed for flip-flops to remember own state
}
type pulse struct {
	w     *wire
	state uint8
}

func bfs(b *module, btnPress uint64, receivers []*wire, found []uint64) {
	var queue []pulse
	for _, output := range b.outputs {
		queue = append(queue, pulse{state: 0, w: output})
	}
	for len(queue) > 0 { // <- while we have pulses to send
		p := queue[0] // <- pop pulse from queue
		queue = queue[1:]
		p.w.state = p.state // <- set pulse state on wire
		if p.state == 0 {
			// check which of the receivers received a low pulse
			if idx := slices.IndexFunc(receivers, func(r *wire) bool {
				return r.sender == p.w.receiver
			}); idx != -1 && found[idx] == 0 {
				// mark receiver as found
				found[idx] = btnPress
			}
		}
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
func gcd(a, b uint64) uint64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
func lcm(a, b uint64, rest ...uint64) uint64 {
	result := a * b / gcd(a, b)
	for i := 0; i < len(rest); i++ {
		result = lcm(result, rest[i])
	}
	return result
}
func parseAndBuildGraph(m map[string]*module) *module {
	file, _ := os.Open("input.txt")
	r := bufio.NewScanner(file)
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
	return broadcaster
}
func main() {
	m := make(map[string]*module)
	broadcaster := parseAndBuildGraph(m)

	// assume and verify multiple independent cycles behind a single inverter
	rx := m["rx"]
	if len(rx.inputs) != 1 {
		panic("rx.inputs != 1")
	}
	rxInput := rx.inputs[0].sender
	if rxInput.ty != '&' {
		panic("rxInput.ty != '&'")
	}
	for _, input := range rxInput.inputs {
		if input.sender.ty != '&' {
			panic("input.sender.ty != '&'")
		}
	}
	btnPresses := make([]uint64, len(rxInput.inputs))
	var btnPress uint64
	for {
		btnPress++
		bfs(broadcaster, btnPress, rxInput.inputs, btnPresses)
		if slices.Index(btnPresses, 0) == -1 {
			// all receivers before the single inverter now received a low pulse
			break
		}
	}

	// calculate lcm of cycle lengths
	println(lcm(btnPresses[0], btnPresses[1], btnPresses[2:]...))
}
