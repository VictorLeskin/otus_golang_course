package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
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

func SetupCommandLineParameters() {
	flag.Usage = Usage
}

func parseCommandLine(args0 []string) (ret CommanLineParameter, err error) {
	fs := flag.NewFlagSet("privitive-telnet", flag.ContinueOnError)

	fs.DurationVar(&ret.timeout, "timeout", 10*time.Second, "connection timeout")
	err = fs.Parse(args0)
	if err != nil {
		return ret, fmt.Errorf("error parsing command line parameters:\n%s", err.Error())
	}

	// ge host and port
	args := fs.Args()
	if len(args) < 2 {
		fs.Usage()
		return ret, fmt.Errorf("host and port are required")
	}

	ret.host = args[0]
	if ret.host == "" {
		return ret, errors.New("missed host address")
	}

	// Check port.
	ret.port, err = strconv.Atoi(args[1])
	if err != nil {
		return ret, fmt.Errorf("port must be a number")
	}

	if ret.port < 1 || ret.port > 65535 {
		return ret, fmt.Errorf("port number must be in range [1,65535]")
	}

	return ret, nil
}

func ParseCommandLine() (ret CommanLineParameter, err error) {
	return parseCommandLine(os.Args[1:])
}
