package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const (
	DEFAULT_BUFFER_SIZE int = 1024
)

type CommanLineParameter struct {
	input, output string
	offset, limit int64
}

type IOCopyData struct {
	src           io.Reader
	dst           io.Writer
	offset, limit int64
	bufSize       int
	buf           []byte

	//	done          <-chan struct{}
	//	progress      <-chan int
}

func (cp *IOCopyData) seekStart() error {
	// good boy
	seeker, _ := cp.src.(io.Seeker)
	n, err := seeker.Seek(cp.offset, io.SeekStart)

	if err != nil && err != io.EOF {
		return err
	}
	if n < cp.offset {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

func (cp *IOCopyData) skipBytes() error {
	r := cp.src
	// copy offset bytest to the Discard ()
	n, err := io.CopyN(io.Discard, r, cp.offset)
	if err != nil && err != io.EOF {
		return err
	}
	if n < cp.offset {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

func (cp *IOCopyData) seek() error {
	if cp.offset != 0 {
		//  Try to cast to Seeker
		if _, ok := cp.src.(io.Seeker); ok {
			return cp.seekStart()
		} else {
			return cp.skipBytes()
		}
	}

	return nil
}

func (cp *IOCopyData) BufferSize() int {
	if cp.bufSize == 0 {
		return DEFAULT_BUFFER_SIZE
	} else {
		return cp.bufSize
	}
}

func (cp *IOCopyData) copy() error {
	cp.buf = make([]byte, cp.BufferSize())

	for {
		// Read to buffer from a input stream
		n, err := cp.src.Read(cp.buf)
		if n > 0 {
			// Write from buffer to output stream.
			toWrite := min(int64(n), cp.limit)
			if _, writeErr := cp.dst.Write(cp.buf[:toWrite]); writeErr != nil {
				return writeErr
			}
			cp.limit -= toWrite
		}

		if cp.limit == 0 {
			break
		}

		if err != nil {
			if err == io.EOF {
				break // Hurra
			}
			return err // Ups....
		}
	}

	return nil
}

func (cp *IOCopyData) main() error {
	if err := cp.seek(); err != nil {
		return err
	}
	if err := cp.copy(); err != nil {
		return err
	}

	return nil
}

func IOCopy(src io.Reader, dst io.Writer, offset, limit int64) error {
	cp := IOCopyData{
		src:    src,
		dst:    dst,
		offset: offset,
		limit:  limit,
	}

	return cp.main()
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Open input file.
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %s", fromPath)
	}
	defer srcFile.Close()

	// Open output file.
	dstFile, err := os.Open(toPath)
	if err != nil {
		return fmt.Errorf("error opening output file: %s", toPath)
	}
	defer dstFile.Close()

	return IOCopy(srcFile, dstFile, offset, limit)
}

func Usage() {
	fmt.Println("Copy a file or a part of it.")
	flag.PrintDefaults()
}

func SetupCommadLineParameters() {
	flag.Usage = Usage
}

func ParseCommadLine() (ret CommanLineParameter, err error) {

	// Copy a part of a file according to the operands.
	flag.StringVar(&ret.input, "from", "", "file to read from")
	flag.StringVar(&ret.output, "to", "", "file to copy")
	flag.Int64Var(&ret.offset, "offset", 0, "skip offset bytes at start of output")
	flag.Int64Var(&ret.limit, "limit", 0, "copy only 'limit' bytes")

	flag.Parse() // проанализировать аргументы

	if ret.input == "" {
		return ret, errors.New("there is not name of the file to read from")
	}
	if ret.output == "" {
		return ret, errors.New("there is not name of the file to copy")
	}
	return ret, nil
}
