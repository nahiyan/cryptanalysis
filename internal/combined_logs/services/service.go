package services

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/samber/mo"
)

type LogFileMapping struct {
	Offset uint32
	Size   uint32
}

type Properties struct {
	LogFiles map[string]*LogFileMapping
}

func (combinedLogsSvc *CombinedLogsService) Generate(workers int) {
	errorSvc := combinedLogsSvc.errorSvc
	configSvc := combinedLogsSvc.configSvc

	combinedLogFile, err := os.OpenFile("all.clog", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	errorSvc.Fatal(err, "Combine: failed to create the combined logs file")
	defer combinedLogFile.Close()
	combinedLogWriter := bufio.NewWriter(combinedLogFile)
	lock := sync.Mutex{}

	pool := pond.New(workers, 1000, pond.IdleTimeout(100*time.Millisecond))
	files, err := os.ReadDir(configSvc.Config.Paths.Logs)
	errorSvc.Fatal(err, "Combine: failed to find the log files")
	numFiles := len(files)
	startTime := time.Now()
	for i, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if path.Ext(fileName) != ".log" {
			continue
		}
		pool.Submit(func(fileName string, i_ int) func() {
			filePath := path.Join(configSvc.Config.Paths.Logs, fileName)
			return func() {
				bytes, err := os.ReadFile(filePath)
				if err != nil {
					errorSvc.Fatal(err, "Combine: couldn't read "+filePath)
				}

				lock.Lock()
				_, err = combinedLogWriter.WriteString("FN:" + fileName + "\n")
				errorSvc.Fatal(err, "Combine: failed to write")
				_, err = combinedLogWriter.Write(bytes)
				errorSvc.Fatal(err, "Combine: failed to write")
				combinedLogWriter.Flush()
				lock.Unlock()
				log.Printf("Combine: [%d/%d] %s\n", i_+1, numFiles, fileName)
			}
		}(fileName, i))
	}
	pool.StopAndWait()
	log.Printf("Took %s to process %d files", time.Since(startTime), numFiles)
}

func (combinedLogsSvc *CombinedLogsService) Load() error {
	log.Println("CombinedLogs: Loading all.clog")
	startTime := time.Now()
	combinedLogsSvc.LogFiles = map[string]*LogFileMapping{}
	file, err := os.OpenFile("all.clog", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	bytesAccumulator := uint32(0)
	var currentFileSizePtr *uint32
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		line := string(lineBytes)
		lineSize := uint32(len(lineBytes) + 1)
		bytesAccumulator += lineSize
		if strings.HasPrefix(line, "FN:") {
			fileName := line[3:]
			entry := LogFileMapping{
				Offset: bytesAccumulator,
				Size:   -lineSize,
			}
			combinedLogsSvc.LogFiles[fileName] = &entry
			currentFileSizePtr = &entry.Size
		}
		if currentFileSizePtr != nil {
			*currentFileSizePtr += lineSize
		}
	}
	log.Printf("CombinedLogs: Load took %s", time.Since(startTime))

	return nil
}

func (combinedLogsSvc *CombinedLogsService) IsLoaded() bool {
	return len(combinedLogsSvc.LogFiles) > 0
}

func (combinedLogsSvc *CombinedLogsService) Get(name string) (mo.Option[string], error) {
	entry, exists := combinedLogsSvc.LogFiles[name]
	if !exists {
		return mo.None[string](), nil
	}

	file, err := os.OpenFile("all.clog", os.O_RDONLY, 0644)
	if err != nil {
		return mo.None[string](), err
	}

	data := make([]byte, entry.Size)
	_, err = file.ReadAt(data, int64(entry.Offset))
	if err != nil && !errors.Is(err, io.EOF) {
		return mo.None[string](), err
	}
	if errors.Is(err, io.EOF) {
		data = data[:len(data)-1]
	}

	return mo.Some(string(data)), nil
}
