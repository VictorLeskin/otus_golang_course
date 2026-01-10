package main

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testIOCopyData struct {
	IOCopyData
}

func Test_IOCopyData_ctor(t *testing.T) {
	var t0 testIOCopyData

	assert.Equal(t, nil, t0.src)
	assert.Equal(t, nil, t0.dst)
	assert.Equal(t, int64(0), t0.offset)
	assert.Equal(t, int64(0), t0.limit)
	assert.Equal(t, 0, len(t0.buf))
}

// MockSeeker для тестирования.
type MockSeeker struct {
	realOffset int64
	err        error
}

// MockSeeker для тестирования.
type MockReader struct {
	realRead int
	err      error
}

func (m MockSeeker) Seek(offset int64, whence int) (int64, error) {
	_ = offset
	if whence != io.SeekStart {
		return 0, errors.New("the unsupported seek parameter: expected a offset only from the start")
	}
	return m.realOffset, m.err
}

func (m MockSeeker) Read(p []byte) (int, error) {
	_ = p
	// io.Reader interface
	return 0, io.EOF
}

func (m MockReader) Read(p []byte) (int, error) {
	_ = p
	// io.Reader interface
	return m.realRead, m.err
}

func Test_IOCopyData_seekStart(t *testing.T) {
	var t0 testIOCopyData

	{
		t0.offset = 5
		m := MockSeeker{realOffset: 5, err: nil}
		t0.src = m

		err := t0.seekStart()
		assert.Nil(t, err)
	}

	{
		t0.offset = 5
		m := MockSeeker{realOffset: 5, err: errors.New("Upsssss")}
		t0.src = m

		err := t0.seekStart()
		assert.NotNil(t, err)
	}

	{
		t0.offset = 5
		m := MockSeeker{realOffset: 5, err: io.EOF}
		t0.src = m

		err := t0.seekStart()
		assert.Nil(t, err)
	}

	{
		t0.offset = 5
		m := MockSeeker{realOffset: 4, err: nil}
		t0.src = m

		err := t0.seekStart()
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
	}
}

func Test_IOCopyData_skipBytes(t *testing.T) {
	var t0 testIOCopyData

	{
		t0.offset = 5
		m := MockReader{realRead: 5, err: nil}
		t0.src = m

		err := t0.skipBytes()
		assert.Nil(t, err)
	}

	{
		t0.offset = 5
		m := MockReader{realRead: 5, err: errors.New("Upsssss")}
		t0.src = m

		err := t0.skipBytes()
		assert.Nil(t, err) // io.CopyN suppress the io.Reader error in such cases
	}

	{
		t0.offset = 5
		m := MockReader{realRead: 5, err: io.EOF}
		t0.src = m

		err := t0.skipBytes()
		assert.Nil(t, err) // io.CopyN suppress the io.Reader error in such cases
	}

	{
		t0.offset = 5
		m := MockReader{realRead: 4, err: io.EOF}
		t0.src = m

		err := t0.skipBytes()
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
	}
}

func Test_IOCopyData_seek(t *testing.T) {
	{
		var t0 testIOCopyData
		t0.offset = 5
		m := MockSeeker{realOffset: 5, err: io.EOF}
		t0.src = m

		err := t0.seek()
		assert.Nil(t, err)
	}

	{
		var t0 testIOCopyData
		t0.offset = 5
		m := MockSeeker{realOffset: 4, err: io.EOF}
		t0.src = m

		err := t0.seek()
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
	}

	{
		var t0 testIOCopyData
		t0.offset = 5
		m := MockReader{realRead: 5, err: io.EOF}
		t0.src = m

		err := t0.seek()
		assert.Nil(t, err)
	}

	{
		var t0 testIOCopyData
		t0.offset = 5
		m := MockReader{realRead: 4, err: io.EOF}
		t0.src = m

		err := t0.seek()
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
	}
}

// MockSeeker для тестирования.
type MockReader1 struct {
	size   int
	readed *int
}

// MockSeeker для тестирования.
type MockWriter1 struct {
	buffer   []byte
	capacity int
}

func (m *MockReader1) Init(sz int) {
	m.size = sz
	m.readed = new(int)
}

func (m MockReader1) Read(p []byte) (int, error) {
	// io.Reader interface
	sz := len(p)
	if sz == 0 {
		panic("MockReader1.Read: zero buffer size")
	}

	for i := 0; i < sz; i++ {
		if *m.readed < m.size {
			p[i] = byte(*m.readed)
			*m.readed++
		} else {
			return i, io.EOF
		}
	}
	return sz, nil
}

func (m *MockWriter1) Write(p []byte) (n int, err error) {
	m.buffer = append(m.buffer, p...)
	return len(p), nil
}

func Test_IOCopyData_BufSize(t *testing.T) {
	var t0 testIOCopyData
	assert.Equal(t, 1024, t0.BufferSize())

	t0.bufSize = 99
	assert.Equal(t, 99, t0.BufferSize())
}

func Test_IOCopyData_copyLimit(t *testing.T) {
	var t0 testIOCopyData
	t0.limit = 25
	m := MockReader1{size: 30}
	m.Init(30)
	d := MockWriter1{capacity: 100}
	t0.src = m
	t0.bufSize = 10
	t0.dst = &d
	t0.buf = make([]byte, 10)

	t0.copyLimit()

	assert.Equal(t, 25, len(d.buffer))
	for i := 0; i < 25; i++ {
		assert.Equal(t, byte(i), d.buffer[i])
	}
}

func Test_IOCopyData_copyNoLimit(t *testing.T) {
	var t0 testIOCopyData
	t0.limit = 0
	m := MockReader1{size: 30}
	m.Init(30)
	d := MockWriter1{capacity: 100}
	t0.src = m
	t0.bufSize = 10
	t0.dst = &d
	t0.buf = make([]byte, 10)

	t0.copyNoLimit()

	assert.Equal(t, 30, len(d.buffer))
	for i := 0; i < 30; i++ {
		assert.Equal(t, byte(i), d.buffer[i])
	}
}

func Test_IOCopyData_copy(t *testing.T) {
	// copy with limits.
	{
		var t0 testIOCopyData
		t0.limit = 25
		m := MockReader1{size: 30}
		m.Init(30)
		d := MockWriter1{capacity: 100}
		t0.src = m
		t0.bufSize = 10
		t0.dst = &d

		t0.copy()

		assert.Equal(t, 10, len(t0.buf))
		assert.Equal(t, 25, len(d.buffer))
		for i := 0; i < 25; i++ {
			assert.Equal(t, byte(i), d.buffer[i])
		}
	}

	// copy with no limit.
	{
		var t0 testIOCopyData
		t0.limit = 0
		m := MockReader1{size: 30}
		m.Init(30)
		d := MockWriter1{capacity: 100}
		t0.src = m
		t0.bufSize = 10
		t0.dst = &d

		t0.copy()

		assert.Equal(t, 10, len(t0.buf))
		assert.Equal(t, 30, len(d.buffer))
		for i := 0; i < 30; i++ {
			assert.Equal(t, byte(i), d.buffer[i])
		}
	}
}
