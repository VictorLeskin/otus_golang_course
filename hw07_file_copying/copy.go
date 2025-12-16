package hw07_file_copying

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

var bufSize int = 1024

type CommanLineParameter struct {
	input, output string
	offset, limit int64
}

type IOCopyData struct {
	src           io.Reader
	dst           io.Writer
	offset, limit int64
	buf           []byte
	done          <-chan struct{}
	progress      <-chan int
}

func IOCopy(src io.Reader, dst io.Writer, offset, limit int64) error {
	buf := make([]byte, bufSize)

	for {
		// Read to buffer from a intput stream
		n, err := src.Read(buf)
		if n > 0 {
			// Write from buffer to output stream.
			if _, writeErr := dst.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
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

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Open input file.
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("Error opening input file: %s", fromPath)
	}
	defer srcFile.Close()

	// Open output file.
	dstFile, err := os.Open(toPath)
	if err != nil {
		return fmt.Errorf("Error opening output file: %s", toPath)
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
