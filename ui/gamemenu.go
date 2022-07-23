package ui

import (
	"fmt"
	"termtyper/settings"
	util "termtyper/utils"
	"time"

	"github.com/gdamore/tcell/v2"
)

type GameMenu struct {
	words [][]*Sprite

	CorrectTextColor   settings.DisplayStyle
	IncorrectTextColor settings.DisplayStyle
	GameContainer      *Container
	GameCursor         *GameCursor
	GameStart          bool
	yIndex             int
	xIndex             int
	wordIndex          int
}

type GameCursor struct {
	Style settings.DisplayStyle
	Xpos  int
	Ypos  int
}

func NewGameMenu(c *Container, s *settings.Settings) *GameMenu {
	words := StringArrayToSprite(util.GenerateWords(), s)

	gc := &GameCursor{
		Style: settings.DisplayStyle(tcell.Style(s.CursorColor).Attributes(tcell.AttrUnderline)),
		Xpos:  c.xBound.begin,
		Ypos:  c.yBound.begin,
	}

	return &GameMenu{
		words:              words,
		CorrectTextColor:   s.CorrectTextColor,
		IncorrectTextColor: s.IncorrectTextColor,
		GameContainer:      c,
		GameCursor:         gc,
		GameStart:          false,
		xIndex:             0,
		yIndex:             0,
		wordIndex:          0,
	}

}

func (gm *GameMenu) StartGameTimer(s *Screen) {
	ticker := time.NewTicker(1 * time.Second)
	go func(s *Screen) {
		i := 60
		for {
			if !gm.GameStart {
				continue
			}
			select {
			case <-ticker.C:
				if i == 60 {
					print("1:00")
				} else {
					fmt.Printf("0:%i", i)
				}
				i--
			}
			if i == -1 {
				break
			}
			// End the Game
			return
		}
	}(s)

}

func (gm *GameMenu) Show(s *Screen) {
	for !gm.GameContainer.isFull {
		gm.GameCursor.Xpos = gm.words[gm.yIndex][gm.xIndex].Xpos
		gm.GameCursor.Ypos = gm.words[gm.yIndex][gm.xIndex].Ypos
		if gm.wordIndex < len(gm.words)-1 {
			s.DrawContentWithinContainer(gm.GameContainer,
				gm.GameContainer.xCursor,
				gm.GameContainer.yCursor,
				gm.words[gm.wordIndex])
			gm.wordIndex++
		}
	}
	gm.GameCursor.Xpos = gm.words[gm.yIndex][gm.xIndex].Xpos
	gm.GameCursor.Ypos = gm.words[gm.yIndex][gm.xIndex].Ypos
	// Debug
	debug := fmt.Sprintf("cx:%d cy:%d val:%s", gm.words[gm.yIndex][gm.xIndex].Xpos, gm.words[gm.yIndex][gm.xIndex].Ypos, string(gm.words[gm.yIndex][gm.xIndex].Char))
	debug3 := fmt.Sprintf("cx_begin:%d  cx_end:%d cy_begin:%d cy_end:%d",
		gm.GameContainer.xBound.begin,
		gm.GameContainer.xBound.end,
		gm.GameContainer.yBound.begin,
		gm.GameContainer.yBound.end)
	r, _, _, _ := s.Window.GetContent(gm.GameCursor.Xpos, gm.GameContainer.yCursor)

	debug2 := fmt.Sprintf("gcx:%d gcy:%d val:%s", gm.GameCursor.Xpos, gm.GameCursor.Ypos, string(r))
	debug4 := fmt.Sprintf("Word Index: %d", gm.wordIndex)
	debug5 := fmt.Sprintf("Word Index: %d", gm.yIndex)
	s.DrawContent(1, 1, DebugStringToSprite(debug))
	s.DrawContent(1, 2, DebugStringToSprite(debug3))
	s.DrawContent(1, 3, DebugStringToSprite(debug2))
	s.DrawContent(1, 4, DebugStringToSprite(debug4))
	s.DrawContent(1, 5, DebugStringToSprite(debug5))

	// Draw cursor
	s.Window.SetContent(gm.GameCursor.Xpos, gm.GameCursor.Ypos, gm.words[gm.yIndex][gm.xIndex].Char, nil, tcell.Style(gm.GameCursor.Style))
}

