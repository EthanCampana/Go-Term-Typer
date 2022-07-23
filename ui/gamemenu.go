package ui

import (
	"fmt"
	"termtyper/settings"
	"termtyper/stats"
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
	GameStats          *stats.Stats
	GameStart          bool
	GameEnd            bool
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
	gs := stats.New()

	return &GameMenu{
		words:              words,
		CorrectTextColor:   s.CorrectTextColor,
		IncorrectTextColor: s.IncorrectTextColor,
		GameContainer:      c,
		GameCursor:         gc,
		GameStats:          gs,
		GameStart:          false,
		GameEnd:            false,
		xIndex:             0,
		yIndex:             0,
		wordIndex:          0,
	}

}

func (gm *GameMenu) StartGameTimer(s *Screen, op *settings.Settings) {
	go func(s *Screen, op *settings.Settings) {
		i := 60
		winX, _ := s.Window.Size()
		for {
			if gm.GameEnd {
				break
			}
			if !gm.GameStart {
				continue
			}

			ticker := time.NewTicker(1 * time.Second)
			select {
			case <-ticker.C:
				if i == 60 {
					s.DrawContent(winX-10, 1, StringToSprite("1:00", op))
				} else if i > 10 {
					s.DrawContent(winX-10, 1, StringToSprite(fmt.Sprintf("0:%d", i), op))
				} else {
					s.DrawContent(winX-10, 1, StringToSprite(fmt.Sprintf("0:0%d", i), op))
				}
				i--
			}
			if i == -1 {
				break
			}
		}
		s.Window.Clear()
		s.CurrentMenu = s.MenuList[0]
		s.CurrentMenu.(*MainMenu).animStopped = false
		s.CurrentMenu.Resize(s)
		return
	}(s, op)

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

	// Draw cursor
	s.Window.SetContent(gm.GameCursor.Xpos, gm.GameCursor.Ypos, gm.words[gm.yIndex][gm.xIndex].Char, nil, tcell.Style(gm.GameCursor.Style))
}

func (gm *GameMenu) KeyListener(ev *tcell.EventKey, s *Screen) {
	key := ev.Key()

	if key == tcell.KeyEsc {
		gm.GameEnd = true
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
			wpmCheck := false
			for _, c := range gm.words[gm.yIndex-2] {
				if c.Correct == false {
					wpmCheck = false
					break
				}
				wpmCheck = true

			}
			if wpmCheck {
				gm.GameStats.Wpm++
			}

			return
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
			gm.words[gm.yIndex][gm.xIndex].Correct = false
			gm.GameStats.Streak = 0
		} else {
			gm.words[gm.yIndex][gm.xIndex].Style = gm.CorrectTextColor
			gm.words[gm.yIndex][gm.xIndex].Correct = true
			gm.GameStats.Streak++
			gm.GameStats = gm.GameStats.UpdateStreak()
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
