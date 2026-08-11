package services

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

func (filesystemSvc *FilesystemService) CountLines(filePath string) (int, error) {
	reader, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return 0, err
	}

	lines, err := filesystemSvc.CountLinesFromReader(reader)
	if err != nil {
		return 0, err
	}

	return lines, nil
}

func (filesystemSvc *FilesystemService) CountLinesFromReader(r io.Reader) (int, error) {
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

func (filesystemSvc *FilesystemService) FileExistsNonEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return info.Size() != 0
}

func (filesystemSvc *FilesystemService) WriteFromPipe(pipe io.Reader, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// TODO: Test if it works for long lines
	_, err = io.Copy(file, pipe)
	if err != nil {
		return err
	}

	return nil
}

func (filesystemSvc *FilesystemService) PrepareDir(name string) error {
	if filesystemSvc.FileExists(name) {
		return nil
	}

	err := os.Mkdir(name, os.ModePerm)
	return err
}

func (filesystemSvc *FilesystemService) PrepareDirs(names []string) error {
	for _, name := range names {
		if err := filesystemSvc.PrepareDir(name); err != nil {
			return err
		}
	}

	return nil
}
