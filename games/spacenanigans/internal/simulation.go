package internal

import (
	"time"

	"github.com/lmas/spacenanigans/wsclient"
)

const (
	cTicksPerSec int     = 60
	cWalkSpeed   float64 = 6 // yeah.. walk.. speeeed...
)

var (
	cMobDistance Vector = Vector{11, 9} // seeing distance before mobs disappear
)

func (s *Server) Simulate() {
	ticker := time.NewTicker(time.Second / time.Duration(cTicksPerSec))
	defer ticker.Stop()
	last := time.Now()
	for {
		select {
		case t := <-ticker.C:
			delta := t.Sub(last)
			last = t
			s.Update(delta.Seconds())
		}
	}
}

func clientInList(c *Client, list []*Client) bool {
	for _, cc := range list {
		if c == cc {
			return true
		}
	}
	return false
}

func (s *Server) Update(delta float64) {
	clients := s.getClients()
	for _, c := range clients {
		newPos := s.newPosition(c, delta)
		if newPos == nil {
			continue
		}

		oldNeighbors := s.playersNearby(c)
		c.position = *newPos
		updates := s.playersNearby(c)

		for _, n := range oldNeighbors {
			if !clientInList(n, updates) {
				// c disappeard from n
				n.Write(wsclient.PacketHideMob, c.id)
				c.Write(wsclient.PacketHideMob, n.id)
				continue
			}
		}

		for _, n := range updates {
			n.Write(wsclient.PacketSeeMob, map[string]interface{}{
				"id":   float64(c.id),
				"name": c.char.FullName(),
				"x":    c.position.X,
				"y":    c.position.Y,
				"d":    c.direction,
				"m":    c.char.IsMale,
			})
			if n == c {
				// avoid sending update to the moving client,
				// twice, with it's own pos
				continue
			}
			c.Write(wsclient.PacketSeeMob, map[string]interface{}{
				"id":   float64(n.id),
				"name": n.char.FullName(),
				"x":    n.position.X,
				"y":    n.position.Y,
				"d":    n.direction,
				"m":    n.char.IsMale,
			})
		}
	}
}

func (s *Server) newPosition(c *Client, delta float64) *Vector {
	if !c.velocity.Zero() {
		vel := c.velocity.Mulf(cWalkSpeed).Mulf(delta)
		pos := c.position.Add(vel)
		if !s.world.IsBlocking(pos) {
			return &pos
		}
		c.velocity = Vector{}
	}
	return nil
}

//func (s *Server) playersNearby1(c *Client) []*Client {
//// Note: will always return itself too
//var nearby []*Client
//clients := s.getClients()
//for _, n := range clients {
//d := c.position.Distance(n.position)
//if d < cMobDistance {
//nearby = append(nearby, n)
//}
//}
//return nearby
//}

func (s *Server) playersNearby(c *Client) []*Client {
	// Note: will always return itself too
	var nearby []*Client
	clients := s.getClients()
	start := c.position.Sub(cMobDistance)
	stop := c.position.Add(cMobDistance)
	for _, n := range clients {
		if (n.position.X >= start.X && n.position.X <= stop.X) && (n.position.Y >= start.Y && n.position.Y <= stop.Y) {
			nearby = append(nearby, n)
		}
	}
	return nearby
}
