package wsgclient

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn   net.Conn
	closed bool
}

func check_error(err error) {
	if err != nil {
		if err != io.EOF {
			log.Fatal("ERROR ", err)
		} else {
			log.Fatal("Connection lost.")
		}
	}
}

func Connect(proto, addr string) *Client {
	c, err := net.Dial(proto, addr)
	check_error(err)
	client := &Client{c, false}

	return client
}

// Disconnect the client.
func (client *Client) Close() {
	if client.closed {
		return
	}
	client.conn.Close()
	client.closed = true
}

// Send a message to the client.
func (client *Client) Send(msg string) {
	if client.closed {
		return
	}
	_, err := client.conn.Write([]byte(msg + "\n"))
	check_error(err)
}

func (client *Client) read_loop() {
	buf := bufio.NewReader(client.conn)
	for {
		if client.closed {
			return
		}

		msg, err := buf.ReadString('\n')
		check_error(err)
		msg = strings.TrimSpace(msg)
		if len(msg) < 1 {
			continue
		}
		log.Println(msg)
	}
}

func (client *Client) write_loop() {
	buf := bufio.NewReader(os.Stdin)
	for {
		if client.closed {
			return
		}

		msg, err := buf.ReadString('\n')
		check_error(err)
		msg = strings.TrimSpace(msg)
		if len(msg) < 1 {
			continue
		}
		client.Send(msg + "\n")
	}
}

func (client *Client) Loop() {
	go client.read_loop()
	client.write_loop()
}
