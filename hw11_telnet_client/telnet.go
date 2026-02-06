package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	ctx, cancel := context.WithCancel(context.Background())
	return &MyTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		ctx:     ctx,
		cancel:  cancel,
		dialer:  net.DialTimeout,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.

type MyTelnetClient struct {
	TelnetClient

	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	ctx    context.Context
	cancel context.CancelFunc

	conn net.Conn

	// by default it is net.DialTimeout.
	// for testing purposes it can be replaced by a function to destroy Universe.
	dialer func(network, address string, timeout time.Duration) (net.Conn, error)
}

func (c *MyTelnetClient) Connect() error {
	conn, err := c.dialer("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	fmt.Printf("Connected to %s\n", c.address)
	return nil
}

func (c *MyTelnetClient) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			text := scanner.Text()
			_, err := c.conn.Write([]byte(text + "\n"))
			if err != nil {
				return fmt.Errorf("send error: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("input scanner error: %w", err)
	}

	// Ctrl+D - end of input.
	c.cancel()
	return nil
}

func (c *MyTelnetClient) Receive() error {
	// nonblocking reading.
	c.conn.SetReadDeadline(time.Time{}) // remove timeout.

	reader := bufio.NewReader(c.conn)
	buf := make([]byte, 1024)

	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			n, err := reader.Read(buf)
			if err != nil {
				c.cancel()
				if err == io.EOF {
					fmt.Println("Connection closed by server")
					return nil
				}
				return fmt.Errorf("server send error: %w", err)
			}

			if n > 0 {
				// Output received data.
				c.out.Write(buf[:n])
			}
		}
	}
}

func (c *MyTelnetClient) Close() error {
	c.cancel()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *MyTelnetClient) Run() error {
	// add processsing Ctrl+C.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n^C")
		c.cancel()
		// close immediately.
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	if err := c.Connect(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer c.Close()
	var wg sync.WaitGroup

	wg.Add(2)

	// start goroutines.
	go func() {
		c.Send()
	}()

	go func() {
		c.Receive()
	}()

	// end finishing WaitGroup.
	wg.Wait()

	return nil
}
