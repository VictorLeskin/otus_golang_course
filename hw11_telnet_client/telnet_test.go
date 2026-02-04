package main

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestMyTelnetClient(t *testing.T) {
	client := NewTelnetClient(
		"74.125.29.102:80",
		5*time.Second,
		nil,
		nil)

	client.Connect()
	client.Send()
	client.Receive()
	client.Close()
}

func TestTelnetClient_Connect(t *testing.T) {
	// Создаем mock сервер
	server := &MockTelnetServer{
		Port:     8080,
		Response: "Welcome!\n",
	}

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	// Ждем, чтобы сервер запустился
	time.Sleep(100 * time.Millisecond)

	// Создаем клиент
	client := NewTelnetClient(
		"localhost:8080",
		5*time.Second,
		NewMockReadCloser(""),
		NewMockWriter(),
	)

	// Пытаемся подключиться
	err := client.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	client.Close()
}

func TestTelnetClient_Send(t *testing.T) {
	server := &MockTelnetServer{
		Port:     8081,
		Response: "OK\n",
	}

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	time.Sleep(100 * time.Millisecond)

	// Mock входные данные
	input := NewMockReadCloser("Hello\nWorld\n")
	output := NewMockWriter()

	client := NewTelnetClient(
		"localhost:8081",
		5*time.Second,
		input,
		output,
	)

	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Запускаем Send в отдельной горутине
	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Send()
	}()

	// Ждем завершения Send (он завершится, когда входные данные закончатся)
	select {
	case err := <-errCh:
		if err != nil && err != io.EOF {
			t.Errorf("Send failed: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Send timeout")
	}

	// Проверяем, что сообщения дошли до сервера
	time.Sleep(100 * time.Millisecond)
	if len(server.Messages) < 2 {
		t.Errorf("Expected at least 2 messages, got %d", len(server.Messages))
	}
}

func TestTelnetClient_Receive(t *testing.T) {
	server := &MockTelnetServer{
		Port:     8082,
		Response: "Server response\n",
	}

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	time.Sleep(100 * time.Millisecond)

	input := NewMockReadCloser("")
	output := NewMockWriter()

	client := NewTelnetClient(
		"localhost:8082",
		5*time.Second,
		input,
		output,
	)

	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Запускаем Receive
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Receive()
	}()

	select {
	case err := <-errCh:
		// Receive может не завершиться, так как сервер продолжает отправлять данные
		if err != nil && err != context.Canceled {
			t.Errorf("Receive failed: %v", err)
		}
	case <-ctx.Done():
		// Ожидаемо - сервер продолжает отправлять данные
	}

	// Проверяем, что получили данные от сервера
	outputStr := output.String()
	if outputStr == "" {
		t.Error("Expected to receive data from server")
	}
}

func TestTelnetClient_CtrlD(t *testing.T) {
	server := &MockTelnetServer{Port: 8083}

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	time.Sleep(100 * time.Millisecond)

	// Ctrl+D - пустой ввод
	input := NewMockReadCloser("") // EOF сразу
	output := NewMockWriter()

	client := NewTelnetClient(
		"localhost:8083",
		5*time.Second,
		input,
		output,
	)

	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Send должен сразу вернуть EOF
	err := client.Send()
	if err != io.EOF && err != context.Canceled {
		t.Errorf("Expected EOF or Canceled, got: %v", err)
	}
}

func TestTelnetClient_ConnectionTimeout(t *testing.T) {
	// Не запускаем сервер - проверяем таймаут подключения

	client := NewTelnetClient(
		"localhost:9999",     // Несуществующий порт
		100*time.Millisecond, // Маленький таймаут
		NewMockReadCloser(""),
		NewMockWriter(),
	)

	err := client.Connect()
	if err == nil {
		t.Error("Expected connection timeout error")
	}

	// Проверяем, что это действительно таймаут
	if _, ok := err.(net.Error); !ok {
		t.Errorf("Expected net.Error, got: %T", err)
	}
}
