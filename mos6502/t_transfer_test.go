package mos6502

import (
	"testing"
)

func Test_op_tax(t *testing.T) {
	//
	//	INS_TAX
	//
	testCases{
		testCase{
			INS_TAX,
			"positive",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0x0f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TAX,
			"negative",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TAX,
			"zero",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x00)
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

func Test_op_tay(t *testing.T) {
	//
	//	INS_TAY
	//
	testCases{
		testCase{
			INS_TAY,
			"positive",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareY(t, c, 0x0f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TAY,
			"negative",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareY(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TAY,
			"zero",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x00)
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

func Test_op_tsx(t *testing.T) {
	//
	//	INS_TSX
	//
	testCases{
		testCase{
			INS_TSX,
			"positive",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.S.Set(0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0x0f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TSX,
			"negative",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.S.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareX(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TSX,
			"zero",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.S.Set(0x00)
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

func Test_op_txa(t *testing.T) {
	//
	//	INS_TXA
	//
	testCases{
		testCase{
			INS_TXA,
			"positive",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x0f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TXA,
			"negative",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TXA,
			"zero",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0x00)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
	}.Run(t)
}

func Test_op_tya(t *testing.T) {
	//
	//	INS_TYA
	//
	testCases{
		testCase{
			INS_TYA,
			"positive",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.Y.Set(0x0f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x0f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TYA,
			"negative",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.Y.Set(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_TYA,
			"zero",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.Y.Set(0x00)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
	}.Run(t)
}

func Test_op_txs(t *testing.T) {
	//
	//	INS_TXS
	//
	testCases{
		testCase{
			INS_TXS,
			"implied",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.X.Set(0x0f)

				c.Registers.P.N = true
				c.Registers.P.Z = true
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareSP(t, c, 0x0f)

				// Check that the status flags have not been affected
				NSet(t, c)
				ZSet(t, c)
			},
		},
	}.Run(t)
}
