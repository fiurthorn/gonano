package nano

import "time"

type blinker interface {
	refresh()
	set()
	clear()
}

type realBlinker struct {
	blinkIsSet bool
	d          *display
}

func (r *realBlinker) refresh() {
	if r.blinkIsSet {
		r.set()
	} else {
		r.clear()
	}
}

func (r *realBlinker) set() {
	r.d.screen.putStr(r.d.getBlinkerX(), r.d.getBlinkerY(), rune('▉'))
	r.d.screen.sync()
}

func (r *realBlinker) clear() {
	if len(r.d.getCurrentEl().data) > r.d.getCurrentEl().pos {
		r.d.screen.putStr(r.d.getBlinkerX(), r.d.getBlinkerY(), r.d.getCurrentEl().getCurrentChar())
	} else {
		// FIXME: exception: it might be on beginning of another line. Fix the case.
		if r.d.getBlinkerX() == 0 && r.d.hasNextEl() && r.d.getBlinkerY() == r.d.getNextEl().startingCoordY {
			if len(r.d.getNextEl().data) > 0 {
				r.d.screen.putStr(r.d.getBlinkerX(), r.d.getBlinkerY(), r.d.getNextEl().data[0])
			} else {
				r.d.screen.clearStr(r.d.getBlinkerX(), r.d.getBlinkerY())
			}
		} else {
			r.d.screen.clearStr(r.d.getBlinkerX(), r.d.getBlinkerY())
		}
	}
	r.d.screen.sync()
}

func newBlinker(e *Editor) blinker {
	b := &realBlinker{d: e.Display, blinkIsSet: false}
	go func(c chan contentOperation) {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			<-ticker.C
			b.blinkIsSet = !b.blinkIsSet
			c <- blinkOperation{blink: b.blinkIsSet}
		}
	}(e.Display.monitorChannel)

	return b
}
