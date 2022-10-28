package mos6502

type Instruction struct {
	Mode   AddrMode
	Bytes  int
	Format string
	F      func(Instruction) error
}

// makeInstructionSet returns the table of opcodes with their metadata
func (c *CPU) makeInstructionSet() map[Opcode]Instruction {
	return map[Opcode]Instruction{
		INS_BRK:     {IMPLIED, 0, "BRK ", c.op_brk},
		INS_ORA_ZP:  {ZERO_PAGE, 1, "ORA $%02x", c.op_ora},
		INS_ASL_ZP:  {ZERO_PAGE, 1, "ASL $%02x", c.op_asl},
		INS_PHP:     {IMPLIED, 0, "PHP ", c.op_php},
		INS_ORA_IM:  {IMMEDIATE, 1, "ORA #$%02x", c.op_ora},
		INS_ASL_AC:  {ACCUMULATOR, 0, "ASL ", c.op_asl},
		INS_BPL_RE:  {RELATIVE, 1, "BPL $%02x", c.op_bpl},
		INS_ORA_IY:  {INDIRECT_Y, 1, "ORA ($%02x),Y", c.op_ora},
		INS_ASL_ZPX: {ZERO_PAGE_X, 1, "ASL $%02x,X", c.op_asl},
		INS_CLC:     {IMPLIED, 0, "CLC ", c.op_clc},
		INS_JSR_AB:  {ABSOLUTE, 2, "JSR $%04x", c.op_jsr},
		INS_BIT_ZP:  {ZERO_PAGE, 1, "BIT $%02x", c.op_bit},
		INS_ROL_ZP:  {ZERO_PAGE, 1, "ROL $%02x", c.op_rol},
		INS_PLP:     {IMPLIED, 0, "PLP ", c.op_plp},
		INS_AND_IM:  {IMMEDIATE, 1, "AND #$%02x", c.op_and},
		INS_ROL_AC:  {ACCUMULATOR, 0, "ROL ", c.op_rol},
		INS_BIT_AB:  {ABSOLUTE, 2, "BIT $%04x", c.op_bit},
		INS_BMI_RE:  {RELATIVE, 1, "BMI $%02x", c.op_bmi},
		INS_SEC:     {IMPLIED, 0, "SEC ", c.op_sec},
		INS_RTI:     {IMPLIED, 0, "RTI ", c.op_rti},
		INS_EOR_ZP:  {ZERO_PAGE, 1, "EOR $%02x", c.op_eor},
		INS_LSR_ZP:  {ZERO_PAGE, 1, "LSR $%02x", c.op_lsr},
		INS_PHA:     {IMPLIED, 0, "PHA ", c.op_pha},
		INS_EOR_IM:  {IMMEDIATE, 1, "EOR #$%02x", c.op_eor},
		INS_LSR_AC:  {ACCUMULATOR, 0, "LSR ", c.op_lsr},
		INS_JMP_AB:  {ABSOLUTE, 2, "JMP $%04x", c.op_jmp},
		INS_LSR_ZPX: {ZERO_PAGE_X, 1, "LSR $%02x,X", c.op_lsr},
		INS_CLI:     {IMPLIED, 0, "CLI ", c.op_cli},
		INS_RTS:     {IMPLIED, 0, "RTS ", c.op_rts},
		INS_ADC_ZP:  {ZERO_PAGE, 1, "ADC $%02x", c.op_adc},
		INS_ROR_ZP:  {ZERO_PAGE, 1, "ROR $%02x", c.op_ror},
		INS_PLA:     {IMPLIED, 0, "PLA ", c.op_pla},
		INS_ADC_IM:  {IMMEDIATE, 1, "ADC #$%02x", c.op_adc},
		INS_ROR_AC:  {ACCUMULATOR, 0, "ROR ", c.op_ror},
		INS_JMP_IN:  {INDIRECT, 2, "JMP ($%04x)", c.op_jmp},
		INS_BVS_RE:  {RELATIVE, 1, "BVS $%02x", c.op_bvs},
		INS_ADC_IY:  {INDIRECT_Y, 1, "ADC ($%02x),Y", c.op_adc},
		INS_ROR_ZPX: {ZERO_PAGE_X, 1, "ROR $%02x,X", c.op_ror},
		INS_ADC_ABY: {ABSOLUTE_Y, 2, "ADC $%04x,Y", c.op_adc},
		INS_SEI:     {IMPLIED, 0, "SEI ", c.op_sei},
		INS_STY_ZP:  {ZERO_PAGE, 1, "STY $%02x", c.op_sty},
		INS_STA_ZP:  {ZERO_PAGE, 1, "STA $%02x", c.op_sta},
		INS_STX_ZP:  {ZERO_PAGE, 1, "STX $%02x", c.op_stx},
		INS_DEY:     {IMPLIED, 0, "DEY ", c.op_dey},
		INS_TXA:     {IMPLIED, 0, "TXA ", c.op_txa},
		INS_STY_AB:  {ABSOLUTE, 2, "STY $%04x", c.op_sty},
		INS_STA_AB:  {ABSOLUTE, 2, "STA $%04x", c.op_sta},
		INS_STX_AB:  {ABSOLUTE, 2, "STX $%04x", c.op_stx},
		INS_BCC_RE:  {RELATIVE, 1, "BCC $%02x", c.op_bcc},
		INS_STA_IY:  {INDIRECT_Y, 1, "STA ($%02x),Y", c.op_sta},
		INS_STY_ZPX: {ZERO_PAGE_X, 1, "STY $%02x,X", c.op_sty},
		INS_STA_ZPX: {ZERO_PAGE_X, 1, "STA $%02x,X", c.op_sta},
		INS_TYA:     {IMPLIED, 0, "TYA ", c.op_tya},
		INS_STA_ABY: {ABSOLUTE_Y, 2, "STA $%04x,Y", c.op_sta},
		INS_TXS:     {IMPLIED, 0, "TXS ", c.op_txs},
		INS_STA_ABX: {ABSOLUTE_X, 2, "STA $%04x,X", c.op_sta},
		INS_LDY_IM:  {IMMEDIATE, 1, "LDY #$%02x", c.op_ldy},
		INS_LDX_IM:  {IMMEDIATE, 1, "LDX #$%02x", c.op_ldx},
		INS_LDY_ZP:  {ZERO_PAGE, 1, "LDY $%02x", c.op_ldy},
		INS_LDA_ZP:  {ZERO_PAGE, 1, "LDA $%02x", c.op_lda},
		INS_LDX_ZP:  {ZERO_PAGE, 1, "LDX $%02x", c.op_ldx},
		INS_TAY:     {IMPLIED, 0, "TAY ", c.op_tay},
		INS_LDA_IM:  {IMMEDIATE, 1, "LDA #$%02x", c.op_lda},
		INS_TAX:     {IMPLIED, 0, "TAX ", c.op_tax},
		INS_LDY_AB:  {ABSOLUTE, 2, "LDY $%04x", c.op_ldy},
		INS_LDA_AB:  {ABSOLUTE, 2, "LDA $%04x", c.op_lda},
		INS_LDX_AB:  {ABSOLUTE, 2, "LDX $%04x", c.op_ldx},
		INS_BCS_RE:  {RELATIVE, 1, "BCS $%02x", c.op_bcs},
		INS_LDA_IY:  {INDIRECT_Y, 1, "LDA ($%02x),Y", c.op_lda},
		INS_LDY_ZPX: {ZERO_PAGE_X, 1, "LDY $%02x,X", c.op_ldy},
		INS_LDA_ZPX: {ZERO_PAGE_X, 1, "LDA $%02x,X", c.op_lda},
		INS_LDX_ZPY: {ZERO_PAGE_Y, 1, "LDX $%02x,Y", c.op_ldx},
		INS_LDA_ABY: {ABSOLUTE_Y, 2, "LDA $%04x,Y", c.op_lda},
		INS_TSX:     {IMPLIED, 0, "TSX ", c.op_tsx},
		INS_LDA_ABX: {ABSOLUTE_X, 2, "LDA $%04x,X", c.op_lda},
		INS_LDX_ABY: {ABSOLUTE_Y, 2, "LDX $%04x,Y", c.op_ldx},
		INS_DEC_ZP:  {ZERO_PAGE, 1, "DEC $%02x", c.op_dec},
		INS_DEC_AB:  {ABSOLUTE, 2, "DEC $%04x", c.op_dec},
		INS_CPY_IM:  {IMMEDIATE, 1, "CPY #$%02x", c.op_cpy},
		INS_CPY_ZP:  {ZERO_PAGE, 1, "CPY $%02x", c.op_cpy},
		INS_CMP_ZP:  {ZERO_PAGE, 1, "CMP $%02x", c.op_cmp},
		INS_INY:     {IMPLIED, 0, "INY ", c.op_iny},
		INS_CMP_IM:  {IMMEDIATE, 1, "CMP #$%02x", c.op_cmp},
		INS_DEX:     {IMPLIED, 0, "DEX ", c.op_dex},
		INS_CMP_AB:  {ABSOLUTE, 2, "CMP $%04x", c.op_cmp},
		INS_BNE_RE:  {RELATIVE, 1, "BNE $%02x", c.op_bne},
		INS_CMP_IY:  {INDIRECT_Y, 1, "CMP ($%02x),Y", c.op_cmp},
		INS_CLD:     {IMPLIED, 0, "CLD ", c.op_cld},
		INS_CMP_ABY: {ABSOLUTE_Y, 2, "CMP $%04x,Y", c.op_cmp},
		INS_CMP_ABX: {ABSOLUTE_X, 2, "CMP $%04x,X", c.op_cmp},
		INS_CPX_IM:  {IMMEDIATE, 1, "CPX #$%02x", c.op_cpx},
		INS_CPX_ZP:  {ZERO_PAGE, 1, "CPX $%02x", c.op_cpx},
		INS_SBC_ZP:  {ZERO_PAGE, 1, "SBC $%02x", c.op_sbc},
		INS_INC_ZP:  {ZERO_PAGE, 1, "INC $%02x", c.op_inc},
		INS_INX:     {IMPLIED, 0, "INX ", c.op_inx},
		INS_SBC_IM:  {IMMEDIATE, 1, "SBC #$%02x", c.op_sbc},
		INS_SBC_IY:  {INDIRECT_Y, 1, "SBC ($%02x),Y", c.op_sbc},
		INS_SBC_ABY: {ABSOLUTE_Y, 2, "SBC $%04x,Y", c.op_sbc},
		INS_NOP:     {IMPLIED, 0, "NOP ", c.op_nop},
		INS_CPX_AB:  {ABSOLUTE, 2, "CPX $%04x", c.op_cpx},
		INS_INC_AB:  {ABSOLUTE, 2, "INC $%04x", c.op_inc},
		INS_BEQ_RE:  {RELATIVE, 1, "BEQ $%02x", c.op_beq},
		INS_BVC_RE:  {RELATIVE, 1, "BVC $%02x", c.op_bvc},
		INS_AND_ZP:  {ZERO_PAGE, 1, "AND $%02x", c.op_and},
		INS_SBC_ABX: {ABSOLUTE_X, 1, "SBC %02x,X", c.op_sbc},
	}
}

