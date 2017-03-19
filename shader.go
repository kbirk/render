package render

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	glslRegex *regexp.Regexp
)

func init() {
	glslRegex, _ = regexp.Compile(`void\s+main\s*\(\s*(void)*\s*\)\s*`)
}

func isGLSL(str string) bool {
	return glslRegex.MatchString(str)
}

// Shader represents a shader program.
type Shader struct {
	id               uint32
	shaders          []uint32
	descriptors      map[string]*UniformDescriptor
	blockDescriptors map[string]*UniformBlockDescriptor
}

// Use activates the shader.
func (s *Shader) Use() {
	gl.UseProgram(s.id)
}

// CreateShader creates an individual shader object.
func (s *Shader) CreateShader(source string, typ uint32) (uint32, error) {
	if !isGLSL(source) {
		// load shader file into memory
		raw, err := ioutil.ReadFile(source)
		if err != nil {
			return 0, err
		}
		source = string(raw)
	}
	// create shader object
	shader := gl.CreateShader(typ)
	// get c string
	cstr, free := gl.Strs(source + "\x00")
	// set source code of shader object
	gl.ShaderSource(shader, 1, cstr, nil)
	// free c string
	free()
	// compile shader
	gl.CompileShader(shader)
	// check error
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		// get info log length
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		// get error message
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		// delete current objects and abort constructor
		gl.DeleteShader(shader)
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	// return shader object
	return shader, nil
}

// AttachShader attaches a shader object to the program.
func (s *Shader) AttachShader(shader uint32) {
	if s.id == 0 {
		s.id = gl.CreateProgram()
	}
	if s.shaders == nil {
		s.shaders = make([]uint32, 0)
		s.shaders = append(s.shaders, shader)
	}
	gl.AttachShader(s.id, shader)
}

// LinkProgram links the shader program.
func (s *Shader) LinkProgram() error {
	// link shader program
	gl.LinkProgram(s.id)
	// error check
	var status int32
	gl.GetProgramiv(s.id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.id, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.id, logLength, nil, gl.Str(log))
		// delete shader objects
		s.deleteShaders()
		return fmt.Errorf("failed to link program: %v", log)
	}
	// delete shader objects
	s.deleteShaders()
	// query uniform information
	s.queryUniforms()
	return nil
}

// SetUniform1i buffers a int32 by value.
func (s *Shader) SetUniform1i(location int32, arg interface{}) error {
	value, ok := arg.(int32)
	if !ok {
		return fmt.Errorf("%v is not of type int32", arg)
	}
	gl.Uniform1i(location, value)
	return nil
}

// SetUniform1ui buffers an uint32 by value.
func (s *Shader) SetUniform1ui(location int32, arg interface{}) error {
	value, ok := arg.(uint32)
	if !ok {
		return fmt.Errorf("%v is not of type uint32", arg)
	}
	gl.Uniform1ui(location, value)
	return nil
}

// SetUniform1f buffers a float32 by value.
func (s *Shader) SetUniform1f(location int32, arg interface{}) error {
	value, ok := arg.(float32)
	if !ok {
		return fmt.Errorf("%v is not of type float32", arg)
	}
	gl.Uniform1f(location, value)
	return nil
}

// SetUniform1iv buffers one or more int32 by address.
func (s *Shader) SetUniform1iv(location int32, count int32, arg interface{}) error {
	value, ok := arg.(*int32)
	if !ok {
		return fmt.Errorf("%v is not of type *int32", arg)
	}
	gl.Uniform1iv(location, count, value)
	return nil
}

// SetUniform1uiv buffers  one or more uint32 by address.
func (s *Shader) SetUniform1uiv(location int32, count int32, arg interface{}) error {
	value, ok := arg.(*uint32)
	if !ok {
		return fmt.Errorf("%v is not of type *uint32", arg)
	}
	gl.Uniform1uiv(location, count, value)
	return nil
}

// SetUniform1fv buffers one or more float32 by address.
func (s *Shader) SetUniform1fv(location int32, count int32, arg interface{}) error {
	value, ok := arg.(*float32)
	if !ok {
		return fmt.Errorf("%v is not of type *float32", arg)
	}
	gl.Uniform1fv(location, count, value)
	return nil
}

