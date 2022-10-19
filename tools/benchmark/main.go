package main

import (
	"benchmark/cmd"
	"benchmark/utils"
)

func main() {
	utils.AggregateLogs()

	cmd.Execute()
}
