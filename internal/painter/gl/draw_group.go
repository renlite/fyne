//go:build !multi_shader && !single_shader

package gl

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	paint "fyne.io/fyne/v2/internal/painter"
)

func (p *painter) drawObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	switch obj := o.(type) {
	case *canvas.Circle:
		p.renderRoundRectangle(o, pos, frame)
		//println("Circle")
	case *canvas.Line:
		if obj.Position1.X == obj.Position2.X || obj.Position1.Y == obj.Position2.Y {
			p.renderRoundRectangle(o, pos, frame)
			println("Line |-")
		} else {
			//p.drawShapes(frame)
			//p.drawLine(obj, pos, frame)
			println("Line /")
		}
	case *canvas.Text, *canvas.Image, *canvas.Raster:
		p.renderTexture(o, pos, frame)
		// max textureBuffer GL=32, GLES=16
		//println("Texture")
	case *canvas.Rectangle:
		p.renderRoundRectangle(o, pos, frame)
		//println("Rectangle")
	case *canvas.LinearGradient:
		//* ?? TODO
		// p.drawGradient(obj, p.newGlLinearGradientTexture, pos, frame)
		//println("LinearGradient")
	case *canvas.RadialGradient:
		//* ?? TODO
		// p.drawGradient(obj, p.newGlRadialGradientTexture, pos, frame)
		//println("RadialGradient")
	case *canvas.Shape:
		p.drawShapes(frame)
	}
}

/*
func (p *painter) renderRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	if (rect.FillColor == color.Transparent || rect.FillColor == nil) && (rect.StrokeColor == color.Transparent || rect.StrokeColor == nil || rect.StrokeWidth == 0) {
		return
	}

}
*/

func (p *painter) renderRoundRectangle(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	switch obj := o.(type) {
	case *canvas.Rectangle:
		rect := obj //Re-Assignment for better readability
		if (rect.FillColor == color.Transparent || rect.FillColor == nil) && (rect.StrokeColor == color.Transparent || rect.StrokeColor == nil || rect.StrokeWidth == 0) {
			return
		}
		rectPoints := p.roundRectCoords(pos, rect, frame)
		p.groupLayers[0].roundRectPoints = append(p.groupLayers[0].roundRectPoints, rectPoints...)
	case *canvas.Circle:
		// Conversion of Circle to Rectangle with CornerRaius = 0.5 * (Width or Height)
		circleDiameter := obj.Position1.X - obj.Position2.X
		rect := canvas.Rectangle{FillColor: obj.FillColor,
			StrokeColor:  obj.StrokeColor,
			StrokeWidth:  obj.StrokeWidth,
			CornerRadius: 0.5 * circleDiameter}
		rect.Move(obj.Position1)
		rect.Resize(fyne.Size{Width: circleDiameter, Height: circleDiameter})
		rectPoints := p.roundRectCoords(pos, &rect, frame)
		p.groupLayers[0].roundRectPoints = append(p.groupLayers[0].roundRectPoints, rectPoints...)
	case *canvas.Line:
		if obj.StrokeColor == color.Transparent || obj.StrokeColor == nil || obj.StrokeWidth == 0 {
			return
		}
		// Conversion of horizontal/vertical Line to Rectangle
		var rectPosition fyne.Position
		var rectSize fyne.Size
		if obj.Position1.X == obj.Position2.X {
			// vertical line
			rectPosition.X = obj.Position1.X - obj.StrokeWidth*0.5
			rectPosition.Y = obj.Position1.Y
			rectSize.Width = obj.StrokeWidth
			rectSize.Height = obj.Position2.Y - obj.Position1.Y
			if rectSize.Height < 0.00 {
				rectSize.Height = rectSize.Height * (-1)
			}
		} else {
			// horizontal line
			rectPosition.X = obj.Position1.X
			rectPosition.Y = obj.Position1.Y - obj.StrokeWidth*0.5
			rectSize.Height = obj.StrokeWidth
			rectSize.Width = obj.Position2.X - obj.Position1.X
			if rectSize.Width < 0.00 {
				rectSize.Width = rectSize.Width * (-1)
			}
		}
		rect := canvas.Rectangle{FillColor: obj.StrokeColor}
		rect.Move(rectPosition)
		rect.Resize(rectSize)
		rectPoints := p.roundRectCoords(pos, &rect, frame)
		p.groupLayers[0].roundRectPoints = append(p.groupLayers[0].roundRectPoints, rectPoints...)
	}
}

