package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fiurthorn/gonano/krypta"
	"github.com/fiurthorn/gonano/nano"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	arg := os.Args[1]
	switch arg {
	case "keygen":
		krypta.Keygen(krypta.IdentitiesFile(), krypta.RecipientsFile())
	default:
		editor := nano.NewEditor(os.Args[1], krypta.New())
		defer editor.Close()
		editor.PollKeyboard(nil)
	}
}
