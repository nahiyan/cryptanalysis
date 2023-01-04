package services

import (
	"os"
	"text/template"
)

func (slurmSvc *SlurmService) GenerateJob(numTasks, nodes, cpuCores, memory, timeout int, command string) (string, error) {
	randomSvc := slurmSvc.randomSvc
	config := slurmSvc.configSvc.Config
	commandTmpl := "#!/bin/bash\n#SBATCH --nodes={{.Nodes}}\n#SBATCH --cpus-per-task={{.CpuCores}}\n#SBATCH --mem={{.Memory}}M\n#SBATCH --time=00:{{.GlobalTimeout}}\n#SBATCH --array=1-{{.NumTasks}}\n\n{{.Command}}\n"

	tmpl, err := template.New("tmpl").Parse(commandTmpl)
	if err != nil {
		return "", err
	}

	// Random name
	name := randomSvc.RandString(10)
	filePath := "/tmp/" + string(name) + ".sh"

	jobFile, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer jobFile.Close()

	if err := tmpl.Execute(jobFile, map[string]interface{}{
		"Nodes":         nodes,
		"CpuCores":      cpuCores,
		"Memory":        memory,
		"Timeout":       timeout,
		"GlobalTimeout": timeout + 5, // 5 extra seconds to gracefully shutdown
		"BenchmarkBin":  config.Paths.Bin.Benchmark,
		"Command":       command,
		"NumTasks":      numTasks,
	}); err != nil {
		return "", err
	}

	return filePath, nil
}
