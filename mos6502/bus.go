package mos6502

import (
	"errors"
)

// Functions to read & write a single Byte from/to a Word address
type (
	ReadByteFunc  func(Word) Byte
	WriteByteFunc func(Word, Byte)
)

// Interface for devices that can read & write Bytes
type ReadWriter interface {
	Read(Word) Byte
	Write(Word, Byte)
}

// WordReadWrite implements methods for reading & writing 16bit Words
type WordReadWrite struct {
	ReadWriter
}

// ReadWord reads a 16bit word using 2 8bit reads
func (w *WordReadWrite) ReadWord(address Word) Word {
	lo := w.Read(address)
	hi := w.Read(address + 1)
	return Word(hi)<<8 | Word(lo)
}

// WriteWord writes a 16bit word to the given address
func (w *WordReadWrite) WriteWord(address Word, data Word) {
	hi := Byte((data >> 8) & 0xFF)
	lo := Byte(data & 0xFF)

	w.Write(address, lo)
	w.Write(address+1, hi)
}

// Implement ReadWriter interface for the CPU as shims on top of the bus I/O
// functions
func (c *CPU) Read(address Word) Byte {
	return c.BusRead(address)
}

func (c *CPU) Write(address Word, data Byte) {
	c.BusWrite(address, data)
}

// ReadByte reads a single 8bit byte
func (c *CPU) ReadByte(address Word) Byte {
	return c.Read(address)
}

// FetchByte reads a single 8bit byte and increments PC
func (c *CPU) FetchByte() Byte {
	b := c.ReadByte(c.PC.Get())
	c.PC.Inc()
	return b
}

// FetchByteMode forwards the Fetch to the appropriate function for the given addressing mode
func (c *CPU) FetchByteMode(m AddrMode) (Byte, error) {
	switch m {
	case IMMEDIATE:
		return c.FetchByteImmediate(), nil
	case ABSOLUTE:
		return c.FetchByteAbsolute(), nil
	case ABSOLUTE_X:
		return c.FetchByteAbsoluteX(), nil
	case ABSOLUTE_Y:
		return c.FetchByteAbsoluteY(), nil
	case ZERO_PAGE:
		return c.FetchByteZeroPage(), nil
	case ZERO_PAGE_X:
		return c.FetchByteZeroPageX(), nil
	case ZERO_PAGE_Y:
		return c.FetchByteZeroPageY(), nil
	case INDIRECT_Y:
		return c.FetchByteIndirectY(), nil
	default:
		return Byte(0), errors.New("unknown or unsupported addressing mode")
	}
}

func (c *CPU) FetchByteImmediate() Byte {
	return c.FetchByte()
}

func (c *CPU) FetchByteAbsolute() Byte {
	addr := c.FetchWord()
	return c.ReadByte(addr)
}

func (c *CPU) FetchByteAbsoluteX() Byte {
	addr := c.FetchWord()
	return c.ReadByte(Word(addr) + Word(c.Registers.X.Get()))
}

func (c *CPU) FetchByteAbsoluteY() Byte {
	addr := c.FetchWord()
	return c.ReadByte(Word(addr) + Word(c.Registers.Y.Get()))
}

func (c *CPU) FetchByteZeroPage() Byte {
	zpa := c.FetchByte()
	return c.ReadByte(Word(zpa))
}

func (c *CPU) FetchByteZeroPageX() Byte {
	zpa := c.FetchByte()
	addr := Word(zpa + c.Registers.X.Get()&0xff)
	return c.ReadByte(addr)
}

func (c *CPU) FetchByteZeroPageY() Byte {
	zpa := c.FetchByte()
	addr := Word(zpa + c.Registers.Y.Get()&0xff)
	return c.ReadByte(addr)
}

func (c *CPU) FetchByteIndirectY() Byte {
	zpa := c.FetchByte()
	addr := c.ReadWord(Word(zpa))
	return c.ReadByte(addr + Word(c.Registers.Y.Get()))
}

// FetchWord reads a 16bit word and increments PC by 2
func (c *CPU) FetchWord() Word {
	w := c.ReadWord(c.PC.Get())
	c.PC.Inc() // low byte
	c.PC.Inc() // high byte
	return w
}

// FetchWordMode forwards the Fetch to the appropriate function for the given addressing mode
func (c *CPU) FetchWordMode(m AddrMode) (Word, error) {
	switch m {
	case ABSOLUTE:
		return c.FetchWordAbsolute(), nil
	case INDIRECT:
		return c.FetchWordIndirect(), nil
	default:
		return Word(0), errors.New("unknown or unsupported addressing mode")
	}
}

func (c *CPU) FetchWordAbsolute() Word {
	return c.FetchWord()
}

func (c *CPU) FetchWordIndirect() Word {
	addr := c.FetchWord()
	return c.ReadWord(addr)
}

// WriteByte writes a single byte to the given address
func (c *CPU) WriteByte(address Word, data Byte) {
	c.Write(address, data)
}

func (c *CPU) WriteByteAbsolute(address Word, data Byte) {
	c.WriteByte(address, data)
}

func (c *CPU) WriteByteAbsoluteX(address Word, data Byte) {
	c.WriteByte(address+Word(c.Registers.X.Get()), data)
}

func (c *CPU) WriteByteAbsoluteY(address Word, data Byte) {
	c.WriteByte(address+Word(c.Registers.Y.Get()), data)
}

func (c *CPU) WriteByteZeroPage(zpa Byte, data Byte) {
	c.WriteByte(Word(zpa), data)
}

func (c *CPU) WriteByteZeroPageX(zpa Byte, data Byte) {
	addr := Word(zpa + c.Registers.X.Get()&0xff)
	c.WriteByte(addr, data)
}

func (c *CPU) WriteByteZeroPageY(zpa Byte, data Byte) {
	addr := Word(zpa + c.Registers.Y.Get()&0xff)
	c.WriteByte(addr, data)
}

func (c *CPU) WriteByteIndirectX(base Byte, data Byte) {
	addr := c.ReadWord(Word(base) + Word(c.Registers.X.Get()))
	c.WriteByte(addr, data)
}

func (c *CPU) WriteByteIndirectY(base Byte, data Byte) {
	addr := c.ReadWord(Word(base))
	c.WriteByte(addr+Word(c.Registers.Y.Get()), data)
}

// PushByte writes an 8bit byte to the stack and increments the stack pointer by 1
func (c *CPU) PushByte(data Byte) {
	addr := STACK_BOTTOM + Word(c.Registers.S.Get())
	c.WriteByte(addr, data)
	c.Registers.S.Dec()
}

// PushWord writes a 16bit word to the stack and decrements the stack pointer by 2
func (c *CPU) PushWord(data Word) {
	addr := STACK_BOTTOM + Word(c.Registers.S.Get()-1)
	c.WriteWord(addr, data)
	c.Registers.S.Set(c.Registers.S.Get() - 2)
}

// PopByte reads an 8bit byte from the stack and increments the stack pointer by 1
func (c *CPU) PopByte() Byte {
	c.Registers.S.Inc()
	return c.ReadByte(STACK_BOTTOM + Word(c.Registers.S.Get()))
}

// PopWord reads a 16bit word from the stack and increments the stack pointer by 2
func (c *CPU) PopWord() Word {
	c.Registers.S.Set(c.Registers.S.Get() + 2)
	return c.ReadWord(STACK_BOTTOM + Word(c.Registers.S.Get()-1))
}
