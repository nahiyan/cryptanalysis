package services

import (
	"fmt"
	"os"
	"text/template"
)

func (slurmSvc *SlurmService) GenerateSlurmJob(nodes, cpuCores, memory int, timeout int, command string) {
	errorSvc := slurmSvc.errorSvc
	randomSvc := slurmSvc.randomSvc
	config := slurmSvc.configSvc.Config
	commandTmpl := "#!/bin/bash\n#SBATCH --nodes={{.Nodes}}\n#SBATCH --cpus-per-task={{.CpuCores}}\n#SBATCH --mem={{.Memory}}M\n#SBATCH --time=00:{{.GlobalTimeout}}\n\n{{.Command}}\n"

	tmpl, err := template.New("tmpl").Parse(commandTmpl)
	errorSvc.Fatal(err, "Solver: failed to generate slurm job")

	// Random name
	name := randomSvc.RandString(10)

	fmt.Println(name)

	jobFile, err := os.OpenFile("/tmp/"+string(name)+".sh", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	errorSvc.Fatal(err, "Solver: failed to create the slurm job file")
	defer jobFile.Close()

	err = tmpl.Execute(jobFile, map[string]interface{}{
		"Nodes":         nodes,
		"CpuCores":      cpuCores,
		"Memory":        memory,
		"Timeout":       timeout,
		"GlobalTimeout": timeout + 5, // 5 extra seconds to gracefully shutdown
		"BenchmarkBin":  config.Paths.Bin.Benchmark,
		"Command":       command,
	})

	errorSvc.Fatal(err, "Solver: failed to write the slurm job file")
}
