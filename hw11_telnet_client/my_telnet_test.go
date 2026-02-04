package main

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockReadCloser имитирует io.ReadCloser

type MockReadCloser1 struct {
	closed bool
	Data   []string
	pos    int
}

func (m *MockReadCloser1) Read(p []byte) (n int, err error) {
	if m.closed {
		return 0, fmt.Errorf("stream had closed")
	}
	if m.pos < len(m.Data) {
		p = []byte(m.Data[m.pos])
		m.pos++
		return len(p), nil
	}

	return 0, io.EOF
}

func (m *MockReadCloser1) Close() error {
	m.closed = true
	return nil
}

// MockWriter имитирует io.Writer и сохраняет записанные данные
type MockWriter1 struct {
	buffer string
}

func (m *MockWriter1) Write(p []byte) (n int, err error) {
	m.buffer += string(p)
	return len(p), nil
}

func (m *MockWriter1) Free() (ret string) {
	ret = m.buffer
	m.buffer = ""
	return ret
}

func Test_MockReadCloser1_Ctor(t *testing.T) {
	var t0 MockReadCloser1
	assert.Equal(t, false, t0.closed)
	assert.Equal(t, 0, len(t0.Data))
	assert.Equal(t, 0, t0.pos)
}

func Test_NewTelnetClient(t *testing.T) {
	in := &MockReadCloser1{}
	out := &MockWriter1{}
	t0 := NewTelnetClient("1.2.3.4:5", 11*time.Second, in, out)

	c, ok := t0.(*MyTelnetClient)

	assert.True(t, ok)
	assert.Equal(t, "1.2.3.4:5", c.address)
	assert.Equal(t, 11*time.Second, c.timeout)

	c1, ok1 := c.in.(*MockReadCloser1)
	assert.True(t, ok1)
	assert.Equal(t, in, c1)

	c2, ok2 := c.out.(*MockWriter1)
	assert.True(t, ok2)
	assert.Equal(t, out, c2)

	assert.NotNil(t, c.ctx)
	assert.NotNil(t, c.cancel)
}

func Test_MyTelnetClient_Connect(t *testing.T) {
	t0 := MyTelnetClient{
		address: "1.2.3.4:5",
		timeout: 1 * time.Second,
		in:      nil,
		out:     nil}
	t0.Connect()
}

//func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
//	// Place your code here.
//	ctx, cancel := context.WithCancel(context.Background())
//	return &MyTelnetClient{
//		address: address,
//		timeout: timeout,
//		ctx:     ctx,
//		cancel:  cancel,
//	}
//}
