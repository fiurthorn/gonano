package nano

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

type screenHandler interface {
	putStr(x, y int, b rune)
	clearStr(x, y int)
	getSize() (int, int)
	pollKeyPress() interface{}
	close()
	sync()
}

type keyEvent struct {
	rn rune
	// Instead copy-pasting and mapping all constants....
	k tcell.Key
}

type resizeEvent struct {
}

type physicalScreenHandler struct {
	screen tcell.Screen
}

func (s *physicalScreenHandler) close() {
	s.screen.Fini()
}

func (s *physicalScreenHandler) sync() {
	s.screen.Show()
}

func (s *physicalScreenHandler) putStr(x, y int, b rune) {
	s.screen.SetContent(x, y, b, []rune{}, tcell.StyleDefault)
}
func (s *physicalScreenHandler) clearStr(x, y int) {
	s.putStr(x, y, 0)
}

func (s *physicalScreenHandler) getSize() (int, int) {
	return s.screen.Size()
}

func (s *physicalScreenHandler) pollKeyPress() interface{} {
	for {
		switch ev := s.screen.PollEvent().(type) {
		case *tcell.EventKey:
			return keyEvent{rn: ev.Rune(), k: ev.Key()}
		case *tcell.EventResize:
			return resizeEvent{}
		}
	}
}

func initPhysicalScreenHandler() screenHandler {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error 1: %v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "Error 2:%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	return &physicalScreenHandler{screen: s}
}
