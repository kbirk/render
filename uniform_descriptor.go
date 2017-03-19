package render

// UniformDescriptor represents a single shader uniforms attributes.
type UniformDescriptor struct {
	Name     string
	Type     uint32
	Count    int32
	Location int32
}
