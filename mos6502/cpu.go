package mos6502

import (
	"fmt"
	"io"
)

// TrapFunc is a callback for entering an emulator trap from a CPU instruction
type TrapFunc func(Byte)

type CPU struct {
	Registers struct {
		A ByteRegister // Accumulater
		X ByteRegister // X index
		Y ByteRegister // Y index
		S ByteRegister // Stack pointer
		P Flags        // Processer state
	}
	PC WordRegister // Program Counter
	IR ByteRegister // Current instruction

	instructionSet map[Opcode]Instruction // Table of opcodes
	insCount       int                    // Number of instructions executed
	isr            bool                   // Is the CPU running the ISR?

	BusRead  ReadByteFunc  // Read a single byte from the bus
	BusWrite WriteByteFunc // Write a single byte to the bus
	Trap     TrapFunc      // Handle emulator trap

	WordReadWrite // Read & Write 16bit words

	Writer io.Writer // io.Writer for log output
}

const (
	VEC_RESET     = 0xfffc   // Reset vector
	VEC_INTERRUPT = 0xfffe   // Interrupt vector
	MAX_ADDR      = 64 * KiB // Maximum addressable memory
	STACK_BOTTOM  = 0x0100   // Address of the bottom of the stack
	STACK_TOP     = 0x0200   // Address of the top of the stack
)

// Create & initialise a new CPU object
func NewCPU(rf ReadByteFunc, wf WriteByteFunc, tf TrapFunc, w io.Writer) *CPU {
	cpu := &CPU{
		BusRead:  rf,
		BusWrite: wf,
		Trap:     tf,
		Writer:   w,
	}
	cpu.ReadWriter = cpu

	return cpu
}

// Log writes a formatted sting to the configured output
func (c *CPU) Log(format string, a ...any) (int, error) {
	if c.Writer != nil {
		return fmt.Fprintf(c.Writer, format, a...)
	}
	return 0, nil
}

func (c *CPU) Dump() {
	c.Log("\n\tPC: $%04x\n", c.PC.Get())
	c.Log("\n\tA: %s\n\tX: %s\n\tY: %s\n\tS: %s\n",
		c.Registers.A,
		c.Registers.X,
		c.Registers.Y,
		c.Registers.S)
	c.Log("\nFlags:\n%s\n", c.Registers.P)
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

	// Initialise opcode table
	c.instructionSet = c.makeInstructionSet()

	// Reset the instruction count
	c.insCount = 0
}

// Step fetches & executes a single instruction
func (c *CPU) Step() error {
	// Fetch next instruction from PC
	opcode := c.FetchByte()
	c.Log("%d:\t0x%.2x:\t0x%.2x:\t(S: %s)\t(A: %s, X: %s: Y: %s)\t",
		c.insCount,
		c.PC.Get()-1,
		opcode,
		c.Registers.S,
		c.Registers.A,
		c.Registers.X,
		c.Registers.Y)
	c.IR.Set(opcode)

	// Decode
	ins, ok := c.instructionSet[Opcode(opcode)]
	if !ok {
		return fmt.Errorf("invalid or unknown instruction 0x%2x", opcode)
	}
	c.insCount++

	// Disasemble & log
	switch ins.Bytes {
	case 0:
		c.Log(ins.Format + "\r\n")
	case 1:
		c.Log(ins.Format+"\r\n", c.ReadByte(c.PC.Get()))
	case 2:
		c.Log(ins.Format+"\r\n", c.ReadWord(c.PC.Get()))
	}

	// Call the instruction implementation
	err := ins.F(ins)
	if err != nil {
		return err
	}

	return nil
}

// Raise interrupt
func (c *CPU) Interrupt() {
	if c.Registers.P.I == false && c.isr == false {
		c.isr = true

		c.PushWord(c.PC.Get() - 1)
		c.PushByte(c.Registers.P.GetByte())

		addr := c.ReadWord(VEC_INTERRUPT)
		c.PC.Set(addr)
	}
}
