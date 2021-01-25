package ggl

import (
	"gobatch/mat"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var blend bool

// SetBlend enables or disables blending
func SetBlend(val bool) {
	if blend == val {
		return
	}
	blend = val
	if val {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
}

// ClearMode is enum type used for clearing a framebuffer
type ClearMode uint32

// Clearing modes
const (
	Color   ClearMode = gl.COLOR_BUFFER_BIT
	Depth   ClearMode = gl.DEPTH_BUFFER_BIT
	Stencil ClearMode = gl.STENCIL_BUFFER_BIT
)

var color mat.RGBA

// Clear clears currently bound framebuffer with given mode
func Clear(c mat.RGBA, mode ClearMode) {
	if color != c {
		gl.ClearColor(float32(c.R), float32(c.G), float32(c.B), float32(c.A))
		color = c
	}

	gl.Clear(uint32(mode))
}
