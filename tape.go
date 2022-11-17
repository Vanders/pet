package main

import (
	"os"
)

const DATA_START = 2

type Prg struct {
	data []byte

	addr Word
	size Word
}

func (p *Prg) Addr() Word {
	return p.addr
}

func (p *Prg) Size() Word {
	return p.size
}

func (p *Prg) Load(data []byte) {
	// Calculate addr from first two bytes
	lo := data[0]
	hi := data[1]
	p.addr = Word(hi)<<8 | Word(lo)

	// Size is the remaining data
	p.size = Word(len(data) - DATA_START)

	// Copy remaining data, minus header
	p.data = make([]byte, p.size)
	copy(p.data, data[DATA_START:])
}

func (p *Prg) Read(addr Word) Byte {
	if addr < p.size {
		return Byte(p.data[addr])
	}
	return Byte(0x00)
}

func (p *Prg) SetAddr(address Word) {
	p.addr = address
}

func (p *Prg) SetSize(size Word) {
	p.size = size
}

func (p *Prg) Save(data []Byte) []byte {
	prg := make([]byte, p.size+DATA_START)

	hi := byte((p.addr >> 8) & 0xff)
	lo := byte(p.addr & 0xff)
	prg[0] = lo
	prg[1] = hi

	for n := Word(0); n < p.size; n++ {
		prg[DATA_START+n] = byte(data[n])
	}
	return prg
}

type Cassette struct {
	filename string
	prg      *Prg
	cb       Word
}

func (c *Cassette) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	prg := &Prg{}
	prg.Load(data)

	c.filename = filename
	c.prg = prg
	c.cb = 0

	return nil
}

func (c *Cassette) Addr() Word {
	return c.prg.Addr()
}

func (c *Cassette) Size() Word {
	return c.prg.Size()
}

func (c *Cassette) FetchByte() Byte {
	b := c.prg.Read(c.cb)
	c.cb++
	return b
}

func (c *Cassette) Save(filename string, address Word, size Word, data []Byte) error {
	prg := Prg{}
	prg.SetAddr(address)
	prg.SetSize(size)

	out := prg.Save(data)
	err := os.WriteFile(filename, out, 0644)
	if err != nil {
		return err
	}
	return nil
}
