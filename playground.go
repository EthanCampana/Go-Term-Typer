package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type cursor struct {
	xpos int
	ypos int
	char string
}
type board struct {
	height    int
	width     int
	board_map [][]string
}

func NewBoard() *board {
	b := &board{}
	return b
}

type words struct {
	r     rune
	style tcell.Style
}

func ResizeTextArea(s tcell.Screen, rwords [][]rune, nx, ny int, c chan bool, cursor *cursor) {
	l := 0
	screen_content := [][]words{}
	_, _, c_style, _ := s.GetContent(cursor.xpos, cursor.ypos)
	s.Sync()
	//GRAB ALL INFO FROM THE SCREEN
	for z := range rwords {
		ct := []words{}
		xpos := (nx / 2) - (len(rwords[l]) / 2) + l
		ypos := (ny / 2) + l
		for _ = range rwords[z] {
			r, _, style, _ := s.GetContent(xpos, ypos)
			word := words{r: r, style: style}
			ct = append(ct, word)
			xpos++
		}
		screen_content = append(screen_content, ct)
		l++
	}
	time.Sleep(1000)
	//POST GATHERING ALL THE INFO FROM THE SCREEN
	s.Clear()
	ox, oy := s.Size()

	l = 0
	// drawText(s, 0, oy-10, ox, oy, tcell.StyleDefault, fmt.Sprintf("", screen_content[1]))
	// drawText(s, 0, oy-20, ox, oy, tcell.StyleDefault, fmt.Sprintf("", rwords[1]))
	for z := range screen_content {
		xpos := (ox / 2) - (len(rwords[l]) / 2) + l
		ypos := (oy / 2) + l
		for q := range screen_content[z] {
			s.SetContent(xpos, ypos, screen_content[z][q].r, nil, screen_content[z][q].style)
			if screen_content[z][q].style == c_style {
				cursor.xpos = xpos
				cursor.ypos = ypos
			}
			xpos++
		}
		l++
	}
	c <- true
}

func AdjustBoard(b *board, s tcell.Screen) {
	s.Sync()
	s.Clear()
	b.width, b.height = s.Size()
	new_map := make([][]string, b.height)
	for i := range new_map {
		new_map[i] = make([]string, b.width)
		for x := range new_map[i] {
			new_map[i][x] = "x"
		}
	}
	b.board_map = new_map
}

