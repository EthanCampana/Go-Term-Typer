package ui

import "github.com/gdamore/tcell/v2"

type cursor struct {
	xpos  int
	ypos  int
	style tcell.Style
}

type UI struct {
	sc tcell.Screen
	c  *cursor
}

//Clear the Screen
func (u *UI) Clear() {
	u.sc.Clear()
}

//Update the screen
func (u *UI) Show() {
	u.sc.Show()
}

func (u *UI) GetSize() (int, int) {
	x, y := u.sc.Size()
	return x, y

}

func (u *UI) drawContent(x, y int, text string, t tcell.Style, isCursor bool) {
	x2, y2 := u.GetSize()
	row := y
	col := x
	for _, r := range []rune(text) {
		if isCursor {
			u.sc.SetContent(col, row, r, nil, u.c.style)
		} else {
			u.sc.SetContent(col, row, r, nil, t)
		}
		col++
		if col >= x2 {
			row++
			col = x
		}
		if row > y2 {
			break
		}
	}
}
