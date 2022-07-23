package ui

import (
	"fmt"
	"termtyper/settings"
	"termtyper/stats"
)

type GameOverMenu struct {
	Title [][]*Sprite
	Stats [][]*Sprite
}

func NewGameOverMenu(stats *stats.Stats, op *settings.Settings) *GameOverMenu {
	title := []string{
		"█▀▀█ █▀▀█ █▀▄▀█ █▀▀    █▀▀▀█ ▀█ █▀ █▀▀ █▀▀█",
		"█ ▄▄ █▄▄█ █ ▀ █ █▀▀    █   █  █▄█  █▀▀ █▄▄▀",
		"█▄▄█ ▀  ▀ ▀   ▀ ▀▀▀    █▄▄▄█   ▀   ▀▀▀ ▀ ▀▀",
	}
	content := []string{
		fmt.Sprintf("Completed Words per Minute: %d", stats.Wpm),
		fmt.Sprintf("Longest Streak without a Mistake: %d characters", stats.GetLongestStreak()),
		fmt.Sprintf("Number of Mistakes: %d", stats.Mistakes),
		fmt.Sprintf("Please Press Any Button to Head back to the Main Menu"),
	}
	return &GameOverMenu{
		Title: StringArrayToSprite(title, op),
		Stats: StringArrayToSprite(content, op),
	}

}

func (gom *GameOverMenu) BackToMainMenu(s *Screen) {
	s.Window.Clear()
	s.CurrentMenu = s.MenuList[0]
	s.CurrentMenu.(*MainMenu).animStopped = false
	s.CurrentMenu.Resize(s)
}

func (gom *GameOverMenu) Show(s *Screen) {
	winX, winY := s.Window.Size()
	for i := range gom.Stats {
		//Draws the stats
		offset := (3 * i) - 3
		s.DrawContent((winX/2)-(len(gom.Stats[1])/2), winY/2+offset, gom.Stats[i])
	}
	for i := range gom.Title {
		s.DrawContent((winX/2)-(len(gom.Title[0])/2), gom.Stats[0][0].Ypos-7+i, gom.Title[i])
	}

	return
}

func (gom *GameOverMenu) Resize(s *Screen) {
	s.Window.Sync()
	s.Window.Clear()
	return
}
