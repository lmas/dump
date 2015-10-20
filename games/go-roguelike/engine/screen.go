package engine

import (
	"os"
	"os/signal"
	"syscall"

	curses "github.com/rthornton128/goncurses"
)

const (
	BLACK int16 = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
	B_BLACK
	B_RED
	B_GREEN
	B_YELLOW
	B_BLUE
	B_MAGENTA
	B_CYAN
	B_WHITE
)

var stdscr *curses.Window

func InitScreen() {
	stdscr = NewScreen()
	err := curses.StartColor()
	check_error(err)

	curses.InitPair(BLACK, curses.C_BLACK, curses.C_BLACK)
	curses.InitPair(RED, curses.C_RED, curses.C_BLACK)
	curses.InitPair(GREEN, curses.C_GREEN, curses.C_BLACK)
	curses.InitPair(YELLOW, curses.C_YELLOW, curses.C_BLACK)
	curses.InitPair(BLUE, curses.C_BLUE, curses.C_BLACK)
	curses.InitPair(MAGENTA, curses.C_MAGENTA, curses.C_BLACK)
	curses.InitPair(CYAN, curses.C_CYAN, curses.C_BLACK)
	curses.InitPair(WHITE, curses.C_WHITE, curses.C_BLACK)

	curses.InitPair(B_BLACK, curses.C_BLACK, curses.C_BLACK)
	curses.InitPair(B_RED, curses.C_BLACK, curses.C_RED)
	curses.InitPair(B_GREEN, curses.C_BLACK, curses.C_GREEN)
	curses.InitPair(B_YELLOW, curses.C_BLACK, curses.C_YELLOW)
	curses.InitPair(B_BLUE, curses.C_BLACK, curses.C_BLUE)
	curses.InitPair(B_MAGENTA, curses.C_BLACK, curses.C_MAGENTA)
	curses.InitPair(B_CYAN, curses.C_BLACK, curses.C_CYAN)
	curses.InitPair(B_WHITE, curses.C_BLACK, curses.C_WHITE)

	// Catch window resize signal and resize the screen
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGWINCH)
	go func() {
		for {
			_ = <-sigc
			ResizeScreen()
		}
	}()
}

func NewScreen() *curses.Window {
	scr, err := curses.Init()
	check_error(err)

	curses.Echo(false) // do not echo keys
	curses.Cursor(0)   // show no cursor
	curses.Raw(true)   // raw input
	scr.Keypad(true)

	return scr
}

func ShutdownScreen() {
	curses.End()
}

func Clear() {
	stdscr.Erase()
}

func GetScreenSize() (x, y int) {
	y, x = stdscr.MaxYX()
	return x, y
}

func ResizeScreen() {
	ShutdownScreen()
	stdscr = NewScreen()
}

func GetKey() curses.Key {
	return stdscr.GetChar()
}

func GetStr(size int) (string, error) {
	return stdscr.GetString(size)
}

func ScreenPrint(x int, y int, color int16, format string, args ...interface{}) {
	stdscr.ColorOn(color)
	stdscr.MovePrintf(y, x, format, args...)
	stdscr.ColorOff(color)
}
