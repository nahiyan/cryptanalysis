package filesystem

import "io"

type FilesystemService interface {
	CountLines(io.Reader) (int, error)
	ReadLine(io.Reader, int) (string, int, error)
	FileExists(string) bool
}
