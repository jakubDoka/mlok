package ggl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Setup2D is basic 2D game rendering setup, if you would like to change some part of it just embed it
// in your setup struct and override methods. Thought setup isn't just that, it brings lot of 2D utility
// to its scope to separate it from 3D setup (comming soon)
type Setup2D struct{}

// Batch creates batch from texture with fragment shader, it does the setup of shader program for you
func (s Setup2D) Batch(texture *Texture, fragmentShader string) (*Batch2D, error) {
	pg, err := NProgramFromSource(s.VertexShader(), fragmentShader)
	if err != nil {
		return nil, err
	}
	return NBatch2D(texture, nil, pg), nil
}

// Canvas creates canvas with custom fragment shader prepared for 2D rendering
func (s Setup2D) Canvas(w, h int, fragmentShader string) (*Canvas, error) {
	pg, err := NProgramFromSource(s.VertexShader(), fragmentShader)
	if err != nil {
		return nil, err
	}

	return NCanvas(*RawTexture(int32(w), int32(h), nil, DefaultTextureConfig...), *pg, s.Buffer()), nil
}

// VertexShader implements Setup interface
func (s Setup2D) VertexShader() string {
	return `
	#version 330
	
	layout (location = 0) in vec2 vert;
	layout (location = 1) in vec2 tex;
	layout (location = 2) in vec4 mask;
	layout (location = 3) in float intensity;
	
	uniform mat3 camera2D;
	uniform vec2 viewportSize;
	uniform vec2 textureSize;
	
	out vec2 fragTex;
	out vec4 fragMask;
	out float fragIntensity;
	void main() {
		fragMask = mask;
		fragIntensity = intensity;
		fragTex = tex/textureSize;
		gl_Position = vec4(camera2D * vec3(vert/viewportSize, 0), 1);
	}
	`
}

// FragmentShader implements Setup interface
func (s Setup2D) FragmentShader() string {
	return `
	#version 330

	uniform sampler2D tex;
	uniform int useTexture;

	in vec2 fragTex;
	in vec4 fragMask;
	in float fragIntensity;

	out vec4 outputColor;

	void main() {
		if(fragIntensity == 1) {
			outputColor = texture(tex, fragTex) * fragMask;
		} else {
			outputColor = fragMask;
		}
	}
	`
}

// Buffer implements Setup interface
func (s Setup2D) Buffer() Buffer {
	return NBuffer(2, 2, 4, 1)
}

// Modify implements Setup interface
func (s Setup2D) Modify(win *Window) {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
