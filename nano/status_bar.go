package nano

type statusBar interface {
	draw(text []string)
}

type physicalStatusBar struct {
	d *display
}

func newStatusBar(d *display) statusBar {
	return &physicalStatusBar{d}
}

func (s *physicalStatusBar) draw(text []string) {
	h := s.d.getHeight()
	w := s.d.getWidth()

	for i := 0; i < w; i++ {
		s.d.screen.putStr(i, h-1-len(text), '—')
	}
	for i := 0; i < len(text); i++ {
		for j, c := range text[i] {
			s.d.screen.putStr(j, h-len(text)+i, rune(c))
		}
	}

	s.d.screen.sync()

}