func New() *cursor {
	c := &cursor{xpos: 1, ypos: 1, char: " "}
	return c
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func moveCursor(mx, my, gx, gy int, c *cursor, s tcell.Screen, style tcell.Style, b *board) {
	if mx < 0 && c.xpos == 0 || mx > 0 && c.xpos == gx-1 {
		return
	} else if my < 0 && c.ypos == 0 || my > 0 && c.ypos == gy-1 {
		return
	} else {
		b.board_map[c.ypos][c.xpos] = "x"
		c.xpos += mx
		c.ypos += my
		b.board_map[c.ypos][c.xpos] = "c"
		r, _, style2, _ := s.GetContent(c.xpos, c.ypos)
		x, y := c.xpos, c.ypos
		s.SetContent(c.xpos, c.ypos, ' ', nil, style)
		s.SetContent(x-mx, y-my, r, nil, style2)
	}
}

func buildBorder(s tcell.Screen, b *board) {
	borderStyle := tcell.StyleDefault.Background(tcell.ColorMistyRose)
	for i := range b.board_map[0] {
		b.board_map[0][i] = "b"
		s.SetContent(i, 0, ' ', nil, borderStyle)
		time.Sleep(100 * time.Millisecond)
	}
	for i := range b.board_map {
		b.board_map[i][0] = "b"
		s.SetContent(0, i, ' ', nil, borderStyle)
		time.Sleep(100 * time.Millisecond)
	}
	for i := range b.board_map[b.height-1] {
		b.board_map[b.height-1][i] = "b"
		s.SetContent(i, b.height-1, ' ', nil, borderStyle)
		time.Sleep(100 * time.Millisecond)
	}
	for i := range b.board_map {
		b.board_map[i][len(b.board_map[0])-1] = "b"
		s.SetContent(len(b.board_map[0])-1, i, ' ', nil, borderStyle)
		time.Sleep(100 * time.Millisecond)
	}

}

func TyperPlayground(s tcell.Screen) {
	s.Clear()
	basicText := tcell.StyleDefault.Foreground(tcell.ColorDarkOrange)
	IncorrectText := tcell.StyleDefault.Foreground(tcell.ColorDarkRed)
	correctText := tcell.StyleDefault.Foreground(tcell.ColorLightGreen)
	cursorColor := tcell.StyleDefault.Attributes(tcell.AttrUnderline)
	ox, oy := s.Size()
	words := "Testing 123 Testing 123 123 Testing\nHello 123 Testing 123"
	val := strings.Split(words, "\n")

	rwords := [][]rune{}
	ch1 := make(chan bool)
	for _, word := range val {

		k := []rune(word)
		rwords = append(rwords, k)
	}

	line := 0
	pos := 0
	progress := 0
	drawText(s, (ox/2)-(len(rwords[line])/2)-line, (oy/2)-line, ox, oy, basicText, string(rwords[line]))
	c := New()
	c.xpos = (ox / 2) - (len(rwords[line]) / 2) - line
	c.ypos = (oy / 2) - line

	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	for {
		ox, oy = s.Size()
		r, _, _, _ := s.GetContent(c.xpos, c.ypos)
		s.SetContent(c.xpos, c.ypos, r, nil, cursorColor)
		p := float64(progress) / float64(len(rwords)) * 100
		progress_bar := fmt.Sprintf("Progress: %f percent", p)
		drawText(s, 0, oy-2, ox, oy, basicText, progress_bar)
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			go ResizeTextArea(s, rwords, ox, oy, ch1, c)
			<-ch1
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				quit()
			} else if ev.Key() == tcell.KeyBackspace2 {
				if pos == len(rwords[line]) {
					pos--
					s.SetContent(c.xpos, c.ypos, ' ', nil, basicText)
					c.xpos--
				}
				x := c.xpos
				p := pos
				if pos > 0 {
					pos--
					c.xpos--
					progress--
				}
				s.SetContent(x, c.ypos, rwords[line][p], nil, basicText)
				s.SetContent(c.xpos, c.ypos, rwords[line][pos], nil, cursorColor)

			} else if ev.Key() == tcell.KeyEnter {
				if pos != len(rwords[line]) {
					continue
				} else {
					s.SetContent(c.xpos, c.ypos, ' ', nil, basicText)
					line++
					pos = 0
					drawText(s, (ox/2)-(len(rwords[line])/2)+line, (oy/2)+line, ox, oy, basicText, string(rwords[line]))
					c.xpos = (ox / 2) - (len(rwords[line]) / 2) + line
					c.ypos = (oy / 2) + line
				}

			} else {
				if pos != len(rwords[line]) {
					u_key := ev.Rune()
					if u_key == rwords[line][pos] {
						s.SetContent(c.xpos, c.ypos, u_key, nil, correctText)
					} else {
						s.SetContent(c.xpos, c.ypos, u_key, nil, IncorrectText)
					}
					pos++
					if pos != len(rwords[line]) {
						c.xpos++
						progress++
						s.SetContent(c.xpos, c.ypos, rwords[line][pos], nil, cursorColor)
					} else {
						c.xpos++
						progress++
						s.SetContent(c.xpos, c.ypos, ' ', nil, cursorColor)
					}

				}
			}

		}
	}

}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorRebeccaPurple).Background(tcell.ColorWheat)
	b := NewBoard()
	framerate := time.NewTicker(1 * time.Second / 60)
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	cursor := New()
	AdjustBoard(b, s)
	s.Clear()
	go func(s tcell.Screen) {
		for _ = range framerate.C {
			s.Show()
		}
	}(s)
	for {
		//You Should always draw your screen before re-sizing so that the screen can resize appropriately
		//Only Draw things that are going to be in static position before re-sizing. Dynamic Things will need extra work to be drawn in the correct place
		ox, oy := s.Size()
		location := fmt.Sprintf("cx:%d cy:%d", cursor.xpos, cursor.ypos)
		size := fmt.Sprintf("x:%d y:%d", b.width, b.height)
		s.SetContent(cursor.xpos, cursor.ypos, ' ', nil, textStyle)
		drawText(s, ox-10, oy-2, ox, oy, textStyle, size)
		drawText(s, ox-10, oy-1, ox, oy, textStyle, location)
		ev := s.PollEvent()
		//s.SetContent(cursor.xpos, cursor.ypos, 'b', nil, textStyle)
		switch ev := ev.(type) {
		case *tcell.EventResize:
			AdjustBoard(b, s)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				quit()
			} else if ev.Key() == tcell.KeyUp {
				moveCursor(0, -1, ox, oy, cursor, s, textStyle, b)
			} else if ev.Key() == tcell.KeyDown {
				moveCursor(0, 1, ox, oy, cursor, s, textStyle, b)
			} else if ev.Key() == tcell.KeyRight {
				moveCursor(1, 0, ox, oy, cursor, s, textStyle, b)
			} else if ev.Key() == tcell.KeyLeft {
				moveCursor(-1, 0, ox, oy, cursor, s, textStyle, b)
			} else if ev.Rune() == 'b' {
				go buildBorder(s, b)
			} else if ev.Rune() == 'k' {
				TyperPlayground(s)
			}

		}
	}
}
