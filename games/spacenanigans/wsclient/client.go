package wsclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gobwas/ws"
)

var (
	gTimeout           = time.Duration(30) * time.Second
	gConnectionTimeout = time.Duration(2) * time.Minute
)

const (
	PacketUnknown int = iota
	PacketID
	PacketMap
	PacketConnect
	PacketDisconnect
	PacketMessage
	PacketMove
	PacketStopMove
	PacketSeeMob
	PacketHideMob
)

type Packet struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

type Client struct {
	conn      net.Conn
	logger    *log.Logger
	lastReply time.Time
	closeOnce sync.Once
	write     chan []byte
}

func New(conn net.Conn, logger *log.Logger) *Client {
	c := &Client{
		conn:      conn,
		logger:    logger,
		lastReply: time.Now(),
		write:     make(chan []byte, 100),
	}
	go c.writeLoop()
	return c
}

////////////////////////////////////////////////////////////////////////////////

func (c *Client) Log(msg string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Printf(msg+"\n", args...)
	}
}

func (c *Client) Read() (*Packet, error) {
	d := time.Now().Sub(c.lastReply)
	if d > gConnectionTimeout {
		return nil, fmt.Errorf("connection timeout")
	}

	b, err := c.readRaw()
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return nil, nil // Ignore timeouts
		}
		return nil, err
	}

	c.lastReply = time.Now()
	if b == nil {
		return nil, nil
	}
	var pkt *Packet
	err = c.unmarshal(b, &pkt)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %s", err)
	}
	return pkt, nil
}

func (c *Client) Write(t int, d interface{}) error {
	b, err := c.marshal(Packet{
		Type: t,
		Data: d,
	})
	if err != nil {
		return err
	}

	// Blocking operation (unless the channel is nil)
	select {
	case c.write <- b:
	default:
	}
	return nil
}

func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		c.conn.Close()
		// we first ensure no other goroutines will try to cause any
		// "send on closed channel" panics or blocked calls
		ch := c.write
		c.write = nil
		// then we close the channel properly, which will cause
		// the read loop goroutine to close
		close(ch)
	})
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (c *Client) marshal(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	return b, err
}

func (c *Client) unmarshal(b []byte, v interface{}) error {
	err := json.Unmarshal(b, &v)
	return err
}

func (c *Client) readRaw() ([]byte, error) {
	c.conn.SetReadDeadline(time.Now().Add(gTimeout))
	h, err := ws.ReadHeader(c.conn)
	if err != nil {
		return nil, err
	}

	switch h.OpCode {
	case ws.OpText:
		b := make([]byte, h.Length)
		c.conn.SetReadDeadline(time.Now().Add(gTimeout))
		_, err = io.ReadFull(c.conn, b)
		if err != nil {
			return nil, err
		}
		ws.Cipher(b, h.Mask, 0)
		return b, nil

	//case ws.OpBinary:

	case ws.OpClose:
		return nil, fmt.Errorf("closed")

	case ws.OpPing:
		c.conn.SetWriteDeadline(time.Now().Add(gTimeout))
		_, err = c.conn.Write(ws.CompiledPong)
		return nil, err

	case ws.OpPong:
		return nil, nil

	}
	return nil, fmt.Errorf("unknown opcode in ws frame: %d", h.OpCode)
}

func (c *Client) writeLoop() {
	ticker := time.NewTicker(gTimeout - time.Second)
	defer ticker.Stop()
	buf := bufio.NewWriter(c.conn)

	for {
		var b []byte
		select {
		case <-ticker.C:
			b = ws.CompiledPing
		case pkt, ok := <-c.write:
			if !ok {
				return
			}
			var err error
			b, err = ws.CompileFrame(ws.NewTextFrame(pkt))
			if err != nil {
				c.Log("ws frame error: %s", err)
				c.Close()
				return
			}
		}

		c.conn.SetWriteDeadline(time.Now().Add(gTimeout))
		//_, err := c.conn.Write(b)
		_, err := buf.Write(b)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				continue // TODO: obviously try send it again
			}
			//c.Log("write error: %s", err)
			c.Close()
			return
		}
		buf.Flush()
	}
}
