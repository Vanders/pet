package main

import (
	"bytes"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowHeight = 480
	windowWidth  = 360
	borderTop    = 30
	borderLeft   = 40
	fontSize     = 8
)

var (
	color = sdl.Color{R: 0, G: 255, B: 0, A: 255}
)

type Video struct {
	num int

	window *sdl.Window
	font   *ttf.Font
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

	return nil
}

func (v *Video) Redraw() error {
	surface, err := v.window.GetSurface()
	if err != nil {
		return err
	}

	// Clear the surface
	surface.FillRect(nil, 0x0)
	v.window.UpdateSurface()

	// Draw 40 x 25
	for y := 0; y < 25; y++ {
		var lineBuffer bytes.Buffer

		// Assemble the line of text
		for x := 0; x < 40; x++ {
			//lineBuffer.WriteRune('*')
			lineBuffer.WriteString(fmt.Sprintf("%d", v.num))
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

	v.num += 1
	if v.num > 9 {
		v.num = 0
	}

	return nil
}

func (v *Video) PollEvent() bool {
	event := sdl.PollEvent()
	if event != nil {
		switch event.(type) {
		case *sdl.QuitEvent:
			return true
		}
	}
	return false
}

func (v *Video) Stop() {
	v.window.Destroy()
	v.font.Close()

	ttf.Quit()
	sdl.Quit()
}

func main() {
	video := Video{}
	err := video.Reset()
	if err != nil {
		panic(err)
	}

	for {
		quit := video.PollEvent()
		if quit {
			break
		}

		// Wait 50ms and then redraw the screen
		delay := 50
		sdl.Delay(uint32(delay))
		err = video.Redraw()
		if err != nil {
			break
		}
	}

	// Clean up
	video.Stop()
}
