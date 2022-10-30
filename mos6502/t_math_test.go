package mos6502

import (
	"testing"
)

/*
The tests in this file were built with the helpful tables at
http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html
*/

func Test_op_adc(t *testing.T) {
	//
	//	INS_ADC_IM
	//	INS_ADC_ZP
	//	INS_ADC_ABY
	//	INS_ADC_IY
	//
	testCases{
		testCase{
			INS_ADC_IM,
			"immediate #1 (positive + positive, positive result, no carry in, no carry out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0x50)

				m.WriteByte(0x10)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x60)

				CClear(t, c)
				NClear(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_ADC_IM,
			"immediate (#2, positive + positive, negative result, no carry in, no carry out, overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0x50)

				m.WriteByte(0x50)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xa0)

				CClear(t, c)
				NSet(t, c)
				ZClear(t, c)
				VSet(t, c)
			},
		},
		testCase{
			INS_ADC_IM,
			"immediate (#3, positive + negative, negative result, no carry in, no carry out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0x50)

				m.WriteByte(0x90)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xe0)

				CClear(t, c)
				NSet(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_ADC_IM,
			"immediate (#4, positive + negative, positive result, no carry in, carry out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0x50)

				m.WriteByte(0xd0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x20)

				CSet(t, c)
				NClear(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_ADC_IM,
			"immediate (#5, negative + positive, negative result, no carry in, no carry out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0xd0)

				m.WriteByte(0x10)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xe0)

				CClear(t, c)
				NSet(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_ADC_IM,
			"immediate (#6, negative + positive, positive result, no carry in, carry out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0xd0)

				m.WriteByte(0x50)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x20)

				CSet(t, c)
				NClear(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_ADC_IM,
			"immediate (#7, negative + negative, positive result, no carry in, carry out, overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0xd0)

				m.WriteByte(0x90)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x60)

				CSet(t, c)
				NClear(t, c)
				ZClear(t, c)
				VSet(t, c)
			},
		},
		testCase{
			INS_ADC_IM,
			"immediate (#8, negative + negative, negative result, no carry in, carry out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0xd0)

				m.WriteByte(0xd0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xa0)

				CSet(t, c)
				NSet(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_ADC_IM,
			"immediate (zero, no carry in, no carry out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = false
				c.Registers.A.Set(0x00)

				m.WriteByte(0x00)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)

				CClear(t, c)
				NClear(t, c)
				ZSet(t, c)
				VClear(t, c)
			},
		},
	}.Run(t)
}

func Test_op_sbc(t *testing.T) {
	//
	//	INS_SBC_IM
	//	INS_SBC_ZP
	//	INS_SBC_ABX
	//	INS_SBC_ABY
	//	INS_SBC_IY
	//
	testCases{
		testCase{
			INS_SBC_IM,
			"immediate #1 (positive - negative, positive result, no borrow in, borrow out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0x50)

				m.WriteByte(0xf0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x60)

				CClear(t, c) // C=0 (^C=B=1)
				NClear(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_SBC_IM,
			"immediate #2 (positive - negative, negative result, no borrow in, borrow out, overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0x50)

				m.WriteByte(0xb0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xa0)

				CClear(t, c) // C=0 (^C=B=1)
				NSet(t, c)
				ZClear(t, c)
				VSet(t, c)
			},
		},
		testCase{
			INS_SBC_IM,
			"immediate #3 (positive - positive, negative result, no borrow in, borrow out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0x50)

				m.WriteByte(0x70)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xe0)

				CClear(t, c) // C=0 (^C=B=1)
				NSet(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_SBC_IM,
			"immediate #4 (positive - positive, positive result, no borrow in, no borrow out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0x50)

				m.WriteByte(0x30)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x20)

				CSet(t, c) // C=1 (^C=B=0)
				NClear(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_SBC_IM,
			"immediate #5 (negative - negative, negative result, no borrow in, borrow out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0xd0)

				m.WriteByte(0xf0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xe0)

				CClear(t, c) // C=0 (^C=B=1)
				NSet(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_SBC_IM,
			"immediate #6 (negative - negative, positive result, no borrow in, no borrow out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0xd0)

				m.WriteByte(0xb0)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x20)

				CSet(t, c) // C=1 (^C=B=0)
				NClear(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_SBC_IM,
			"immediate #7 (negative - positive, positive result, no borrow in, no borrow out, overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0xd0)

				m.WriteByte(0x70)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x60)

				CSet(t, c) // C=1 (^C=B=0)
				NClear(t, c)
				ZClear(t, c)
				VSet(t, c)
			},
		},
		testCase{
			INS_SBC_IM,
			"immediate #8 (negative - positive, negative result, no borrow in, no borrow out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0xd0)

				m.WriteByte(0x30)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0xa0)

				CSet(t, c) // C=1 (^C=B=0)
				NSet(t, c)
				ZClear(t, c)
				VClear(t, c)
			},
		},
		testCase{
			INS_SBC_IM,
			"immediate (zero, no borrow in, no borrow out, no overflow)",
			// Setup
			func(t *testing.T, c *CPU, m *fakeMem) {
				c.Registers.P.C = true // C=1 (^C=B=0)
				c.Registers.A.Set(0x01)

				m.WriteByte(0x01)
			},
			// Check
			func(t *testing.T, c *CPU, m *fakeMem) {
				CompareA(t, c, 0x00)

				CSet(t, c) // C=1 (^C=B=0)
				NClear(t, c)
				ZSet(t, c)
				VClear(t, c)
			},
		},
	}.Run(t)
}
