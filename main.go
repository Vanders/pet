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

	// Load ROMs
	edit := &ROM{
		Base: 0xe000,
		Size: 0x800, // 2k
	}
	edit.Reset()
	edit.Load("roms/edit-2-n.901447-24.bin")
	mem.Map(edit)

	kernal := &ROM{
		Base: 0xf000,
		Size: 0x1000, // 4k
	}
	kernal.Reset()
	kernal.Load("roms/kernal-2.901465-03.bin")
	mem.Map(kernal)

	// Create PIAs & VIA
	pia1 := &PIA{
		Base: 0xe810,
	}
	mem.Map(pia1)

	pia2 := &PIA{
		Base: 0xe820,
	}
	mem.Map(pia2)

	via := &VIA{
		Base: 0xe840,
	}
	mem.Map(via)

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
