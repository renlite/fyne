package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Shape)(nil)

// Rectangle describes a colored rectangle primitive in a Fyne canvas
type Shape struct {
	baseObject
}

// Hide will set this rectangle to not be visible
func (r *Shape) Hide() {
	r.baseObject.Hide()

	repaint(r)
}

// Move the rectangle to a new position, relative to its parent / canvas
func (r *Shape) Move(pos fyne.Position) {
	r.baseObject.Move(pos)

	repaint(r)
}

// Refresh causes this rectangle to be redrawn with its configured state.
func (r *Shape) Refresh() {
	Refresh(r)
}

// Resize on a rectangle updates the new size of this object.
// If it has a stroke width this will cause it to Refresh.
func (r *Shape) Resize(s fyne.Size) {
	if s == r.Size() {
		return
	}

	r.baseObject.Resize(s)
	Refresh(r)
}

// NewRectangle returns a new Rectangle instance
func NewShape(color color.Color) *Shape {
	return &Shape{}
}
