package main

import (
	"benchmark/cmd"
	"benchmark/utils"
	"fmt"
	"os/exec"
)

func main() {
	cmd_ := exec.Command("bash", "-c", "rm results/*.log")
	if err := cmd_.Run(); err != nil {
		fmt.Println("Failed to remove the logs file", err.Error())
	}

	utils.AggregateLogs()

	cmd.Execute()
}
