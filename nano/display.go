package nano

import (
	"container/list"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// display ss
type display struct {
	data           *list.List // list of *Line
	currentElement *list.Element
	screen         screenHandler
	monitorChannel chan contentOperation
	blinker        blinker
	statusBar      statusBar
	offsetY        int
}

func (c *display) dump() {
	fmt.Printf("Current: x:%v, y:%v", c.getBlinkerX(), c.getBlinkerY())
	fmt.Printf("width :%v, height:%v", c.getWidth(), c.getHeight())
	fmt.Println("Dumping lines:")
	for i, e := 0, c.data.Front(); e != nil; i, e = i+1, e.Next() {
		l := e.Value.(*Line)
		fmt.Printf("Line %v: data %v startY %v height %v pos %v", i, string(l.data), l.startingCoordY, l.calculateHeight(), l.pos)
	}
}

func (c *display) getWidth() int {
	w, _ := c.screen.getSize()
	return w
}

func (c *display) getHeight() int {
	_, h := c.screen.getSize()
	return h
}

func (c *display) getCurrentEl() *Line {
	return c.currentElement.Value.(*Line)
}

func (c *display) hasPrevEl() bool {
	return c.currentElement.Prev() != nil
}

func (c *display) getPrevEl() *Line {
	return c.currentElement.Prev().Value.(*Line)
}

func (c *display) hasNextEl() bool {
	return c.currentElement.Next() != nil
}

func (c *display) getNextEl() *Line {
	return c.currentElement.Next().Value.(*Line)
}

func createDisplay(handler screenHandler) *display {
	channel := make(chan contentOperation)
	lst := list.New()
	d := display{screen: handler, data: lst, monitorChannel: channel}
	lst.PushBack(d.newLine())
	d.currentElement = lst.Back()
	return &d
}

func (c *display) Close() {
	c.screen.close()
}

func (c *display) startLoop() {
	for op := range c.monitorChannel {
		c.blinker.clear()

		switch decoded := op.(type) {
		case typeOperation:
			{
				c.handleKeyPress(decoded)
				c.blinker.refresh()
				if decoded.resp != nil {
					decoded.resp <- true
				}
			}
		case blinkOperation:
			{
				c.blinker.refresh()
			}
		case announcementOperation:
			{
				c.statusBar.draw(decoded.text)
				if decoded.resp != nil {
					decoded.resp <- true
				}
			}
		}
	}
}

type contentOperation interface{}

type typeOperation struct {
	rn   rune
	key  tcell.Key
	resp chan bool
}

type blinkOperation struct {
	blink bool
}

type announcementOperation struct {
	text []string
	resp chan bool
}
