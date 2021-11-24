package game

import (
	"termtyper/ui"

	"github.com/gdamore/tcell/v2"
)

type game struct {
	ui *ui.UI
}

func New() *game {
	c := &ui.cursor{x: 0, y: 0, style: tcell.StyleDefault}
	ui := &ui.UI{sc: tcell.Screen, c: c}
	g := &game{ui: ui}
	return g

}

func (g *game) Start() {
	//Do frame rate stufff
	//Do Styling Stuff
	for {
		g.drawMainMenu()
	}

}

func (g *game) drawMainMenu() {
	cursor := "->"
	menu1 := "Start"
	menu2 := "Options"
	menu3 := "Exit"
	uix, uiy := g.ui.GetSize()
	//Drawing Start 3 rows up From the Middle of the page
	g.ui.drawContent((uix/2)-(len(menu2)/2), (uiy/2)-3, menu1, nil, false)
	//Drawing Options to  the  middle of the UI
	g.ui.drawContent((uix/2)-(len(menu2)/2), uiy/2, menu2, nil, false)
	//Drawing Exit 3 rows down from the middle of the page
	g.ui.drawContent((uix/2)-(len(menu2)/2), (uiy/2)+3, menu1, nil, false)
	// If theere is no last location set the cursor to location of start
	if g.ui.c.xpos == 0 {
		g.ui.c.xpos = (uix / 2) - (len(menu2) / 2) - 4
		g.ui.c.xpos = (uiy / 2) - 3
	}
	//Draw last location of the cursor
	g.ui.drawContent(g.ui.c.xpos, g.ui.c.ypos, menu1, g.ui.c.style, false)
}
