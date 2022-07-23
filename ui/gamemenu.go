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
	// Debug
	debug := fmt.Sprintf("cx:%d cy:%d val:%s", gm.words[gm.yIndex][gm.xIndex].Xpos, gm.words[gm.yIndex][gm.xIndex].Ypos, string(gm.words[gm.yIndex][gm.xIndex].Char))
	debug3 := fmt.Sprintf("cx_begin:%d  cx_end:%d cy_begin:%d cy_end:%d",
		gm.GameContainer.xBound.begin,
		gm.GameContainer.xBound.end,
		gm.GameContainer.yBound.begin,
		gm.GameContainer.yBound.end)
	r, _, _, _ := s.Window.GetContent(gm.GameCursor.Xpos, gm.GameContainer.yCursor)
	debug2 := fmt.Sprintf("gcx:%d gcy:%d val:%s", gm.GameCursor.Xpos, gm.GameCursor.Ypos, string(r))
	s.DrawContent(1, 1, DebugStringToSprite(debug))
	s.DrawContent(1, 2, DebugStringToSprite(debug3))
	s.DrawContent(1, 3, DebugStringToSprite(debug2))

	// Draw cursor
	s.Window.SetContent(gm.GameCursor.Xpos, gm.GameCursor.Ypos, gm.words[gm.yIndex][gm.xIndex].Char, nil, tcell.Style(gm.GameCursor.Style))
}

/*
Resize notes:
To find the cursor position after resize:
new x : cur x * change in x of the grid (oldx / new x) + any modifiers (screen wrapping etc)
new y : same as x ^
We need to account for whether the positions round up or down.
ideally the rune values should be the same  before_rune == after_rune
This is another check to make sure we ended up at the same position


Static Items can be resized the same way as the other menu's

*/
/*
Game Design Notes:
- There will be no new lines just spaces Woohooo We can just have one massive [][]*Sprite
	- The way things are drawn to the screen will handle the wrapping of content
	- We drawing words to a screen we should make sure it fits len(word) < len(empty space in line) if not draw it on the next lin

- Game area bounds
- We will need a new function DrawContentInBounds()
	- This function should fill the game area space with as much content as possible
	- Once the Player finishes all the information in the area it should be  called again and draw new content
	- The cursor should be moved to the top-left bound
  - We check the last word to see if there is any wrong characters. If so word is not counted and other statitics
  - Space and Enter should be disabled on letters

- BackSpacing a.k.a erasing errors:
	- []*Sprite is 1 whole word
	- [][]*Sprite will be a line of words
	- We will need two counters that we will need to keep track internally of where we are at.
	- Using array bounds we don't have to worry about overflow.
	- If we update the UI with when they type the wrong character with the wrong character.
	- all we have to do is move the counter back one space to get the data that is present


*/

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
		//Enter new line keys stuff

		// if gm.words[gm.yIndex][gm.xIndex].Char == ' ' && len(gm.words[gm.yIndex+1]) > gm.GameContainer.xBound.end-gm.GameCursor.Xpos {
		// 	gm.GameCursor.Xpos = gm.GameContainer.xBound.begin
		// 	gm.GameCursor.Ypos++
		// 	gm.xIndex = 0
		// 	gm.yIndex++

		// }
		// if gm.words[gm.yIndex][gm.xIndex].Char == ' ' && gm.words[gm.yIndex+1][gm.xIndex].Char != ' ' {
		// 	gm.GameCursor.Xpos++
		// 	gm.xIndex = 0
		// 	gm.yIndex++
		// }

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
		// Check if we exit the bounds of the container
		// if gm.GameCursor.Xpos+1 > gm.GameContainer.xBound.end {
		// 	gm.GameCursor.Xpos = gm.GameContainer.xBound.begin
		// 	gm.GameCursor.Ypos++
		// } else {
		// 	gm.GameCursor.Xpos++
		// }
		// Check if we exited the word
		if gm.xIndex+1 > len(gm.words[gm.yIndex])-1 {
			gm.xIndex = 0
			gm.yIndex++
		} else {
			gm.xIndex++
		}
		gm.GameCursor.Xpos = gm.words[gm.yIndex][gm.xIndex].Xpos
		gm.GameCursor.Ypos = gm.words[gm.yIndex][gm.xIndex].Ypos

	}
	return
}

func (gm *GameMenu) Resize(s *Screen) {
	// oldx, oldy := s.Window.Size()
	s.Window.Sync()
	s.Window.Clear()
	winX, winY := s.Window.Size()
	new_Container := NewContainer(winX, winY)
	// xDiff := winX / oldx
	// yDiff := winY / oldy
	// To Get the new Cursor location we multiply its current location with the differance in change
	// gm.GameCursor.Xpos *= xDiff
	// gm.GameCursor.Ypos *= yDiff
	if gm.GameContainer.Size > new_Container.Size {
		gm.wordIndex = gm.yIndex
		gm.GameContainer = new_Container
		gm.GameCursor.Xpos = gm.GameContainer.xBound.begin + gm.xIndex
		gm.GameCursor.Ypos = gm.GameContainer.yBound.begin
	} else {
		// I Already Know this is going to cause a bunch of problems... Need better Math here
		gm.wordIndex = 0
		gm.GameContainer = new_Container
	}
	s.Window.Sync()
	s.Window.Clear()

	return
}
