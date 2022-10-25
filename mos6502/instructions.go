package mos6502

import (
	"errors"
)

var (
	errUnimplemented   = errors.New("unimplemented instruction")
	errUnsupportedMode = errors.New("unknown or unsupported addressing mode")
)

// Shift Left One Bit (Memory or Accumulator)
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

// Set PC to the relative branch address
func (c *CPU) op_branch_relative(addr Byte) {
	if addr < 128 {
		c.PC.Set(c.PC.Get() + Word(addr))
	} else {
		c.PC.Set(c.PC.Get() + Word(int(addr)-256))
	}
}

// Branch on Carry Clear
func (c *CPU) op_bcc(i Instruction) error {
	addr := c.FetchByte()
	if c.Registers.P.C == true {
		c.op_branch_relative(addr)
	}

	return nil
}

func (c *CPU) op_bcs(i Instruction) error {
	return errUnimplemented
}

// Branch on Result Zero
func (c *CPU) op_beq(i Instruction) error {
	addr := c.FetchByte()
	if c.Registers.P.Z == true {
		c.op_branch_relative(addr)
	}

	return nil
}

// Branch on Result not Zero
func (c *CPU) op_bne(i Instruction) error {
	addr := c.FetchByte()
	if c.Registers.P.Z == false {
		c.op_branch_relative(addr)
	}

	return nil
}

// Branch on Result Minus
func (c *CPU) op_bmi(i Instruction) error {
	addr := c.FetchByte()
	if c.Registers.P.N == true {
		c.op_branch_relative(addr)
	}

	return nil
}

// Branch on Result Plus (Positive)
func (c *CPU) op_bpl(i Instruction) error {
	addr := c.FetchByte()
	if c.Registers.P.N == false {
		c.op_branch_relative(addr)
	}

	return nil
}

func (c *CPU) op_bvs(i Instruction) error {
	return errUnimplemented
}

func (c *CPU) op_brk(i Instruction) error {
	return errUnimplemented
}

func (c *CPU) op_php(i Instruction) error {
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
func (c *CPU) op_adc(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_pla(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_ror(i Instruction) error {
	return errUnimplemented
}

// Return from Subroutine
func (c *CPU) op_rts(i Instruction) error {
	addr := c.PopWord()
	c.PC.Set(addr + 1)

	return nil
}

func (c *CPU) op_sei(i Instruction) error {
	return errUnimplemented
}

// Store Accumulator in Memory
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

// Store Index X in Memory
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

// Store Index Y in Memory
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

func (c *CPU) op_txa(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_tya(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_txs(i Instruction) error {
	c.Registers.S.Set(c.Registers.X.Get())

	return nil
}

func (c *CPU) op_tay(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_tax(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_tsx(i Instruction) error {
	return errUnimplemented
}

// Compare Memory with Accumulator
func (c *CPU) op_cmp(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	a := c.Registers.A.Get()

	c.Registers.P.SetCarry(a >= data)
	c.Registers.P.SetZero(a == data)
	c.Registers.P.SetNegative((a-data)&BIT_7 != 0)

	return nil
}

// Clear Interrupt Disable bit
func (c *CPU) op_cli(i Instruction) error {
	c.Registers.P.SetInterrupt(false)

	return nil
}

// Clear Carry Flag
func (c *CPU) op_clc(i Instruction) error {
	c.Registers.P.SetCarry(false)

	return nil
}

// Clear Decimal Mode
func (c *CPU) op_cld(i Instruction) error {
	c.Registers.P.D = false

	return nil
}

// Compare Memory and Index X
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

func (c *CPU) op_cpy(i Instruction) error {
	return errUnimplemented
}
func (c *CPU) op_dec(i Instruction) error {
	return errUnimplemented
}

// Decrement Index X by One
func (c *CPU) op_dex(i Instruction) error {
	c.Registers.X.Dec()
	c.Registers.P.Update(c.Registers.X.Get())

	return nil
}

// Decrement Index Y by One
func (c *CPU) op_dey(i Instruction) error {
	c.Registers.Y.Dec()
	c.Registers.P.Update(c.Registers.Y.Get())

	return nil
}

func (c *CPU) op_inc(i Instruction) error {
	return errUnimplemented
}

// Increment Index X by One
func (c *CPU) op_inx(i Instruction) error {
	c.Registers.X.Inc()
	c.Registers.P.Update(c.Registers.X.Get())

	return nil
}

// Increment Index Y by One
func (c *CPU) op_iny(i Instruction) error {
	c.Registers.Y.Inc()
	c.Registers.P.Update(c.Registers.Y.Get())

	return nil
}

// Load Accumulator with Memory
func (c *CPU) op_lda(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	c.Registers.A.Set(data)
	c.Registers.P.Update(data)

	return nil
}

// Load Index X with Memory
func (c *CPU) op_ldx(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	c.Registers.X.Set(data)

	return nil
}

// Load Index Y with Memory
func (c *CPU) op_ldy(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	c.Registers.Y.Set(data)

	return nil
}

func (c *CPU) op_nop(i Instruction) error {
	return errUnimplemented
}

// OR Memory with Accumulator
func (c *CPU) op_ora(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}

	a := c.Registers.A.Get()
	a |= data
	c.Registers.A.Set(a)
	c.Registers.P.Update(a)

	return nil
}

// Subtract Memory from Accumulator with Borrow
func (c *CPU) op_sbc(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}

	var res Word

	// add A to NOT(memory), plus 1 for the carry flag (if set)
	a := c.Registers.A.Get()
	if c.Registers.P.C {
		res = Word(a) + Word(^data) + 1
	} else {
		res = Word(a) + Word(^data)
	}

	// check sign bits match
	c.Registers.P.SetOverflow((data & BIT_7) != (a & BIT_7))

	// check overflow & set carry flag if needed
	c.Registers.P.SetCarry(res > 0xff)

	// update accumulator
	c.Registers.A.Set(Byte(res & 0xff))
	c.Registers.P.Update(c.Registers.A.Get())

	return nil
}