// Addressing modes
type AddrMode int

const (
	IMPLIED     = iota // Instruction requires no address
	IMMEDIATE          // Immediate
	RELATIVE           // Relative
	ACCUMULATOR        // Accumulator
	ABSOLUTE           // Absolute
	ABSOLUTE_X         // Absolute indexed X
	ABSOLUTE_Y         // Absolute indexed Y
	ZERO_PAGE          // Zero Page
	ZERO_PAGE_X        // Zero Page indexed X
	ZERO_PAGE_Y        // Zero page indexed Y
	INDIRECT           // Indirect
	INDIRECT_X         // Indirect indexed X
	INDIRECT_Y         // Indirect indexed Y
)

// Supported opcodes
type Opcode Byte

const (
	INS_BRK    = 0x00 // break
	INS_ORA_ZP = 0x05 // inclusive OR zero page
	INS_ASL_ZP = 0x06 // arithmatic shift left zero page
	INS_PHP    = 0x08 // push status
	INS_ORA_IM = 0x09 // inclusive OR immediate
	INS_ASL_AC = 0x0a // arithmatic shift left accumulator

	INS_BPL_RE  = 0x10 // branch if positive relative
	INS_ORA_IY  = 0x11 // inclusive OR indirect y
	INS_ASL_ZPX = 0x16 // arithmatic shift left zero page indexed
	INS_CLC     = 0x18 // clear carry

	INS_JSR_AB = 0x20 // jump subroutine absolute
	INS_BIT_ZP = 0x24 // test bit zero page
	INS_ROL_ZP = 0x26 // rotate left zero page
	INS_PLP    = 0x28 // pull status
	INS_AND_IM = 0x29 // AND immediate
	INS_AND_ZP = 0x25 // AND zero page
	INS_ROL_AC = 0x2a // rotate left accumulator
	INS_BIT_AB = 0x2c // test bit absolute

	INS_BMI_RE = 0x30 // branch if minus relative
	INS_SEC    = 0x38 // set carry flag

	INS_RTI    = 0x40 // return from interrupt
	INS_EOR_ZP = 0x45 // exclusive OR zero page
	INS_LSR_ZP = 0x46 // logical shift right zero page
	INS_PHA    = 0x48 // push accumulator
	INS_EOR_IM = 0x49 // exclusive OR immediate
	INS_LSR_AC = 0x4a // logical shift right accumulator
	INS_JMP_AB = 0x4c // jump absolute

	INS_LSR_ZPX = 0x56 // logical shift right zero page indexed
	INS_CLI     = 0x58 // clear interrupt disable

	INS_RTS    = 0x60 // return from subroutine (implicit)
	INS_ADC_ZP = 0x65 // add with carry zero page
	INS_ROR_ZP = 0x66 // rotate right zero page
	INS_PLA    = 0x68 // pull accumulator
	INS_ADC_IM = 0x69 // add with carry immediate
	INS_ROR_AC = 0x6a // rotate right accumulator
	INS_JMP_IN = 0x6c // jump indirect

	INS_BVS_RE  = 0x70 // branch if overflow set relative
	INS_ADC_IY  = 0x71 // add with carry indirect y
	INS_ROR_ZPX = 0x76 // rotate right zero page indexed
	INS_ADC_ABY = 0x79 // add with carry absolute indexed y
	INS_SEI     = 0x78 // set interrupt disable

	INS_STY_ZP = 0x84 // store y zero page
	INS_STA_ZP = 0x85 // store accumulator zero page
	INS_STX_ZP = 0x86 // store x zero page
	INS_DEY    = 0x88 // decrement y
	INS_TXA    = 0x8a // transfer x to accumulator
	INS_STY_AB = 0x8c // store y absolute
	INS_STA_AB = 0x8d // store accumulator absolute
	INS_STX_AB = 0x8e // store x absolute

	INS_BCC_RE  = 0x90 // branch if carry clear relative
	INS_STA_IY  = 0x91 // store accumulator indirect y
	INS_STY_ZPX = 0x94 // store y zero page indexed
	INS_STA_ZPX = 0x95 // store accumulator zero page indexed
	INS_TYA     = 0x98 // transfer y to accumulator
	INS_STA_ABY = 0x99 // store accumulator absolute indexed y
	INS_TXS     = 0x9a // transfer x to sp
	INS_STA_ABX = 0x9d // store accumulator absolute indexed x

	INS_LDY_IM = 0xa0 // load y immediate
	INS_LDX_IM = 0xa2 // load x immediate
	INS_LDY_ZP = 0xa4 // load y zero page
	INS_LDA_ZP = 0xa5 // load accumulator zero page
	INS_LDX_ZP = 0xa6 // load x zero page
	INS_TAY    = 0xa8 // transfer accumulator to y
	INS_LDA_IM = 0xa9 // load acumulator immediate
	INS_TAX    = 0xaa // transfer accumulator to x
	INS_LDY_AB = 0xac // load y absolute
	INS_LDA_AB = 0xad // load accumulator absolute
	INS_LDX_AB = 0xae // load x absolute

	INS_BCS_RE  = 0xb0 // branch if carry set relative
	INS_LDA_IY  = 0xb1 // load acumulator indirect y
	INS_LDY_ZPX = 0xb4 // load y zero page indexed
	INS_LDA_ZPX = 0xb5 // load accumulator zero page indexed
	INS_LDX_ZPY = 0xb6 // load x zero page indexed y
	INS_LDA_ABY = 0xb9 // load accumulator absolute y
	INS_TSX     = 0xba // transfer sp to x
	INS_LDA_ABX = 0xbd // load accumulator absolute x
	INS_LDX_ABY = 0xbe // load x absolute indexed y

	INS_DEC_ZP = 0xc6 // decrement zero page
	INS_DEC_AB = 0xce // decrement absolute

	INS_CPY_IM = 0xc0 // compare y immediate
	INS_CPY_ZP = 0xc4 // compare y zero page
	INS_CMP_ZP = 0xc5 // compare zero page
	INS_INY    = 0xc8 // increment y
	INS_CMP_IM = 0xc9 // compare immediate
	INS_DEX    = 0xca // decrement x
	INS_CMP_AB = 0xcd // compare absolute
	INS_BNE_RE = 0xd0 // branch not equal relative
	INS_CMP_IY = 0xd1 // compare indirect y

	INS_CLD     = 0xd8 // clear decimal
	INS_CMP_ABY = 0xd9 // compare absolute indexed y
	INS_CMP_ABX = 0xdd // compare absolute indexed x

	INS_CPX_IM = 0xe0 // compare x immediate
	INS_CPX_ZP = 0xe4 // compare x zero page
	INS_SBC_ZP = 0xe5 // subtract with carry zero page
	INS_INC_ZP = 0xe6 // increment zero page
	INS_INX    = 0xe8 // increment x
	INS_SBC_IM = 0xe9 // subtract with carry immediate
	INS_NOP    = 0xea // no-op
	INS_CPX_AB = 0xec // compare x absolute
	INS_INC_AB = 0xee // increment absolute

	INS_SBC_IY  = 0xf1 // subtract with carry indirect y
	INS_SBC_ABY = 0xf9 // subtract with carry absolute y
	INS_SBC_ABX = 0xfd // subtract with carry zero page indexed x

	INS_BEQ_RE = 0xf0 // branch if equal relative
	INS_BVC_RE = 0x50 // branch if overflow relative
)
