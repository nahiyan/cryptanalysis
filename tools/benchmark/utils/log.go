package utils

import (
	"benchmark/constants"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

func AppendLog(satSolver string, instanceName string, duration time.Duration, messages []string, exitCode int, validity string) {
	// Open the file
	filePath := path.Join(constants.LogsDirPath, "logs.csv")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to write logs: " + err.Error())
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{
		fmt.Sprintf("%s_%s", satSolver, instanceName),
		fmt.Sprintf("%.2f", duration.Seconds()),
		strings.Join(messages, "; "),
		fmt.Sprintf("%d", exitCode),
		validity,
	})
	csvWriter.Flush()

	file.Close()
}
