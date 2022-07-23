package ui

import (
	"os"
	"termtyper/settings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type MainMenu struct {
	KeyPressListener map[tcell.Key]keyListenerFn
	MenuOptions      [3][]*Sprite
	Title            [3][]*Sprite
	Cursor           []*Sprite
	ClearCursor      []*Sprite
	Credits          []*Sprite
	CurrentOption    int8
	animActive       bool
	animStopped      bool
}

func NewMainMenu(settings *settings.Settings) *MainMenu {
	options := [3][]*Sprite{}
	title := [3][]*Sprite{}
	title[0] = StringToSprite("▀▀█▀▀ █▀▀ █▀▀█ █▀▄▀█  ▀▀█▀▀ █  █ █▀▀█ █▀▀ █▀▀█ ", settings)
	title[1] = StringToSprite("  █   █▀▀ █▄▄▀ █ ▀ █    █   █▄▄█ █  █ █▀▀ █▄▄▀ ", settings)
	title[2] = StringToSprite("  ▀   ▀▀▀ ▀ ▀▀ ▀   ▀    ▀   ▄▄▄█ █▀▀▀ ▀▀▀ ▀ ▀▀ ", settings)
	Credits := StringToSprite("Created By: Ethan Campana", settings)

	options[0] = StringToSprite("Start Game", settings)
	options[1] = StringToSprite("Options", settings)
	options[2] = StringToSprite("Exit", settings)
	cursor := CursorToSprite("->", settings)
	clearCursor := StringToSprite("    ", settings)
	mm := &MainMenu{
		Cursor:        cursor,
		Title:         title,
		Credits:       Credits,
		ClearCursor:   clearCursor,
		MenuOptions:   options,
		CurrentOption: 0,
		animActive:    false,
		animStopped:   false,
	}
	mm.KeyPressListener = make(map[tcell.Key]keyListenerFn)
	mm.KeyPressListener[tcell.KeyESC] = mm.keyEscListener
	mm.KeyPressListener[tcell.KeyUp] = mm.keyUpListener
	mm.KeyPressListener[tcell.KeyDown] = mm.keyDownListener
	mm.KeyPressListener[tcell.KeyEnter] = mm.keyEnterListener
	return mm

}

func (mm *MainMenu) AnimStart(s *Screen) {
	framerate := time.NewTicker(1 * time.Second / 5)
	go func(s *Screen) {
		i := 0
		for range framerate.C {
			if !mm.animActive && mm.animStopped {
				break
			}
			if i == 0 {
				s.DrawContent(mm.Cursor[0].Xpos, mm.Cursor[0].Ypos, mm.ClearCursor)
				mm.Cursor[0].Xpos -= 3
				s.DrawContent(mm.Cursor[0].Xpos, mm.Cursor[0].Ypos, mm.Cursor)
				i++
			} else {
				s.DrawContent(mm.Cursor[0].Xpos, mm.Cursor[0].Ypos, mm.ClearCursor)
				mm.Cursor[0].Xpos += 3
				s.DrawContent(mm.Cursor[0].Xpos, mm.Cursor[0].Ypos, mm.Cursor)
				i--

			}
		}
		return
	}(s)
}

func (mm *MainMenu) Resize(s *Screen) {
	s.Window.Sync()
	s.Window.Clear()
	winX, winY := s.Window.Size()
	offset := (3 * mm.CurrentOption) - 3
	s.DrawContent(((winX / 2) - (len(mm.MenuOptions[mm.CurrentOption]) / 2) - 4), winY/2+int(offset), mm.Cursor)
}

func (mm *MainMenu) Show(s *Screen) {
	winX, winY := s.Window.Size()

	//Checking if the cursor was just initialized
	if mm.Cursor[0].Xpos == 0 {
		//Move the Cursor to the correct Location
		s.DrawContent(((winX / 2) - (len(mm.MenuOptions[1]) / 2) - 4), winY/2-3, mm.Cursor)
	}

	if !mm.animActive && !mm.animStopped {
		mm.AnimStart(s)
		mm.animActive = true
	}

	for i := range mm.MenuOptions {
		//Draws the Menu Options
		offset := (3 * i) - 3
		s.DrawContent((winX/2)-(len(mm.MenuOptions[1])/2), winY/2+offset, mm.MenuOptions[i])
	}
	for i := range mm.Title {
		s.DrawContent((winX/2)-(len(mm.Title[0])/2), mm.MenuOptions[0][0].Ypos-7+i, mm.Title[i])
	}
	s.DrawContent((winX/2)+5, mm.Title[2][0].Ypos+2, mm.Credits)

}

func (mm *MainMenu) keyUpListener(s *Screen) {

	s.DrawContent(mm.Cursor[0].Xpos, mm.Cursor[0].Ypos, mm.ClearCursor)
	if mm.CurrentOption == 0 {
		mm.CurrentOption = 2
		mm.Cursor[0].Ypos += 6
		return
	}
	mm.CurrentOption--
	mm.Cursor[0].Ypos -= 3
}

func (mm *MainMenu) keyDownListener(s *Screen) {

	s.DrawContent(mm.Cursor[0].Xpos, mm.Cursor[0].Ypos, mm.ClearCursor)
	if mm.CurrentOption == 2 {
		mm.CurrentOption = 0
		mm.Cursor[0].Ypos -= 6
		return
	}
	mm.CurrentOption++
	mm.Cursor[0].Ypos += 3
}

func (mm *MainMenu) keyEscListener(s *Screen) {
	os.Exit(0)
}

func (mm *MainMenu) keyEnterListener(s *Screen) {
	mm.animActive = false
	mm.animStopped = true
	switch mm.CurrentOption {
	case 0:
		// Create the Main Game Window
		s.Window.Clear()
		s.StartGame <- true
		return
	case 1:
		// Open Settings Menu
		s.Window.Clear()
		s.CurrentMenu = s.MenuList[1]
		s.CurrentMenu.(*SettingsMenu).animStopped = false
		s.CurrentMenu.Resize(s)
	case 2:
		os.Exit(0)
	}
}
