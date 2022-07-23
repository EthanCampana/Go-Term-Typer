package settings

import "github.com/gdamore/tcell/v2"

type FrameRate int

type DisplayStyle tcell.Style

type Settings struct {
	CursorColor        DisplayStyle
	TextColor          DisplayStyle
	IncorrectTextColor DisplayStyle
	CorrectTextColor   DisplayStyle
	FrameRate          FrameRate
	ColorMap           map[tcell.Color]string
}

func DefaultSettings() *Settings {
	return &Settings{
		CursorColor:        DisplayStyle(tcell.StyleDefault),
		TextColor:          DisplayStyle(tcell.StyleDefault),
		IncorrectTextColor: DisplayStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkRed)),
		CorrectTextColor:   DisplayStyle(tcell.StyleDefault.Foreground(tcell.ColorLightGreen)),
		FrameRate:          FrameRate(60),
		ColorMap:           ReverseColorMap(),
	}
}

func ReverseColorMap() map[tcell.Color]string {
	n := make(map[tcell.Color]string)
	for k, v := range tcell.ColorNames {
		n[v] = k
	}
	n[0] = "Default"
	return n
}
