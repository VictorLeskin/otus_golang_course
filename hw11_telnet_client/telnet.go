package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
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

func ParseCommadLine() (ret CommanLineParameter, err error) {
	flag.StringVar(&ret.host, "host", "", "server host (required)")
	flag.IntVar(&ret.port, "port", -1, "server port (required)")
	ret.timeout = *flag.Duration("timeout", 10*time.Second, "connection timeout")

	flag.Parse()

	// ge host and port
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		return ret, fmt.Errorf("Host and port are required")
	}

	ret.host = args[0]
	if ret.host == "" {
		return ret, errors.New("Missed host name")
	}

	// Check port.
	ret.port, err = strconv.Atoi(args[1])
	if err != nil {
		return ret, fmt.Errorf("Port must be a number\n")
	}

	// Check result
	fmt.Printf("Host: %s\n", ret.host)
	fmt.Printf("Port: %d\n", ret.port)
	fmt.Printf("Timeout: %v\n", ret.timeout)

	return ret, nil
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
