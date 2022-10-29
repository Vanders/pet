package mos6502

import (
	"testing"
)

func Test_op_dec(t *testing.T) {
	//
	//	INS_DEC_AB
	//	INS_DEC_ZP
	//
	testCases{
		testCase{
			INS_DEC_AB,
			"absolute (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(dataStart, 0x0f) // $0300 = $0f
				m.WriteWord(dataStart)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, dataStart, 0x0e)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_DEC_AB,
			"absolute (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(dataStart, 0xff) // $0300 = $ff
				m.WriteWord(dataStart)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, dataStart, 0xfe)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_DEC_AB,
			"absolute (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(dataStart, 0x01) // $0300 = $01
				m.WriteWord(dataStart)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, dataStart, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			INS_DEC_ZP,
			"zero page",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(0x01, 0x0f) // ZP $01=$0f
				m.WriteByte(0x01)     // DEC ZP $01
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, 0x01, 0x0e)
				NClear(t, c)
				ZClear(t, c)
			},
		},
	}.Run(t)
}

func Test_op_dex(t *testing.T) {
	//
	//	INS_DEX
	//
	testCases{
		testCase{
			INS_DEX,
			"implied (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0x0e)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_DEX,
			"implied (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0xfe)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_DEX,
			"implied (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0x01)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
	}.Run(t)
}

func Test_op_dey(t *testing.T) {
	//
	//	INS_DEY
	//
	testCases{
		testCase{
			INS_DEY,
			"implied (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.Y.Set(0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareY(t, c, 0x0e)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_DEY,
			"implied (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.Y.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareY(t, c, 0xfe)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_DEY,
			"implied (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.Y.Set(0x01)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareY(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
	}.Run(t)
}

func Test_op_inc(t *testing.T) {
	//
	//	INS_INC_AB
	//	INS_INC_ZP
	//
	testCases{
		testCase{
			INS_INC_AB,
			"absolute (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(dataStart, 0x0e) // $0300 = $0e
				m.WriteWord(dataStart)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, dataStart, 0x0f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_INC_AB,
			"absolute (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(dataStart, 0xfe) // $0300 = $fe
				m.WriteWord(dataStart)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, dataStart, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_INC_AB,
			"absolute (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(dataStart, 0xff) // $0300 = $ff
				m.WriteWord(dataStart)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, dataStart, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			INS_INC_ZP,
			"zero page",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(0x01, 0x0e) // ZP $01=$0e
				m.WriteByte(0x01)     // DEC ZP $01
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, 0x01, 0x0f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
	}.Run(t)
}

func Test_op_inx(t *testing.T) {
	//
	//	INS_INX
	//
	testCases{
		testCase{
			INS_INX,
			"implied (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0x0e)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0x0f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_INX,
			"implied (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0xfe)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_INX,
			"implied (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
	}.Run(t)
}
