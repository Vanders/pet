package mos6502

import (
	"fmt"
)

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
	case INDIRECT_Y:
		return c.FetchByteIndirectY(), nil
	default:
		return Byte(0), fmt.Errorf("Unknown or unsupported addressing mode %d")
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

func (c *CPU) FetchByteIndirectY() Byte {
	zpa := c.FetchByte()
	addr := c.ReadWord(Word(zpa))
	return c.ReadByte(addr + Word(c.Registers.Y.Get()))
}

func (c *CPU) WriteByte(address Word, data Byte) {
	c.Write(address, data)
}

// ReadWord reads a 16bit word using 2 8bit reads
func (c *CPU) ReadWord(address Word) Word {
	lo := c.Read(address)
	hi := c.Read(address + 1)
	return Word(hi)<<8 | Word(lo)
}

// FetchWord reads a 16bit word and increments PC by 2
func (c *CPU) FetchWord() Word {
	w := c.ReadWord(c.PC.Get())
	c.PC.Inc() // low byte
	c.PC.Inc() // high byte
	return w
}

func (c *CPU) WriteWord(address Word, data Word) {
	hi := Byte((data >> 8) & 0xFF)
	lo := Byte(data & 0xFF)

	c.WriteByte(address, lo)
	c.WriteByte(address+1, hi)
}

func (c *CPU) PushWord(data Word) {
	addr := STACK_BOTTOM + Word(c.Registers.S.Get()-1)
	c.WriteWord(addr, data)
	c.Registers.S.Set(c.Registers.S.Get() - 2)
}
