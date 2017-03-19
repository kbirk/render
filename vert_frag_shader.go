package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// NewVertFragShader instantiates a new shader object.
func NewVertFragShader(vert, frag string) (*Shader, error) {
	// create shader
	shader := &Shader{}
	// create vertex shader
	vertex, err := shader.CreateShader(vert, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	// create fragment shader
	fragment, err := shader.CreateShader(frag, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	// attach shaders
	shader.AttachShader(vertex)
	shader.AttachShader(fragment)
	// link program
	err = shader.LinkProgram()
	if err != nil {
		return nil, err
	}
	return shader, nil
}
