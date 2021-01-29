package client

import (
	"bufio"
	"game/common/logger"
	"net"
)

const (
	writeBacklog = 16
)

type Message struct {
	data string
}

type Client struct {
	Connection net.Conn
	chWrite    chan *Message
	chClose    chan struct{}
}

func NewClient(c net.Conn) *Client {
	return &Client{
		Connection: c,
		chWrite:    make(chan *Message, writeBacklog),
		chClose:    make(chan struct{}),
	}
}

func (c *Client) BeginWrite() {
	defer func() {
		close(c.chWrite)
		c.Connection.Close()
	}()

	for c.chClose != nil {
		select {
		case msg := <-c.chWrite:
			_, err := c.Connection.Write([]byte(msg.data))
			if err != nil {
				panic(err)
			}
			logger.Debugf("[echo] %v", msg.data)
			break
		case <-c.chClose:
			logger.Info("close client")
			c.chClose = nil
			break
		default:
			break
		}
	}
	logger.Debugf("write exit normal")
}

func (c *Client) HandleConnection() {
	logger.Debugf("Serving %s\n", c.Connection.RemoteAddr().String())

	go c.BeginWrite()

	r := bufio.NewReader(c.Connection)
	for {
		netData, err := r.ReadString('\n')
		if err != nil {
			logger.Debug("begin close client:", err)
			close(c.chClose)
			break
		}

		result := string(netData)

		m := &Message{
			data: result,
		}
		c.chWrite <- m
	}
}