func (p *painter) renderTexture(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	switch obj := o.(type) {
	case *canvas.Text:
		text := obj //Re-Assignment for better readability
		if text.Text == "" || text.Text == " " {
			return
		}
		size := text.MinSize()
		containerSize := text.Size()
		switch text.Alignment {
		case fyne.TextAlignTrailing:
			pos = fyne.NewPos(pos.X+containerSize.Width-size.Width, pos.Y)
		case fyne.TextAlignCenter:
			pos = fyne.NewPos(pos.X+(containerSize.Width-size.Width)/2, pos.Y)
		}

		if containerSize.Height > size.Height {
			pos = fyne.NewPos(pos.X, pos.Y+(containerSize.Height-size.Height)/2)
		}

		// text size is sensitive to position on screen
		size, _ = roundToPixelCoords(size, text.Position(), p.pixScale)
		size.Width += roundToPixel(paint.VectorPad(text), p.pixScale)
		texturePoints := p.renderTextureWithDetails(text, p.newGlTextTexture, pos, size, frame, canvas.ImageFillStretch, 1.0, 0)
		p.groupLayers[0].texturePoints = append(p.groupLayers[0].texturePoints, texturePoints...)

	case *canvas.Image:
		img := obj //Re-Assignment for better readability
		imgPoints := p.renderTextureWithDetails(img, p.newGlImageTexture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0)
		p.groupLayers[0].texturePoints = append(p.groupLayers[0].texturePoints, imgPoints...)

	case *canvas.Raster:
		img := obj //Re-Assignment for better readability
		imgPoints := p.renderTextureWithDetails(img, p.newGlRasterTexture, pos, img.Size(), frame, canvas.ImageFillStretch, float32(img.Alpha()), 0)
		p.groupLayers[0].texturePoints = append(p.groupLayers[0].texturePoints, imgPoints...)

	}
}

func (p *painter) renderTextureWithDetails(
	o fyne.CanvasObject,
	creator func(canvasObject fyne.CanvasObject) Texture,
	pos fyne.Position,
	size, frame fyne.Size,
	fill canvas.ImageFill,
	alpha float32,
	pad float32,
) []float32 {

	texture, err := p.getTexture(o, creator)
	if err != nil {
		return nil
	}
	p.textures = append(p.textures, texture)

	aspect := float32(0)
	if img, ok := o.(*canvas.Image); ok {
		aspect = img.Aspect()
		if aspect == 0 {
			aspect = 1 // fallback, should not occur - normally an image load error
		}
	}
	return p.textureCoords(size, pos, frame, fill, aspect, pad)

}

func (p *painter) roundRectCoords(pos fyne.Position, rect *canvas.Rectangle, frame fyne.Size) []float32 {
	size := rect.Size()
	pos1 := rect.Position()

	xPosDiff := pos.X - pos1.X
	yPosDiff := pos.Y - pos1.Y
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	size.Width = roundToPixel(size.Width, p.pixScale)
	size.Height = roundToPixel(size.Height, p.pixScale)

	x1Pos := pos1.X
	x2Pos := (pos1.X + size.Width)
	y1Pos := pos1.Y
	y2Pos := (pos1.Y + size.Height)

	// constants for the same rectangle
	strokeWidthScaled := roundToPixel(rect.StrokeWidth*p.pixScale, 1.0)
	radiusScaled := roundToPixel(rect.CornerRadius*p.pixScale, 1.0)
	fillR, fillG, fillB, fillA := getFragmentColor(rect.FillColor)
	strokeColor := rect.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	strokeR, strokeG, strokeB, strokeA := getFragmentColor(strokeColor)
	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(x1Pos, x2Pos, y1Pos, y2Pos)

	coords := []float32{
		// #1
		x1Pos, y1Pos, // att_vert
		strokeWidthScaled, radiusScaled, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #2
		x2Pos, y1Pos,
		strokeWidthScaled, radiusScaled, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #3
		x1Pos, y2Pos,
		strokeWidthScaled, radiusScaled, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #4
		x1Pos, y2Pos,
		strokeWidthScaled, radiusScaled, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #5
		x2Pos, y1Pos,
		strokeWidthScaled, radiusScaled, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #6
		x2Pos, y2Pos,
		strokeWidthScaled, radiusScaled, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
	}

	return coords
}

func (p *painter) textureCoords(size fyne.Size, pos fyne.Position, frame fyne.Size,
	fill canvas.ImageFill, aspect float32, pad float32) []float32 {
	size, pos = rectInnerCoords(size, pos, fill, aspect)
	size, pos = roundToPixelCoords(size, pos, p.pixScale)

	x1Pos := pos.X - pad
	x2Pos := pos.X + size.Width + pad
	y1Pos := pos.Y - pad
	y2Pos := pos.Y + size.Height + pad

	if p.textrureIdx == textureBuffer {
		p.textrureIdx = 0
	}
	p.textrureIdx = p.textrureIdx + 1
	//println(p.textrureIdx)
	return []float32{
		/* TEMPLATE OLD for triangleStrip
		// coord x, y, z texture x, y
		x1, y2, 0, 0.0, 1.0, // top left
		x1, y1, 0, 0.0, 0.0, // bottom left
		x2, y2, 0, 1.0, 1.0, // top right
		x2, y1, 0, 1.0, 0.0, // bottom right
		*/

		// #1
		x1Pos, y2Pos, // att_vert
		// (texture x | texture y | texIdx | UNUSED )
		0.0, 1.0, float32(p.textrureIdx - 1), 0.0, // att_type

		// #2
		x1Pos, y1Pos, // att_vert
		// (texture x | texture y | texIdx | UNUSED )
		0.0, 0.0, float32(p.textrureIdx - 1), 0.0, // att_type

		// #3
		x2Pos, y2Pos, // att_vert
		// (texture x | texture y | texIdx | UNUSED )
		1.0, 1.0, float32(p.textrureIdx - 1), 0.0, // att_type

		// #4
		x2Pos, y1Pos, // att_vert
		// (texture x | texture y | texIdx | UNUSED )
		1.0, 0.0, float32(p.textrureIdx - 1), 0.0, // att_type

		// #5
		x1Pos, y1Pos, // att_vert
		// (texture x | texture y | texIdx | UNUSED )
		0.0, 0.0, float32(p.textrureIdx - 1), 0.0, // att_type

		// #6
		x2Pos, y2Pos, // att_vert
		// (texture x | texture y | texIdx | UNUSED )
		1.0, 1.0, float32(p.textrureIdx - 1), 0.0, // att_type

	}
}

