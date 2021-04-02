package spacenanigans

import (
	"log"

	"github.com/lmas/spacenanigans/internal"
)

// This is just a nice wrapper around the internals

type Conf struct {
	Addr     string
	Database string
	Map      string
	Logger   *log.Logger
}

type Server struct {
	srv *internal.Server
}

func (s *Server) Run() error {
	return s.srv.Run()
}

func New(conf Conf) (*Server, error) {
	c := internal.Conf{
		Addr:     conf.Addr,
		Database: conf.Database,
		Map:      conf.Map,
		Logger:   conf.Logger,
	}

	s, err := internal.NewServer(c)
	if err != nil {
		return nil, err
	}
	return &Server{s}, nil
}
