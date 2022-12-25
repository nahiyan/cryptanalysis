package services

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/samber/do"
)

type FilesystemService struct {
}

func NewFilesystemService(i *do.Injector) (*FilesystemService, error) {
	return &FilesystemService{}, nil
}

func (filesystemSvc *FilesystemService) CountLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func (filesystemSvc *FilesystemService) ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

func (filesystemSvc *FilesystemService) FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
