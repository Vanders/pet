package mos6502

import (
	"errors"
)

var (
	errUnimplemented   = errors.New("unimplemented instruction")
	errUnsupportedMode = errors.New("unknown or unsupported addressing mode")
)

func (c *CPU) op_brk(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_ora(i Instruction) error {
	return errUnimplemented
}

func (c *CPU) op_asl(i Instruction) error {
	carryAndShift := func(data Byte) Byte {
		c.Registers.P.SetCarry(data&BIT_7 != 0)
		return (data << 1) & 0xfe
	}

	switch i.Mode {
	case ACCUMULATOR:
		data := c.Registers.A.Get()

		data = carryAndShift(data)

		c.Registers.A.Set(data)
		c.Registers.P.Update(data)
	case ZERO_PAGE:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data = carryAndShift(data)

		c.WriteByteZeroPage(zpa, data)
		c.Registers.P.Update(data)
	case ZERO_PAGE_X:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data = carryAndShift(data)

		c.WriteByteZeroPageX(zpa, data)
		c.Registers.P.Update(data)
	case ABSOLUTE:
		return errUnsupportedMode
	case ABSOLUTE_X:
		return errUnsupportedMode
	default:
		return errUnsupportedMode
	}

	return nil
}

func (c *CPU) op_php(i Instruction) error {
	return errUnimplemented
}

func (c *CPU) op_bpl(i Instruction) error {
	switch i.Mode {
	case RELATIVE:
		addr := c.FetchByte()
		if c.Registers.P.N == false {
			if addr < 128 {
				c.PC.Set(c.PC.Get() + Word(addr))
			} else {
				c.PC.Set(c.PC.Get() + Word(int(addr)-256))
			}
		}
	default:
		return errUnsupportedMode
	}

	return nil
}

func (c *CPU) op_clc(i Instruction) error {
	return errUnimplemented
}

func (c *CPU) op_jsr(i Instruction) error {
	addr := c.FetchWord()
	c.PushWord(c.PC.Get() - 1)
	c.PC.Set(addr)

	return nil
}

func (c *CPU) op_bit(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}

	c.Registers.P.SetZero(c.Registers.A.Get()&data == 0)
	c.Registers.P.SetOverflow(data&BIT_6 == 0)
	c.Registers.P.SetNegative(data&BIT_7 == 0)

	return nil
}

func (c *CPU) op_rol(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_plp(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_and(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_bmi(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_sec(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_rti(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_eor(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_lsr(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_pha(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_jmp(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_cli(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_rts(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_adc(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_pla(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_bvs(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_sei(i Instruction) error {
	return errUnimplemented
}

func (c *CPU) op_sta(i Instruction) error {
	data := c.Registers.A.Get()

	switch i.Mode {
	case ABSOLUTE:
		addr := c.FetchWord()
		c.WriteByteAbsolute(addr, data)
	case ABSOLUTE_X:
		addr := c.FetchWord()
		c.WriteByteAbsoluteX(addr, data)
	case ABSOLUTE_Y:
		addr := c.FetchWord()
		c.WriteByteAbsoluteY(addr, data)
	case ZERO_PAGE:
		zpa := c.FetchByte()
		c.WriteByteZeroPage(zpa, data)
	case ZERO_PAGE_X:
		zpa := c.FetchByte()
		c.WriteByteZeroPageX(zpa, data)
	case INDIRECT_X:
		base := c.FetchByte()
		c.WriteByteIndirectX(base, data)
	case INDIRECT_Y:
		base := c.FetchByte()
		c.WriteByteIndirectY(base, data)
	default:
		return errUnsupportedMode
	}
	return nil
}

func (c *CPU) op_stx(i Instruction) error {
	data := c.Registers.X.Get()
	switch i.Mode {
	case ABSOLUTE:
		addr := c.FetchWord()
		c.WriteByte(addr, data)
	case ZERO_PAGE:
		zpa := c.FetchByte()
		c.WriteByteZeroPage(zpa, data)
	case ZERO_PAGE_X:
		zpa := c.FetchByte()
		c.WriteByteZeroPageX(zpa, data)
	default:
		return errUnsupportedMode
	}
	return nil
}

func (c *CPU) op_sty(i Instruction) error {
	data := c.Registers.Y.Get()
	switch i.Mode {
	case ABSOLUTE:
		addr := c.FetchWord()
		c.WriteByte(addr, data)
	case ZERO_PAGE:
		zpa := c.FetchByte()
		c.WriteByteZeroPage(zpa, data)
	case ZERO_PAGE_X:
		zpa := c.FetchByte()
		c.WriteByteZeroPageX(zpa, data)
	default:
		return errUnsupportedMode
	}
	return nil
}

func (c *CPU) op_dey(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_ror(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_txa(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_bcc(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_tya(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_txs(i Instruction) error {
	c.Registers.S.Set(c.Registers.X.Get())

	return nil
}

func (c *CPU) op_ldy(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	c.Registers.Y.Set(data)

	return nil
}

func (c *CPU) op_ldx(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	c.Registers.X.Set(data)

	return nil
}

func (c *CPU) op_lda(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	c.Registers.A.Set(data)
	c.Registers.P.Update(data)

	return nil
}

func (c *CPU) op_tay(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_tax(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_bcs(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_tsx(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_dec(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_cpy(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_iny(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_cmp(i Instruction) error {
	return errUnimplemented
}

func (c *CPU) op_dex(i Instruction) error {
	c.Registers.X.Dec()
	c.Registers.P.Update(c.Registers.X.Get())

	return nil
}

func (c *CPU) op_bne(i Instruction) error {
	return errUnimplemented
}

func (c *CPU) op_cld(i Instruction) error {
	c.Registers.P.D = false

	return nil
}

func (c *CPU) op_cpx(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	x := c.Registers.X.Get()

	c.Registers.P.SetCarry(x >= data)
	c.Registers.P.SetZero(x == data)
	c.Registers.P.SetNegative((x-data)&BIT_7 != 0)

	return nil
}

func (c *CPU) op_sbc(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_inc(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_inx(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_nop(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_beq(i Instruction) error {
	return errUnimplemented
}
