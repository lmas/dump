package main

import (
	"flag"

	"github.com/lmas/websocketgame/wsgserver"
)

var (
	// Command line flags
	proto   = flag.String("proto", "tcp", "Listening socket protocol")
	addr    = flag.String("addr", "127.0.0.1:9000", "Server's listening address")
	verbose = flag.Bool("verbose", false, "Show incomming and sent messages")
)

func main() {
	flag.Parse()
	wsgserver.VERBOSE = *verbose
	sock := wsgserver.Listen(*proto, *addr)
	defer sock.Close()

	wsgserver.RunLoop(sock)
}
