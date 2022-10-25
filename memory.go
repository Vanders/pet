package main

import (
	"fmt"
	"io"

	"github.com/vanders/pet/mos6502"
)

type ReadWriter interface {
	Read(mos6502.Word) mos6502.Byte
	Write(mos6502.Word, mos6502.Byte)
}

type Device interface {
	GetBase() mos6502.Word
	GetSize() mos6502.Word
	ReadWriter
}

type Memory struct {
	Devices []Device

	Writer io.Writer // io.Writer for log output

	mem [mos6502.MAX_ADDR]mos6502.Byte
}

func (m *Memory) Reset() {
	for n := mos6502.STACK_TOP; n < mos6502.MAX_ADDR; n++ {
		m.mem[n] = 0x00
	}
}

func (m *Memory) debug(format string, a ...any) (int, error) {
	if m.Writer != nil {
		return fmt.Fprintf(m.Writer, format, a...)
	} else {
		return 0, nil
	}
}

func (m *Memory) Read(address mos6502.Word) mos6502.Byte {
	m.debug("read $%04x\n", address)

	for n, d := range m.Devices {
		base := d.GetBase()
		size := d.GetSize()
		top := base + (size - 1)

		m.debug("device %d at $%04x:$%04x\n", n, base, base+top)

		if address >= base && address <= top {
			m.debug("selected device %d at $%04x\n", n, base)
			return d.Read(address)
		}
	}
	return m.mem[address]
}

func (m *Memory) Write(address mos6502.Word, data mos6502.Byte) {
	m.debug("write $%04x\n", address)

	for n, d := range m.Devices {
		base := d.GetBase()
		size := d.GetSize()
		top := base + (size - 1)

		m.debug("device %d at $%04x:$%04x\n", n, base, base+top)
		if address >= base && address <= top {
			m.debug("selected device %d at $%04x\n", n, base)
			d.Write(address, data)
		}
	}

	m.mem[address] = data
}

func (m *Memory) Map(device Device) {
	m.Devices = append(m.Devices, device)
}
