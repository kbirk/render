package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// IndexBuffer represents an indexbuffer.
type IndexBuffer struct {
	id uint32
}

// BufferUint8 allocates uint8 buffer data.
func (i *IndexBuffer) BufferUint8(data []uint8) {
	if i.id == 0 {
		gl.GenBuffers(1, &i.id)
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, i.id)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data), gl.Ptr(data), gl.STATIC_DRAW)
}

// BufferUint16 allocates uint16 buffer data.
func (i *IndexBuffer) BufferUint16(data []uint16) {
	if i.id == 0 {
		gl.GenBuffers(1, &i.id)
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, i.id)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data)*2, gl.Ptr(data), gl.STATIC_DRAW)
}

// BufferUint32 allocates uint32 buffer data.
func (i *IndexBuffer) BufferUint32(data []uint32) {
	if i.id == 0 {
		gl.GenBuffers(1, &i.id)
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, i.id)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
}

// Bind binds the indexbuffer.
func (i *IndexBuffer) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, i.id)
}

// Unbind unbinds the indexbuffer.
func (i *IndexBuffer) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

// Draw renders the indexbuffer.
func (i *IndexBuffer) Draw(mode uint32, count int32, typ uint32, byteOffset int) {
	gl.DrawElements(mode, count, typ, gl.PtrOffset(byteOffset))
}

// DrawInstanced renders multiple instances of the indexbuffer.
func (i *IndexBuffer) DrawInstanced(mode uint32, count int32, typ uint32, byteOffset int, primcount int32) {
	gl.DrawElementsInstanced(mode, count, typ, gl.PtrOffset(byteOffset), primcount)
}

// Destroy deallocates the indexbuffer.
func (i *IndexBuffer) Destroy() {
	if i.id != 0 {
		gl.DeleteBuffers(1, &i.id)
		i.id = 0
	}
}
