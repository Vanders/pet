package main

import (
	"fmt"
	"os"

	"github.com/vanders/pet/mos6502"
)

func dumpAndExit(cpu *mos6502.CPU, err error) {
	fmt.Println(err)
	cpu.Dump()
	os.Exit(1)
}

func main() {
	var mem Memory
	mem.Reset()

	// load some real data
	mem.Load(0xe000, "roms/edit-2-n.901447-24.bin")
	mem.Load(0xf000, "roms/kernal-2.901465-03.bin")

	cpu := mos6502.CPU{
		Read:   mem.Read,
		Write:  mem.Write,
		Writer: os.Stderr,
	}
	cpu.Reset()

	// Execute as many instructions as possible
	for {
		err := cpu.Step()
		if err != nil {
			dumpAndExit(&cpu, fmt.Errorf("execution stopped: %s", err))
		}
	}
}
