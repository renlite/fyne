package main

import (
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-loremipsum/loremipsum"
)

type row struct {
	id                int
	name, description string
}

var rows []row

func init() {
	l := loremipsum.New()

	for i := 0; i < 200; i++ {
		name := l.Word() + " " + l.Word()
		desc := l.Word() + " " + l.Word() + " " + l.Word() + " " + l.Word() + " " + l.Word()
		r := row{i, name, desc}
		rows = append(rows, r)
	}
}

type dir int

const (
	sortOff dir = iota
	sortAsc
	sortDesc
)

var sorts = [3]dir{}

func main() {
	a := app.New()
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(800, 600))

	t := widget.NewTableWithHeaders(func() (int, int) {
		return len(rows), 3
	},
		func() fyne.CanvasObject {
			l := widget.NewLabel("John Smith")
			return l
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			l.Truncation = fyne.TextTruncateEllipsis
			switch id.Col {
			case 0:
				l.Truncation = fyne.TextTruncateOff
				l.SetText(strconv.Itoa(rows[id.Row].id))
			case 1:
				l.SetText(rows[id.Row].name)
			case 2:
				l.SetText(rows[id.Row].description)
			}
		})
	t.SetColumnWidth(0, 40)
	t.SetColumnWidth(1, 125)
	t.SetColumnWidth(2, 450)

	t.CreateHeader = func() fyne.CanvasObject {
		return widget.NewButton("000", func() {})
	}
	t.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			switch id.Col {
			case 0:
				b.SetText("ID")
				switch sorts[0] {
				case sortAsc:
					b.Icon = theme.MoveUpIcon()
				case sortDesc:
					b.Icon = theme.MoveDownIcon()
				default:
					b.Icon = nil
				}
			case 1:
				b.SetText("Name")
				switch sorts[1] {
				case sortAsc:
					b.Icon = theme.MoveUpIcon()
				case sortDesc:
					b.Icon = theme.MoveDownIcon()
				default:
					b.Icon = nil
				}
			case 2:
				b.SetText("Description")
				switch sorts[2] {
				case sortAsc:
					b.Icon = theme.MoveUpIcon()
				case sortDesc:
					b.Icon = theme.MoveDownIcon()
				default:
					b.Icon = nil
				}
			}
			b.Importance = widget.MediumImportance
			b.OnTapped = func() {
				applySort(id.Col, t)
			}
			b.Enable()
			b.Refresh()
		}
	}

	w.SetContent(t)
	w.ShowAndRun()
}

func applySort(col int, t *widget.Table) {
	order := sorts[col]
	order++
	if order > sortDesc {
		order = sortOff
	}
	// reset all and assign tapped sort
	for i := 0; i < 3; i++ {
		sorts[i] = sortOff
	}
	sorts[col] = order

	sort.Slice(rows, func(i, j int) bool {
		a := rows[i]
		b := rows[j]

		// re-sort with no sort selected
		if order == sortOff {
			return a.id < b.id
		}

		switch col {
		case 1:
			if order == sortAsc {
				return strings.Compare(a.name, b.name) < 0
			}
			return strings.Compare(a.name, b.name) > 0
		case 2:
			if order == sortAsc {
				return strings.Compare(a.description, b.description) < 0
			}
			return strings.Compare(a.description, b.description) > 0
		default:
			if order == sortDesc {
				return a.id > b.id
			}
			return a.id < b.id
		}
	})

	t.Refresh()
}
