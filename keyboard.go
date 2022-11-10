package main

// Keypress contains the data from a key press
type Keypress struct {
	Scancode Byte
	State    KeyState
	Shifted  bool
}

type KeyState int

const (
	KEY_DOWN = 0
	KEY_UP   = 1
)

// Key defines the key scan data for an individual key
type Key struct {
	row uint8 // keyboard row
	bit uint8 // keyboard bit (column)
}

// Matrix defines the map of the keyboard scan matrix
type Matrix struct {
	Rows [10]Byte
}

// Reset all the rows active high (off)
func (m *Matrix) Reset() {
	for n := 0; n < 10; n++ {
		m.Rows[n] = 0xff
	}
}

// Set a single bit in a row active low (on)
func (m *Matrix) Set(row, bit uint8) {
	m.Rows[row] &= ^(0x01 << bit)
}

// Set a single bit in a row active high (off)
func (m *Matrix) Clear(row, bit uint8) {
	m.Rows[row] |= 0x01 << bit
}

// Get the current status of a row
func (m *Matrix) Get(row uint8) Byte {
	if row < 10 {
		return m.Rows[row]
	}
	return 0xff
}

type Keyboard struct {
	Buffer chan<- (Key) // Keyboard "buffer"

	matrix Matrix // Keyboard scan matrix
	keys   map[Byte]Key
}

func (kbd *Keyboard) Reset() {
	kbd.matrix.Reset()

	kbd.keys = make(map[Byte]Key)
	kbd.keys[0x3d] = Key{9, 7}
	kbd.keys[0x2e] = Key{9, 6}
	kbd.keys[0xff] = Key{9, 5} // UNUSED
	kbd.keys[0x03] = Key{9, 4} // ^C/STOP (Escape)
	kbd.keys[0x3c] = Key{9, 3}
	kbd.keys[0x20] = Key{9, 2}
	kbd.keys[0x5b] = Key{9, 1}
	kbd.keys[0x12] = Key{9, 0} // ^R/REVERSE ON

	kbd.keys[0x2d] = Key{8, 7}
	kbd.keys[0x30] = Key{8, 6}
	kbd.keys[0x00] = Key{8, 5} // RIGHT SHIFT
	kbd.keys[0x3e] = Key{8, 4}
	kbd.keys[0xff] = Key{8, 3} // UNUSED
	kbd.keys[0x5d] = Key{8, 2}
	kbd.keys[0x40] = Key{8, 1}
	kbd.keys[0x00] = Key{8, 0} // LEFT SHIFT

	kbd.keys[0x2b] = Key{7, 7}
	kbd.keys[0x32] = Key{7, 6}
	kbd.keys[0xff] = Key{7, 5} // UNUSED
	kbd.keys[0x3f] = Key{7, 4}
	kbd.keys[0x2c] = Key{7, 3}
	kbd.keys[0x4e] = Key{7, 2}
	kbd.keys[0x56] = Key{7, 1}
	kbd.keys[0x58] = Key{7, 0}

	kbd.keys[0x33] = Key{6, 7}
	kbd.keys[0x31] = Key{6, 6}
	kbd.keys[0x0d] = Key{6, 5} // RETURN
	kbd.keys[0x3b] = Key{6, 4}
	kbd.keys[0x4d] = Key{6, 3}
	kbd.keys[0x42] = Key{6, 2}
	kbd.keys[0x43] = Key{6, 1}
	kbd.keys[0x5a] = Key{6, 0}

	kbd.keys[0x2a] = Key{5, 7}
	kbd.keys[0x35] = Key{5, 6}
	kbd.keys[0xff] = Key{5, 5} // UNUSED
	kbd.keys[0x3a] = Key{5, 4}
	kbd.keys[0x4b] = Key{5, 3}
	kbd.keys[0x48] = Key{5, 2}
	kbd.keys[0x46] = Key{5, 1}
	kbd.keys[0x53] = Key{5, 0}

	kbd.keys[0x36] = Key{4, 7}
	kbd.keys[0x34] = Key{4, 6}
	kbd.keys[0xff] = Key{4, 5} // UNUSED
	kbd.keys[0x4c] = Key{4, 4}
	kbd.keys[0x4a] = Key{4, 3}
	kbd.keys[0x47] = Key{4, 2}
	kbd.keys[0x44] = Key{4, 1}
	kbd.keys[0x41] = Key{4, 0}

	kbd.keys[0x2f] = Key{3, 7}
	kbd.keys[0x38] = Key{3, 6}
	kbd.keys[0xff] = Key{3, 5} // UNUSED
	kbd.keys[0x50] = Key{3, 4}
	kbd.keys[0x49] = Key{3, 3}
	kbd.keys[0x59] = Key{3, 2}
	kbd.keys[0x52] = Key{3, 1}
	kbd.keys[0x57] = Key{3, 0}

	kbd.keys[0x39] = Key{2, 7}
	kbd.keys[0x37] = Key{2, 6}
	kbd.keys[0x5e] = Key{2, 5}
	kbd.keys[0x4f] = Key{2, 4}
	kbd.keys[0x55] = Key{2, 3}
	kbd.keys[0x54] = Key{2, 2}
	kbd.keys[0x45] = Key{2, 1}
	kbd.keys[0x51] = Key{2, 0}

	kbd.keys[0x14] = Key{1, 7} // DEL
	kbd.keys[0x11] = Key{1, 6} // cursor down
	kbd.keys[0xff] = Key{1, 5} // UNUSED
	kbd.keys[0x29] = Key{1, 4}
	kbd.keys[0x5c] = Key{1, 3}
	kbd.keys[0x27] = Key{1, 2}
	kbd.keys[0x24] = Key{1, 1}
	kbd.keys[0x22] = Key{1, 0}

	kbd.keys[0x1d] = Key{0, 7} // cursor right
	kbd.keys[0x13] = Key{0, 6} // home
	kbd.keys[0x5f] = Key{0, 5} // backspace; unused, mapped to DEL
	kbd.keys[0x28] = Key{0, 4}
	kbd.keys[0x26] = Key{0, 3}
	kbd.keys[0x25] = Key{0, 2}
	kbd.keys[0x23] = Key{0, 1}
	kbd.keys[0x21] = Key{0, 0}
}

func (kbd *Keyboard) Scan(k Keypress) {
	key, ok := kbd.keys[k.Scancode]
	if ok {
		kbd.Buffer <- key

		switch k.State {
		case KEY_DOWN:
			kbd.matrix.Set(key.row, key.bit)
			if k.Shifted {
				kbd.matrix.Set(8, 5) // Right shift
			}
		case KEY_UP:
			kbd.matrix.Clear(key.row, key.bit)
			if k.Shifted {
				kbd.matrix.Clear(8, 5) // Right shift
			}
		}
	}
}

func (kbd *Keyboard) Get(row Byte) Byte {
	return kbd.matrix.Get(uint8(row))
}
