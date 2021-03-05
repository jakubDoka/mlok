package ggl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Setup2D is basic  game rendering setup, if you would like to change some part of it just embed it
// in your setup struct and override methods. Thought setup isn't just that, it brings lot of  utility
// to its scope to separate it from 3D setup (comming soon)
type Setup2D struct{}

// Batch creates batch from texture with fragment shader, it does the setup of shader program for you
func (s Setup2D) Batch(texture *Texture, fragmentShader string) (*Batch, error) {
	pg, err := NProgramFromSource(s.VertexShader(), fragmentShader)
	if err != nil {
		return nil, err
	}
	return &Batch{Data{}, nil, pg, texture}, nil
}

// Canvas creates canvas with custom fragment shader prepared for  rendering
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
	
	uniform mat3 camera;
	uniform vec2 viewportSize;
	uniform vec2 textureSize;
	
	out vec2 fragTex;
	out vec4 fragMask;
	out float fragIntensity;
	void main() {
		fragMask = mask;
		fragIntensity = intensity;
		fragTex = tex/textureSize;
		gl_Position = vec4((camera * vec3(vert, 1)).xy / viewportSize, 0, 1);
	}
	`
}

// FragmentShader implements Setup interface
func (s Setup2D) FragmentShader() string {
	return `
	#version 330
	#define WHITE vec4(1, 1, 1, 1)

	uniform sampler2D tex;
	uniform int useTexture;

	in vec2 fragTex;
	in vec4 fragMask;
	in float fragIntensity;

	out vec4 outputColor;

	void main() {
		if(fragIntensity == 1) {
			outputColor = texture(tex, fragTex) * fragMask;
		} else if(fragIntensity == 0) {
			outputColor = fragMask;
		} else {
			vec4 col = texture(tex, fragTex);
			outputColor = col + (WHITE - col) * (1 - fragIntensity);
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
	SetBlend(true)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