// SetUniform2fv buffers one or more 2-component float32 by address.
func (s *Shader) SetUniform2fv(location int32, count int32, arg interface{}) error {
	value, ok := arg.(*float32)
	if !ok {
		return fmt.Errorf("%v is not of type *float32", arg)
	}
	gl.Uniform2fv(location, count, value)
	return nil
}

// SetUniform3fv buffers one or more 3-component float32 by address.
func (s *Shader) SetUniform3fv(location int32, count int32, arg interface{}) error {
	value, ok := arg.(*float32)
	if !ok {
		return fmt.Errorf("%v is not of type *float32", arg)
	}
	gl.Uniform3fv(location, count, value)
	return nil
}

// SetUniform4fv buffers one or more 4-component float32 by address.
func (s *Shader) SetUniform4fv(location int32, count int32, arg interface{}) error {
	value, ok := arg.(*float32)
	if !ok {
		return fmt.Errorf("%v is not of type *float32", arg)
	}
	gl.Uniform4fv(location, count, value)
	return nil
}

// SetUniformMatrix3fv buffers one or more 9-component float32 by address.
func (s *Shader) SetUniformMatrix3fv(location int32, count int32, arg interface{}) error {
	value, ok := arg.(*float32)
	if !ok {
		return fmt.Errorf("%v is not of type *float32", arg)
	}
	gl.UniformMatrix3fv(location, count, false, value)
	return nil
}

// SetUniformMatrix4fv buffers one or more 16-component float32 by address.
func (s *Shader) SetUniformMatrix4fv(location int32, count int32, arg interface{}) error {
	value, ok := arg.(*float32)
	if !ok {
		return fmt.Errorf("%v is not of type *float32", arg)
	}
	gl.UniformMatrix4fv(location, count, false, value)
	return nil
}

// SetUniform buffers one or more uniforms.
func (s *Shader) SetUniform(name string, arg interface{}) error {
	// check descriptors
	descriptor, ok := s.descriptors[name]
	if !ok {
		return fmt.Errorf("uniform `%s` was not recognized", name)
	}
	// buffer uniform data
	switch descriptor.Type {
	case gl.SAMPLER_2D:
		return s.SetUniform1i(descriptor.Location, arg)
	case gl.SAMPLER_CUBE:
		return s.SetUniform1i(descriptor.Location, arg)
	case gl.INT:
		if descriptor.Count > 1 {
			return s.SetUniform1iv(descriptor.Location, descriptor.Count, arg)
		}
		return s.SetUniform1i(descriptor.Location, arg)
	case gl.UNSIGNED_INT:
		if descriptor.Count > 1 {
			return s.SetUniform1uiv(descriptor.Location, descriptor.Count, arg)
		}
		return s.SetUniform1ui(descriptor.Location, arg)
	case gl.FLOAT:
		if descriptor.Count > 1 {
			return s.SetUniform1fv(descriptor.Location, descriptor.Count, arg)
		}
		return s.SetUniform1f(descriptor.Location, arg)
	case gl.FLOAT_VEC2:
		return s.SetUniform2fv(descriptor.Location, descriptor.Count, arg)
	case gl.FLOAT_VEC3:
		return s.SetUniform3fv(descriptor.Location, descriptor.Count, arg)
	case gl.FLOAT_VEC4:
		return s.SetUniform4fv(descriptor.Location, descriptor.Count, arg)
	case gl.FLOAT_MAT3:
		return s.SetUniformMatrix3fv(descriptor.Location, descriptor.Count, arg)
	case gl.FLOAT_MAT4:
		return s.SetUniformMatrix4fv(descriptor.Location, descriptor.Count, arg)
	}
	return nil
}

// Destroy deallocates the shader program.
func (s *Shader) Destroy() {
	if s.id != 0 {
		gl.DeleteProgram(s.id)
		s.id = 0
	}
}

