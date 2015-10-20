package wsgserver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Client struct {
	// Addr contains the client's network address.
	Addr   string
	conn   net.Conn
	closed bool
}

var Clients []*Client

func check_client_error(err error, client *Client) {
	if err != nil {
		if err != io.EOF {
			log.Printf("ERROR %s: %s", client.Addr, err)
		}
		client.Close()
	}
}

// Start up a new client, using an existing socket that's already connected.
func NewClient(c net.Conn) *Client {
	addr := c.RemoteAddr().String()
	client := &Client{addr, c, false}
	Clients = append(Clients, client)
	return client
}

// Disconnect the client.
func (client *Client) Close() {
	if client.closed {
		return
	}
	client.conn.Close()
	client.closed = true

	// Remove the client from the list
	for i, c := range Clients {
		if c == client {
			Clients = append(Clients[:i], Clients[i+1:]...)
			break
		}
	}

	SendAll(fmt.Sprintf("%s disconnected.\n", client.Addr))
}

// Send a message to the client.
func (client *Client) Send(msg string) {
	if client.closed {
		return
	}
	_, err := client.conn.Write([]byte(msg + "\n"))
	check_client_error(err, client)
}

// Send a message to all available clients.
func SendAll(msg string) {
	msg = strings.TrimSpace(msg)
	if len(msg) < 1 {
		return
	}

	if VERBOSE {
		log.Printf("ALL: %s", msg)
	}
	for _, client := range Clients {
		client.Send(msg)
	}
}

func (client *Client) Loop() {
	buf := bufio.NewReader(client.conn)
	for {
		if client.closed {
			return
		}

		msg, err := buf.ReadString('\n')
		check_client_error(err, client)
		msg = strings.TrimSpace(msg)
		if len(msg) < 1 {
			continue
		}

		if VERBOSE {
			log.Printf("IN %s: %s", client.Addr, msg)
		}
		SendAll(fmt.Sprintf("%s: %s", client.Addr, msg))
	}
}
