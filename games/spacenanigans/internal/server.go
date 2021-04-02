package internal

import (
	"log"
	"net/http"
	"sync"
)

type Conf struct {
	Addr     string
	Database string
	Map      string
	Logger   *log.Logger
}

type Server struct {
	conf Conf

	web    *http.Server
	db     *DB
	static map[string]staticFile

	lock    *sync.RWMutex
	clients []*Client
	world   *Map
	lastID  int64
}

func NewServer(conf Conf) (*Server, error) {
	logger = conf.Logger
	s := &Server{
		conf: conf,
		lock: &sync.RWMutex{},
	}

	s.web = &http.Server{
		ReadTimeout:  gTimeout,
		WriteTimeout: gTimeout,
		IdleTimeout:  gTimeout,
		Addr:         s.conf.Addr,
		Handler:      s,
	}

	var err error
	//s.db, err = OpenDB(`host=localhost user=spacenanigans dbname=spacenanigans connect_timeout=60`)
	s.db, err = OpenDB(conf.Database)
	if err != nil {
		return nil, err
	}

	s.static, err = loadStatic("static/")
	if err != nil {
		return nil, err
	}

	s.world, err = LoadMap(s.conf.Map)
	if err != nil {
		return nil, err
	}

	go s.Simulate()
	return s, nil
}

func (s *Server) Run() error {
	Log("Server running on %s", s.conf.Addr)
	return s.web.ListenAndServe()
}
