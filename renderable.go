package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// AttributePointer represents a vertex attribute pointer.
type AttributePointer struct {
	Index      uint32
	Size       int32
	Type       uint32
	ByteStride int32
	ByteOffset int
}

// Renderable represents a renderable object.
type Renderable struct {
	id           uint32
	vertexbuffer *VertexBuffer
	indexbuffer  *IndexBuffer
	pointers     map[uint32]*AttributePointer
	instanced    map[uint32]bool
	// draw params
	mode       uint32
	count      int32
	first      int32
	typ        uint32
	byteOffset int
	primcount  int32
}

// SetVertexBuffer sets the vertexbuffer of the renderable.
func (r *Renderable) SetVertexBuffer(vb *VertexBuffer) {
	r.vertexbuffer = vb
}

// SetIndexBuffer sets the indexbuffer of the renderable.
func (r *Renderable) SetIndexBuffer(ib *IndexBuffer) {
	r.indexbuffer = ib
}

// SetPointer sets a vertex attribute pointer of the renderable.
func (r *Renderable) SetPointer(index uint32, pointer *AttributePointer) {
	if r.pointers == nil {
		r.pointers = make(map[uint32]*AttributePointer)
	}
	r.pointers[index] = pointer
}

// SetDrawArrays sets the params to render the underlying vertexbuffer.
func (r *Renderable) SetDrawArrays(mode uint32, first int32, count int32) {
	r.mode = mode
	r.first = first
	r.count = count
}

// SetDrawElements sets the instancing params to render the underlying
// vertexbuffer.
func (r *Renderable) SetDrawElements(mode uint32, count int32, typ uint32, byteOffset int) {
	r.mode = mode
	r.count = count
	r.typ = typ
	r.byteOffset = byteOffset
}

// SetDrawArraysInstanced sets the instancing params to render the underlying
// vertexbuffer.
func (r *Renderable) SetDrawArraysInstanced(mode uint32, first int32, count int32, primcount int32) {
	r.mode = mode
	r.first = first
	r.count = count
	r.primcount = primcount
}

// SetDrawElementsInstanced sets the params to render the underlying vertexbuffer.
func (r *Renderable) SetDrawElementsInstanced(mode uint32, count int32, typ uint32, byteOffset int, primcount int32) {
	r.mode = mode
	r.count = count
	r.typ = typ
	r.byteOffset = byteOffset
	r.primcount = primcount
}

// SetInstancedAttributes flags provided attributes for instancing.
func (r *Renderable) SetInstancedAttributes(instancedIndices []uint32) {
	if r.instanced == nil {
		r.instanced = make(map[uint32]bool)
	}
	for _, index := range instancedIndices {
		r.instanced[index] = true
	}
}

// Upload allocates the renderable to the GPU.
func (r *Renderable) Upload() {
	// create underlying vao
	gl.GenVertexArrays(1, &r.id)
	// bind
	gl.BindVertexArray(r.id)
	// bind vbo
	r.vertexbuffer.Bind()
	// set attribute pointers
	for index, pointer := range r.pointers {
		gl.EnableVertexAttribArray(index)
		gl.VertexAttribPointer(
			index,
			pointer.Size,
			pointer.Type,
			false,
			pointer.ByteStride,
			gl.PtrOffset(pointer.ByteOffset))
		// check if the attribute is instanced
		_, instanced := r.instanced[index]
		if instanced {
			gl.VertexAttribDivisor(index, 1)
		}
	}
	// bind EABO
	if r.indexbuffer != nil {
		r.indexbuffer.Bind()
	}
	// unbind
	gl.BindVertexArray(0)
}

// Bind binds the renderable.
func (r *Renderable) Bind() {
	gl.BindVertexArray(r.id)
}

// Unbind ubinds the renderable.
func (r *Renderable) Unbind() {
	gl.BindVertexArray(0)
}

// Draw renders the renderable.
func (r *Renderable) Draw() {
	if r.indexbuffer != nil {
		if r.primcount > 0 {
			r.indexbuffer.DrawInstanced(r.mode, r.count, r.typ, r.byteOffset, r.primcount)
		} else {
			r.indexbuffer.Draw(r.mode, r.count, r.typ, r.byteOffset)
		}
	} else {
		if r.primcount > 0 {
			r.vertexbuffer.DrawInstanced(r.mode, r.first, r.count, r.primcount)
		} else {
			r.vertexbuffer.Draw(r.mode, r.first, r.count)
		}
	}
}

// Destroy deallocates the renderable.
func (r *Renderable) Destroy() {
	if r.id != 0 {
		gl.DeleteVertexArrays(1, &r.id)
		r.id = 0
	}
}
