package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"

	. "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("Box Layout")

	text1 := canvas.NewText("Column left", color.Black)
	text2 := canvas.NewText("Column right", color.Black)
	btn1 := widget.NewButton("click me", func() {
		log.Println("tapped")
	})
	btn1.Resize((fyne.Size{Width: 200.0}))
	btn2 := widget.NewButton("Home", func() {
		log.Println("tapped home")
	})
	btn3 := widget.NewButton("Done", func() {
		log.Println("tapped done")
	})
	btn4 := widget.NewButton("New", func() {
		log.Println("tapped ...")
	})
	btn5 := widget.NewButton("Save", func() {
		log.Println("tapped ...")
	})
	btn6 := widget.NewButton("Delete", func() {
		log.Println("tapped ...")
	})
	btn7 := widget.NewButton("Execute", func() {
		log.Println("tapped ...")
	})
	btn8 := widget.NewButton("Execute", func() {
		log.Println("tapped ...")
	})
	contentLeft := container.New(layout.NewVBoxLayout(), text1, btn1, btn2, btn3, btn4, btn5, btn6, btn7)

	shape := Shape{}
	contentRight := container.New(layout.NewVBoxLayout(), text2, btn8, &shape)
	gridCont := container.NewGridWithColumns(2, contentLeft, contentRight)

	myWindow.SetContent(gridCont)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()

}

/*
package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	. "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Renlite - Flex RoundRectangle Prototype")
	green := color.NRGBA{R: 0, G: 180, B: 0, A: 150}
	green_blue := color.NRGBA{R: 0, G: 180, B: 50, A: 150}
	orange := color.NRGBA{R: 255, G: 120, B: 0, A: 255}
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	//blue := color.NRGBA{R: 0, G: 0, B: 255, A: 100}
	purple := color.NRGBA{R: 150, G: 0, B: 205, A: 255}
	//blue_gray := color.NRGBA{R: 83, G: 140, B: 162, A: 150}
	//blue_gray1 := color.NRGBA{R: 134, G: 174, B: 189, A: 255}
	yellow := color.NRGBA{R: 255, G: 200, B: 0, A: 180}
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255.0}

	// RRect1
	rr1 := Rectangle{FillColor: red, StrokeColor: orange, StrokeWidth: 20.0}
	rr1.Resize((fyne.NewSize(300, 150)))
	rr1.Move(fyne.NewPos(10, 0))
	// RRect2
	rr2 := Rectangle{FillColor: purple, StrokeColor: yellow, StrokeWidth: 10.0, CornerRadius: 20}
	rr2.Resize((fyne.NewSize(300, 150)))
	rr2.Move(fyne.NewPos(10, 180))
	// RRect3
	rr3 := Rectangle{FillColor: white, StrokeColor: red, StrokeWidth: 5.0, CornerRadius: 20.0}
	rr3.Resize((fyne.NewSize(300, 150)))
	rr3.Move(fyne.NewPos(520, 0))
	// RRect4
	rr4 := Rectangle{FillColor: green, CornerRadius: 10.0}
	rr4.Resize((fyne.NewSize(150, 150)))
	rr4.Move(fyne.NewPos(360, 0))
	// RRect5
	rr5 := Rectangle{FillColor: purple, CornerRadius: 30.0}
	rr5.Resize((fyne.NewSize(300, 150)))
	rr5.Move(fyne.NewPos(400, 180))

	// >>BEG: composition
	// RRect6
	rr6 := Rectangle{FillColor: orange, CornerRadius: 10.0}
	rr6.Resize(fyne.NewSize(300, 150))
	rr6.Move(fyne.NewPos(360, 360))
	// RRect7
	rr7 := Rectangle{FillColor: yellow, CornerRadius: 10.0}
	rr7.Resize(fyne.NewSize(290, 140))
	rr7.Move(fyne.NewPos(365, 365))

	txtSeg1 := &widget.TextSegment{Text: "TEST"}
	txt1 := widget.NewRichText(txtSeg1)
	txt1.Move(fyne.NewPos(10, 10))
	// >>END: composition

	// RRect8
	rr8 := Rectangle{FillColor: yellow, StrokeColor: green_blue, StrokeWidth: 5.0}
	rr8.Resize((fyne.NewSize(200, 100)))
	rr8.Move(fyne.NewPos(255, 460))

	txtSeg2 := &widget.TextSegment{Text: "one GL stream to GPU"}
	txt2 := widget.NewRichText(txtSeg2)
	txt2.Move(fyne.NewPos(60, 60))

	// RRect9
	rr9 := Rectangle{FillColor: yellow, StrokeColor: orange, StrokeWidth: 5.0, CornerRadius: 8.0}
	rr9.Resize((fyne.NewSize(300, 150)))
	rr9.Move(fyne.NewPos(50, 360))

	btn1 := widget.NewButton("Button", nil)
	btn1.Resize(fyne.NewSize(80.0, 40.0))
	btn1.Move(fyne.NewPos(110, 110))

	egText := NewText("egText", color.White)
	//egText.Resize((fyne.NewSize(200, 50)))
	egText.Move(fyne.NewPos(160, 160))

	//shape := Shape{}

	cont := container.NewWithoutLayout(
		&rr1,
		&rr2,
		txt1,
		txt2,
		&rr3,
		btn1,
		egText,
		&rr7,
		//&shape,
	)
	myWindow.SetContent(cont)
	myWindow.Resize(fyne.NewSize(900, 600))

	myWindow.ShowAndRun()
}

*/
