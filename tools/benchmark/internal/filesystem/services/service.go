package services

import (
	"benchmark/internal/consts"
	"bufio"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
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

func (filesystemSvc *FilesystemService) FileExistsNonEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	fmt.Println(info.Size())

	return info.Size() != 0
}

func (filesystemSvc *FilesystemService) WriteFromPipe(pipe io.Reader, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		file.WriteString(scanner.Text() + "\n")
	}

	return nil
}

// TODO: Move this to the log module
func (filesystemSvc *FilesystemService) LogInfo(messages ...string) {
	filesystemSvc.Log(consts.Info, messages...)
}

func (filesystemSvc *FilesystemService) LogDebug(messages ...string) {
	filesystemSvc.Log(consts.Debug, messages...)
}

func (filesystemSvc *FilesystemService) Log(type_ string, messages ...string) {
	filePath := ""
	switch type_ {
	case consts.Info:
		filePath = "info.log"
	case consts.Debug:
		filePath = "debug.log"
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("FS log: failed to open log file of type " + type_)
	}

	if _, err := file.WriteString(strings.Join(messages, " ") + "\n"); err != nil {
		log.Println("FS log: failed to write log of type " + type_)
	}
}

func (filesystemSvc *FilesystemService) Checksum(filePath string) (string, error) {
	startTime := time.Now()
	defer filesystemSvc.LogInfo("Checksum: took", time.Since(startTime).String())

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return "", os.ErrNotExist
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
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

func (filesystemSvc *FilesystemService) PrepareTempDir() error {
	err := filesystemSvc.PrepareDir("tmp")
	return err
}
