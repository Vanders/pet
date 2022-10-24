package mos6502

import (
	"fmt"
	"io"
)

// XXX Move to memory?
type ReadWriter interface {
	Read() Word
	Write(Word)
}

type CPU struct {
	Registers struct {
		A WordRegister // Accumulater
		X WordRegister // X index
		Y WordRegister // Y index
		S ByteRegister // Stack pointer
		P Flags        // Processer state
	}
	PC WordRegister // Program Counter
	IR ByteRegister // Current instruction

	insCount int // Number of instructions executed

	Read  func(address Word) Byte        // Read a single byte from the bus
	Write func(address Word, value Byte) // Write a single byte to the bus

	Writer io.Writer // io.Writer for log output
}

const (
	VEC_RESET     = 0xfffc   // Reset vector
	VEC_INTERRUPT = 0xfffe   // Interrupt vector
	MAX_ADDR      = 64 * KiB // Maximum addressable memory
	STACK_BOTTOM  = 0x0100   // Address of the bottom of the stack
	STACK_TOP     = 0x0200   // Address of the top of the stack
)

// Log writes a formatted sting to the configured output
func (c *CPU) Log(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(c.Writer, format, a...)
}

// ReadByte reads a single 8bit byte
func (c *CPU) ReadByte(address Word) Byte {
	return c.Read(address)
}

// ReadWord reads a 16bit word using 2 8bit reads
func (c *CPU) ReadWord(address Word) Word {
	lo := c.Read(address)
	hi := c.Read(address + 1)
	return Word(hi)<<8 | Word(lo)
}

// Reset resets the CPU as though RST had been asserted
func (c *CPU) Reset() {
	// Reset the registers
	c.Registers.A.Set(0xaa)
	c.Registers.X.Set(0x00)
	c.Registers.Y.Set(0x00)

	// Set the stack pointer & reset PC from the reset vector
	c.Registers.S.Set(0xff)
	c.PC.Set(c.ReadWord(VEC_RESET))
	c.IR.Set(0x00)

	// Clear flags
	c.Registers.P.Reset()

	// Reset the instruction count
	c.insCount = 0
}

// Step fetches & executes a single instruction
func (c *CPU) Step() {
	// Read next instruction from & increment PC
	c.insCount++

	ins := c.ReadByte(c.PC.Get())
	c.Log("%d:\t0x%.2x:\t0x%.2x\n", c.insCount, c.PC.Get(), ins)

	c.IR.Set(ins)
	c.PC.Inc()

	// XXX Decode

	// XXX Disasemble & log
}