func (s *Shader) deleteShaders() {
	if s.shaders != nil {
		for _, shader := range s.shaders {
			gl.DeleteShader(shader)
		}
		s.shaders = nil
	}
}

func (s *Shader) queryUniforms() {
	// query all necessary uniform information
	uniformIndices := s.queryUniformIndices()
	uniformNames := s.queryUniformNames(uniformIndices)
	uniformTypes := s.queryUniformTypes(uniformIndices)
	uniformCounts := s.queryUniformCounts(uniformIndices)
	parentBlockIndices := s.queryParentBlockIndices(uniformIndices)
	uniformOffsets := s.queryUniformOffsets(uniformIndices)
	uniformLocations := s.queryUniformLocations(uniformNames)

	// create descriptor maps
	s.descriptors = make(map[string]*UniformDescriptor)
	s.blockDescriptors = make(map[string]*UniformBlockDescriptor)

	// for each uniform index
	for _, index := range uniformIndices {

		// check if part of a block or not
		if parentBlockIndices[index] == -1 {
			// not in a block, add as uniform descriptor
			name := uniformNames[index]
			typ := uniformTypes[index]
			count := uniformCounts[index]
			location := uniformLocations[index]

			s.descriptors[name] = &UniformDescriptor{
				Name:     name[:len(name)-1],
				Type:     typ,
				Count:    count,
				Location: location,
			}
		}
	}

	// query all necessary uniform block information
	blockIndices := s.queryUniformBlockIndices()
	blockNames := s.queryUniformBlockNames(blockIndices)
	blockSizes := s.queryUniformBlockSizes(blockIndices)
	bufferAlignment := s.queryUniformBufferAlignment()

	// for each block
	for _, index := range blockIndices {

		// get all uniform offsets that are part of this block
		offsets := make(map[string]int32)
		for i, parentIndex := range parentBlockIndices {
			if parentIndex == int32(index) {
				// uniform is part of this block
				offsets[uniformNames[i]] = uniformOffsets[i]
			}
		}

		blockName := blockNames[index]
		blockIndex := blockIndices[index]
		blockSize := blockSizes[index]

		// add block descriptor
		s.blockDescriptors[blockName] = &UniformBlockDescriptor{
			Name:      blockName,
			Index:     blockIndex,
			Size:      blockSize,
			Offsets:   offsets,
			Alignment: bufferAlignment,
		}

		// set binding point for block index and shader
		// TODO: make this configurable
		gl.UniformBlockBinding(s.id, blockIndex, blockIndex)
	}
}

// UniformDescriptors returns the map of uniform descriptors.
func (s *Shader) UniformDescriptors() map[string]*UniformDescriptor {
	return s.descriptors
}

// UniformBlockDescriptors returns the map of uniform block descriptors.
func (s *Shader) UniformBlockDescriptors() map[string]*UniformBlockDescriptor {
	return s.blockDescriptors
}

func toString(buff []uint8) string {
	b := make([]byte, len(buff))
	for i, v := range buff {
		b[i] = byte(v)
	}
	return string(b[:len(b)-1]) // trim null terminator
}

func toUint32(arr []int32) []uint32 {
	res := make([]uint32, len(arr))
	for i, val := range arr {
		res[i] = uint32(val)
	}
	return res
}

func (s *Shader) queryUniformIndices() []uint32 {
	// get the number of uniforms
	var numActiveUniforms int32
	gl.GetProgramiv(s.id, gl.ACTIVE_UNIFORMS, &numActiveUniforms)
	// get uniform indices from 0 to gl.ACTIVE_UNIFORMS
	indices := make([]uint32, numActiveUniforms)
	for i := int32(0); i < numActiveUniforms; i++ {
		indices[i] = uint32(i)
	}
	return indices
}

