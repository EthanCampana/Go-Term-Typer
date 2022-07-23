package ui

import (
	"os"
	"termtyper/settings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type SettingsMenu struct {
	KeyPressListener          map[tcell.Key]keyListenerFn
	ColorOptions              [][]*Sprite
	MenuOptions               [7][]*Sprite
	FrameRateOptions          [3][]*Sprite
	Title                     [3][]*Sprite
	Cursor                    []*Sprite
	ClearCursor               []*Sprite
	CurrentTextColor          int16
	CurrentCursorColor        int16
	CurrentIncorrectTextColor int16
	CurrentCorrectTextColor   int16
	CurrentOption             int8
	CurrentFrameRateOption    int8
	animActive                bool
	animStopped               bool
}

func getColors() []string {
	colors := []string{}
	colors = append(colors, "Default")
	for k := range tcell.ColorNames {
		colors = append(colors, k)
	}
	return colors
}

func findColorId(color string, colors []string) int {
	for i, k := range colors {
		if k == color {
			return i
		}
	}
	return 0
}
func findFrameRateOption(s *settings.Settings) int8 {
	switch s.FrameRate {
	case settings.FrameRate(30):
		return 0
	case settings.FrameRate(60):
		return 1
	case settings.FrameRate(90):
		return 2
	}
	return 0
}

func NewSettingsMenu(settings *settings.Settings) *SettingsMenu {
	options := [7][]*Sprite{}
	foptions := [3][]*Sprite{}
	title := [3][]*Sprite{}
	colors := getColors()
	coptions := StringArrayToSprite(colors, settings)
	tc, _, _ := tcell.Style(settings.TextColor).Decompose()
	cc, _, _ := tcell.Style(settings.CursorColor).Decompose()

	incc, _, _ := tcell.Style(settings.IncorrectTextColor).Decompose()
	corc, _, _ := tcell.Style(settings.CorrectTextColor).Decompose()

	title[0] = StringToSprite("█▀▀▀█ █▀▀█ ▀▀█▀▀  ▀  █▀▀█ █▀▀▄ █▀▀", settings)
	title[1] = StringToSprite("█   █ █  █   █   ▀█▀ █  █ █  █ ▀▀█", settings)
	title[2] = StringToSprite("█▄▄▄█ █▀▀▀   ▀   ▀▀▀ ▀▀▀▀ ▀  ▀ ▀▀▀", settings)

	foptions[0] = StringToSprite(" 30 >", settings)
	foptions[1] = StringToSprite("< 60 >", settings)
	foptions[2] = StringToSprite("< 90 ", settings)

	options[0] = StringToSprite("Text Color:", settings)
	options[1] = StringToSprite("Cursor Color:", settings)
	options[2] = StringToSprite("Incorrect Text Color:", settings)
	options[3] = StringToSprite("Correct Text Color:", settings)
	options[4] = StringToSprite("FrameRate:", settings)
	options[5] = StringToSprite("Save and Apply", settings)
	options[6] = StringToSprite("Back", settings)
	cursor := CursorToSprite("->", settings)
	clearCursor := StringToSprite("    ", settings)
	menu := &SettingsMenu{
		MenuOptions:               options,
		FrameRateOptions:          foptions,
		Title:                     title,
		ColorOptions:              coptions,
		Cursor:                    cursor,
		ClearCursor:               clearCursor,
		CurrentOption:             0,
		CurrentFrameRateOption:    findFrameRateOption(settings),
		CurrentTextColor:          int16(findColorId(settings.ColorMap[tc], colors)),
		CurrentIncorrectTextColor: int16(findColorId(settings.ColorMap[incc], colors)),
		CurrentCorrectTextColor:   int16(findColorId(settings.ColorMap[corc], colors)),
		CurrentCursorColor:        int16(findColorId(settings.ColorMap[cc], colors)),
		animActive:                false,
		animStopped:               false,
	}
	menu.KeyPressListener = make(map[tcell.Key]keyListenerFn)
	menu.KeyPressListener[tcell.KeyEsc] = menu.keyEscListener
	menu.KeyPressListener[tcell.KeyUp] = menu.keyUpListener
	menu.KeyPressListener[tcell.KeyDown] = menu.keyDownListener
	menu.KeyPressListener[tcell.KeyEnter] = menu.keyEnterListener
	menu.KeyPressListener[tcell.KeyLeft] = menu.keyLeftListener
	menu.KeyPressListener[tcell.KeyRight] = menu.keyRightListener
	return menu
}

