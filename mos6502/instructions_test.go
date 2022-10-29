package mos6502

import (
	"testing"
)

func Test_op_ora(t *testing.T) {
	/*
		INS_ORA_IM
		INS_ORA_ZP
		INS_ORA_IY
	*/

	testCases := []testCase{
		/* INS_ORA_IM */
		testCase{
			"immediate (positive)",
			Instruction{
				Mode: IMMEDIATE,
			},
			/* Setup */
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x0f)
				m.WriteByte(0x70)
			},
			/* Check */
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x7f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			"immediate (negative)",
			Instruction{
				Mode: IMMEDIATE,
			},
			/* Setup */
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.WriteByte(0x55)
			},
			/* Check */
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			"immediate (zero)",
			Instruction{
				Mode: IMMEDIATE,
			},
			/* Setup */
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x00)
				m.WriteByte(0x00)
			},
			/* Check */
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			/* INS_ORA_ZP */
			"zero page",
			Instruction{
				Mode: ZERO_PAGE,
			},
			/* Setup */
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				m.SetByte(0x01, 0x55) // ZP $01=$55
				m.WriteByte(0x01)     // Read ZP $01
			},
			/* Check */
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			/* INS_ORA_IY */
			"indirect, y",
			Instruction{
				Mode: INDIRECT_Y,
			},
			/* Setup */
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0xaa)
				c.Registers.Y.Set(0x01) // Y = 0x01

				m.WriteWord(0x01)         // ZP $01
				m.SetWord(0x0001, 0x0300) // ...points to 0x0300
				m.SetByte(0x0301, 0x55)   // ...plus Y, reads 0x0301 = 0x55
			},
			/* Check */
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xff)
				NSet(t, c)
				ZClear(t, c)
			},
		},
	}

	m := newMem()
	c := newCPU(m)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t, c, m)

			err := c.op_ora(tc.ins)
			if err != nil {
				t.Error(err)
			}

			tc.check(t, c, m)
		})

		m.Reset()
		c.Reset()
	}
}

func Test_op_eor(t *testing.T) {
	/*
		INS_EOR_IM
		INS_EOR_ZP
	*/

	testCases := []testCase{
		/* INS_ORA_IM */
		testCase{
			"immediate (positive)",
			Instruction{
				Mode: IMMEDIATE,
			},
			/* Setup */
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x0f)
				m.WriteByte(0x70)
			},
			/* Check */
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x7f)
				NClear(t, c)
				ZClear(t, c)
			},
		},
	}

	m := newMem()
	c := newCPU(m)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t, c, m)

			err := c.op_eor(tc.ins)
			if err != nil {
				t.Error(err)
			}

			tc.check(t, c, m)
		})

		m.Reset()
		c.Reset()
	}
}
