package nano

import (
	"container/list"
	"strconv"
)

func oneCharOperation(c *display, f func()) {
	oldH := c.getCurrentEl().getOnScreenLineEndingY()
	oldCursorY := c.getCurrentEl().getRelativeCursorY()
	f()
	newH := c.getCurrentEl().getOnScreenLineEndingY()
	newCursorY := c.getCurrentEl().getRelativeCursorY()

	if oldCursorY != newCursorY {
		c.resyncNewCursorY()
	} else {
		if oldH != newH {
			c.resyncBelow(c.currentElement)
		} else {
			c.getCurrentEl().resync()
		}
	}
}

func (c *display) insert(char rune) {
	if strconv.IsPrint(char) {
		oneCharOperation(c, func() {
			c.getCurrentEl().data = insertInSlice(c.getCurrentEl().data, char, c.getCurrentEl().pos)
			c.getCurrentEl().pos++
		})
	}
}

func (c *display) remove() {
	if c.getCurrentEl().pos == 0 {
		if !c.hasPrevEl() {
			return
		}

		// Remove current line
		p := c.currentElement.Prev()
		p.Value.(*Line).pos = len(p.Value.(*Line).data)
		p.Value.(*Line).data = append(p.Value.(*Line).data, c.getCurrentEl().data...)
		c.data.Remove(c.currentElement)
		c.currentElement = p
		c.recalcBelow(c.currentElement)

		// Fix Y!
		if c.getCurrentEl().getOnScreenCursorY() < 0 {
			c.offsetY--
		}
		c.resyncBelow(c.data.Front())
	} else {
		oneCharOperation(c, func() {
			c.getCurrentEl().pos--
			c.getCurrentEl().data = removeFromSlice(c.getCurrentEl().data, c.getCurrentEl().pos)
		})
	}
}

func (c *display) delete() {
	endPos := len(c.getCurrentEl().data)
	if c.getCurrentEl().pos == endPos {
		if !c.hasNextEl() {
			return
		}

		// Remove next line
		p := c.currentElement
		p.Value.(*Line).data = append(p.Value.(*Line).data, c.getNextEl().data...)
		c.data.Remove(c.currentElement.Next())
		c.recalcBelow(c.currentElement)

		// Fix Y!
		c.resyncBelow(c.data.Front())
	} else {
		oneCharOperation(c, func() {
			c.getCurrentEl().data = removeFromSlice(c.getCurrentEl().data, c.getCurrentEl().pos)
		})
	}
}

// Current line should have correct startingY !
func (c *display) resyncBelow(from *list.Element) {
	c.recalcBelow(from)
	for ; from != nil && from.Value.(*Line).startingCoordY-c.offsetY < c.getHeight(); from = from.Next() {
		from.Value.(*Line).resync()
	}

	// Clean at startingY
	startingY := 0
	if from != nil {
		startingY = from.Value.(*Line).startingCoordY + from.Value.(*Line).calculateHeight()
	} else {
		startingY = c.data.Back().Value.(*Line).startingCoordY + c.data.Back().Value.(*Line).calculateHeight()
	}

	for ; startingY-c.offsetY < c.getHeight(); startingY++ {
		for i := 0; i < c.getWidth(); i++ {
			c.screen.clearStr(i, startingY-c.offsetY)
		}
	}
	c.screen.sync()
}

func (c *display) recalcBelow(from *list.Element) {
	startingY := from.Value.(*Line).startingCoordY
	for ; from != nil; from = from.Next() {
		line := from.Value.(*Line)
		line.startingCoordY = startingY
		startingY += line.calculateHeight()
	}
}

func (c *display) resyncNewCursorY() {
	onScreenCursorY := c.getCurrentEl().getOnScreenCursorY()
	// If cursor jumped below screen
	if onScreenCursorY >= c.getHeight() {
		c.offsetY++
		c.resyncBelow(c.data.Front())
	} else if onScreenCursorY < 0 {
		c.offsetY--
		c.resyncBelow(c.data.Front())
	} else {
		c.resyncBelow(c.currentElement)
	}
}
