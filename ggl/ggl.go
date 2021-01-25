// Package ggl defines essential abstraction over opengl types, you should use this and dodge gl calls to
// make you code cleaner. If struct has methdod Drop it has to be called in order to prevent memory leak.
// All main datatypes that has Drop method contain pointer to underlining gl object so droping one struct
// makes all copies dengling and using them can produce undefined behavior. All constructors starts with N
// folwed by name of a struct. If struct ha a constructor it has to be used or you will not get expected
// behavior. Other methods that are repetitive are Start() and End<StructName>(), this is for binding and
// unbinding objects. Because calling gl functions introduces 120 nanoseconds overhead so package uses global
// state to skip redundant start and end calls, using of global state si justified by the fact that YOU HAVE
// TO USE ALL STRUCTS INTERFACING OPENGL OBJECTS FROM THREAD WHERE YOU CREATED A WINDOW anyway.
//
// Package also utilizes gogen code generator (github.com/jakubDoka/gogen)
package ggl

import (
	"fmt"
	"gobatch/mat"
	"gobatch/refl"
	"image"
	"image/draw"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// VertexData is what you use to hand data to opengl, its very important for you to pass
// a Slice containing doubles or structs composed only from doubles. If you are not sure
// whether your data is valid call VerifyVertexData with instance of your data, it will
// return helpfull error if behavior is not satisfied, though function is slow thats why
// its not used to infer parameters that you can provide with lot lower overhead.
type VertexData interface {
	// Len simply returns length of slice
	Len() int
	// VertexSize is amount of floats one element of VertexData contains
	VertexSize() int
}

// VerifyVertexData should be used in UnitTest to verify that type you are using for VertexData is safe
// for use. Function uses heavy reflection and is not suited for ordinary usage.
func VerifyVertexData(v VertexData) error {
	tp := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	if tp.Kind() != reflect.Slice {
		return fmt.Errorf("unexpected data kind, expected: %v got: %v", reflect.Slice.String(), tp.Kind().String())
	}

	if val.Len() < 2 {
		return fmt.Errorf("Tested data has to have at least len 2")
	}

	if v.Len() != val.Len() {
		return fmt.Errorf("Len does not match, expacted: %v, got: %v", val.Len(), v.Len())
	}

	if tp.Elem().Size()/componentByteSize != uintptr(v.VertexSize()) {
		return fmt.Errorf("the VertexSize returns invalid value, expected: %v, got: %v", tp.Elem().Size()/componentByteSize, v.VertexSize())
	}

	if err := refl.AssertHomogeneity(reflect.TypeOf(v), reflect.TypeOf(float64(0))); err != nil {
		return fmt.Errorf("slice element contains other datatypes other then float64: %v", err)
	}

	return nil
}

// Indices are drawing indices passed to opengl to lower data traffic. Theier use is optional
// you can always pass nil if you don't want to use them. Though mind that if Batch already contain
// indices produced by other drawer you have to use them otherwise your triangles will not get rendered.
type Indices []uint32

// Clear clears indices
func (i *Indices) Clear() {
	*i = (*i)[:0]
}

// VAO abstracts over opengl vertex array
type VAO struct {
	Ptr
}

// NVAO ...
func NVAO() VAO {
	v := VAO{}
	gl.GenVertexArrays(1, &v.ptr)
	return v
}

// Start ...
func (v VAO) Start() {
	setVao(v.ptr)
}

// EndVao ...
func EndVao() {
	setVao(0)
}

var vao uint32

func setVao(nc uint32) {
	if vao == nc {
		return
	}
	vao = nc

	gl.BindVertexArray(vao)
}

// Drop ...
func (v VAO) Drop() {
	gl.DeleteVertexArrays(1, &v.ptr)
}

const componentByteSize = 8

// Static is used when you want to draw once and reuse multiple times
const Static = gl.STATIC_DRAW

// Stream is used when you want to update a lot and use at most few times
const Stream = gl.STREAM_DRAW

// Dynamic combines Static and stream
const Dynamic = gl.DYNAMIC_DRAW

// VBO abstracts over opengl array buffer
type VBO struct {
	Ptr
	vertexSize int
}

// NVBO initializes buffer structure:
//
// 	ggl.NVBO(2, 2, 4) // initializes buffer with 3 pointers with given sizes
//
// Then in your vertex shader you access them like:
//
//	layout (location = 0) in vec2 vert;
//	layout (location = 1) in vec2 tex;
//	layout (location = 2) in vec4 mask;
//
// bdw this is default 2D buffer setup
func NVBO(vertexSizes ...int32) VBO {
	v := VBO{}
	gl.GenBuffers(1, &v.ptr)

	v.Start()
	v.setup(vertexSizes...)

	return v
}

// SetData loads data into VBO, for mode you can use gl version, or redefined package modes
// that also have some comments, Stream is mostly the best option.
//
// panics if buffer vertexSize does not match with data vertex size
func (v *VBO) SetData(data VertexData, mode uint32) {

	v.Start()
	vertSz := data.VertexSize()
	if vertSz != v.vertexSize {
		panic(fmt.Errorf("unexpected vertex size, buffer expects %v, but vertex data with vertexSize %v was inputted", v.vertexSize, vertSz))
	}

	if data.Len() == 0 {
		gl.BufferData(gl.ARRAY_BUFFER, 0, nil, mode)
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, data.Len()*vertSz*componentByteSize, gl.Ptr(data), mode)
	}
}

