package mos6502

import (
	"testing"
)

func Test_op_asl(t *testing.T) {
	//
	//	INS_ASL_AC
	//	INS_ASL_ZP
	//	INS_ASL_ZPX
	//
	testCases{
		testCase{
			INS_ASL_AC,
			"accumulator (positive, no carry)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x3f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x7e)

				CClear(t, c)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ASL_AC,
			"accumulator (negative, no carry)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x7f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xfe)

				CClear(t, c)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ASL_AC,
			"accumulator (negative, carry)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xfe)

				CSet(t, c)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ASL_AC,
			"accumulator (zero, carry)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x80)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)

				CSet(t, c)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			INS_ASL_ZP,
			"zero page",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(0x01, 0x3f) // ZP $01=$3f
				m.WriteByte(0x01)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, 0x01, 0x7e)

				CClear(t, c)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ASL_ZPX,
			"zero page, indexed x",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(0x02, 0x3f) // ZP $02=$3f

				c.Registers.X.Set(0x01)
				m.WriteByte(0x01) // $01 + X($01) = ZP $02
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, 0x02, 0x7e)

				CClear(t, c)
				NClear(t, c)
				ZClear(t, c)
			},
		},
	}.Run(t)
}

func Test_op_rol(t *testing.T) {
	//
	//	INS_ROL_AC
	//	INS_ROL_ZP
	//
	testCases{
		testCase{
			INS_ROL_AC,
			"accumulator (positive, no carry in, no carry out)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x3f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x7e)

				CClear(t, c)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ROL_AC,
			"accumulator (negative, no carry in, no carry out)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x7f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xfe)

				CClear(t, c)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ROL_AC,
			"accumulator (positive, carry in, no carry out)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x3f)
				c.Registers.P.C = true
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x7f)

				CClear(t, c)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ROL_AC,
			"accumulator (negative, no carry in, carry out)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xfe)

				CSet(t, c)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_ROL_AC,
			"accumulator (zero, no carry out)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x00)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)

				CClear(t, c)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			INS_ROL_ZP,
			"zero page",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				m.SetByte(0x01, 0x3f) // ZP $01=$3f
				m.WriteByte(0x01)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareMem(t, m, 0x01, 0x7e)

				CClear(t, c)
				NClear(t, c)
				ZClear(t, c)
			},
		},
	}.Run(t)
}
