package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const (
	DefaultBufferSize int = 1024
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

	pb ProgressBar

	progressChan chan int64 // to show a progreas
	cancelChan   chan error // to finish
}

func (cp *IOCopyData) seekStart() error {
	// good boy
	seeker, _ := cp.src.(io.Seeker)
	n, err := seeker.Seek(cp.offset, io.SeekStart)

	if err != nil && !errors.Is(err, io.EOF) {
		return err // pass up the nonEOF error.
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
	if err != nil && !errors.Is(err, io.EOF) {
		return err // pass up the nonEOF error.
	}
	if n < cp.offset {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

func (cp IOCopyData) getStreamSize() (int64, error) {
	// First: check is it a file.
	if file, ok := cp.src.(*os.File); ok {
		stat, err := file.Stat()
		if err == nil {
			return stat.Size(), nil
		}

		return 0, err
	}

	// Second: check is it supporting Seeker.
	if _, ok := cp.src.(io.Seeker); ok {
		// store current position.
		sz, err := cp.getSeekerSize()
		if err != nil {
			return sz, err
		}
	}

	return 0, nil
}

func (cp IOCopyData) getSeekerSize() (int64, error) {
	seeker, _ := cp.src.(io.Seeker)

	// store current position.
	current, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// go to the end.
	end, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// go back.
	_, err = seeker.Seek(current, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return end, nil
}

func (cp *IOCopyData) seek() error {
	if cp.offset != 0 {
		//  Try to cast to Seeker
		if _, ok := cp.src.(io.Seeker); ok {
			return cp.seekStart()
		}
		return cp.skipBytes()
	}

	return nil
}

func (cp *IOCopyData) BufferSize() int {
	if cp.bufSize == 0 {
		return DefaultBufferSize
	}
	return cp.bufSize
}

func (cp *IOCopyData) copyNoLimit() error {
	for {
		// Read to buffer from a input stream
		n, err := cp.src.Read(cp.buf)
		if n > 0 {
			// Write from buffer to output stream.
			if _, writeErr := cp.dst.Write(cp.buf[:n]); writeErr != nil {
				return writeErr
			}
			if cp.pb != nil {
				cp.progressChan <- int64(n)
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

func (cp *IOCopyData) copyLimit() error {
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
			if cp.pb != nil {
				cp.progressChan <- toWrite
			}
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

func (cp *IOCopyData) runProgressUpdater() {
	var total int64

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	updater := func() {
		// There is a modes when progres bar isn't active.
		if cp.pb != nil {
			cp.pb.Update(total)
			cp.pb.Render()
		}
	}

	for {
		select {
		case bytes, ok := <-cp.progressChan:
			// read from the channel till the last value and update when channel is empty
			if !ok {
				updater()
				return
			}
			total += bytes

			// update by time
		case <-ticker.C:
			updater()

		case <-cp.cancelChan:
			return
		}
	}
}

func (cp *IOCopyData) copy() error {
	// Init channels.
	cp.progressChan = make(chan int64, 100)
	cp.cancelChan = make(chan error, 1)

	// Start progress bar.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cp.runProgressUpdater()
	}()

	// start copying....
	cp.buf = make([]byte, cp.BufferSize())

	var err error
	if cp.limit == 0 {
		err = cp.copyNoLimit()
	} else {
		err = cp.copyLimit()
	}

	close(cp.progressChan)

	return err
}

func (cp *IOCopyData) setupProgressBar() error {
	copyCnt := cp.limit

	sz, err := cp.getStreamSize()
	if copyCnt == 0 {
		// Reading till the end of the stream.
		if err != nil {
			return err
		}
		// If the stream provides a size evalute the count for the progress indicator.
		if sz != 0 {
			copyCnt = sz - cp.offset
		}
	} else if err == nil {
		// Reading thje disired count of bytes. If it excees the rest  of the stream side, we should
		// copy only a tail of the stream.
		realCnt := sz - cp.offset
		copyCnt = min(copyCnt, realCnt)
	}

	// setup a progress bar to show a progress in percents.
	if copyCnt != 0 {
		cp.pb = NewTxtProgressBar(copyCnt, 100)
	}

	return nil
}

func (cp *IOCopyData) main() error {
	if err := cp.setupProgressBar(); err != nil {
		return err
	}

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
	dstFile, err := os.Create(toPath)
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