func (sm *SettingsMenu) AnimStart(s *Screen) {
	framerate := time.NewTicker(1 * time.Second / 5)
	go func(s *Screen) {
		i := 0
		for range framerate.C {
			if !sm.animActive && sm.animStopped {
				break
			}
			if i == 0 {
				s.DrawContent(sm.Cursor[0].Xpos, sm.Cursor[0].Ypos, sm.ClearCursor)
				sm.Cursor[0].Xpos -= 3
				s.DrawContent(sm.Cursor[0].Xpos, sm.Cursor[0].Ypos, sm.Cursor)
				i++
			} else {
				s.DrawContent(sm.Cursor[0].Xpos, sm.Cursor[0].Ypos, sm.ClearCursor)
				sm.Cursor[0].Xpos += 3
				s.DrawContent(sm.Cursor[0].Xpos, sm.Cursor[0].Ypos, sm.Cursor)
				i--

			}
		}
		return
	}(s)
}

func (sm *SettingsMenu) keyEscListener(s *Screen) {
	os.Exit(0)
}

func (sm *SettingsMenu) Resize(s *Screen) {
	s.Window.Sync()
	s.Window.Clear()
	winX, winY := s.Window.Size()
	offset := (3 * sm.CurrentOption) - 6
	s.DrawContent(((winX / 2) - (len(sm.MenuOptions[1]) / 2) - 4), winY/2+int(offset), sm.Cursor)
}

func (sm *SettingsMenu) Show(s *Screen) {

	winX, winY := s.Window.Size()
	//Checking if the cursor was just initialized
	if sm.Cursor[0].Xpos == 0 {
		//Move the Cursor to the correct Location
		s.DrawContent(((winX / 2) - (len(sm.MenuOptions[1]) / 2) - 4), winY/2-12, sm.Cursor)
	}
	for i := range sm.MenuOptions {
		//Draws all the menu options
		offset := (3 * i) - 6
		s.DrawContent((winX/2)-(len(sm.MenuOptions[1])/2), winY/2+offset, sm.MenuOptions[i])
	}

	if !sm.animActive && !sm.animStopped {
		sm.AnimStart(s)
		sm.animActive = true
	}

	for i := range sm.Title {
		s.DrawContent((winX/2)-(len(sm.Title[0])/2), sm.MenuOptions[0][0].Ypos-6+i, sm.Title[i])
	}

	s.DrawContent((winX/2)+(len(sm.MenuOptions[1])/2)+10, winY/2-6, sm.ColorOptions[sm.CurrentTextColor])
	s.DrawContent((winX/2)+(len(sm.MenuOptions[1])/2)+10, winY/2-3, sm.ColorOptions[sm.CurrentCursorColor])
	s.DrawContent((winX/2)+(len(sm.MenuOptions[1])/2)+10, winY/2, sm.ColorOptions[sm.CurrentIncorrectTextColor])
	s.DrawContent((winX/2)+(len(sm.MenuOptions[1])/2)+10, winY/2+3, sm.ColorOptions[sm.CurrentCorrectTextColor])
	s.DrawContent((winX/2)+(len(sm.MenuOptions[1])/2), winY/2+6, sm.FrameRateOptions[sm.CurrentFrameRateOption])

	//Draw the last Location Of the Cursor Sprite
	s.DrawContent(sm.Cursor[0].Xpos, sm.Cursor[0].Ypos, sm.Cursor)
}

func (sm *SettingsMenu) keyUpListener(s *Screen) {
	//Might need to -2 the cursor xpos because of Animation that plays
	//This function will not handle the actually drawing of the cursor
	s.DrawContent(sm.Cursor[0].Xpos, sm.Cursor[0].Ypos, sm.ClearCursor)
	if sm.CurrentOption == 0 {
		sm.CurrentOption = 6
		sm.Cursor[0].Ypos += 18
		return
	}
	sm.CurrentOption--
	sm.Cursor[0].Ypos -= 3
}

