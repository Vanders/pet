package mos6502

import (
	"errors"
)

var (
	errUnimplemented   = errors.New("unimplemented instruction")
	errUnsupportedMode = errors.New("unknown or unsupported addressing mode")
)

// Add Memory to Accumulator with Carry
func (c *CPU) op_adc(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}

	var res Word

	a := c.Registers.A.Get()
	sameSign := (data & BIT_7) == (a & BIT_7)

	// add A to memory, plus 1 for the carry flag (if set)
	res = Word(a) + Word(data)
	if c.Registers.P.C == true {
		res += 1
	}

	// if the inputs were the same sign check for overflow
	if sameSign {
		// does the result sign match the input sign? if not we overflowed
		c.Registers.P.SetOverflow((data & BIT_7) != (a & BIT_7))
	} else {
		c.Registers.P.SetOverflow(false)
	}

	// update accumulator
	a = Byte(res & 0xff)
	c.Registers.A.Set(a)

	// check overflow & set carry flag if needed
	c.Registers.P.SetCarry(res > 0xff)

	c.Registers.P.Update(a)

	return nil
}

// AND Memory with Accumulator
func (c *CPU) op_and(i Instruction) error {
	data := c.FetchByteImmediate()
	a := c.Registers.A.Get() & data
	c.Registers.A.Set(a)
	c.Registers.P.Update(a)

	return nil
}

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
	if c.Registers.P.C == false {
		c.op_branch_relative(addr)
	}

	return nil
}

// Branch on Carry Set
func (c *CPU) op_bcs(i Instruction) error {
	addr := c.FetchByte()
	if c.Registers.P.C == true {
		c.op_branch_relative(addr)
	}

	return nil
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

func (c *CPU) op_bit(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}

	c.Registers.P.SetZero(c.Registers.A.Get()&data == 0)
	c.Registers.P.SetOverflow(data&BIT_6 != 0)
	c.Registers.P.SetNegative(data&BIT_7 != 0)

	return nil
}

// Push Accumulator
func (c *CPU) op_pha(i Instruction) error {
	c.PushByte(c.Registers.A.Get())
	return nil
}

// Push Processor Status on Stack
func (c *CPU) op_php(i Instruction) error {
	c.PushByte(c.Registers.P.GetByte())
	return nil
}

// Pull Accumulator from Stack
func (c *CPU) op_pla(i Instruction) error {
	data := c.PopByte()
	c.Registers.A.Set(data)
	c.Registers.P.Update(data)

	return nil
}

// Pull Processor Status from Stack
func (c *CPU) op_plp(i Instruction) error {
	data := c.PopByte()
	c.Registers.P.SetByte(data)

	return nil
}

