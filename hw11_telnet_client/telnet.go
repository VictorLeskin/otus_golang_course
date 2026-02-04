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
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.

type MyTelnetClient struct {
	TelnetClient
	address string
	timeout time.Duration
	conn    net.Conn
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewMyTelnetClient(address string, timeout time.Duration) *MyTelnetClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &MyTelnetClient{
		address: address,
		timeout: timeout,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (c *MyTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	fmt.Printf("Connected to %s\n", c.address)
	return nil
}

func (c *MyTelnetClient) Send() error {
	defer c.wg.Done()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			text := scanner.Text()
			_, err := fmt.Fprintf(c.conn, "%s\n", text)
			if err != nil {
				return err
			}
		}
	}

	// Обработка Ctrl+D
	fmt.Println("^D")
	c.cancel()
	return nil
}

func (c *MyTelnetClient) Receive() error {
	defer c.wg.Done()

	reader := bufio.NewReader(c.conn)
	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					c.cancel()
					return nil
				}
				return err
			}
			fmt.Print(line)
		}
	}
}

func (c *MyTelnetClient) Close() error {
	c.cancel()
	c.wg.Wait()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *MyTelnetClient) Run() error {
	// Обработка Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n^C")
		c.cancel()
	}()

	if err := c.Connect(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer c.Close()

	c.wg.Add(2)

	go func() {
		if err := c.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "Send error: %v\n", err)
		}
	}()

	go func() {
		if err := c.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "Receive error: %v\n", err)
		}
	}()

	c.wg.Wait()
	return nil
}