func (sm *SettingsMenu) keyDownListener(s *Screen) {
	//Might need to -2 the cursor xpos because of Animation that plays
	//This function will not handle the actually drawing of the cursor
	s.DrawContent(sm.Cursor[0].Xpos, sm.Cursor[0].Ypos, sm.ClearCursor)
	if sm.CurrentOption == 6 {
		sm.CurrentOption = 0
		sm.Cursor[0].Ypos -= 18
		return
	}
	sm.CurrentOption++
	sm.Cursor[0].Ypos += 3
}

func (sm *SettingsMenu) keyEnterListener(s *Screen) {

	switch sm.CurrentOption {
	case 5:

		tColor := tcell.ColorNames[SpriteToString(sm.ColorOptions[sm.CurrentTextColor])]
		cColor := tcell.ColorNames[SpriteToString(sm.ColorOptions[sm.CurrentCursorColor])]
		incColor := tcell.ColorNames[SpriteToString(sm.ColorOptions[sm.CurrentIncorrectTextColor])]
		corColor := tcell.ColorNames[SpriteToString(sm.ColorOptions[sm.CurrentCorrectTextColor])]

		var framerate int
		switch sm.CurrentFrameRateOption {
		case 0:
			framerate = 30
		case 1:
			framerate = 60
		case 2:
			framerate = 90

		}
		tStyle := tcell.StyleDefault.Foreground(tColor)
		cStyle := tcell.StyleDefault.Foreground(cColor)
		incStyle := tcell.StyleDefault.Foreground(incColor)
		corStyle := tcell.StyleDefault.Foreground(corColor)

		s.UpdateChannel <- &settings.Settings{
			FrameRate:          settings.FrameRate(framerate),
			CursorColor:        settings.DisplayStyle(cStyle),
			TextColor:          settings.DisplayStyle(tStyle),
			IncorrectTextColor: settings.DisplayStyle(incStyle),
			CorrectTextColor:   settings.DisplayStyle(corStyle),
			ColorMap:           settings.ReverseColorMap(),
		}
	case 6:
		sm.animActive = false
		sm.animStopped = true
		s.Window.Clear()
		s.CurrentMenu = s.MenuList[0]
		s.CurrentMenu.(*MainMenu).animStopped = false
		s.CurrentMenu.Resize(s)
	}
}
func (sm *SettingsMenu) keyLeftListener(s *Screen) {
	switch sm.CurrentOption {
	case 0:
		if sm.CurrentTextColor != 0 {
			sm.CurrentTextColor--
			s.Window.Clear()
		}
	case 1:
		if sm.CurrentCursorColor != 0 {
			sm.CurrentCursorColor--
			s.Window.Clear()
		}
	case 2:
		if sm.CurrentIncorrectTextColor != 0 {
			sm.CurrentIncorrectTextColor--
			s.Window.Clear()
		}
	case 3:
		if sm.CurrentCorrectTextColor != 0 {
			sm.CurrentCorrectTextColor--
			s.Window.Clear()
		}
	case 4:
		if sm.CurrentFrameRateOption != 0 {
			sm.CurrentFrameRateOption--
			s.Window.Clear()
		}
	}
}

func (sm *SettingsMenu) keyRightListener(s *Screen) {
	switch sm.CurrentOption {
	case 0:
		if int(sm.CurrentTextColor) != len(sm.ColorOptions)-1 {
			sm.CurrentTextColor++
			s.Window.Clear()
		}
	case 1:
		if int(sm.CurrentCursorColor) != len(sm.ColorOptions)-1 {
			sm.CurrentCursorColor++
			s.Window.Clear()
		}
	case 2:
		if int(sm.CurrentIncorrectTextColor) != len(sm.ColorOptions)-1 {
			sm.CurrentIncorrectTextColor++
			s.Window.Clear()
		}
	case 3:
		if int(sm.CurrentCorrectTextColor) != len(sm.ColorOptions)-1 {
			sm.CurrentCorrectTextColor++
			s.Window.Clear()
		}
	case 4:
		if int(sm.CurrentFrameRateOption) != len(sm.FrameRateOptions)-1 {
			sm.CurrentFrameRateOption++
			s.Window.Clear()
		}

	}
}
