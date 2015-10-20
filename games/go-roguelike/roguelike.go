package main

import (
	rl "github.com/lmas/go-rl/engine"

	curses "github.com/rthornton128/goncurses"
)

var (
	Running = true
)

const (
	COMP_PLAYER_INPUT = iota
	COMP_DOG_AI
)

func main() {
	rl.OpenLog("debug.log")
	defer rl.CloseLog()

	rl.InitScreen()
	defer rl.ShutdownScreen()

	rl.InitSeeds(1)
	rl.LoadMap()

	player := rl.NewMob("Player", "@", 0, 0)
	dog := rl.NewMob("Dog", "D", 10, 10)

	em := rl.NewEntityManager()
	em.NewComponent(COMP_PLAYER_INPUT, *new(map[string]interface{}))

	playerid := em.NewEntity()
	em.AddComponent(playerid, COMP_PLAYER_INPUT)

	// System to update screen and handle player input
	em.NewSystem(func(entities []int) {
		rl.Clear()
		rl.MapDraw(player.Pos, 10)
		x, y := rl.GetMapScreenSize()
		rl.ScreenPrint(x+1, 0, rl.WHITE, "Plr: %s", player.Pos.Str())
		rl.ScreenPrint(x+1, 1, rl.WHITE, "Dog: %s", dog.Pos.Str())
		rl.ScreenPrint(0, y, rl.WHITE, "Arrow keys to move, ESC/Q to quit")

		key := rl.GetKey()
		//log.Println(key)
		switch key {
		case 27, 113:
			Running = false
		case curses.KEY_UP:
			player.Move(0, -1)
		case curses.KEY_DOWN:
			player.Move(0, 1)
		case curses.KEY_LEFT:
			player.Move(-1, 0)
		case curses.KEY_RIGHT:
			player.Move(1, 0)
		}
	}, []int{COMP_PLAYER_INPUT})

	for {
		if !Running {
			break
		}
		em.RunSystems()
		dog.Track(&player)

	}
}

func PlayerInputSystem(entities []int) {
}
