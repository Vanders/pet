package main

import (
	"context"
	"image/color"
	"unicode/utf8"

	"github.com/sqweek/dialog"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowHeight = 480
	windowWidth  = 360
)

var scancodes = map[sdl.Keycode]Byte{
	/* row 9 */
	sdl.K_EQUALS:      0x3d,
	sdl.K_PERIOD:      0x2e,
	sdl.K_ESCAPE:      0x03,
	sdl.K_LESS:        0x3c,
	sdl.K_SPACE:       0x20,
	sdl.K_LEFTBRACKET: 0x5b,
	sdl.K_LCTRL:       0x12, // Reverse

	/* row 8 */
	sdl.K_MINUS:        0x2d,
	sdl.K_0:            0x30,
	sdl.K_RSHIFT:       0x00,
	sdl.K_GREATER:      0x3e,
	sdl.K_RIGHTBRACKET: 0x5d,
	sdl.K_AT:           0x40,
	// We're going to need a non-mapped shift key for shifted input characters
	// sdl.K_LSHIFT:       0x00,

	/* row 7 */
	sdl.K_PLUS:     0x2b,
	sdl.K_2:        0x32,
	sdl.K_QUESTION: 0x3f,
	sdl.K_COMMA:    0x2c,
	sdl.K_n:        0x4e,
	sdl.K_v:        0x56,
	sdl.K_x:        0x58,

	/* row 6 */
	sdl.K_3:         0x33,
	sdl.K_1:         0x31,
	sdl.K_RETURN:    0x0d,
	sdl.K_SEMICOLON: 0x3b,
	sdl.K_m:         0x4d,
	sdl.K_b:         0x42,
	sdl.K_c:         0x43,
	sdl.K_z:         0x5a,

	/* row 5 */
	sdl.K_ASTERISK: 0x2a, // *
	sdl.K_5:        0x35,
	sdl.K_COLON:    0x3a,
	sdl.K_k:        0x4b,
	sdl.K_h:        0x48,
	sdl.K_f:        0x46,
	sdl.K_s:        0x53,

	/* row 4 */
	sdl.K_6: 0x36,
	sdl.K_4: 0x34,
	sdl.K_l: 0x4c,
	sdl.K_j: 0x4a,
	sdl.K_g: 0x47,
	sdl.K_d: 0x44,
	sdl.K_a: 0x41,

	/* row 3 */
	sdl.K_SLASH: 0x2f,
	sdl.K_8:     0x38,
	sdl.K_p:     0x50,
	sdl.K_i:     0x49,
	sdl.K_y:     0x59,
	sdl.K_r:     0x52,
	sdl.K_w:     0x57,

	/* row 2 */
	sdl.K_9:     0x39,
	sdl.K_7:     0x37,
	sdl.K_CARET: 0x5e, // ^
	sdl.K_o:     0x4f,
	sdl.K_u:     0x55,
	sdl.K_t:     0x54,
	sdl.K_e:     0x45,
	sdl.K_q:     0x51,

	/* row 1 */
	sdl.K_DELETE:     0x14,
	sdl.K_UP:         0x11, // UP & DOWN are the same key, shifted
	sdl.K_DOWN:       0x11,
	sdl.K_RIGHTPAREN: 0x29, // )
	sdl.K_BACKSLASH:  0x5c,
	sdl.K_QUOTE:      0x27, // '
	sdl.K_DOLLAR:     0x24,
	sdl.K_QUOTEDBL:   0x22, // "

	/* row 0 */
	sdl.K_LEFT:      0x1d, // LEFT & RIGHT are the same key, shifted
	sdl.K_RIGHT:     0x1d,
	sdl.K_HOME:      0x13,
	sdl.K_BACKSPACE: 0x14,
	sdl.K_LEFTPAREN: 0x28, // (
	sdl.K_AMPERSAND: 0x26, // &
	sdl.K_PERCENT:   0x25, // %
	sdl.K_HASH:      0x23, // #
	sdl.K_EXCLAIM:   0x21, // !
}

type Remapper struct {
	runeToScan map[rune]Byte
	remapped   map[Byte]Byte
}

