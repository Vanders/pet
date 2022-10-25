package mos6502

import (
	"errors"
)

var unimpl = errors.New("unimplemented instruction")

func (c *CPU) op_brk(i Instruction) error {
	return unimpl
}
func (c *CPU) op_ora(i Instruction) error {
	return unimpl
}
func (c *CPU) op_asl(i Instruction) error {
	return unimpl
}
func (c *CPU) op_php(i Instruction) error {
	return unimpl
}
func (c *CPU) op_bpl(i Instruction) error {
	return unimpl
}
func (c *CPU) op_clc(i Instruction) error {
	return unimpl
}
func (c *CPU) op_jsr(i Instruction) error {
	addr := c.FetchWord()
	c.PushWord(c.PC.Get() - 1)
	c.PC.Set(addr)

	return nil
}
func (c *CPU) op_bit(i Instruction) error {
	return unimpl
}
func (c *CPU) op_rol(i Instruction) error {
	return unimpl
}
func (c *CPU) op_plp(i Instruction) error {
	return unimpl
}
func (c *CPU) op_and(i Instruction) error {
	return unimpl
}
func (c *CPU) op_bmi(i Instruction) error {
	return unimpl
}
func (c *CPU) op_sec(i Instruction) error {
	return unimpl
}
func (c *CPU) op_rti(i Instruction) error {
	return unimpl
}
func (c *CPU) op_eor(i Instruction) error {
	return unimpl
}
func (c *CPU) op_lsr(i Instruction) error {
	return unimpl
}
func (c *CPU) op_pha(i Instruction) error {
	return unimpl
}
func (c *CPU) op_jmp(i Instruction) error {
	return unimpl
}
func (c *CPU) op_cli(i Instruction) error {
	return unimpl
}
func (c *CPU) op_rts(i Instruction) error {
	return unimpl
}
func (c *CPU) op_adc(i Instruction) error {
	return unimpl
}
func (c *CPU) op_pla(i Instruction) error {
	return unimpl
}
func (c *CPU) op_bvs(i Instruction) error {
	return unimpl
}
func (c *CPU) op_sei(i Instruction) error {
	return unimpl
}
func (c *CPU) op_sty(i Instruction) error {
	return unimpl
}
func (c *CPU) op_sta(i Instruction) error {
	return unimpl
}
func (c *CPU) op_stx(i Instruction) error {
	return unimpl
}
func (c *CPU) op_dey(i Instruction) error {
	return unimpl
}
func (c *CPU) op_ror(i Instruction) error {
	return unimpl
}
func (c *CPU) op_txa(i Instruction) error {
	return unimpl
}
func (c *CPU) op_bcc(i Instruction) error {
	return unimpl
}
func (c *CPU) op_tya(i Instruction) error {
	return unimpl
}
func (c *CPU) op_txs(i Instruction) error {
	c.Registers.S.Set(c.Registers.X.Get())

	return nil
}
func (c *CPU) op_ldy(i Instruction) error {
	return unimpl
}
func (c *CPU) op_ldx(i Instruction) error {
	b, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	c.Registers.X.Set(b)

	return nil
}
func (c *CPU) op_lda(i Instruction) error {
	b, err := c.FetchByteMode(i.Mode)
	if err != nil {
		return err
	}
	c.Registers.A.Set(b)
	// XXX update P

	return nil
}
func (c *CPU) op_tay(i Instruction) error {
	return unimpl
}
func (c *CPU) op_tax(i Instruction) error {
	return unimpl
}
func (c *CPU) op_bcs(i Instruction) error {
	return unimpl
}
func (c *CPU) op_tsx(i Instruction) error {
	return unimpl
}
func (c *CPU) op_dec(i Instruction) error {
	return unimpl
}
func (c *CPU) op_cpy(i Instruction) error {
	return unimpl
}
func (c *CPU) op_iny(i Instruction) error {
	return unimpl
}
func (c *CPU) op_cmp(i Instruction) error {
	return unimpl
}
func (c *CPU) op_dex(i Instruction) error {
	return unimpl
}
func (c *CPU) op_bne(i Instruction) error {
	return unimpl
}
func (c *CPU) op_cld(i Instruction) error {
	c.Registers.P.D = false

	return nil
}
func (c *CPU) op_cpx(i Instruction) error {
	return unimpl
}
func (c *CPU) op_sbc(i Instruction) error {
	return unimpl
}
func (c *CPU) op_inc(i Instruction) error {
	return unimpl
}
func (c *CPU) op_inx(i Instruction) error {
	return unimpl
}
func (c *CPU) op_nop(i Instruction) error {
	return unimpl
}
func (c *CPU) op_beq(i Instruction) error {
	return unimpl
}
