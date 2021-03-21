package main

import (
	"log"

	engine "github.com/lmas/marsman/engine"

	curses "github.com/rthornton128/goncurses"
)

var (
	Running = true
	Player  *engine.Mob
)

func main() {
	engine.OpenLog("./tmp/debug.log")
	defer engine.CloseLog()

	engine.InitScreen()
	defer engine.ShutdownScreen()

	engine.InitSeeds(1)
	engine.LoadMap()
	engine.LoadItems()
	Player = engine.NewMob("Player", "@", 0, 0)
	Player.Tool = engine.NewItem(engine.ITEM_WRENCH)

	for Running == true {
		engine.Clear()
		engine.MapDraw(Player.Pos, 10)
		x, y := engine.GetMapScreenSize()
		engine.ScreenPrint(x+1, 0, engine.WHITE, "Pos: %s", Player.Pos.Str())
		engine.ScreenPrint(x+1, 1, engine.WHITE, "Dir: %d", Player.Dir)
		if Player.Tool != nil {
			engine.ScreenPrint(x+1, 2, engine.WHITE, "Tool: %s", Player.Tool.Name)
		}
		if Player.Block != nil {
			engine.ScreenPrint(x+1, 3, engine.WHITE, "Block: %s", Player.Block.Name)
		}
		engine.ScreenPrint(0, y, engine.WHITE, "Arrow keys to move, ESC/Q to quit")

		key := engine.GetKey()
		log.Println(key)
		switch key {
		case 27, 113:
			Running = false
		case curses.KEY_UP:
			Player.Move(0, -1)
			Player.Dir = engine.DIR_NORTH
		case curses.KEY_DOWN:
			Player.Move(0, 1)
			Player.Dir = engine.DIR_SOUTH
		case curses.KEY_LEFT:
			Player.Move(-1, 0)
			Player.Dir = engine.DIR_WEST
		case curses.KEY_RIGHT:
			Player.Move(1, 0)
			Player.Dir = engine.DIR_EAST
		case 117: // u
			Player.Use()
		case 112: // p
			Player.Place()
		}
	}
}