func (s *Shader) queryUniformNames(indices []uint32) []string {
	// check if no uniforms
	if len(indices) == 0 {
		return make([]string, 0)
	}
	// get uniform name lengths
	nameLengths := make([]int32, len(indices))
	gl.GetActiveUniformsiv(s.id, int32(len(indices)), &indices[0], gl.UNIFORM_NAME_LENGTH, &nameLengths[0])
	// for each uniform index
	names := make([]string, len(indices))
	for _, index := range indices {
		// get uniform name
		nameLength := nameLengths[index]
		// create name slice
		name := make([]uint8, nameLength)
		// get name bytes
		gl.GetActiveUniformName(s.id, index, nameLength, nil, &name[0])
		// cast from uint8 to string
		names[index] = toString(name)
	}
	return names
}

func (s *Shader) queryUniformTypes(indices []uint32) []uint32 {
	// check if no uniforms
	if len(indices) == 0 {
		return make([]uint32, 0)
	}
	// get uniform types
	types := make([]int32, len(indices))
	gl.GetActiveUniformsiv(s.id, int32(len(indices)), &indices[0], gl.UNIFORM_TYPE, &types[0])
	// cast to uint32
	return toUint32(types)
}

func (s *Shader) queryUniformCounts(indices []uint32) []int32 {
	// check if no uniforms
	if len(indices) == 0 {
		return make([]int32, 0)
	}
	// get uniform types
	sizes := make([]int32, len(indices))
	gl.GetActiveUniformsiv(s.id, int32(len(indices)), &indices[0], gl.UNIFORM_SIZE, &sizes[0])
	return sizes
}

func (s *Shader) queryParentBlockIndices(indices []uint32) []int32 {
	// check if no uniforms
	if len(indices) == 0 {
		return make([]int32, 0)
	}
	// get uniform block indices (-1 is not part of a block)
	blockIndices := make([]int32, len(indices))
	gl.GetActiveUniformsiv(s.id, int32(len(indices)), &indices[0], gl.UNIFORM_BLOCK_INDEX, &blockIndices[0])
	return blockIndices
}

func (s *Shader) queryUniformOffsets(indices []uint32) []int32 {
	// check if no uniforms
	if len(indices) == 0 {
		return make([]int32, 0)
	}
	// get uniform offsets
	offsets := make([]int32, len(indices))
	gl.GetActiveUniformsiv(s.id, int32(len(indices)), &indices[0], gl.UNIFORM_OFFSET, &offsets[0])
	return offsets
}

func (s *Shader) queryUniformLocations(names []string) []int32 {
	locations := make([]int32, len(names))
	for i, name := range names {
		locations[i] = gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	}
	return locations
}

func (s *Shader) queryUniformBlockIndices() []uint32 {
	// get the number of active uniform blocks
	var numActiveBlocks int32
	gl.GetProgramiv(s.id, gl.ACTIVE_UNIFORM_BLOCKS, &numActiveBlocks)
	// get uniform indices from 0 to gl.ACTIVE_UNIFORMS
	indices := make([]uint32, numActiveBlocks)
	for i := int32(0); i < numActiveBlocks; i++ {
		indices[i] = uint32(i)
	}
	return indices
}

func (s *Shader) queryUniformBlockNames(indices []uint32) []string {
	names := make([]string, len(indices))
	for _, index := range indices {
		// get the length of the name
		var nameLength int32
		gl.GetActiveUniformBlockiv(s.id, index, gl.UNIFORM_BLOCK_NAME_LENGTH, &nameLength)
		// get the block name
		name := make([]uint8, nameLength)
		// get name bytes
		gl.GetActiveUniformBlockName(s.id, index, nameLength, nil, &name[0])
		// cast from uint8 to string
		names[index] = toString(name)
	}
	return names
}

func (s *Shader) queryUniformBlockSizes(indices []uint32) []int32 {
	sizes := make([]int32, len(indices))
	for _, index := range indices {
		var blockSize int32
		gl.GetActiveUniformBlockiv(s.id, index, gl.UNIFORM_BLOCK_DATA_SIZE, &blockSize)
		sizes[index] = blockSize
	}
	return sizes
}

func (s *Shader) queryUniformBufferAlignment() int32 {
	var uniformBufferAlignment int32
	gl.GetIntegerv(gl.UNIFORM_BUFFER_OFFSET_ALIGNMENT, &uniformBufferAlignment)
	return uniformBufferAlignment
}
