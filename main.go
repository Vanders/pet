package main

import (
	"os"

	"github.com/vanders/pet/mos6502"
)

func main() {
	var mem Memory
	mem.Reset()

	// load some real data
	mem.Load(0xf000, "roms/kernal-2.901465-03.bin")

	cpu := mos6502.CPU{
		Read:   mem.Read,
		Write:  mem.Write,
		Writer: os.Stderr,
	}
	cpu.Reset()

	// XXX Attempt to execute the 1st 100 instructions
	for n := 0; n < 100; n++ {
		cpu.Step()
	}
}
