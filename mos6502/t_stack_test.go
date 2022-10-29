package mos6502

import (
	"testing"
)

func Test_stack(t *testing.T) {
	//
	//	INS_PHA
	//	INS_PHP
	//	INS_PLA
	//	INS_PLP
	//
	testCases{
		testCase{
			INS_PHA,
			"PHA",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				// Check the byte at the top of the stack is 0x00
				CompareMem(t, m, STACK_TOP-1, 0xaa)
				// Check the stack pointer has decremented by one
				CompareSP(t, c, 0xfe)
			},
		},
		testCase{
			INS_PHP,
			"PHP",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.I = true
				c.Registers.P.C = false
				c.Registers.P.Z = true
				c.Registers.P.V = false
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				expected := c.Registers.P.GetByte()
				// Check the byte at the top of the stack matches SP
				CompareMem(t, m, STACK_TOP-1, expected)
				// Check the stack pointer has decremented by one
				CompareSP(t, c, 0xfe)
			},
		},
		testCase{
			INS_PLA,
			"PLA (positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.S.Set(0xfe)
				m.SetByte(STACK_TOP-1, 0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x0f)

				NClear(t, c)
				ZClear(t, c)

				// Check the stack pointer has incremented by one
				CompareSP(t, c, 0xff)
			},
		},
		testCase{
			INS_PLA,
			"PLA (negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.S.Set(0xfe)
				m.SetByte(STACK_TOP-1, 0xf0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xf0)

				NSet(t, c)
				ZClear(t, c)

				// Check the stack pointer has incremented by one
				CompareSP(t, c, 0xff)
			},
		},
		testCase{
			INS_PLA,
			"PLA (zero)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.S.Set(0xfe)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)

				NClear(t, c)
				ZSet(t, c)

				// Check the stack pointer has incremented by one
				CompareSP(t, c, 0xff)
			},
		},
		testCase{
			INS_PLP,
			"PLP",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.S.Set(0xfe)
				m.SetByte(STACK_TOP-1, 0xff) // All flags set
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				NSet(t, c)
				ZSet(t, c)
				CSet(t, c)
				if c.Registers.P.I != true {
					t.Error("interrupt flag is not set")
				}
				if c.Registers.P.D != true {
					t.Error("decimal flag is not set")
				}
				if c.Registers.P.V != true {
					t.Error("overflow flag is not set")
				}

				// Ensure bits 6 & 5 are clear
				P := c.Registers.P.GetByte()
				if P != 0xcf {
					t.Errorf("bits 5 & 6 are not clear: got 0x%02x", P)
				}

				// Check the stack pointer has incremented by one
				CompareSP(t, c, 0xff)
			},
		},
	}.Run(t)
}
