package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

var (
	timeout               time.Duration
	ErrArgsNumberMismatch = errors.New("number of arguments should be equal 2")
	ErrInvalidPort        = errors.New("port should be in range 0-65535")
	ErrUnableToConnect    = errors.New("unable to connect")
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", ErrArgsNumberMismatch)
		os.Exit(1)
	}
	host := args[0]
	port, err := checkPort(args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	ctx := context.Background()
	err = work(ctx, host, *port, timeout)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func work(ctx context.Context, host string, port uint16, timeout time.Duration) (err error) {
	addr := net.JoinHostPort(host, strconv.Itoa(int(port)))
	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		return ErrUnableToConnect
	}
	defer func() {
		err = client.Close()
	}()
	_, _ = fmt.Fprintf(os.Stderr, "...connected to %s\n", addr)
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		err := client.Send()
		if err != nil {
			if errors.Is(err, io.EOF) {
				_, _ = fmt.Fprintln(os.Stderr, "...EOF")
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "send error: %s\n", err)
			}
		}
		cancel()
	}()
	go func() {
		err := client.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				_, _ = fmt.Fprintln(os.Stderr, "...EOF")
			} else {
				_, _ = fmt.Fprintln(os.Stderr, "...connection was closed by peer")
			}
		}
		cancel()
	}()
	<-ctx.Done()
	return nil
}

func checkPort(s string) (*uint16, error) {
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return nil, ErrInvalidPort
	}
	port := uint16(n)
	return &port, nil
}
