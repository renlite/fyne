//go:build multi_shader

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
		p.renderCanvasObject(o, pos, frame)
	case *canvas.Line:
		if obj.Position1.X == obj.Position2.X || obj.Position1.Y == obj.Position2.Y {
			p.renderCanvasObject(o, pos, frame)
			println("Line |-")
		} else {
			p.drawShapes(frame)
			p.drawLine(obj, pos, frame)
			println("Line /")
		}
	case *canvas.Text, *canvas.Image, *canvas.Raster:
		p.renderCanvasObject(o, pos, frame)
		// max textureBuffer GL=32, GLES=16
		if len(p.textures) == textureBuffer {
			p.drawShapes(frame)
			//println("Texture: ", textureBuffer)
		}
	case *canvas.Rectangle:
		p.renderCanvasObject(o, pos, frame)
	case *canvas.LinearGradient:
		//* ?? TODO
		//p.drawGradient(obj, p.newGlLinearGradientTexture, pos, frame)
	case *canvas.RadialGradient:
		//* ?? TODO
		//p.drawGradient(obj, p.newGlRadialGradientTexture, pos, frame)
	case *canvas.Shape:
		p.drawShapes(frame)
	}
}

func (p *painter) renderCanvasObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	var shapeType float32
	switch obj := o.(type) {
	case *canvas.Rectangle:
		rect := obj //Re-Assignment for better readability
		if (rect.FillColor == color.Transparent || rect.FillColor == nil) && (rect.StrokeColor == color.Transparent || rect.StrokeColor == nil || rect.StrokeWidth == 0) {
			return
		}
		roundedCorners := rect.CornerRadius != 0
		if roundedCorners {
			shapeType = 2.0
		} else {
			shapeType = 1.0
		}
		//println("ShapeType: ", shapeType)

		rectPoints := p.roundRectCoords(pos, rect, frame)
		p.multiPoints = append(p.multiPoints, rectPoints...)
	//println(p.multiPoints)
	/*
		NOT USED
		for i := 0; i < 6; i++ {
			points[i*18+2] = shapeType
			points[i*18+3] = rect.StrokeWidth
			points[i*18+4] = rect.CornerRadius
			points[i*18+5] = 0.0 // NOT USED

		}
	*/

	case *canvas.Circle:
		// Conversion of Circle to Rectangle with CornerRaius = 0.5 * (Width or Height)
		circleDiameter := obj.Position1.X - obj.Position2.X
		rect := canvas.Rectangle{FillColor: obj.FillColor,
			StrokeColor:  obj.StrokeColor,
			StrokeWidth:  obj.StrokeWidth,
			CornerRadius: 0.5 * circleDiameter}
		rect.Move(obj.Position1)
		rect.Resize(fyne.Size{Width: circleDiameter, Height: circleDiameter})
		shapeType = 2.0 // Rectangle with CornerRadius
		rectPoints := p.roundRectCoords(pos, &rect, frame)
		p.multiPoints = append(p.multiPoints, rectPoints...)
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
		shapeType = 1.0 // Rectangle
		rectPoints := p.roundRectCoords(pos, &rect, frame)
		p.multiPoints = append(p.multiPoints, rectPoints...)
	case *canvas.Text:
		text := obj //Re-Assignment for better readability
		shapeType = 3.0
		//println("ShapeType: ", shapeType)
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
		texturePoints := p.renderTextureWithDetails(shapeType, text, p.newGlTextTexture, pos, size, frame, canvas.ImageFillStretch, 1.0, 0)
		p.multiPoints = append(p.multiPoints, texturePoints...)
		//p.textureIdx += 1
	case *canvas.Image:
		img := obj //Re-Assignment for better readability
		shapeType = 3.0
		imgPoints := p.renderTextureWithDetails(shapeType, img, p.newGlImageTexture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0)
		p.multiPoints = append(p.multiPoints, imgPoints...)
	case *canvas.Raster:
		img := obj //Re-Assignment for better readability
		shapeType = 3.0
		imgPoints := p.renderTextureWithDetails(shapeType, img, p.newGlRasterTexture, pos, img.Size(), frame, canvas.ImageFillStretch, float32(img.Alpha()), 0)
		p.multiPoints = append(p.multiPoints, imgPoints...)
	}
}

func (p *painter) renderTextureWithDetails(
	shapeType float32,
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
	return p.textureCoords(shapeType, size, pos, frame, fill, aspect, pad)

}

