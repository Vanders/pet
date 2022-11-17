package main

import "github.com/vanders/pet/mos6502"

var (
	LOADPATCH_v2 = []Byte{
		mos6502.INS_PHP,    // PHP: Store current flags
		mos6502.INS_PHA,    // PHA: Store current A register contents
		mos6502.INS_LDA_IM, // LDA: Load A register...
		TRAP_LOAD,          // #$01: with trap selector for LOAD
		mos6502.INS_TRAP,   // TRAP: Emulator trap
		mos6502.INS_PLA,    // PLA: Restore contents of A register
		mos6502.INS_PLP,    // PLP: Restore flags
		mos6502.INS_RTS,    // RTS: Return from vector subroutine
	}

	LOADPATCH_v4 = []Byte{
		mos6502.INS_PHP, // PHP: Store current flags
		mos6502.INS_PHA, // PHA: Store current A register contents
		mos6502.INS_TYA, // TYA: Transfer Y to A
		mos6502.INS_PHA, // PHY: Store current A (Y) register contents
		mos6502.INS_LDY_IM,
		0x41,               // press play
		mos6502.INS_JSR_AB, // print message
		0x85,
		0xf1,
		mos6502.INS_LDY_IM,
		0x56,               // on tape #
		mos6502.INS_JSR_AB, // print message
		0x85,
		0xf1,
		mos6502.INS_LDA_ZP,
		0xd4, // current tape #
		mos6502.INS_ORA_IM,
		0x30,
		mos6502.INS_JSR_AB, // print character
		0x02,
		0xe2,
		mos6502.INS_LDA_IM, // LDA: Load A register...
		TRAP_LOAD,          // #$01: with trap selector for LOAD
		mos6502.INS_TRAP,   // TRAP: Emulator trap
		mos6502.INS_PLA,    // PLA: Restore contents of A (Y) register
		mos6502.INS_TAY,    // TAY: Transfer A to Y
		mos6502.INS_PLA,    // PLA: Restore contents of A register
		mos6502.INS_PLP,    // PLP: Restore flags
		mos6502.INS_RTS,    // RTS: Return from vector subroutine
	}

	SAVEPATCH_v2 = []Byte{
		mos6502.INS_PHP,    // PHP: Store current flags
		mos6502.INS_PHA,    // PHA: Store current A register contents
		mos6502.INS_LDA_IM, // LDA: Load A register...
		TRAP_SAVE,          // #$02: with trap selector for SAVE
		mos6502.INS_TRAP,   // TRAP: Emulator trap
		mos6502.INS_PLA,    // PLA: Restore contents of A register
		mos6502.INS_PLP,    // PLP: Restore flags
		mos6502.INS_RTS,    // RTS: Return from vector subroutine
	}

	SAVEPATCH_v4 = []Byte{
		mos6502.INS_PHP, // PHP: Store current flags
		mos6502.INS_PHA, // PHA: Store current A register contents
		mos6502.INS_TYA, // TYA: Transfer Y to A
		mos6502.INS_PHA, // PHY: Store current A (Y) register contents
		mos6502.INS_LDY_IM,
		0x41,               // press play
		mos6502.INS_JSR_AB, // print message
		0x85,
		0xf1,
		mos6502.INS_LDY_IM,
		0x4d,               // & record
		mos6502.INS_JSR_AB, // print message
		0x85,
		0xf1,
		mos6502.INS_LDY_IM,
		0x56,               // on tape #
		mos6502.INS_JSR_AB, // print message
		0x85,
		0xf1,
		mos6502.INS_LDA_ZP,
		0xd4, // current tape #
		mos6502.INS_ORA_IM,
		0x30,
		mos6502.INS_JSR_AB, // print character
		0x02,
		0xe2,
		mos6502.INS_LDY_IM,
		0x64,               // writing
		mos6502.INS_JSR_AB, // print message
		0x85,
		0xf1,
		mos6502.INS_LDA_IM, // LDA: Load A register...
		TRAP_SAVE,          // #$02: with trap selector for SAVE
		mos6502.INS_TRAP,   // TRAP: Emulator trap
		mos6502.INS_PLA,    // PLA: Restore contents of A (Y) register
		mos6502.INS_TAY,    // TAY: Transfer A to Y
		mos6502.INS_PLA,    // PLA: Restore contents of A register
		mos6502.INS_PLP,    // PLP: Restore flags
		mos6502.INS_RTS,    // RTS: Return from vector subroutine
	}
)
