package ratatat

import (
	"bufio"
	"bytes"
	"io"
	"time"
)

const (
	EOL          = '\r'
	WaitInterval = time.Second * 3
)

type Commander struct {
	dev     *bufio.ReadWriter
	lastCmd time.Time
}

// Creates a new Commander for communicating with an AT device
// over and readable/writable source
func NewCommander(dev io.ReadWriter) *Commander {
	return &Commander{
		bufio.NewReadWriter(bufio.NewReader(dev), bufio.NewWriter(dev)),
		time.Time{},
	}
}

func (c *Commander) readResponse() (string, error) {
	out := bytes.NewBuffer([]byte{})
	for {
		b, err := c.dev.ReadByte()
		if err != nil {
			return "", err
		}
		if b == EOL {
			break
		}
		out.WriteByte(b)
	}
	return out.String(), nil
}

// Enters command mode only if the last command was run over
// three seconds ago. Otherwise, short-circuits
func (c *Commander) EnterCommandMode() error {
	if time.Since(c.lastCmd) < WaitInterval {
		return nil
	}
	_, err := c.dev.Write([]byte("+++"))
	c.dev.Flush()
	if err != nil {
		return err
	}
	res, err := c.readResponse()
	if err != nil {
		return err
	}
	if res != "OK" {
		panic("Non-OK response from serial port")
	}
	return nil
}

// Sends an AT command and receives the response. Note that the outgoing
// command should not have a <CR> at the end
func (c *Commander) SendAndRecv(cmd string) (string, error) {
	c.EnterCommandMode()
	out := []byte(cmd[:len(cmd)-1] + "\r\n")
	_, err := c.dev.Write(out)
	c.dev.Flush()
	if err != nil {
		return "", err
	}
	res, err := c.readResponse()
	c.lastCmd = time.Now()
	return res, err
}