func (p *painter) drawShapes(frame fyne.Size) {
	var program Program
	var vbo Buffer
	var values []int32
	for idx := 0; idx < len(p.groupLayers); idx++ {
		if p.groupLayers[idx].roundRectPoints != nil {
			program = p.groupRoundRectProgram
			p.ctx.UseProgram(program)
			vbo = p.createBuffer(p.groupLayers[idx].roundRectPoints)
			p.defineVertexArray(program, "att_vert", 2, 16, 0)
			p.defineVertexArray(program, "att_type", 2, 16, 2)
			p.defineVertexArray(program, "att_fill_color", 4, 16, 4)
			p.defineVertexArray(program, "att_stroke_color", 4, 16, 8)
			p.defineVertexArray(program, "att_rect_coords", 4, 16, 12)

			p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
			p.logError()
			frameSizeUniform := p.ctx.GetUniformLocation(program, "frame_size")
			frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
			p.ctx.Uniform2f(frameSizeUniform, frameWidthScaled, frameHeightScaled)
			p.logError()

			p.ctx.DrawArrays(triangles, 0, (len(p.groupLayers[idx].roundRectPoints) / 16))

			p.ctx.DisableVertexAttribArray(p.ctx.GetAttribLocation(program, "att_fill_color"))
			p.ctx.DisableVertexAttribArray(p.ctx.GetAttribLocation(program, "att_stroke_color"))
			p.ctx.DisableVertexAttribArray(p.ctx.GetAttribLocation(program, "att_rect_coords"))

			p.logError()
			p.freeBuffer(vbo)
		}

		if p.groupLayers[idx].linePoints != nil {
			println("no lines")
		}

		if p.groupLayers[idx].texturePoints != nil {
			var texturePoints []float32
			var textures []Texture
			textureBufferCount := len(p.groupLayers[idx].texturePoints) / 6 / 6 / textureBuffer
			textureCount := len(p.groupLayers[idx].texturePoints) / 6 / 6
			lastTextureBuffer := textureCount - (textureBufferCount * textureBuffer)
			// BEG: OpenGL
			program = p.groupTextureProgram
			p.ctx.UseProgram(program)
			// END: OpenGL
			for i := 0; i <= (textureBufferCount); i++ {
				if i < textureBufferCount {
					texturePoints = p.groupLayers[idx].texturePoints[(i * (6 * 6 * textureBuffer)):((i + 1) * (6 * 6 * (textureBuffer)))]
					textures = p.textures[(i * textureBuffer):((i + 1) * (textureBuffer))]
				} else {
					texturePoints = p.groupLayers[idx].texturePoints[(i * (6 * 6 * textureBuffer)):((i * (6 * 6 * textureBuffer)) + (6 * 6 * lastTextureBuffer))]
					textures = p.textures[(i * textureBuffer):((i * textureBuffer) + lastTextureBuffer)]
				}
				//println(texturePoints[(len(texturePoints) - 2)])
				vbo = p.createBuffer(texturePoints)
				p.defineVertexArray(program, "att_vert", 2, 6, 0)
				p.defineVertexArray(program, "att_type", 4, 6, 2)

				p.ctx.BlendFunc(one, oneMinusSrcAlpha)
				p.logError()
				frameSizeUniform := p.ctx.GetUniformLocation(program, "frame_size")
				frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
				p.ctx.Uniform2f(frameSizeUniform, frameWidthScaled, frameHeightScaled)
				p.logError()

				for j, texture := range textures {
					values = append(values, int32(j))
					p.ctx.ActiveTexture(texture0 + uint32(j))
					p.ctx.BindTexture(texture2D, texture)
					p.logError()
				}
				samplers := p.ctx.GetUniformLocation(program, "textures")
				p.ctx.Uniform1iv(samplers, values)
				p.logError()
				p.ctx.DrawArrays(triangles, 0, (len(texturePoints) / 6))
				p.logError()
				p.freeBuffer(vbo)
				// delete
				texturePoints = nil
				textures = nil
				values = nil
			}
		}
		// delete
		p.textures = nil
		p.groupLayers[idx].roundRectPoints = nil
		p.groupLayers[idx].linePoints = nil
		p.groupLayers[idx].texturePoints = nil

	} //for
	p.textrureIdx = 0
}
