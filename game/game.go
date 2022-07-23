package Game

import (
	"termtyper/settings"
	"termtyper/ui"
	"time"

	"github.com/gdamore/tcell/v2"
)

type game struct {
	UI       *ui.Screen
	Settings *settings.Settings
	ticker   *time.Ticker
}

func New() *game {
	settings := settings.DefaultSettings()
	framerate := time.NewTicker(1 * time.Second / time.Duration(settings.FrameRate))
	return &game{
		UI:       ui.NewScreen(settings),
		Settings: settings,
		ticker:   framerate,
	}

}

func (g *game) Start() {
	g.startFPS()
	g.Run()
}

func (g *game) startFPS() {
	go func() {
		for _ = range g.ticker.C {
			g.UI.Window.Show()
			g.UI.CurrentMenu.Show(g.UI)
		}
		return
	}()
}

func (g *game) updateUI(s *settings.Settings) *game {
	g.ticker.Stop()
	framerate := time.NewTicker(1 * time.Second / time.Duration(s.FrameRate))
	return &game{
		UI:       ui.NewScreen(s),
		Settings: s,
		ticker:   framerate,
	}
}

func (g *game) Run() {
	for {
		ev := g.UI.Window.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			g.UI.CurrentMenu.Resize(g.UI)
		case *tcell.EventKey:
			switch menu := g.UI.CurrentMenu.(type) {
			case *ui.GameMenu:
				menu.KeyListener(ev, g.UI)
			case *ui.MainMenu:
				val, ok := menu.KeyPressListener[ev.Key()]
				if ok {
					val(g.UI)
				}
			case *ui.SettingsMenu:
				val, ok := menu.KeyPressListener[ev.Key()]
				if ok {
					val(g.UI)
				}

			}
		}
		select {
		case s := <-g.UI.UpdateChannel:
			g = g.updateUI(s)
			g.startFPS()
			break
		case <-time.After(5 * time.Millisecond):
			break
		case <-g.UI.StartGame:
			g.UI.NewGame(g.Settings)
		}
	}
}
