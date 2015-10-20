package main

import (
	"flag"
	"log"

	"github.com/lmas/websocketgame/wsgclient"
)

var (
	proto = flag.String("proto", "tcp", "Socket protocol")
	addr  = flag.String("addr", "127.0.0.1:9000", "Server's address")
)

func check_error(err error) {
	if err != nil {
		log.Fatal("ERROR ", err)
	}
}

func main() {
	flag.Parse()
	client := wsgclient.Connect(*proto, *addr)
	defer client.Close()

	client.Loop()
}
