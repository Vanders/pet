package main

import (
	"bytes"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowHeight = 480
	windowWidth  = 360
	borderTop    = 30
	borderLeft   = 40
	fontSize     = 8
	VID_MEM      = 0x8000
)

var (
	petscii = [256]rune{
		0x00: '@',
		0x01: 'A',
		0x02: 'B',
		0x03: 'C',
		0x04: 'D',
		0x05: 'E',
		0x06: 'F',
		0x07: 'G',
		0x08: 'H',
		0x09: 'I',
		0x0a: 'J',
		0x0b: 'K',
		0x0c: 'L',
		0x0d: 'M',
		0x0e: 'N',
		0x0f: 'O',
		0x10: 'P',
		0x11: 'Q',
		0x12: 'R',
		0x13: 'S',
		0x14: 'T',
		0x15: 'U',
		0x16: 'V',
		0x17: 'W',
		0x18: 'X',
		0x19: 'Y',
		0x1a: 'Z',
		0x1b: '[',
		0x1c: '\\',
		0x1d: ']',
		0x1e: ' ',
		0x1f: ' ',
		0x20: ' ',
		0x21: '!',
		0x22: '"',
		0x23: '#',
		0x24: '$',
		0x25: '%',
		0x26: '&',
		0x27: '\'',
		0x28: '(',
		0x29: ')',
		0x2a: '*',
		0x2b: '+',
		0x2c: ',',
		0x2d: '-',
		0x2e: '.',
		0x2f: '/',
		0x30: '0',
		0x31: '1',
		0x32: '2',
		0x33: '3',
		0x34: '4',
		0x35: '5',
		0x36: '6',
		0x37: '7',
		0x38: '8',
		0x39: '9',
		0x3a: ':',
		0x3b: ';',
		0x3c: '<',
		0x3d: '=',
		0x3e: '>',
		0x3f: '?',
		/*
		   0x40: '',
		   0x41: '',
		   0x42: '',
		   0x43: '',
		   0x44: '',
		   0x45: '',
		   0x46: '',
		   0x47: '',
		   0x48: '',
		   0x49: '',
		   0x4a: '',
		   0x4b: '',
		   0x4c: '',
		   0x4d: '',
		   0x4e: '',
		   0x4f: '',
		   0x50: '',
		   0x51: '',
		   0x52: '',
		   0x53: '',
		   0x54: '',
		   0x55: '',
		   0x56: '',
		   0x57: '',
		   0x58: '',
		   0x59: '',
		   0x5a: '',
		   0x5b: '',
		   0x5c: '',
		   0x5d: '',
		   0x5e: '',
		   0x5f: '',
		   0x60: '',
		   0x61: '',
		   0x62: '',
		   0x63: '',
		   0x64: '',
		   0x65: '',
		   0x66: '',
		   0x67: '',
		   0x68: '',
		   0x69: '',
		   0x6a: '',
		   0x6b: '',
		   0x6c: '',
		   0x6d: '',
		   0x6e: '',
		   0x6f: '',
		   0x70: '',
		   0x71: '',
		   0x72: '',
		   0x73: '',
		   0x74: '',
		   0x75: '',
		   0x76: '',
		   0x77: '',
		   0x78: '',
		   0x79: '',
		   0x7a: '',
		   0x7b: '',
		   0x7c: '',
		   0x7d: '',
		   0x7e: '',
		   0x7f: '',
		*/
		0x80: '@',
		0x81: 'A',
		0x82: 'B',
		0x83: 'C',
		0x84: 'D',
		0x85: 'E',
		0x86: 'F',
		0x87: 'G',
		0x88: 'H',
		0x89: 'I',
		0x8a: 'J',
		0x8b: 'K',
		0x8c: 'L',
		0x8d: 'M',
		0x8e: 'N',
		0x8f: 'O',
		0x90: 'P',
		0x91: 'Q',
		0x92: 'R',
		0x93: 'S',
		0x94: 'T',
		0x95: 'U',
		0x96: 'V',
		0x97: 'W',
		0x98: 'X',
		0x99: 'Y',
		0x9a: 'Z',
		0x9b: '[',
		0x9c: '\\',
		0x9d: ']',
		0x9e: ' ',
		0x9f: ' ',
		0xa0: ' ',
		0xa1: '!',
		0xa2: '"',
		0xa3: '#',
		0xa4: '$',
		0xa5: '%',
		0xa6: '&',
		0xa7: '\'',
		0xa8: '(',
		0xa9: ')',
		0xaa: '*',
		0xab: '+',
		0xac: ',',
		0xad: '-',
		0xae: '.',
		0xaf: '/',
		0xb0: '0',
		0xb1: '1',
		0xb2: '2',
		0xb3: '3',
		0xb4: '4',
		0xb5: '5',
		0xb6: '6',
		0xb7: '7',
		0xb8: '8',
		0xb9: '9',
		0xba: ':',
		0xbb: ';',
		0xbc: '<',
		0xbd: '=',
		0xbe: '>',
		0xbf: '?',
		0xc0: '-',
		0xc1: 'a',
		0xc2: 'b',
		0xc3: 'c',
		0xc4: 'd',
		0xc5: 'e',
		0xc6: 'f',
		0xc7: 'g',
		0xc8: 'h',
		0xc9: 'i',
		0xca: 'j',
		0xcb: 'k',
		0xcc: 'l',
		0xcd: 'm',
		0xce: 'n',
		0xcf: 'o',
		0xd0: 'p',
		0xd1: 'q',
		0xd2: 'r',
		0xd3: 's',
		0xd4: 't',
		0xd5: 'u',
		0xd6: 'v',
		0xd7: 'w',
		0xd8: 'x',
		0xd9: 'y',
		0xda: 'z',
		/*
		   0xdb: '',
		   0xdc: '',
		   0xdd: '',
		   0xde: '',
		   0xdf: '',
		   0xe0: '',
		   0xe1: '',
		   0xe2: '',
		   0xe3: '',
		   0xe4: '',
		   0xe5: '',
		   0xe6: '',
		   0xe7: '',
		   0xe8: '',
		   0xe9: '',
		   0xea: '',
		   0xeb: '',
		   0xec: '',
		   0xed: '',
		   0xee: '',
		   0xef: '',
		   0xf0: '',
		   0xf1: '',
		   0xf2: '',
		   0xf3: '',
		   0xf4: '',
		   0xf5: '',
		   0xf6: '',
		   0xf7: '',
		   0xf8: '',
		   0xf9: '',
		   0xfa: '',
		   0xfb: '',
		   0xfc: '',
		   0xfd: '',
		   0xfe: '',
		   0xff: '',
		*/
	}
	color = sdl.Color{R: 0, G: 255, B: 0, A: 255}
)