func (r *Remapper) Init() {
	r.runeToScan = map[rune]Byte{
		'!':  0x21,
		'"':  0x22,
		'#':  0x23,
		'$':  0x24,
		'%':  0x25,
		'&':  0x26,
		'\'': 0x27,
		'(':  0x28,
		')':  0x29,
		'*':  0x2a,
		'+':  0x2b,
		',':  0x2c,
		'-':  0x2d,
		'.':  0x2e,
		'/':  0x2f,
		':':  0x3a,
		'<':  0x3c,
		'>':  0x3e,
		'?':  0x3f,
		'@':  0x40,
		'\\': 0x5c,
		'^':  0x5e,
	}
	r.remapped = make(map[Byte]Byte, 0xff)
}

func (r *Remapper) Down(c rune, scancode Byte) Byte {
	newScancode, ok := r.runeToScan[c]
	if ok {
		// Remap the old scancode to the replacement & remember for the subsequent key up
		r.remapped[scancode] = newScancode
		return newScancode
	}
	return scancode
}

func (r *Remapper) Up(scancode Byte) Byte {
	if r.remapped[scancode] != 0x00 {
		// This scancode was remapped on key down, so send the remapped scancode for key up
		oldScancode := r.remapped[scancode]
		r.remapped[scancode] = 0x00
		return oldScancode
	}
	return scancode
}

type GUI struct {
	Video *Video

	remapper *Remapper
	window   *sdl.Window
}

func (g *GUI) Init() error {
	// Initialize SDL & create a window
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return err
	}

	g.window, err = sdl.CreateWindow("pet",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		windowHeight,
		windowWidth,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}

	g.remapper = &Remapper{}
	g.remapper.Init()

	return nil
}

func (g *GUI) redraw() error {
	surface, err := g.window.GetSurface()
	if err != nil {
		return err
	}

	// Clear the surface & lock it
	surface.FillRect(nil, 0x0)
	surface.Lock()

	// Perform video refresh
	green := color.RGBA{G: uint8(255)}

	g.Video.Redraw(func(x, y int) {
		surface.Set(x, y, green)
	})

	// Unlock the surface & update the window
	surface.Unlock()
	g.window.UpdateSurface()

	return nil
}

// EventLoop handles any events generated by the GUI
func (g *GUI) EventLoop(ctx context.Context, events chan<- Event) {
	lastTicks := sdl.GetTicks()
	currentTicks := lastTicks

	for {
		select {
		case <-ctx.Done():
			return
		default:
			break
		}

		sdlEvent := sdl.PollEvent()
		switch event := sdlEvent.(type) {
		case *sdl.KeyboardEvent:
			sym := event.Keysym.Sym
			scancode, ok := scancodes[sym]
			if !ok {
				break
			}
			k := Keypress{
				Scancode: scancode,
			}

			/* fuck you SDL2 */
			nextEvents := make([]sdl.Event, 10)
			inChar := rune(0)

			sdl.PumpEvents()
			sdl.PeepEvents(nextEvents, sdl.GETEVENT, sdl.TEXTINPUT, sdl.TEXTINPUT)
			nextEvent := nextEvents[0]
			if nextEvent != nil {
				inputEvent := nextEvent.(*sdl.TextInputEvent)
				text := inputEvent.GetText()
				inChar, _ = utf8.DecodeRuneInString(text[0:])
			}

			if event.State == sdl.PRESSED {
				// Lookup inChar and replace the scancode if we have a match
				k.Scancode = g.remapper.Down(inChar, scancode)
				k.State = KEY_DOWN
			} else if event.State == sdl.RELEASED {
				// If this scancode was remapped on KEY_DOWN then replace the scancode for KEY_UP
				k.Scancode = g.remapper.Up(scancode)
				k.State = KEY_UP
			}

			// Some keys are shifted and we need to make that explicit
			switch sym {
			case sdl.K_UP,
				sdl.K_LEFT,
				sdl.K_LSHIFT,
				sdl.K_RSHIFT:
				k.Shifted = true
			}

			// Send new key press event
			events <- EventKeypress{
				Key: k,
			}
		case *sdl.QuitEvent:
			events <- EventQuit{}
		}

		// Wait 50ms and then redraw the screen
		currentTicks = sdl.GetTicks()
		if currentTicks > lastTicks+50 {
			err := g.redraw()
			if err != nil {
				break
			}

			lastTicks = currentTicks
		}

		sdl.Delay(10)
	}
}

func (g *GUI) Stop() {
	g.window.Destroy()

	sdl.Quit()
}

func (g *GUI) LoadDialog(title, filter, ext string) (string, error) {
	return dialog.File().Filter(filter, ext).Title(title).Load()
}

func (g *GUI) SaveDialog(title, filter, ext string) (string, error) {
	return dialog.File().Filter(filter, ext).Title(title).Save()
}