func (p *painter) drawShapes(frame fyne.Size) {
	println(len(p.multiPoints) / 18)
	if p.multiPoints == nil {
		println("return")
		return
	}
	//println("ShapeType: 10.0 (Draw shapes NOW!)")
	var program = p.multiProgram
	p.ctx.UseProgram(program)
	vbo := p.createBuffer(p.multiPoints)
	p.defineVertexArray(program, "att_vert", 2, 18, 0)
	p.defineVertexArray(program, "att_type", 4, 18, 2)
	p.defineVertexArray(program, "att_fill_color", 4, 18, 6)
	p.defineVertexArray(program, "att_stroke_color", 4, 18, 10)
	p.defineVertexArray(program, "att_rect_coords", 4, 18, 14)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()

	frameSizeUniform := p.ctx.GetUniformLocation(program, "frame_size")
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.ctx.Uniform2f(frameSizeUniform, frameWidthScaled, frameHeightScaled)
	p.logError()

	// Submit Textures if to draw
	/*
		sampler0 := p.ctx.GetUniformLocation(program, "texture0")
		p.ctx.Uniform1i(sampler0, 0)
		sampler1 := p.ctx.GetUniformLocation(program, "texture1")
		p.ctx.Uniform1i(sampler1, 1)
		sampler2 := p.ctx.GetUniformLocation(program, "texture2")
		p.ctx.Uniform1i(sampler2, 2)
		sampler3 := p.ctx.GetUniformLocation(program, "texture3")
		p.ctx.Uniform1i(sampler3, 3)
	*/
	p.logError()
	//println(len(p.textures))
	var values []int32
	for idx, texture := range p.textures {
		//println(idx)
		/*
			switch idx {
			case 0:
				p.ctx.ActiveTexture(texture0)
			case 1:
				p.ctx.ActiveTexture(texture1)
			case 2:
				p.ctx.ActiveTexture(texture2)
			case 3:
				p.ctx.ActiveTexture(texture3)
			}
		*/
		values = append(values, int32(idx))
		p.ctx.ActiveTexture(texture0 + uint32(idx))
		p.ctx.BindTexture(texture2D, texture)
		p.logError()
	}
	samplers := p.ctx.GetUniformLocation(program, "textures")
	p.ctx.Uniform1iv(samplers, values)

	//println(len(p.multiPoints) / 18)
	p.ctx.DrawArrays(triangles, 0, (len(p.multiPoints) / 18))
	p.logError()
	p.freeBuffer(vbo)
	p.multiPoints = nil
	p.textures = nil
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
	var shapeType float32
	if rect.CornerRadius == 0.0 {
		shapeType = 1.0
	} else {
		shapeType = 2.0
	}
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
		shapeType, strokeWidthScaled, radiusScaled, 0.0, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #2
		x2Pos, y1Pos,
		shapeType, strokeWidthScaled, radiusScaled, 0.0, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #3
		x1Pos, y2Pos,
		shapeType, strokeWidthScaled, radiusScaled, 0.0, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #4
		x1Pos, y2Pos,
		shapeType, strokeWidthScaled, radiusScaled, 0.0, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #5
		x2Pos, y1Pos,
		shapeType, strokeWidthScaled, radiusScaled, 0.0, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
		// #6
		x2Pos, y2Pos,
		shapeType, strokeWidthScaled, radiusScaled, 0.0, // att_type
		fillR, fillG, fillB, fillA, // att_fill_color;
		strokeR, strokeG, strokeB, strokeA, //att_stroke_color;
		x1Scaled, x2Scaled, y1Scaled, y2Scaled, // att_rect_coords;
	}

	return coords
}

func (p *painter) textureCoords(shapeType float32, size fyne.Size, pos fyne.Position, frame fyne.Size,
	fill canvas.ImageFill, aspect float32, pad float32) []float32 {
	size, pos = rectInnerCoords(size, pos, fill, aspect)
	size, pos = roundToPixelCoords(size, pos, p.pixScale)

	x1Pos := pos.X - pad
	x2Pos := pos.X + size.Width + pad
	y1Pos := pos.Y - pad
	y2Pos := pos.Y + size.Height + pad

	//println("in groupRectCoord")
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
		// (shape = 3.0      | texture x    | texture y | texIdx )
		shapeType, 0.0, 1.0, float32(len(p.textures) - 1), // att_type
		// NOT USED
		0.0, 0.0, 0.0, 0.0, // att_fill_color;
		0.0, 0.0, 0.0, 0.0, //att_stroke_color;
		0.0, 0.0, 0.0, 0.0, // att_rect_coords;

		// #2
		x1Pos, y1Pos, // att_vert
		// (shape = 3.0      | texture x    | texture y | texIdx )
		shapeType, 0.0, 0.0, float32(len(p.textures) - 1), // att_type
		// NOT USED
		0.0, 0.0, 0.0, 0.0, // att_fill_color;
		0.0, 0.0, 0.0, 0.0, //att_stroke_color;
		0.0, 0.0, 0.0, 0.0, // att_rect_coords;

		// #3
		x2Pos, y2Pos, // att_vert
		// (shape = 3.0      | texture x    | texture y | texIdx )
		shapeType, 1.0, 1.0, float32(len(p.textures) - 1), // att_type
		// NOT USED
		0.0, 0.0, 0.0, 0.0, // att_fill_color;
		0.0, 0.0, 0.0, 0.0, //att_stroke_color;
		0.0, 0.0, 0.0, 0.0, // att_rect_coords;

		// #4
		x2Pos, y1Pos, // att_vert
		// (shape = 3.0      | texture x    | texture y | texIdx )
		shapeType, 1.0, 0.0, float32(len(p.textures) - 1), // att_type
		// NOT USED
		0.0, 0.0, 0.0, 0.0, // att_fill_color;
		0.0, 0.0, 0.0, 0.0, //att_stroke_color;
		0.0, 0.0, 0.0, 0.0, // att_rect_coords;

		// #5
		x1Pos, y1Pos, // att_vert
		// (shape = 3.0      | texture x    | texture y | texIdx )
		shapeType, 0.0, 0.0, float32(len(p.textures) - 1), // att_type
		// NOT USED
		0.0, 0.0, 0.0, 0.0, // att_fill_color;
		0.0, 0.0, 0.0, 0.0, //att_stroke_color;
		0.0, 0.0, 0.0, 0.0, // att_rect_coords;

		// #6
		x2Pos, y2Pos, // att_vert
		// (shape = 3.0      | texture x    | texture y | texIdx )
		shapeType, 1.0, 1.0, float32(len(p.textures) - 1), // att_type
		// NOT USED
		0.0, 0.0, 0.0, 0.0, // att_fill_color;
		0.0, 0.0, 0.0, 0.0, //att_stroke_color;
		0.0, 0.0, 0.0, 0.0, // att_rect_coords;
	}
}
