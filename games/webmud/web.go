package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	gTimeout  = time.Duration(60) * time.Second
	gUpgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: gTimeout,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func webClientAddr(r *http.Request) string {
	addr := r.Header.Get("x-forwarded-for") // Possible proxy addr
	if addr == "" {
		addr = r.RemoteAddr // fallback
	}
	return addr
}

////////////////////////////////////////////////////////////////////////////////

type Server struct {
	pool *WorkPool
}

func NewServer() *Server {
	s := &Server{
		pool: NewWorkPool(1, 10),
	}
	return s
}

func (s *Server) Log(msg string, args ...interface{}) {
	Log(msg, args...)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Log("%s\t %s\t %s\t",
		webClientAddr(r),
		r.Method,
		r.URL.Path,
	)

	var f func(http.ResponseWriter, *http.Request) error
	switch r.URL.Path {
	case "/ws":
		f = s.newWebsock
	default:
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}

	if err := f(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) newWebsock(w http.ResponseWriter, r *http.Request) error {
	conn, err := gUpgrader.Upgrade(w, r, nil)
	if err != nil {
		// The upgrader auto sends a http error already
		return nil
	}
	c := NewClient(s, conn)
	c.Write([]byte(gLogo))
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type Client struct {
	server *Server
	conn   *websocket.Conn
}

func NewClient(s *Server, conn *websocket.Conn) *Client {
	c := &Client{
		server: s,
		conn:   conn,
	}
	go c.readLoop()
	return c
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Write(msg []byte) {
	c.server.pool.Do(func() {
		if err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			//return err
			c.server.Log("ws write error: %s", err)
			c.Close()
		}
	})
}

func (c *Client) readLoop() {
	for {
		msgType, b, err := c.conn.ReadMessage()
		if err != nil {
			//return err
			c.server.Log("ws read error: %s", err)
			break
		}
		if msgType != websocket.TextMessage {
			//return fmt.Errorf("invalid ws message type: %d", msgType)
			c.server.Log("invalid ws message type: %d", msgType)
			break
		}
		//fmt.Println(string(b))
		c.Write(b)
	}
	c.Close()
}
