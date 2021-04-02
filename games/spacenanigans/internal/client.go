package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lmas/spacenanigans/wsclient"
)

type Client struct {
	conn *wsclient.Client

	addr      string
	id        int64
	char      *Character
	position  Vector
	velocity  Vector
	direction float64
}

func (c *Client) Connect(w http.ResponseWriter, r *http.Request) error {
	addr := webClientAddr(r)
	if c.IsConnected() {
		return fmt.Errorf("client already connected: %s", addr)
	}
	conn, _, _, err := gUpgrader.Upgrade(r, w)
	if err != nil {
		// the upgrader already sent a http error to the client
		// TODO: NOT SURE ANYMORE
		//Log("ws upgrade error: %s", err)
		return fmt.Errorf("ws upgrade error: %s", err)
	}
	c.addr = addr
	c.conn = wsclient.New(conn, logger)
	Log("%s connected from %s", c.char, c.addr)
	return nil
}

func (c *Client) Disconnect() {
	Log("%s disconnected from %s", c.char, c.addr)
	c.conn.Close()
	c.addr = ""
	c.conn = nil
}

func (c *Client) IsConnected() bool {
	return c.conn != nil
}

func (c *Client) Read() (*wsclient.Packet, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("disconnected")
	}
	return c.conn.Read()
}

func (c *Client) Write(t int, data interface{}) error {
	if !c.IsConnected() {
		return fmt.Errorf("disconnected")
	}
	return c.conn.Write(t, data)
}

////////////////////////////////////////////////////////////////////////////////

func (s *Server) newID() int64 {
	// No need to make this func thread safe. it should only be called from
	// addClinet(), which is already made thread safe
	s.lastID += 1
	return s.lastID
}

func (s *Server) newClient(char *Character) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i := range s.clients {
		if s.clients[i].char.Username == char.Username {
			return false
		}
	}

	c := &Client{
		id:        s.newID(),
		char:      char,
		position:  Vector{float64(s.world.SpawnX), float64(s.world.SpawnY)},
		direction: 0,
	}
	s.clients = append(s.clients, c)
	return true
}

func (s *Server) getClient(user string) *Client {
	// TODO: this might be the place to check if client was kicked/banned,
	// since it's run by both charPlay() and connectClient().
	var c *Client
	s.lock.RLock()
	for i := range s.clients {
		if s.clients[i].char.Username == user {
			c = s.clients[i]
			break
		}
	}
	s.lock.RUnlock()
	return c
}

func (s *Server) getClients() []*Client {
	s.lock.RLock()
	clients := s.clients
	s.lock.RUnlock()
	return clients
}

////////////////////////////////////////////////////////////////////////////////

const maxMsgSize = 200

func (s *Server) connectClient(w http.ResponseWriter, r *http.Request) error {
	user := getCookie(r)
	c := s.getClient(user)
	if c == nil {
		return fmt.Errorf("client not logged in")
	}
	err := c.Connect(w, r)
	if err != nil {
		return err
	}

	c.Write(wsclient.PacketID, c.id)
	c.Write(wsclient.PacketMap, s.world)
	nearby := s.playersNearby(c)
	for _, n := range nearby {
		c.Write(wsclient.PacketSeeMob, map[string]interface{}{
			"id":   float64(n.id),
			"name": n.char.FullName(),
			"x":    n.position.X,
			"y":    n.position.Y,
			"d":    n.direction,
			"m":    n.char.IsMale,
		})
		// TODO: skip this here and do it inside addClient()
		n.Write(wsclient.PacketSeeMob, map[string]interface{}{
			"id":   float64(c.id),
			"name": c.char.FullName(),
			"x":    c.position.X,
			"y":    c.position.Y,
			"d":    c.direction,
			"m":    c.char.IsMale,
		})
	}

	for {
		pkt, err := c.Read()
		if err != nil {
			//Log("%s read error: %s", c.user, err)
			break
		}
		if pkt == nil {
			continue
		}
		err = s.parse(c, pkt)
		if err != nil {
			Log("%s parse error: %s", c.char, err)
			break
		}
	}

	c.Disconnect()
	return nil
}

func (s *Server) sendNearby(c *Client, t int, data interface{}) {
	nearby := s.playersNearby(c)
	for _, n := range nearby {
		n.Write(t, data)
	}
}

func (s *Server) sendAll(t int, data interface{}) {
	clients := s.getClients()
	for _, c := range clients {
		c.Write(t, data)
	}
}

func (s *Server) parse(c *Client, pkt *wsclient.Packet) error {
	switch pkt.Type {
	case wsclient.PacketMessage:
		msg := pkt.Data.(string) // TODO: check for error?
		l := len(msg)
		if l > maxMsgSize {
			l = maxMsgSize
		}
		msg = strings.TrimSpace(msg[:l])
		if msg == "" {
			break
		}
		Log("%s: %s", c.char, msg)
		s.sendAll(wsclient.PacketMessage, map[string]interface{}{
			"id":   c.id,
			"name": c.char.FullName(),
			"text": msg,
		})

	case wsclient.PacketMove:
		dir, ok := pkt.Data.(float64)
		if ok {
			switch dir {
			case 0: // down
				c.velocity = Vector{0, 1}
				c.direction = dir
			case 1: // up
				c.velocity = Vector{0, -1}
				c.direction = dir
			case 2: // right
				c.velocity = Vector{1, 0}
				c.direction = dir
			case 3: // left
				c.velocity = Vector{-1, 0}
				c.direction = dir
			}
		}

	case wsclient.PacketStopMove:
		c.velocity = Vector{}

	default:
		return fmt.Errorf("unknown packet type: %v", pkt.Type)
	}
	return nil
}
