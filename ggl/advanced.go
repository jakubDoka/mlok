package ggl

import (
	"fmt"
	"gobatch/mt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Target is something that you can draw to, last tree parameters can be optional and not even be used
// by a target, though if you don't provide tham when you should Target can fall back to defaults or panic
type Target interface {
	Accept(data VertexData, indices Indices, texture *Texture, program *Program, buffer *Buffer)
}

// Canvas allows of screen drawing, drawing to canvas produces draw calls. Its the abstraction
// over opengl framebuffer. It stores drawn image in given texture, if you want to capture it use
// Image method on texture. Program that canvas uses is applied on resulting image but also on triangles
// drawn by batch that does not have custom program. Same goes for Buffer. you can also set canvas state,
// see SetPixels method.
type Canvas struct {
	Ptr

	Texture
	Program
	Buffer Buffer

	ClearColor mt.RGBA

	data2D   Data2D
	sprite2D Sprite2D
}

// NCanvas creates new framebuffer, all three arguments has to e valid instances of
// gl objects for canvas to work.
//
// 	buff := ggl.NBuffer(2, 2, 4) // see buffer doc
//	prog := ggl.LoadProgram("yourVertex.glsl", "yourFragment.glsl")
//	texture := ggl.RawTexture(canvasInitialWidth, canvasInitialHeight, nil, DefaultTextureConfig)
//  canvas := NCanvas(*texture, *program, buffer)
//
// the texture you are creating is a drawing target, you can of corse use existing texture,draw to it
// and then capture the result via image:
//
// 	img := canvas.Image() // returning savable image
//
func NCanvas(texture Texture, program Program, buffer Buffer) *Canvas {
	c := &Canvas{
		Texture: texture,
		Program: program,
		Buffer:  buffer,
	}

	gl.GenFramebuffers(1, &c.ptr)
	c.Start()

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture.ID(), 0)
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		panic("failed to create framebuffer")
	}

	c.Resize(texture.Frame())
	program.SetCamera2D(&mt.IM2)

	return c
}

// Accept implements Target interface
func (c *Canvas) Accept(data VertexData, indices Indices, texture *Texture, program *Program, buffer *Buffer) {
	c.Start()
	p := &c.Program
	if program != nil {
		p = program
	}

	p.Start()

	if texture != nil {
		texture.Start()
		p.SetTextureSize(texture.W, texture.H)
		p.SetUseTexture(true)
	} else {
		p.SetUseTexture(false)
	}

	p.SetViewportSize(c.W, c.H)

	b := &c.Buffer
	if buffer != nil {
		b = buffer
	}

	b.Draw(data, indices, Stream)
}

// Clear2D clears canvas in 2D mode
func (c *Canvas) Clear2D(color mt.RGBA) {
	c.Clear(color, Color)
}

// Clear clears canvas with given color
func (c *Canvas) Clear(color mt.RGBA, mode ClearMode) {
	c.Start()
	Clear(color, mode)
}

var canvas uint32

// Start ...
func (c *Canvas) Start() {
	setCanvas(c.ptr)
}

// EndCanvas unbinds current canvas
func EndCanvas() {
	setCanvas(0)
}

func setCanvas(nc uint32) {
	/*if canvas == nc {
		return
	}*/
	canvas = nc

	gl.BindFramebuffer(gl.FRAMEBUFFER, canvas)
}

// Draw2D draws canvas to another target as a 2D sprite
//
// 	c.Render2D(t, mt.IM2, mt.Alpha(1)) //draws framebuffer to the center of a screen as it is to t
// 	c.Render2D(t, mt.IM2.Scaled(mt.V2{}, 2), mt.RGB(1, 0, 0)) //draws framebuffer scaled up with red mask to t
//
// method makes draw call
func (c *Canvas) Draw2D(t Target, mat mt.Mat2, mask mt.RGBA) {
	c.data2D.Clear()
	c.sprite2D.Draw(&c.data2D, mat, mask)

	t.Accept(c.data2D.Vertexes, c.data2D.indices, &c.Texture, &c.Program, &c.Buffer)
}

// Render2D renders canvas to main framebuffer (window framebuffer) as a 2D sprite:
//
// 	c.Render2D(mt.IM2, mt.Alpha(1)) //draws framebuffer to the center of a screen as it is
// 	c.Render2D(mt.IM2.Scaled(mt.V2{}, 2), mt.RGB(1, 0, 0)) //draws framebuffer scaled up with red mask
//
// method makes draw call
func (c *Canvas) Render2D(mat mt.Mat2, mask mt.RGBA, w, h int32) {
	c.data2D.Clear()
	c.sprite2D.Draw(&c.data2D, mat, mask)
	EndCanvas()

	c.Texture.Start()
	c.Program.Start()

	c.Program.SetViewportSize(w, h)
	c.Program.SetTextureSize(c.W, c.H)
	c.Program.SetUseTexture(true)

	c.Buffer.Draw(c.data2D.Vertexes, c.data2D.indices, Stream)
}

