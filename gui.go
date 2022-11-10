package main

import (
	"context"
	"image/color"

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

type GUI struct {
	Video *Video

	lastChar rune
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

	g.lastChar = rune(0)
	sdl.StartTextInput()

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
			if event.State == sdl.PRESSED {
				k.State = KEY_DOWN
			} else if event.State == sdl.RELEASED {
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
