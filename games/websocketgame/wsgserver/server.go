// The WebSocket Game Server.
package wsgserver

import (
	"fmt"
	"log"
	"net"
)

var (
	// If VERBOSE is set to true it will start to print extra log messages.
	VERBOSE bool
)

func check_error(err error) {
	if err != nil {
		log.Fatal("ERROR ", err)
	}
}

func handle_client(c net.Conn) {
	client := NewClient(c)
	defer client.Close()
	SendAll(fmt.Sprintf("New client from %s", client.Addr))
	client.Loop()
}

// Start listening on a socket.
func Listen(proto, addr string) net.Listener {
	listen, err := net.Listen(proto, addr)
	check_error(err)

	if VERBOSE {
		log.Printf("Now accepting new connections from %s (%s)", addr, proto)
	}
	return listen
}

// Run the main loop and accept new clients.
func RunLoop(l net.Listener) {
	for {
		conn, err := l.Accept()
		check_error(err)
		go handle_client(conn)

	}
}
