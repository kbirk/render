package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	prevBlendFunc   *blendFunc
	prevCullFace    *cullFace
	prevDepthMask   *depthMask
	prevDepthFunc   *depthFunc
	prevViewport    *Viewport
	prevShader      *Shader
	prevFrameBuffer *FrameBuffer
	prevEnables     = make(map[uint32]bool)
)

type blendFunc struct {
	sfactor uint32
	dfactor uint32
}

func (b *blendFunc) Equals(other *blendFunc) bool {
	return other != nil &&
		b.sfactor == other.sfactor &&
		b.dfactor == other.dfactor
}

type cullFace struct {
	mode uint32
}

func (c *cullFace) Equals(other *cullFace) bool {
	return other != nil &&
		c.mode == other.mode
}

type depthMask struct {
	flag bool
}

func (d *depthMask) Equals(other *depthMask) bool {
	return other != nil &&
		d.flag == other.flag
}

type depthFunc struct {
	xfunc uint32
}

func (d *depthFunc) Equals(other *depthFunc) bool {
	return other != nil &&
		d.xfunc == other.xfunc
}

type clearColor struct {
	r float32
	g float32
	b float32
	a float32
}

// Technique represents a render technique.
type Technique struct {
	enables     []uint32
	shader      *Shader
	viewport    *Viewport
	framebuffer *FrameBuffer
	blendFunc   *blendFunc
	cullFace    *cullFace
	depthMask   *depthMask
	depthFunc   *depthFunc
	clearColor  *clearColor
}

// NewTechnique instantiates and returns a new technique instance.
func NewTechnique() *Technique {
	return &Technique{
		blendFunc: &blendFunc{
			sfactor: gl.ONE,
			dfactor: gl.ZERO,
		},
		cullFace: &cullFace{
			mode: gl.BACK,
		},
		depthMask: &depthMask{
			flag: true,
		},
		depthFunc: &depthFunc{
			xfunc: gl.LESS,
		},
	}
}

// Enable enables the rendering states for the technique.
func (t *Technique) Enable(enable uint32) {
	t.enables = append(t.enables, enable)
}

// Shader sets the shader for the technique.
func (t *Technique) Shader(shader *Shader) {
	t.shader = shader
}

// Viewport sets the viewport for the technique.
func (t *Technique) Viewport(viewport *Viewport) {
	t.viewport = viewport
}

// BlendFunc sets the blend func for the technique.
func (t *Technique) BlendFunc(sfactor uint32, dfactor uint32) {
	t.blendFunc = &blendFunc{
		sfactor: sfactor,
		dfactor: dfactor,
	}
}

// CullFace sets the cull face mode for the technique.
func (t *Technique) CullFace(mode uint32) {
	t.cullFace = &cullFace{
		mode: mode,
	}
}

// DepthMask sets the depth mask for the technique.
func (t *Technique) DepthMask(flag bool) {
	t.depthMask = &depthMask{
		flag: flag,
	}
}

// DepthFunc sets the depth mask for the technique.
func (t *Technique) DepthFunc(xfunc uint32) {
	t.depthFunc = &depthFunc{
		xfunc: xfunc,
	}
}

// ClearColor sets the clear color for the frame.
func (t *Technique) ClearColor(r, g, b, a float32) {
	t.clearColor = &clearColor{
		r: r,
		g: g,
		b: b,
		a: a,
	}
}

// Draw renders all commands using the technique.
func (t *Technique) Draw(commands []*Command) {
	t.setup()
	for _, command := range commands {
		command.Execute(t.shader)
	}
}

func (t *Technique) setup() {

	// bind framebuffer
	if t.framebuffer == nil && prevFrameBuffer != nil {
		prevFrameBuffer.Unbind()
	}
	if t.framebuffer != nil && t.framebuffer != prevFrameBuffer {
		t.framebuffer.Bind()
		prevFrameBuffer = t.framebuffer
	}

	// use shader
	if prevShader != t.shader {
		t.shader.Use()
		prevShader = t.shader
	}

	// track previous enables to determine which are stale
	staleEnables := make(map[uint32]bool)
	for state := range prevEnables {
		staleEnables[state] = true
	}

	// enable state
	for _, state := range t.enables {
		if !prevEnables[state] {
			gl.Enable(state)
			prevEnables[state] = true
		}
		delete(staleEnables, state)
	}

	// disable stale state
	for state := range staleEnables {
		gl.Disable(state)
		delete(prevEnables, state)
	}

	// update state functions
	if t.blendFunc != nil && !t.blendFunc.Equals(prevBlendFunc) {
		gl.BlendFunc(t.blendFunc.sfactor, t.blendFunc.dfactor)
		prevBlendFunc = t.blendFunc
	}
	if t.cullFace != nil && !t.cullFace.Equals(prevCullFace) {
		gl.CullFace(t.cullFace.mode)
		prevCullFace = t.cullFace
	}
	if t.depthMask != nil && !t.depthMask.Equals(prevDepthMask) {
		gl.DepthMask(t.depthMask.flag)
		prevDepthMask = t.depthMask
	}
	if t.depthFunc != nil && !t.depthFunc.Equals(prevDepthFunc) {
		gl.DepthFunc(t.depthFunc.xfunc)
		prevDepthFunc = t.depthFunc
	}

	// update viewport
	if t.viewport != nil && !t.viewport.Equals(prevViewport) {
		gl.Viewport(
			t.viewport.X,
			t.viewport.Y,
			t.viewport.Width,
			t.viewport.Height)
		prevViewport = t.viewport
	}
}
