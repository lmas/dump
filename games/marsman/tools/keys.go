package main

import (
	engine "github.com/lmas/marsman/engine"
	"github.com/rthornton128/goncurses"
)

var (
	Running = true
)

func main() {
	engine.InitScreen()
	defer engine.ShutdownScreen()

	var key goncurses.Key
	for Running == true {
		engine.Clear()
		engine.ScreenPrint(0, 0, engine.WHITE, "Press any key. Press ESC to quit.")
		engine.ScreenPrint(0, 1, engine.WHITE, "Key pressed: %d (%+q, %q)", key, key, key)

		key = engine.GetKey()
		switch key {
		case 27:
			Running = false
		}
	}
}