func (gm *GameMenu) KeyListener(ev *tcell.EventKey, s *Screen) {
	key := ev.Key()

	if key == tcell.KeyEsc {

		s.Window.Clear()
		s.CurrentMenu = s.MenuList[0]
		s.CurrentMenu.(*MainMenu).animStopped = false
		s.CurrentMenu.Resize(s)
		return

	} else if key == tcell.KeyBackspace2 {

		if gm.xIndex-1 < 0 {
			if gm.yIndex-1 < 0 {
				return
			}

			if gm.GameCursor.Xpos-1 < gm.GameContainer.xBound.begin {
				r := gm.words[gm.yIndex-1][len(gm.words[gm.yIndex])-1].Char

				// If the word we are trying to go back to is a WhiteSpace we return
				if r == rune(' ') {
					return
				}

				// If not move to that letter in the word.. Ideally this should never happen
				s.Window.SetContent(
					gm.words[gm.yIndex][gm.xIndex].Xpos,
					gm.words[gm.yIndex][gm.xIndex].Ypos,
					gm.words[gm.yIndex][gm.xIndex].Char,
					nil,
					tcell.Style(gm.words[gm.yIndex][gm.xIndex].Style).Attributes(tcell.AttrNone),
				)
				gm.yIndex--
				gm.xIndex = len(gm.words[gm.yIndex]) - 1
				gm.GameCursor.Xpos = gm.GameContainer.xBound.end
				gm.GameCursor.Ypos--
				return
			} else if gm.words[gm.yIndex][gm.xIndex].Char == ' ' {
				s.Window.SetContent(
					gm.words[gm.yIndex][gm.xIndex].Xpos,
					gm.words[gm.yIndex][gm.xIndex].Ypos,
					gm.words[gm.yIndex][gm.xIndex].Char,
					nil,
					tcell.Style(gm.words[gm.yIndex][gm.xIndex].Style).Attributes(tcell.AttrNone),
				)
				gm.yIndex--
				gm.xIndex = len(gm.words[gm.yIndex]) - 1
				gm.GameCursor.Xpos--
				return
			}
			return
		}
		// r, _, _, _ := s.Window.GetContent(gm.GameCursor.Xpos-1, gm.GameContainer.yCursor)
		r := gm.words[gm.yIndex][gm.xIndex-1].Char
		s.Window.SetContent(
			gm.words[gm.yIndex][gm.xIndex].Xpos,
			gm.words[gm.yIndex][gm.xIndex].Ypos,
			gm.words[gm.yIndex][gm.xIndex].Char,
			nil,
			tcell.Style(gm.words[gm.yIndex][gm.xIndex].Style).Attributes(tcell.AttrNone),
		)
		if r == rune(' ') {
			return
		}
		gm.xIndex--
		gm.GameCursor.Xpos--
		return

	} else if key == tcell.KeyEnter {

		if gm.words[gm.yIndex][gm.xIndex].Char == ' ' {

			if gm.xIndex+1 > len(gm.words[gm.yIndex])-1 {
				gm.xIndex = 0
				gm.yIndex++
			} else {
				gm.xIndex++
			}
			gm.GameCursor.Xpos = gm.words[gm.yIndex][gm.xIndex].Xpos
			gm.GameCursor.Ypos = gm.words[gm.yIndex][gm.xIndex].Ypos

		}

	} else {

		gm.GameStart = true
		r := ' '
		if gm.GameCursor.Xpos == gm.words[gm.yIndex][gm.xIndex].Xpos && gm.GameCursor.Ypos == gm.words[gm.yIndex][gm.xIndex].Ypos {
			r = gm.words[gm.yIndex][gm.xIndex].Char
		}
		if r == ' ' {
			return
		}

		if ev.Rune() != r {
			gm.words[gm.yIndex][gm.xIndex].Style = gm.IncorrectTextColor
		} else {
			gm.words[gm.yIndex][gm.xIndex].Style = gm.CorrectTextColor
		}
		// Draw the Option to the Screen
		s.Window.SetContent(
			gm.words[gm.yIndex][gm.xIndex].Xpos,
			gm.words[gm.yIndex][gm.xIndex].Ypos,
			gm.words[gm.yIndex][gm.xIndex].Char,
			nil,
			tcell.Style(gm.words[gm.yIndex][gm.xIndex].Style),
		)

		if gm.xIndex+1 > len(gm.words[gm.yIndex])-1 {
			gm.xIndex = 0
			gm.yIndex++
		} else {
			gm.xIndex++
		}
		gm.GameCursor.Xpos = gm.words[gm.yIndex][gm.xIndex].Xpos
		gm.GameCursor.Ypos = gm.words[gm.yIndex][gm.xIndex].Ypos

	}
	if gm.wordIndex-1 == gm.yIndex {
		gm.wordIndex--
		gm.GameContainer.xCursor = gm.GameContainer.xBound.begin
		gm.GameContainer.yCursor = gm.GameContainer.yBound.begin
		gm.GameContainer.isFull = false
		s.Window.Clear()
	}

	return
}

func (gm *GameMenu) Resize(s *Screen) {
	s.Window.Sync()
	s.Window.Clear()
	winX, winY := s.Window.Size()
	new_Container := NewContainer(winX, winY)
	if gm.GameContainer.Size > new_Container.Size {
		gm.wordIndex = gm.yIndex
		gm.GameContainer = new_Container
		gm.GameCursor.Xpos = gm.GameContainer.xBound.begin + gm.xIndex
		gm.GameCursor.Ypos = gm.GameContainer.yBound.begin
	} else {
		gm.wordIndex = 0
		gm.GameContainer = new_Container
	}
	s.Window.Sync()
	s.Window.Clear()

	return
}
