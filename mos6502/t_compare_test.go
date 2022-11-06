package mos6502

import (
	"testing"
)

func Test_op_cmp(t *testing.T) {
	//
	//	INS_CMP_IM
	//	INS_CMP_AB
	//	INS_CMP_ZP
	//	INS_CMP_IY
	//	INS_CMP_ABY
	//	INS_CMP_ABX
	//
	testCases{
		testCase{
			INS_CMP_IM,
			"immediate (less than, positive)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x01)
				m.WriteByte(0xff)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x01)

				CClear(t, c)
				NClear(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_CMP_IM,
			"immediate (less than, negative)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x01)
				m.WriteByte(0x7f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x01)

				CClear(t, c)
				NSet(t, c)
				ZClear(t, c)
			},
		},
		testCase{
			INS_CMP_IM,
			"immediate (equals)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x7f)
				m.WriteByte(0x7f)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x7f)

				CSet(t, c)
				NClear(t, c)
				ZSet(t, c)
			},
		},
		testCase{
			INS_CMP_IM,
			"immediate (more than)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.A.Set(0x7f)
				m.WriteByte(0x01)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x7f)

				CSet(t, c)
				NClear(t, c)
				ZClear(t, c)
			},
		},
	}.Run(t)
}
