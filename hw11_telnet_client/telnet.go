package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *TelnetClient {
	return &TelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *TelnetClient) Close() error {
	return t.conn.Close()
}

func (t *TelnetClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *TelnetClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}