type Video struct {
	Read func(address Word) Byte // Read a single byte from the bus

	lastChar rune
	window   *sdl.Window
	font     *ttf.Font
}

func (v *Video) Reset() error {
	// Load the font
	err := ttf.Init()
	if err != nil {
		return err
	}

	v.font, err = ttf.OpenFont("Quinquefive.ttf", fontSize)
	if err != nil {
		return err
	}

	// Initialize SDL & create a window
	err = sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return err
	}

	v.window, err = sdl.CreateWindow("pet",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		windowHeight,
		windowWidth,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}

	v.lastChar = rune(0)

	return nil
}

func (v *Video) Redraw() error {
	surface, err := v.window.GetSurface()
	if err != nil {
		return err
	}

	// Clear the surface
	surface.FillRect(nil, 0x0)

	// Draw 40 x 25
	for y := 0; y < 25; y++ {
		var lineBuffer bytes.Buffer

		// Assemble the line of text
		for x := 0; x < 40; x++ {
			addr := Word(VID_MEM) + Word(y*40) + Word(x)
			data := v.Read(addr)
			lineBuffer.WriteString(string(petscii[data]))
		}

		// Render & draw the line
		line, err := v.font.RenderUTF8Blended(lineBuffer.String(), color)
		if err != nil {
			return err
		}

		line.Blit(nil, surface, &sdl.Rect{X: borderLeft, Y: int32(borderTop + (y * 12)), W: 0, H: 0})
		if err != nil {
			return err
		}

		line.Free()
	}
	v.window.UpdateSurface()

	return nil
}

// Event returns any events that have been generated by the GUI
func (v *Video) Event() Event {
	sdlEvent := sdl.PollEvent()
	if sdlEvent != nil {
		switch event := sdlEvent.(type) {
		case *sdl.TextInputEvent:
			// Remember what was typed for the subsequent KEY_UP event
			v.lastChar = rune(event.Text[0])

			// Send a KEY_DOWN event for this printable key
			k := Keypress{
				Char:  v.lastChar,
				State: KEY_DOWN,
			}

			return EventKeypress{
				Key: k,
			}

		case *sdl.KeyboardEvent:
			k := Keypress{
				Keycode: int(event.Keysym.Sym),
			}

			// Only send KEY_DOWN events for non-printable keys are interesting
			if event.State == sdl.PRESSED {
				k.State = KEY_DOWN

				switch event.Keysym.Sym {
				case sdl.K_RETURN,
					sdl.K_ESCAPE,
					sdl.K_BACKSPACE:

					//fmt.Printf("%d (0x%02x)\n", k.Keycode, k.Keycode)

					return EventKeypress{
						Key: k,
					}
				default:
					// Ignore everything else; TextInputEvent will send KEY_DOWN events for those
					break
				}
			} else if event.State == sdl.RELEASED {
				k.State = KEY_UP

				// Send the last printable character that was typed with this KEY_UP and reset the key
				k.Char = v.lastChar
				v.lastChar = rune(0)

				return EventKeypress{
					Key: k,
				}
				/*
					switch event.Keysym.Mod {
					case sdl.KMOD_LSHIFT,
						sdl.KMOD_RSHIFT:
						k.Modifiers.Shift = true
					}
				*/

			}
		case *sdl.QuitEvent:
			return EventQuit{}
		}
	}
	return EventNone{}
}

func (v *Video) Stop() {
	v.window.Destroy()
	v.font.Close()

	ttf.Quit()
	sdl.Quit()
}
