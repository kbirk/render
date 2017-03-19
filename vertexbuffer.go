package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// VertexBuffer represents a vertexbuffer.
type VertexBuffer struct {
	id uint32
}

// AllocateBuffer allocates the size of the underlying buffer.
func (v *VertexBuffer) AllocateBuffer(numBytes int) {
	if v.id == 0 {
		gl.GenBuffers(1, &v.id)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, v.id)
	gl.BufferData(gl.ARRAY_BUFFER, numBytes, gl.Ptr(nil), gl.STATIC_DRAW)
}

// BufferFloat32 buffers a float32 slice.
func (v *VertexBuffer) BufferFloat32(data []float32) {
	if v.id == 0 {
		gl.GenBuffers(1, &v.id)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, v.id)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
}

// BufferSubFloat32 buffers a float32 slice into a portion of the underlying
// buffer.
func (v *VertexBuffer) BufferSubFloat32(data []float32, offset int) {
	if v.id == 0 {
		gl.GenBuffers(1, &v.id)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, v.id)
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, len(data)*4, gl.Ptr(data))
}

// Bind binds the vertexbuffer.
func (v *VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, v.id)
}

// Unbind unbinds the vertexbuffer.
func (v *VertexBuffer) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// Draw renders the vertexbuffer.
func (v *VertexBuffer) Draw(mode uint32, first int32, count int32) {
	gl.DrawArrays(mode, first, count)
}

// DrawInstanced renders multiple instances of the vertexbuffer.
func (v *VertexBuffer) DrawInstanced(mode uint32, first int32, count int32, primcount int32) {
	gl.DrawArraysInstanced(mode, first, count, primcount)
}

// Destroy deallocates the vertexbuffer.
func (v *VertexBuffer) Destroy() {
	if v.id != 0 {
		gl.DeleteBuffers(1, &v.id)
		v.id = 0
	}
}
