package mos6502

import (
	"testing"
)

func Test_flags(t *testing.T) {
	//
	//	INS_CLC
	//	INS_CLD
	//	INS_CLI
	//	INS_SEC
	//	INS_SEI
	//
	testCases{
		testCase{
			INS_CLC,
			"CLC",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CClear(t, c)
			},
		},
		testCase{
			INS_CLD,
			"CLD",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.D = true
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				if c.Registers.P.D != false {
					t.Error("D flag not cleared")
				}
			},
		},
		testCase{
			INS_CLI,
			"CLI",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.I = true
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				if c.Registers.P.I != false {
					t.Error("I flag not cleared")
				}
			},
		},
		testCase{
			INS_SEC,
			"SEC",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CSet(t, c)
			},
		},
		testCase{
			INS_SEI,
			"SEI",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.I = false
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				if c.Registers.P.I != true {
					t.Error("I flag not set")
				}
			},
		},
	}.Run(t)
}
