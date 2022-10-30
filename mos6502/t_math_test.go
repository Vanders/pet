package mos6502

import (
	"testing"
)

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
