//go:build single_shader

package gl

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func (p *painter) drawObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	switch obj := o.(type) {
	case *canvas.Circle:
		p.drawCircle(obj, pos, frame)
	case *canvas.Line:
		p.drawLine(obj, pos, frame)
	case *canvas.Image:
		p.drawImage(obj, pos, frame)
	case *canvas.Raster:
		p.drawRaster(obj, pos, frame)
	case *canvas.Rectangle:
		p.drawRectangle(obj, pos, frame)
	case *canvas.Text:
		p.drawText(obj, pos, frame)
	case *canvas.LinearGradient:
		p.drawGradient(obj, p.newGlLinearGradientTexture, pos, frame)
	case *canvas.RadialGradient:
		p.drawGradient(obj, p.newGlRadialGradientTexture, pos, frame)
	}
}

func (p *painter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	if (rect.FillColor == color.Transparent || rect.FillColor == nil) && (rect.StrokeColor == color.Transparent || rect.StrokeColor == nil || rect.StrokeWidth == 0) {
		return
	}

	roundedCorners := rect.CornerRadius != 0
	var program Program
	if roundedCorners {
		program = p.roundRectangleProgram
	} else {
		program = p.rectangleProgram
	}

	// Vertex: BEG
	bounds, points := p.vecRectCoords(pos, rect, frame)
	p.ctx.UseProgram(program)
	vbo := p.createBuffer(points)
	p.defineVertexArray(program, "vert", 2, 4, 0)
	p.defineVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameSizeUniform := p.ctx.GetUniformLocation(program, "frame_size")
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.ctx.Uniform2f(frameSizeUniform, frameWidthScaled, frameHeightScaled)

	rectCoordsUniform := p.ctx.GetUniformLocation(program, "rect_coords")
	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.ctx.Uniform4f(rectCoordsUniform, x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	strokeWidthScaled := roundToPixel(rect.StrokeWidth*p.pixScale, 1.0)
	if roundedCorners {
		strokeUniform := p.ctx.GetUniformLocation(program, "stroke_width_half")
		p.ctx.Uniform1f(strokeUniform, strokeWidthScaled*0.5)

		rectSizeUniform := p.ctx.GetUniformLocation(program, "rect_size_half")
		rectSizeWidthScaled := x2Scaled - x1Scaled - strokeWidthScaled
		rectSizeHeightScaled := y2Scaled - y1Scaled - strokeWidthScaled
		p.ctx.Uniform2f(rectSizeUniform, rectSizeWidthScaled*0.5, rectSizeHeightScaled*0.5)

		radiusUniform := p.ctx.GetUniformLocation(program, "radius")
		radiusScaled := roundToPixel(rect.CornerRadius*p.pixScale, 1.0)
		p.ctx.Uniform1f(radiusUniform, radiusScaled)
	} else {
		strokeUniform := p.ctx.GetUniformLocation(program, "stroke_width")
		p.ctx.Uniform1f(strokeUniform, strokeWidthScaled)
	}

	var r, g, b, a float32
	fillColorUniform := p.ctx.GetUniformLocation(program, "fill_color")
	r, g, b, a = getFragmentColor(rect.FillColor)
	p.ctx.Uniform4f(fillColorUniform, r, g, b, a)

	strokeColorUniform := p.ctx.GetUniformLocation(program, "stroke_color")
	strokeColor := rect.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.ctx.Uniform4f(strokeColorUniform, r, g, b, a)
	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
	p.freeBuffer(vbo)
}

func (p *painter) vecRectCoords(pos fyne.Position, rect *canvas.Rectangle, frame fyne.Size) ([4]float32, []float32) {
	size := rect.Size()
	pos1 := rect.Position()

	//println(pos.X, "<>", pos1.X)
	xPosDiff := pos.X - pos1.X
	yPosDiff := pos.Y - pos1.Y
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	size.Width = roundToPixel(size.Width, p.pixScale)
	size.Height = roundToPixel(size.Height, p.pixScale)

	x1Pos := pos1.X
	x1Norm := -1 + x1Pos*2/frame.Width
	x2Pos := pos1.X + size.Width
	x2Norm := -1 + x2Pos*2/frame.Width
	y1Pos := pos1.Y
	y1Norm := 1 - y1Pos*2/frame.Height
	y2Pos := pos1.Y + size.Height
	y2Norm := 1 - y2Pos*2/frame.Height

	// output a norm for the fill and the vert is unused, but we pass 0 to avoid optimisation issues
	coords := []float32{
		0, 0, x1Norm, y1Norm, // first triangle
		0, 0, x2Norm, y1Norm, // second triangle
		0, 0, x1Norm, y2Norm,
		0, 0, x2Norm, y2Norm}

	return [4]float32{x1Pos, y1Pos, x2Pos, y2Pos}, coords
}

// for compatibility reasons im painter.go "FinishDrawing()"
func (p *painter) drawShapes(frame fyne.Size) {
}
