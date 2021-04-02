package main

import (
	"flag"
	"log"
	"os"

	"github.com/lmas/spacenanigans"
)

var (
	gAddr     = flag.String("addr", ":8000", "web server address to listen on")
	gDatabase = flag.String("db", "server.db", "database option")
	gMap      = flag.String("map", "data/maps/walltest.txt", "game map to run")
)

func main() {
	flag.Parse()

	conf := spacenanigans.Conf{
		Addr:     *gAddr,
		Database: *gDatabase,
		Map:      *gMap,
		Logger:   log.New(os.Stderr, "", log.Ldate|log.Ltime),
	}
	s, err := spacenanigans.New(conf)
	if err != nil {
		panic(err)
	}
	err = s.Run()
	if err != nil {
		panic(err)
	}
}