// Rotate One Bit Left (Memory or Accumulator)
func (c *CPU) op_rol(i Instruction) error {
	rol := func(data Byte) Byte {
		// old carry becomes bit 0, carry is set to bit 7
		carry := c.Registers.P.C
		c.Registers.P.SetCarry(data&BIT_7 != 0)

		// shift & handle carry bit
		data = data << 1

		if carry {
			// old carry was 1, set bit 0
			data |= BIT_0
		} else {
			// old carry was 0, clear bit 0
			data &= 0xfe
		}
		return data
	}

	switch i.Mode {
	case ACCUMULATOR:
		data := c.Registers.A.Get()

		data = rol(data)

		c.Registers.A.Set(data)
		c.Registers.P.Update(data)
	case ZERO_PAGE:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data = rol(data)

		c.WriteByteZeroPage(zpa, data)
		c.Registers.P.Update(data)
	case ZERO_PAGE_X:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data = rol(data)

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

// Rotate One Bit Right (Memory or Accumulator)
func (c *CPU) op_ror(i Instruction) error {
	ror := func(data Byte) Byte {
		// old carry becomes bit 7, carry is set to bit 0
		carry := c.Registers.P.C
		c.Registers.P.SetCarry(data&BIT_0 != 0)

		// shift & handle carry bit
		data = data >> 1

		if carry {
			// old carry was 1, set bit 7
			data |= BIT_7
		} else {
			// old carry was 0, clear bit 7
			data &= 0x7f
		}
		return data
	}

	switch i.Mode {
	case ACCUMULATOR:
		data := c.Registers.A.Get()

		data = ror(data)

		c.Registers.A.Set(data)
		c.Registers.P.Update(data)
	case ZERO_PAGE:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data = ror(data)

		c.WriteByteZeroPage(zpa, data)
		c.Registers.P.Update(data)
	case ZERO_PAGE_X:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data = ror(data)

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

// Return from Interrupt
func (c *CPU) op_rti(i Instruction) error {
	p := c.PopByte()
	c.Registers.P.SetByte(p)

	// Force clear BCD mode flag
	c.Registers.P.B = false

	addr := c.PopWord()
	c.PC.Set(addr + 1)

	// CPU is no longer servicing an interrupt
	c.isr = false

	return nil
}

// Return from Subroutine
func (c *CPU) op_rts(i Instruction) error {
	addr := c.PopWord()
	c.PC.Set(addr + 1)

	return nil
}

// Set Carry Flag
func (c *CPU) op_sec(i Instruction) error {
	c.Registers.P.SetCarry(true)

	return nil
}

func (c *CPU) op_sei(i Instruction) error {
	c.Registers.P.SetNegative(true)

	return nil
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

// Transfer Accumulator to Index X
func (c *CPU) op_tax(i Instruction) error {
	a := c.Registers.A.Get()
	c.Registers.X.Set(a)
	c.Registers.P.Update(a)

	return nil
}

// Transfer Accumulator to Index Y
func (c *CPU) op_tay(i Instruction) error {
	a := c.Registers.A.Get()
	c.Registers.Y.Set(a)
	c.Registers.P.Update(a)

	return nil
}

// Transfer Stack Pointer to Index X
func (c *CPU) op_tsx(i Instruction) error {
	s := c.Registers.S.Get()
	c.Registers.X.Set(s)
	c.Registers.P.Update(s)

	return nil
}

// Transfer Index X to Stack Pointer
func (c *CPU) op_txs(i Instruction) error {
	c.Registers.S.Set(c.Registers.X.Get())

	return nil
}

// Transfer Index X to Accumulator
func (c *CPU) op_txa(i Instruction) error {
	x := c.Registers.X.Get()
	c.Registers.A.Set(x)
	c.Registers.P.Update(x)

	return nil
}

// Transfer Index Y to Accumulator
func (c *CPU) op_tya(i Instruction) error {
	y := c.Registers.Y.Get()
	c.Registers.A.Set(y)
	c.Registers.P.Update(y)

	return nil
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

// Compare Memory and Index Y
func (c *CPU) op_cpy(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	y := c.Registers.Y.Get()

	c.Registers.P.SetCarry(y >= data)
	c.Registers.P.SetZero(y == data)
	c.Registers.P.SetNegative((y-data)&BIT_7 != 0)

	return nil
}

// Decrement Memory by One
func (c *CPU) op_dec(i Instruction) error {

	switch i.Mode {
	case ZERO_PAGE:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data -= 1

		c.WriteByteZeroPage(zpa, data)
		c.Registers.P.Update(data)
	case ZERO_PAGE_X:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data -= 1

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

// Exclusive-OR Memory with Accumulator
func (c *CPU) op_eor(i Instruction) error {
	data, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	a := c.Registers.A.Get()
	a ^= data

	c.Registers.A.Set(a)
	c.Registers.P.Update(a)

	return nil
}

// Increment Memory by One
func (c *CPU) op_inc(i Instruction) error {
	switch i.Mode {
	case ZERO_PAGE:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data += 1

		c.WriteByteZeroPage(zpa, data)
		c.Registers.P.Update(data)
	case ZERO_PAGE_X:
		zpa := c.FetchByte()
		data := c.ReadByte(Word(zpa))

		data += 1

		c.WriteByteZeroPageX(zpa, data)
		c.Registers.P.Update(data)
	case ABSOLUTE:
		addr := c.FetchWord()
		data := c.ReadByte(addr)

		data += 1

		c.WriteByteAbsolute(addr, data)
		c.Registers.P.Update(data)
	case ABSOLUTE_X:
		addr := c.FetchWord()
		data := c.ReadByte(Word(addr) + Word(c.Registers.X.Get()))

		data += 1

		c.WriteByteAbsoluteX(addr, data)
		c.Registers.P.Update(data)
	default:
		return errUnsupportedMode
	}

	return nil
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

// Jump to New Location
func (c *CPU) op_jmp(i Instruction) error {
	addr, err := c.FetchWordMode(i.Mode)
	if err != nil {
		return err
	}
	c.PC.Set(addr)

	return nil
}

// Jump to New Location Saving Return Address
func (c *CPU) op_jsr(i Instruction) error {
	addr := c.FetchWord()
	c.PushWord(c.PC.Get() - 1)
	c.PC.Set(addr)

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

// Shift One Bit Right (Memory or Accumulator)
func (c *CPU) op_lsr(i Instruction) error {
	carryAndShift := func(data Byte) Byte {
		c.Registers.P.SetCarry(data&BIT_0 != 0)
		return (data >> 1) & 0x7f
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
