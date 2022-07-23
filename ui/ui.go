package ui

import (
	"termtyper/settings"

	"github.com/gdamore/tcell/v2"
)

type (
	keyListenerFn func(*Screen)
)

type Menu interface {
	Show(*Screen)
	Resize(*Screen)
}

type Sprite struct {
	Style   settings.DisplayStyle
	Char    rune
	Xpos    int
	Ypos    int
	Correct bool
}

type Screen struct {
	Window        tcell.Screen
	CurrentMenu   Menu
	MenuList      []Menu
	UpdateChannel chan *settings.Settings
	StartGame     chan bool
}

func StringArrayToSprite(s []string, settings *settings.Settings) [][]*Sprite {
	sprites := [][]*Sprite{}
	for _, k := range s {
		sprites = append(sprites, StringToSprite(k, settings))
	}
	return sprites
}

func CursorToSprite(input string, settings *settings.Settings) []*Sprite {
	out := []*Sprite{}
	for _, r := range []rune(input) {
		out = append(out, &Sprite{
			Style:   settings.CursorColor,
			Char:    r,
			Xpos:    0,
			Ypos:    0,
			Correct: false,
		})
	}
	return out
}

func DebugStringToSprite(input string) []*Sprite {
	out := []*Sprite{}
	for _, r := range []rune(input) {
		out = append(out, &Sprite{
			Style: settings.DisplayStyle(tcell.StyleDefault),
			Char:  r,
			Xpos:  0,
			Ypos:  0,
		})
	}
	return out
}

func StringToSprite(input string, settings *settings.Settings) []*Sprite {
	out := []*Sprite{}
	for _, r := range []rune(input) {
		out = append(out, &Sprite{
			Style: settings.TextColor,
			Char:  r,
			Xpos:  0,
			Ypos:  0,
		})
	}
	return out
}

func SpriteToString(c []*Sprite) string {
	s := []rune{}
	for _, sp := range c {
		s = append(s, sp.Char)
	}
	return string(s)
}

func NewScreen(options *settings.Settings) *Screen {
	Window, _ := tcell.NewScreen()
	Window.Init()
	Window.SetStyle(tcell.Style(options.TextColor))
	menus := []Menu{}
	mainMenu := NewMainMenu(options)
	settingsMenu := NewSettingsMenu(options)
	menus = append(menus, mainMenu)
	menus = append(menus, settingsMenu)
	UpdateChannel := make(chan *settings.Settings, 1)
	StartGame := make(chan bool, 1)
	return &Screen{Window: Window, CurrentMenu: mainMenu, MenuList: menus, UpdateChannel: UpdateChannel, StartGame: StartGame}
}

func (s *Screen) NewGame(options *settings.Settings) {
	x, y := s.Window.Size()
	gameContainer := NewContainer(x, y)
	GameMenu := NewGameMenu(gameContainer, options)
	s.CurrentMenu = GameMenu
	s.CurrentMenu.Resize(s)
}

func (s *Screen) DrawContentWithinContainer(c *Container, x, y int, content []*Sprite) {
	// Checks if the given word fits on the line... If it does it draws it
	if len(content) < c.xBound.end-x {
		for i, sp := range content {
			sp.Xpos = x + i
			sp.Ypos = y
			s.Window.SetContent(sp.Xpos, sp.Ypos, sp.Char, nil, tcell.Style(sp.Style))
			// These cursors are used to keep track how much data is written to a container and where to write data
			c.xCursor++
		}
		// Checks If there is a line below the current one and draws the
	} else if y+2 < c.yBound.end {
		c.yCursor += 2
		c.xCursor = c.xBound.begin
		for i, sp := range content {
			sp.Xpos = c.xBound.begin + i
			sp.Ypos = c.yCursor
			s.Window.SetContent(sp.Xpos, sp.Ypos, sp.Char, nil, tcell.Style(sp.Style))
			c.xCursor++
		}

	} else {
		c.isFull = true
	}

}

func (s *Screen) DrawContent(x, y int, content []*Sprite) {
	winX, winY := s.Window.Size()
	for i, sp := range content {
		sp.Xpos = x + i
		sp.Ypos = y
		s.Window.SetContent(sp.Xpos, sp.Ypos, sp.Char, nil, tcell.Style(sp.Style))
		if x+i > winX {
			y++
		}
		if y > winY {
			break
		}
	}

}
