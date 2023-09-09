//go:build multi_shader

package gl

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func (p *painter) drawObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	switch obj := o.(type) {
	case *canvas.Circle:
		p.renderCanvasObject(o, pos, frame)
	case *canvas.Line:
		if obj.Position1.X == obj.Position2.X || obj.Position1.Y == obj.Position2.Y {
			p.renderCanvasObject(o, pos, frame)
			println("Line |-")
		} else {
			p.drawShapes(10.0, frame)
			p.drawLine(obj, pos, frame)
			println("Line /")
		}
	case *canvas.Text, *canvas.Image, *canvas.Raster:
		p.renderCanvasObject(o, pos, frame)
		// max textureBuffer GL=32, GLES=16
		if len(p.textures) == textureBuffer {
			p.drawShapes(10.0, frame)
			println(textureBuffer)
		}
	case *canvas.Rectangle:
		p.renderCanvasObject(o, pos, frame)
	case *canvas.LinearGradient:
		//* ?? TODO
		p.drawGradient(obj, p.newGlLinearGradientTexture, pos, frame)
	case *canvas.RadialGradient:
		//* ?? TODO
		p.drawGradient(obj, p.newGlRadialGradientTexture, pos, frame)
	case *canvas.Shape:
		p.drawShapes(10.0, frame)
		println("multi_shader")
	}
}
