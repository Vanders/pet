package main

const (
	scr_w      = 40
	scr_h      = 25
	borderTop  = 30
	borderLeft = 40
	VID_MEM    = 0x8000
)

type Video struct {
	Read    func(address Word) Byte // Read a single byte from the bus
	VIA_CB2 func() Byte             // Returns the current status of the VIA CB2 line
	PIA_CB1 func(bool)              // Notify PIA of retrace via. the CB1 line

	ROM *ROM // Character generator ROM
}

func (v *Video) Redraw(drawPixel func(x, y int)) {
	// Start retrace interrupt
	v.PIA_CB1(true)

	// Draw 40 x 25 characters
	var scr_y, scr_x int

	var line [scr_w]Byte
	for y := int32(borderTop); scr_y < scr_h; y += 10 {
		y_addr := Word(VID_MEM + (scr_y * scr_w))
		// read 40 characters for the row
		for n := range line {
			line[n] = v.Read(y_addr + Word(n))
		}

		// For each row, draw 8 scanlines of 25 characters
		for l := int32(0); l < 8; l++ {
			scr_x = 0
			for x := int32(borderLeft); scr_x < scr_w; x += 10 {
				char := line[scr_x]

				// If high bit of vmem is set, invert the video
				invert := char & 0x80

				romAddr := Word(char&0x7f)<<3 | Word(l&0x07)
				// If VIA CB2 is set, set the high bit (A10) of the video ROM address
				if v.VIA_CB2() != 0 {
					romAddr |= 0x400
				}
				bits := v.ROM.Read(romAddr)

				for p := int32(0); p < 8; p++ {
					if (bits<<p)&0x80 != invert {
						drawPixel(int(x+p+1), int(y+l+1))
					}
				}
				scr_x++
			}
		}
		scr_y++
	}

	// End retrace interrupt
	v.PIA_CB1(false)
}
