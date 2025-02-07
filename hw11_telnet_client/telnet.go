package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address:   address,
		timeout:   timeout,
		out:       out,
		inScanner: bufio.NewScanner(in),
	}
}

type Telnet struct {
	address     string
	timeout     time.Duration
	conn        net.Conn
	out         io.Writer
	inScanner   *bufio.Scanner
	connScanner *bufio.Scanner
}

func (t *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err == nil {
		t.conn = conn
		t.connScanner = bufio.NewScanner(conn)
	}
	return err
}

func (t *Telnet) Close() error {
	if t.conn != nil {
		if err := t.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Telnet) Send() error {
	if !t.inScanner.Scan() {
		return errors.New("no data to send")
	}

	line := t.inScanner.Bytes()
	if _, err := t.conn.Write([]byte(fmt.Sprintf("%s\n", line))); err != nil {
		return err
	}

	return nil
}

func (t *Telnet) Receive() error {
	if !t.connScanner.Scan() {
		return errors.New("no data received")
	}

	line := t.connScanner.Bytes()
	if _, err := t.out.Write([]byte(fmt.Sprintf("%s\n", line))); err != nil {
		return err
	}
	return nil
}
