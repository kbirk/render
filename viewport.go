package render

// Viewport represents a viewport.
type Viewport struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// Equals returns true if the viewports are equal.
func (v *Viewport) Equals(other *Viewport) bool {
	return other != nil &&
		v.X == other.X &&
		v.Y == other.Y &&
		v.Width == other.Width &&
		v.Height == other.Height
}
