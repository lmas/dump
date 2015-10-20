package main

import (
	"flag"

	"github.com/lmas/websocketgame/server"
)

var (
	// Command line flags
	addr  = flag.String("addr", "127.0.0.1:8000", "Server's listening address")
	debug = flag.Bool("debug", false, "Run in debug mode")
)

func main() {
	flag.Parse()
	server.DEBUG = *debug

	server.StartGame(*addr)

}
