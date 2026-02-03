package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type CommanLineParameter struct {
	host    string
	port    int
	timeout time.Duration
}

func Usage() {
	fmt.Println(
		`Реализация крайне примитивного TELNET клиента.
* Программа подключется к указанному хосту (IP или доменное имя) и порту по протоколу TCP.
* После подключения STDIN программы записыватся в сокет, а данные, полученные из сокета, выводятся в STDOUT.
* Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout) - по умолчанию 10s.
* При нажатии Ctrl+D программа закрывает сокет и завершается с сообщением.
* При получении SIGINT программа завершает свою работу.
* Если сокет закрылся со стороны сервера, то при следующей попытке отправить сообщение программа должна завершаться.
* При подключении к несуществующему серверу, программа завершается с ошибкой соединения/таймаута.`)

	flag.PrintDefaults()
}

func SetupCommadLineParameters() {
	flag.Usage = Usage
}

func parseCommandLine(args0 []string) (ret CommanLineParameter, err error) {
	fs := flag.NewFlagSet("privitive-telnet", flag.ContinueOnError)

	fs.DurationVar(&ret.timeout, "timeout", 10*time.Second, "connection timeout")
	err = fs.Parse(args0)
	if err != nil {
		return ret, fmt.Errorf("Error parsing command line parameters:\n%s", err.Error())
	}

	// ge host and port
	args := fs.Args()
	if len(args) < 2 {
		fs.Usage()
		return ret, fmt.Errorf("Host and port are required")
	}

	ret.host = args[0]
	if ret.host == "" {
		return ret, errors.New("Missed host address")
	}

	if net.ParseIP(ret.host) == nil { // Не IP
		return ret, errors.New("Wrong host address")
	}

	// Check port.
	ret.port, err = strconv.Atoi(args[1])
	if err != nil {
		return ret, fmt.Errorf("Port must be a number")
	}

	if ret.port < 1 || ret.port > 65535 {
		return ret, fmt.Errorf("Port number must be in range [1,65535]")
	}

	// Check result
	fmt.Printf("Host: %s\n", ret.host)
	fmt.Printf("Port: %d\n", ret.port)
	fmt.Printf("Timeout: %f\n", ret.timeout.Seconds())

	return ret, nil
}

func ParseCommandLine() (ret CommanLineParameter, err error) {
	return parseCommandLine(os.Args[1:])
}

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return nil
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.

type MyTelnetClient struct {
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