// Setup sets up the buffer structure, can be used only once
func (v *VBO) setup(vertexSizes ...int32) {
	var t int32
	for _, s := range vertexSizes {
		t += s
	}

	v.vertexSize = int(t)

	var off int32
	for i, s := range vertexSizes {
		i := uint32(i)
		gl.EnableVertexAttribArray(i)
		gl.VertexAttribLPointer(i, s, gl.DOUBLE, t*componentByteSize, gl.PtrOffset(int(off)*componentByteSize))
		off += s
	}
}

// Start ...
func (v *VBO) Start() {
	setVbo(v.ptr)
}

// EndVbo ...
func EndVbo() {
	setVbo(0)
}

var vbo uint32

func setVbo(nc uint32) {
	if vbo == nc {
		return
	}
	vbo = nc

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
}

// Drop ...
func (v *VBO) Drop() {
	gl.DeleteBuffers(1, &v.ptr)
}

// Shader represents gl sager object
type Shader struct {
	Ptr
}

// LoadVertexShader ...
func LoadVertexShader(p string) (*Shader, error) {
	return LoadShader(p, gl.VERTEX_SHADER)
}

// LoadFragmentShader ...
func LoadFragmentShader(p string) (*Shader, error) {
	return LoadShader(p, gl.FRAGMENT_SHADER)
}

// NVertexShader ...
func NVertexShader(source string) (*Shader, error) {
	return NShader(source, gl.VERTEX_SHADER)
}

// NFragmentShader ...
func NFragmentShader(source string) (*Shader, error) {
	return NShader(source, gl.FRAGMENT_SHADER)
}

// LoadShader loads shader from disk
func LoadShader(p string, shaderType uint32) (*Shader, error) {
	bts, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return NShader(string(bts), shaderType)
}

