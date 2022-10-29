package mos6502

import (
	"testing"
)

func Test_op_ora(t *testing.T) {
	//
	//	INS_ORA_IM
	//	INS_ORA_ZP
	//	INS_ORA_IY
	//
	testCases{
		testCase{
			INS_ORA_IM,
			"immediate (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x0f)
				m.WriteByte(0x70)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x7f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ORA_IM,
			"immediate (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.WriteByte(0x55)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ORA_IM,
			"immediate (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x00)
				m.WriteByte(0x00)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			INS_ORA_ZP,
			"zero page",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.SetByte(0x01, 0x55) // ZP $01=$55
				m.WriteByte(0x01)     // Read ZP $01
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ORA_IY,
			"indirect, y",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				c.Registers.Y.Set(0x01) // Y = 0x01

				m.WriteWord(0x01)            // ZP $01
				m.SetWord(0x0001, dataStart) // ...points to 0x0300
				m.SetByte(dataStart+1, 0x55) // ...plus Y, reads 0x0301 = 0x55
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
	}.Run(t)
}

func Test_op_eor(t *testing.T) {
	//
	//	INS_EOR_IM
	//	INS_EOR_ZP
	//
	testCases{
		testCase{
			INS_EOR_IM,
			"immediate (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.WriteByte(0xa0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x0a)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_EOR_IM,
			"immediate (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.WriteByte(0x0a)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xa0)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_EOR_IM,
			"immediate (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.WriteByte(0xaa)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			INS_EOR_ZP,
			"zero page",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.SetByte(0x01, 0xa0) // ZP $01=$a0
				m.WriteByte(0x01)     // Read ZP $01
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x0a)
				NClear(t, c)
				ZClear(t, c)
			},
		},
	}.Run(t)
}

func Test_op_and(t *testing.T) {
	//
	//	INS_AND_IM
	//	INS_AND_ZP
	//
	testCases{
		testCase{
			INS_AND_IM,
			"immediate (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.WriteByte(0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x0a)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_AND_IM,
			"immediate (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.WriteByte(0xf0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xa0)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_AND_IM,
			"immediate (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.WriteByte(0x55)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			INS_AND_ZP,
			"zero page",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.SetByte(0x01, 0x0f) // ZP $01=$0f
				m.WriteByte(0x01)     // Read ZP $01
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x0a)
				NClear(t, c)
				ZClear(t, c)
			},
		},
	}.Run(t)
}
