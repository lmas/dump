package server

import (
	"html"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/olahol/melody"
)

type ClientHandler struct {
	melody *melody.Melody
}

func NewClientHandler() *ClientHandler {
	m := melody.New()

	if DEBUG {
		upgrader := &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}
		m.Upgrader = upgrader
	}

	cl := ClientHandler{
		melody: m,
	}

	m.HandleError(cl.OnError)
	m.HandleConnect(cl.OnConnect)
	m.HandleDisconnect(cl.OnDisconnect)
	m.HandleMessage(cl.OnMessage)

	return &cl
}

func (cl *ClientHandler) OnRequest(c *gin.Context) {
	// HACK: Nowhere else to set the best ip for the melody session...
	c.Request.RemoteAddr = c.ClientIP()
	cl.melody.HandleRequest(c.Writer, c.Request)
}

func (cl *ClientHandler) OnError(s *melody.Session, err error) {
	if err != nil && err != io.EOF {
		log.Println("ERROR ", err)
	}
}

func (cl *ClientHandler) OnConnect(s *melody.Session) {
	cl.Broadcast(s.Request.RemoteAddr + " connected.") // TODO
}

func (cl *ClientHandler) OnDisconnect(s *melody.Session) {
	cl.Broadcast(s.Request.RemoteAddr + " disconnected.") // TODO
}

func (cl *ClientHandler) OnMessage(s *melody.Session, msg []byte) {
	cl.Broadcast(s.Request.RemoteAddr + ": " + html.EscapeString(string(msg)))
}

func (cl *ClientHandler) Broadcast(msg string) {
	cl.melody.Broadcast([]byte(msg + "\n"))
}

func (cl *ClientHandler) BroadcastExcept(msg string, s *melody.Session) {
	cl.melody.BroadcastOthers([]byte(msg+"\n"), s)
}

func (cl *ClientHandler) BroadcastFilter(msg string, fn func(*melody.Session) bool) {
	cl.melody.BroadcastFilter([]byte(msg+"\n"), fn)
}
