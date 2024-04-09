package nano

import (
	"log"
	"strings"

	"github.com/fiurthorn/gonano/krypta"
)

// Editor is main editor structure.
type Editor struct {
	Display *display

	filename string
	modified bool
	mode     mode

	krypta *krypta.Krypta
}

func NewEditor(filename string, krypta *krypta.Krypta) *Editor {
	handler := initPhysicalScreenHandler()
	display := createDisplay(handler)
	editor := &Editor{Display: display, modified: false, krypta: krypta}
	editor.mode = newNormalMode(editor)

	editor.initData(filename)
	blinkr := newBlinker(editor)
	statusBar := newStatusBar(editor.Display)
	editor.setBlinker(blinkr)
	editor.setStausBar(statusBar)

	go editor.loop()

	return editor
}

func (e *Editor) loop() {
	e.Display.startLoop()
}

func (e *Editor) Close() {
	e.Display.Close()
}

func (e *Editor) initData(filename string) {
	e.filename = filename

	data, err := e.krypta.DecryptFile(e.filename)
	if err != nil {
		log.Fatalf("Error: failed reading file: %v", err)
	}

	// `(\r\n|\n|\r)`
	fields := strings.Split(data, "\n")
	for i, field := range fields {
		if i == 0 {
			e.Display.getCurrentEl().data = []rune(field)
			e.Display.getCurrentEl().pos = 0
		} else {
			newItem := Line{data: []rune(field), startingCoordY: -1, pos: 0, display: e.Display}
			e.Display.data.InsertAfter(&newItem, e.Display.currentElement)
			e.Display.currentElement = e.Display.currentElement.Next()
		}
	}

	e.Display.currentElement = e.Display.data.Front()
	e.Display.resyncBelow(e.Display.currentElement)
}

func (e *Editor) saveData() error {
	data := []rune{}
	for it := e.Display.data.Front(); it != nil; it = it.Next() {
		if chars := len(it.Value.(*Line).data); chars == 0 &&
			it == e.Display.data.Back() {
			continue
		}

		data = append(data, it.Value.(*Line).data...)
		data = append(data, rune(10))
	}

	if err := e.krypta.EncryptFile(e.filename, strings.NewReader(string(data))); err != nil {
		return err
	}

	return nil
}

func (e *Editor) PollKeyboard(resp chan bool) {
	for {
		ev := e.Display.screen.pollKeyPress()

		switch t := ev.(type) {
		case keyEvent:
			exit := e.mode.handleKeyPress(t, resp)
			if exit {
				return
			}
		case resizeEvent:
			e.Display.resyncBelow(e.Display.data.Front())
		}

	}
}

func (e *Editor) setMode(mode mode) {
	e.mode = mode
	e.mode.init()
}

func (e *Editor) setBlinker(b blinker) {
	e.Display.blinker = b
}

func (e *Editor) setStausBar(s statusBar) {
	e.Display.statusBar = s
}