// Resize resizes the canvas to given frame, canvas viewport is also set, that why you are passing frame,
func (c *Canvas) Resize(frame mt.AABB) {
	c.Start()
	c.Texture.Resize(int32(frame.W()), int32(frame.H()), nil)
	c.sprite2D = NSprite2D(frame.Moved(frame.Min.Inv()))
	gl.Viewport(int32(frame.Min.X), int32(frame.Min.Y), int32(frame.W()), int32(frame.H()))
}

// Drop ...
func (c *Canvas) Drop() {
	gl.DeleteFramebuffers(1, &c.ptr)
}

// Buffer combines VAO and VBO to finally draw to current frame buffer
type Buffer struct {
	VAO VAO
	VBO VBO
}

// NBuffer setups a buffer, sizes are passed to VBO constructor, indices determine
// whether you want to use EBO or not
func NBuffer(sizes ...int32) Buffer {
	b := Buffer{}

	b.VAO = NVAO()
	b.VAO.Start()
	b.VBO = NVBO(sizes...)

	return b
}

// Draw draws data with optional indices, if indices are nil or with length 0 classic
// drawcall will be triggered
func (b *Buffer) Draw(data VertexData, indices Indices, mode uint32) {
	b.VAO.Start()
	b.VBO.Start()
	b.VBO.SetData(data, mode)
	if len(indices) != 0 {
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.Ptr(indices))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(data.Len()))
	}
}

// Drop ...
func (b *Buffer) Drop() {
	b.VAO.Drop()
	b.VBO.Drop()
}

// Program is handle to opengl shader program
type Program struct {
	Ptr
	viewpot mt.V2
	texture bool
}

// LoadProgram loads program from disk
func LoadProgram(vertexPath, fragmentPath string) (*Program, error) {
	vertex, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		return nil, err
	}

	fragment, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		return nil, err
	}

	return NProgramFromSource(string(vertex), string(fragment))
}

// NProgramFromSource ...
func NProgramFromSource(vertexSource, fragmentSource string) (*Program, error) {
	vertex, err := NVertexShader(vertexSource)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse vertex shader: %v", err)
	}

	fragment, err := NFragmentShader(fragmentSource)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse fragment shader: %v", err)
	}

	p, err := NProgram(*vertex, *fragment)

	vertex.Drop()
	fragment.Drop()

	return p, err
}

// NProgram links vertex and fragment shader into program
func NProgram(vertex, fragment Shader) (*Program, error) {
	p := &Program{}
	p.ptr = gl.CreateProgram()

	gl.AttachShader(p.ptr, vertex.ptr)
	gl.AttachShader(p.ptr, fragment.ptr)
	gl.LinkProgram(p.ptr)

	var status int32
	gl.GetProgramiv(p.ptr, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(p.ptr, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p.ptr, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	return p, nil
}

// SetTextureSize sets "textureSize" in vertex shader, if the size is already equal to
// given values it does nothing
func (p *Program) SetTextureSize(w, h int32) {
	sz := mt.NV2(float64(w), float64(h))
	if p.viewpot == sz {
		return
	}
	p.viewpot = sz

	p.SetV2("textureSize", sz)
}

// SetViewportSize sets "viewportSize" in vertex shader, if the size is already equal to
// given values it does nothing
func (p *Program) SetViewportSize(w, h int32) {
	sz := mt.NV2(float64(w), float64(h)).Scaled(.5)
	if p.viewpot == sz {
		return
	}
	p.viewpot = sz

	p.SetV2("viewportSize", sz)
}

// SetCamera2D sets "camera2D" field in fragment shader
func (p *Program) SetCamera2D(mat *mt.Mat2) {
	p.SetMat2("camera2D", mat)
}

// SetUseTexture sets "useTexture" in fragment shader, its noop if
// value is already in given state
func (p *Program) SetUseTexture(b bool) {
	if b == p.texture {
		return
	}
	p.texture = b
	var a int32
	if b {
		a = 1
	}
	p.SetInt("useTexture", a)
}

// SetInt sets int uniform
func (p *Program) SetInt(name string, i int32) {
	gl.ProgramUniform1i(p.ptr, p.adr(name), i)
}

// SetV2 sets vec2 uniform
func (p *Program) SetV2(name string, v mt.V2) {
	gl.ProgramUniform2f(p.ptr, p.adr(name), float32(v.X), float32(v.Y))
}

// SetMat2 sets mat3 uniform
func (p *Program) SetMat2(name string, m *mt.Mat2) {
	mat := m.Raw()
	gl.ProgramUniformMatrix3fv(p.ptr, p.adr(name), 1, false, &mat[0])
}

func (p *Program) adr(name string) int32 {
	return gl.GetUniformLocation(p.ptr, gl.Str(name+"\x00"))
}

// Start ...
func (p *Program) Start() {
	setProgram(p.ptr)
}

// EndProgram ...
func EndProgram() {
	setProgram(0)
}

var program uint32

func setProgram(nc uint32) {
	if program == nc {
		return
	}
	program = nc

	gl.UseProgram(program)
}

// Drop ...
func (p *Program) Drop() {
	gl.DeleteProgram(p.ptr)
}
