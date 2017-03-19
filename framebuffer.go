package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// FrameBuffer represents a framebuffer object
type FrameBuffer struct {
	id       uint32
	textures map[uint32]*Texture
}

// NewFrameBuffer instantiates and returns a new framebuffer instance.
func NewFrameBuffer() *FrameBuffer {
	var id uint32
	gl.GenFramebuffers(1, &id)
	return &FrameBuffer{
		id:       id,
		textures: make(map[uint32]*Texture),
	}
}

// Bind binds the framebuffer object.
func (f *FrameBuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
}

// Unbind unbinds the framebuffer object.
func (f *FrameBuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// BindForDraw binds the framebuffer object for drawing.
func (f *FrameBuffer) BindForDraw() {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, f.id)
}

// UnbindForDraw unbinds the framebuffer object for drawing.
func (f *FrameBuffer) UnbindForDraw() {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
}

// BindForRead binds the framebuffer object for reading.
func (f *FrameBuffer) BindForRead() {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, f.id)
}

// UnbindForRead unbinds the framebuffer object for reading.
func (f *FrameBuffer) UnbindForRead() {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
}

// SetDrawBuffers sets the draw buffers for the framebuffer object.
func (f *FrameBuffer) SetDrawBuffers(buffers []uint32) {
	gl.DrawBuffers(int32(len(buffers)), &buffers[0])
}

// AttachTexture attaches the provided texture to the provided attachment id.
func (f *FrameBuffer) AttachTexture(attachment uint32, texture *Texture) error {
	_, ok := f.textures[attachment]
	if ok {
		return fmt.Errorf("texture already attached to attachment `%d`",
			attachment)
	}
	f.Bind()
	gl.FramebufferTexture2D(
		gl.FRAMEBUFFER,
		attachment,
		gl.TEXTURE_2D,
		texture.ID(),
		0)
	err := f.checkAttachmentError()
	f.Unbind()
	if err == nil {
		f.textures[attachment] = texture
	}
	return err
}

// Texture returns the texture for the provided attachment id.
func (f *FrameBuffer) Texture(attachment uint32) (*Texture, bool) {
	tex, ok := f.textures[attachment]
	return tex, ok
}

// Resize will resize all attached textures.
func (f *FrameBuffer) Resize(width uint32, height uint32) {
	for _, texture := range f.textures {
		texture.Resize(width, height)
	}
}

// Destroy deallocates the framebuffer object.
func (f *FrameBuffer) Destroy() {
	gl.DeleteFramebuffers(1, &f.id)
	f.id = 0
}

func (f *FrameBuffer) checkAttachmentError() error {
	// check for errors
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status == gl.FRAMEBUFFER_COMPLETE {
		return nil
	}

	switch status {
	case gl.FRAMEBUFFER_UNDEFINED:
		return fmt.Errorf("target is the default framebuffer, but the " +
			"default framebuffer does not exist")

	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		return fmt.Errorf("any of the framebuffer attachment points are " +
			"framebuffer incomplete")

	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		return fmt.Errorf("the framebuffer does not have at least one image " +
			"attached to it.")

	case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
		return fmt.Errorf("the value of " +
			"gl.FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE is gl.NONE for any color " +
			"attachment point(s) named by gl.DRAW_BUFFERi.")

	case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
		return fmt.Errorf("gl.READ_BUFFER is not gl.NONE and the value of " +
			"gl.FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE is gl.NONE for the color " +
			"attachment point named by gl.READ_BUFFER.")

	case gl.FRAMEBUFFER_UNSUPPORTED:
		return fmt.Errorf("the combination of internal formats of the " +
			"attached images violates an implementation-dependent set of " +
			"restrictions.")

	case gl.FRAMEBUFFER_INCOMPLETE_MULTISAMPLE:
		return fmt.Errorf("the value of gl.RENDERBUFFER_SAMPLES is not the " +
			"same for all attached renderbuffers; if the value of " +
			"gl.TEXTURE_SAMPLES is the not same for all attached textures; " +
			"or, if the attached images are a mix of renderbuffers and " +
			"textures, the value of gl.RENDERBUFFER_SAMPLES does not match " +
			"the value of gl.TEXTURE_SAMPLES. Or the value of " +
			"gl.TEXTURE_FIXED_SAMPLE_LOCATIONS is not the same for all " +
			"attached textures; or, if the attached images are a mix of " +
			"renderbuffers and textures, the value of " +
			"gl.TEXTURE_FIXED_SAMPLE_LOCATIONS is not gl.TRUE for all " +
			"attached textures.")

	case gl.FRAMEBUFFER_INCOMPLETE_LAYER_TARGETS:
		return fmt.Errorf("any framebuffer attachment is layered, and any " +
			"populated attachment is not layered, or if all populated color " +
			"attachments are not from textures of the same target.")

	}
	return fmt.Errorf("unrecognized framebuffer error")
}
