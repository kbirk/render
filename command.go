package render

// Command represents a render command.
type Command struct {
	uniforms   map[string]interface{}
	textures   map[uint32]*Texture
	renderable *Renderable
}

// Uniform sets a uniform to be buffered.
func (c *Command) Uniform(name string, value interface{}) {
	if c.uniforms == nil {
		c.uniforms = make(map[string]interface{})
	}
	c.uniforms[name] = value
}

// Texture sets a texture to be bound.
func (c *Command) Texture(location uint32, texture *Texture) {
	if c.textures == nil {
		c.textures = make(map[uint32]*Texture)
	}
	c.textures[location] = texture
}

// Renderable sets a renderable to be drawn.
func (c *Command) Renderable(renderable *Renderable) {
	c.renderable = renderable
}

// Execute executes the render command.
func (c *Command) Execute(shader *Shader) {
	// bind textures
	for location, texture := range c.textures {
		texture.Bind(location)
	}
	// set uniforms
	for name, value := range c.uniforms {
		shader.SetUniform(name, value)
	}
	// draw
	c.renderable.Bind()
	c.renderable.Draw()
	c.renderable.Unbind()
}
