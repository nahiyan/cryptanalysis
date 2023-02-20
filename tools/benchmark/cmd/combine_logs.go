package cmd

import (
	services2 "benchmark/internal/config/services"
	services1 "benchmark/internal/error/services"
	"benchmark/internal/injector"
	"bufio"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/spf13/cobra"
)

func initCombineLogsCmd() *cobra.Command {
	var workers int

	cmd := &cobra.Command{
		Use:   "combine-logs",
		Short: "Combine log files into one .clog file",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			errorSvc, err := services1.NewErrorService(injector)
			if err != nil {
				panic(err)
			}
			configSvc, err := services2.NewConfigService(injector)
			if err != nil {
				panic(err)
			}

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
		},
	}

	cmd.Flags().IntVarP(&workers, "workers", "w", 100, "Number of workers to read the log files in parallel")

	return cmd
}