// NShader creates shader from source, provided type of shader (gl.VERTEX_SHADER or gl.FRAGMENT_SHADER)
func NShader(source string, shaderType uint32) (*Shader, error) {
	s := &Shader{}
	s.ptr = gl.CreateShader(shaderType)

	src, free := gl.Strs(source + "\x00")
	gl.ShaderSource(s.ptr, 1, src, nil)
	free()
	gl.CompileShader(s.ptr)

	var status int32
	gl.GetShaderiv(s.ptr, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(s.ptr, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(s.ptr, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return s, nil
}

// Drop ...
func (s Shader) Drop() {
	gl.DeleteShader(s.ptr)
}

// Texture is a handle to opengl texture object
type Texture struct {
	Ptr
	W, H int32
}

// LoadTexture loads texture from disk and creates gl texture from it:
//
//	tex, err := LoadTexture("maPath.png")
//
// Use NTexture instead if you have to modify it before turning it into gl object,
// Note that texture is fed with default params and converted ti *image.RGBA if format
// does not match
func LoadTexture(p string) (*Texture, error) {
	imgFile, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", p, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	return NTexture(img), nil
}

// NTexture creates texture with default params
func NTexture(img image.Image) *Texture {
	return NTextureWithParams(img, DefaultTextureConfig...)
}

// NTextureWithParams allows you to specify additional parameters when creating texture,
// paramethers are then in code set like:
//
// 	for _, p := range params {
//		gl.TexParameteri(gl.TEXTURE_2D, p.First, p.Second)
//	}
//
// Note that even though this function takes image.Image it will always convert it to image.RGBA
func NTextureWithParams(img image.Image, params ...TextureParam) *Texture {
	var rgba *image.RGBA

	if val, ok := img.(*image.RGBA); ok {
		rgba = val
	} else {
		rgba = image.NewRGBA(img.Bounds())

		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	}
	return RawTexture(int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), rgba.Pix, params...)
}

// RawTexture creates texture out of raw parts, if you do not provide bytes slice texture will be empty
// but opengl will allocate memory
func RawTexture(w, h int32, pixels []byte, params ...TextureParam) *Texture {
	t := &Texture{}

	gl.GenTextures(1, &t.ptr)
	t.Start()

	for _, p := range params {
		gl.TexParameteri(gl.TEXTURE_2D, p.First, p.Second)
	}

	t.Resize(w, h, pixels)

	EndTexture()

	return t
}

// Resize Resizes the texture resizing does not maintain texture content.
// If texture is already in given size and you don't
// input new pixels, this function does nothing
func (t *Texture) Resize(w, h int32, pixels []byte) {
	var ptr unsafe.Pointer
	if len(pixels) == 0 {
		if t.W == w && t.H == h {
			return
		}
	} else {
		ptr = gl.Ptr(pixels)
	}

	t.W = w
	t.H = h

	t.Start()
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		w,
		h,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		ptr,
	)
}

// Image withdraws texture data from gpu, this is mainly usefull for capturing framebuffer state
// it basically makes recording possible
func (t *Texture) Image() *image.RGBA {
	tex := image.NewRGBA(image.Rect(0, 0, int(t.W), int(t.H)))

	t.Start()
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(tex.Pix))

	return tex
}

// SubImage returns specified part of texture.
//
// panics if region does not fit into texture (only partially overlaps with texture)
func (t *Texture) SubImage(region image.Rectangle) *image.RGBA {
	if region.Min.X < 0 || region.Min.Y < 0 || region.Max.X > int(t.W) || region.Max.Y > int(t.H) {
		panic(fmt.Errorf("texture is of size [%v, %v] so region (%v) does not fit", t.W, t.H, region))
	}

	tex := image.NewRGBA(region)
	const rgbaLen = 4

	t.Start()
	gl.GetTextureSubImage(
		gl.TEXTURE_2D,
		0,
		int32(region.Min.X),
		int32(region.Min.Y),
		0,
		int32(region.Dx()),
		int32(region.Dy()),
		1,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		int32(region.Dx()*region.Dy()*rgbaLen),
		gl.Ptr(tex.Pix),
	)

	return tex
}

// Frame returns texture frame, origin is always at [0, 0]
func (t *Texture) Frame() mat.AABB {
	return mat.NAABB(0, 0, float64(t.W), float64(t.H))
}

// Start ...
func (t *Texture) Start() {
	setTexture(t.ptr)
}

// EndTexture ...
func EndTexture() {
	setTexture(0)
}

var texture uint32

func setTexture(nc uint32) {
	if texture == nc {
		return
	}
	texture = nc

	gl.BindTexture(gl.TEXTURE_2D, texture)
}

// Drop ...
func (t Texture) Drop() {
	gl.DeleteTextures(1, &t.ptr)
}

// TextureParam specifies what is passed to gl.TextureParametri
type TextureParam struct {
	First  uint32
	Second int32
}

// FlipRGBA flips image along the Y-axis
func FlipRGBA(r *image.RGBA) {
	height := r.Rect.Dy()
	hh := height / 2
	row := r.Rect.Dx() * 4
	tmp := make([]byte, row)
	for i := 0; i < hh; i++ {
		nw := row * i
		iv := row * (height - 1 - i)
		a, b := r.Pix[nw:nw+row], r.Pix[iv:iv+row]

		copy(tmp, a)
		copy(a, b)
		copy(b, tmp)
	}
}

// DefaultTextureConfig use this if you don't know what to use
var DefaultTextureConfig = []TextureParam{
	{gl.TEXTURE_MIN_FILTER, gl.NEAREST},
	{gl.TEXTURE_MAG_FILTER, gl.NEAREST},
}

// Ptr is convenience struct, is sores pointer and only allows reading it
type Ptr struct {
	ptr uint32
}

// ID returns pointer valuse
func (p Ptr) ID() uint32 {
	return p.ptr
}
