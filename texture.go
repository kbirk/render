package render

import (
	"fmt"
	"image"
	"image/draw"
	// register png decoder
	_ "image/png"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	// DefaultWrapS is the default wrap S parameter.
	DefaultWrapS = gl.CLAMP_TO_EDGE
	// DefaultWrapT is the default wrap T parameter.
	DefaultWrapT = gl.CLAMP_TO_EDGE
	// DefaultMinFilter is the default min filter parameter.
	DefaultMinFilter = gl.NEAREST
	// DefaultMagFilter is the default mag filter parameter.
	DefaultMagFilter = gl.NEAREST
)

// Texture represents a 2D texture object.
type Texture struct {
	id             uint32
	width          uint32
	height         uint32
	format         uint32
	internalFormat int32
	typ            uint32
}

// TextureParams represents parameters for a 2D texture object.
type TextureParams struct {
	WrapS     int32
	WrapT     int32
	MinFilter int32
	MagFilter int32
}

// LoadRGBATexture loads an image file into an RGBA texture.
func LoadRGBATexture(filename string) (*Texture, error) {
	// load / decode image
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("texture file `%s` not found on disk: %v", filename, err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	// get RGBA
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	// create texture
	return NewRGBATexture(
		rgba.Pix,
		uint32(rgba.Rect.Size().X),
		uint32(rgba.Rect.Size().Y),
		&TextureParams{
			WrapS:     gl.CLAMP_TO_EDGE,
			WrapT:     gl.CLAMP_TO_EDGE,
			MinFilter: gl.LINEAR_MIPMAP_LINEAR,
			MagFilter: gl.LINEAR,
		}), nil
}

// NewRGBATexture returns a new RGBA texture.
func NewRGBATexture(rgba []uint8, width uint32, height uint32, params *TextureParams) *Texture {
	texture := &Texture{
		width:          width,
		height:         height,
		typ:            gl.UNSIGNED_BYTE,
		format:         gl.RGBA,
		internalFormat: gl.RGBA,
	}
	gl.GenTextures(1, &texture.id)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	// default params
	if params == nil {
		params = &TextureParams{}
	}
	if params.WrapS == 0 {
		params.WrapS = DefaultWrapS
	}
	if params.WrapT == 0 {
		params.WrapT = DefaultWrapT
	}
	if params.MinFilter == 0 {
		params.MinFilter = DefaultMinFilter
	}
	if params.MagFilter == 0 {
		params.MagFilter = DefaultMagFilter
	}
	// set params
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, params.MinFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, params.MagFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, params.WrapS)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, params.WrapT)

	// get pointer
	var data unsafe.Pointer
	if rgba != nil {
		data = gl.Ptr(rgba)
	}

	// buffer texture
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		texture.internalFormat,
		int32(texture.width),
		int32(texture.height),
		0,
		texture.format,
		texture.typ,
		data)

	// generate mipmaps
	if params.MinFilter == gl.LINEAR_MIPMAP_LINEAR ||
		params.MinFilter == gl.LINEAR_MIPMAP_NEAREST ||
		params.MinFilter == gl.NEAREST_MIPMAP_LINEAR ||
		params.MinFilter == gl.NEAREST_MIPMAP_NEAREST {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	}
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

// Width returns the width of the texture.
func (t *Texture) Width() uint32 {
	return t.width
}

// Height returns the height of the texture.
func (t *Texture) Height() uint32 {
	return t.height
}

// ID returns the ID of the texture.
func (t *Texture) ID() uint32 {
	return t.id
}

// Bind activates the provided texture unit and binds the texture.
func (t *Texture) Bind(location uint32) {
	gl.ActiveTexture(location)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

// Unbind will unbind the texture.
func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Resize will resize the texture, removing it's current buffer.
func (t *Texture) Resize(width uint32, height uint32) {
	t.width = width
	t.height = height
	gl.BindTexture(gl.TEXTURE_2D, t.id)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		t.internalFormat,
		int32(t.width),
		int32(t.height),
		0,
		t.format,
		t.typ,
		nil)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Destroy deallocates the texture buffer.
func (t *Texture) Destroy() {
	gl.DeleteTextures(1, &t.id)
	t.id = 0
}
