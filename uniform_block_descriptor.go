package render

import (
	"fmt"
)

// UniformBlockDescriptor represents a shader uniform blocks attributes.
type UniformBlockDescriptor struct {
	Name      string
	Index     uint32
	Size      int32
	Offsets   map[string]int32
	Alignment int32
}

// AlignedSize returns the aligned block size of the uniform block.
func (u *UniformBlockDescriptor) AlignedSize() int32 {
	return u.Size + u.Alignment - (u.Size % u.Alignment)
}

// UnAlignedSize returns the unaligned block size of the uniform block.
func (u *UniformBlockDescriptor) UnAlignedSize() int32 {
	return u.Size
}

// Offset returns the offset of the uniform block.
func (u *UniformBlockDescriptor) Offset(name string) (int32, error) {
	offset, ok := u.Offsets[name]
	if !ok {
		return -1, fmt.Errorf("name `%s` not recognizied in the block", name)
	}
	return offset, nil
}
