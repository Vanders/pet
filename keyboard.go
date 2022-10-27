package main

// Keypress contains the data from a key press
type Keypress struct {
	Keycode   int
	Modifiers struct {
		Shift bool
	}
}

type Key struct {
	row uint8 // keyboard row
	bit uint8 // keyboard bit (column)
}

type Keyboard struct {
	Buffer chan<- (Key) // Keyboard "buffer"

	keys map[rune]Key
}

func (kbd *Keyboard) Reset() {
	kbd.keys = make(map[rune]Key)
	kbd.keys['='] = Key{9, 7}
	kbd.keys['.'] = Key{9, 6}
	//kbd.keys[''] = Key{9, 5}
	kbd.keys['\033'] = Key{9, 4} // ^C/STOP (Escape)
	kbd.keys['<'] = Key{9, 3}
	kbd.keys[' '] = Key{9, 2}
	kbd.keys['['] = Key{9, 1}
	//kbd.keys[''] = Key{9, 0} // ^R/REVERSE ON

	kbd.keys['-'] = Key{8, 7}
	kbd.keys['0'] = Key{8, 6}
	//kbd.keys[''] = Key{8,5} // RIGHT SHIFT
	kbd.keys['>'] = Key{8, 4}
	//kbd.keys[''] = Key{8,3}
	kbd.keys[']'] = Key{8, 2}
	kbd.keys['@'] = Key{8, 1}
	//kbd.keys[''] = Key{8,0} // LEFT SHIFT

	kbd.keys['+'] = Key{7, 7}
	kbd.keys['2'] = Key{7, 6}
	//kbd.keys[''] = Key{7,5}
	kbd.keys['?'] = Key{7, 4}
	kbd.keys[','] = Key{7, 3}
	kbd.keys['n'] = Key{7, 2}
	kbd.keys['v'] = Key{7, 1}
	kbd.keys['x'] = Key{7, 0}

	kbd.keys['3'] = Key{6, 7}
	kbd.keys['1'] = Key{6, 6}
	/*
		if EnableTrace == false {
			kbd.keys['\r'] = Key{6, 5} //RETURN
		} else {
			kbd.keys['`'] = Key{6, 5} //RETURN
		}
	*/
	kbd.keys['\r'] = Key{6, 5} //RETURN

	kbd.keys[';'] = Key{6, 4}
	kbd.keys['m'] = Key{6, 3}
	kbd.keys['b'] = Key{6, 2}
	kbd.keys['c'] = Key{6, 1}
	kbd.keys['z'] = Key{6, 0}

	kbd.keys['*'] = Key{5, 7}
	kbd.keys['5'] = Key{5, 6}
	//kbd.keys[''] = Key{5,5}
	kbd.keys[':'] = Key{5, 4}
	kbd.keys['k'] = Key{5, 3}
	kbd.keys['h'] = Key{5, 2}
	kbd.keys['f'] = Key{5, 1}
	kbd.keys['s'] = Key{5, 0}

	kbd.keys['6'] = Key{4, 7}
	kbd.keys['4'] = Key{4, 6}
	//kbd.keys[''] = Key{4,5}
	kbd.keys['l'] = Key{4, 4}
	kbd.keys['j'] = Key{4, 3}
	kbd.keys['g'] = Key{4, 2}
	kbd.keys['d'] = Key{4, 1}
	kbd.keys['a'] = Key{4, 0}
	kbd.keys['A'] = Key{4, 0}

	kbd.keys['/'] = Key{3, 7}
	kbd.keys['8'] = Key{3, 6}
	//kbd.keys[''] = Key{3,5}
	kbd.keys['p'] = Key{3, 4}
	kbd.keys['i'] = Key{3, 3}
	kbd.keys['y'] = Key{3, 2}
	kbd.keys['r'] = Key{3, 1}
	kbd.keys['w'] = Key{3, 0}

	kbd.keys['9'] = Key{2, 7}
	kbd.keys['7'] = Key{2, 6}
	kbd.keys['^'] = Key{2, 5}
	kbd.keys['o'] = Key{2, 4}
	kbd.keys['u'] = Key{2, 3}
	kbd.keys['t'] = Key{2, 2}
	kbd.keys['e'] = Key{2, 1}
	kbd.keys['q'] = Key{2, 0}

	//kbd.keys['\b'] = Key{1, 7} // DEL
	kbd.keys['\102'] = Key{1, 6} // cursor down
	//kbd.keys[''] = Key{1,5}
	kbd.keys[')'] = Key{1, 4}
	kbd.keys['\\'] = Key{1, 3}
	kbd.keys['\''] = Key{1, 2}
	kbd.keys['$'] = Key{1, 1}
	kbd.keys['"'] = Key{1, 0}

	kbd.keys['\103'] = Key{0, 7} // cursor right
	kbd.keys['\176'] = Key{0, 6} // home
	kbd.keys['\177'] = Key{0, 5} // backspace
	kbd.keys['('] = Key{0, 4}
	kbd.keys['&'] = Key{0, 3}
	kbd.keys['%'] = Key{0, 2}
	kbd.keys['#'] = Key{0, 1}
	kbd.keys['!'] = Key{0, 0}
}

func (kbd *Keyboard) Scan(k Keypress) {
	key, ok := kbd.keys[rune(k.Keycode)]
	if ok {
		kbd.Buffer <- key
	}
}
