package tools

import (
	"log"
	"os"

	"game-of-life/internal/resources"

	"github.com/gdamore/tcell/v2"
)

func InitializeScreen() (s tcell.Screen){
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	return s
}

func UpdateScreen(s tcell.Screen, w *resources.World, boxStyle tcell.Style) {
	drawBox(s, w, boxStyle)
	s.Show()
}

func drawBox(s tcell.Screen, w *resources.World, style tcell.Style) {
	x, y := w.GetHeight()-1, w.GetWidth()-1

	// Fill background
	bst := tcell.StyleDefault.Background(tcell.ColorRed)
	st := tcell.StyleDefault.Background(tcell.ColorWhite)

	for row := 0; row <= x; row++ {
		for col := 0; col <= y; col++ {
			if w.GetCellState(row, col){
				s.SetContent(3*col, row, ' ', nil, bst)
				s.SetContent(3*col+1, row, ' ', nil, bst)
				s.SetContent(3*col+2, row, ' ', nil, bst)
			} else {
				s.SetContent(3*col, row, ' ', nil, st)
				s.SetContent(3*col+1, row, ' ', nil, st)
				s.SetContent(3*col+2, row, ' ', nil, st)
			}
		}
	}

	// Draw borders
	for col := 0; col <= y; col++ {
		s.SetContent(3*col, 0, tcell.RuneHLine, nil, style)
		s.SetContent(3*col+1, 0, tcell.RuneHLine, nil, style)
		s.SetContent(3*col+2, 0, tcell.RuneHLine, nil, style)
		s.SetContent(3*col, x, tcell.RuneHLine, nil, style)
		s.SetContent(3*col+1, x, tcell.RuneHLine, nil, style)
		s.SetContent(3*col+2, x, tcell.RuneHLine, nil, style)
	}
	for row := 1; row < x; row++ {
		s.SetContent(0, row, tcell.RuneVLine, nil, style)
		s.SetContent(3*y+2, row, tcell.RuneVLine, nil, style)
	}

	// Draw corners
	s.SetContent(0, 0, tcell.RuneULCorner, nil, style)
	s.SetContent(3*y+2, 0, tcell.RuneURCorner, nil, style)
	s.SetContent(0, x, tcell.RuneLLCorner, nil, style)
	s.SetContent(3*y+2, x, tcell.RuneLRCorner, nil, style)
}

func CancelRoutine(s tcell.Screen) {

	// Create anonymous function to close screen and kill program
	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	for {
		// Poll event
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}
}
