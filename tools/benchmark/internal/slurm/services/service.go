package services

import (
	"benchmark/internal/slurm"
	"os"
	"strconv"
	"text/template"
)

type Properties struct {
	Bucket string
}

func (slurmSvc *SlurmService) Init() {
	slurmSvc.Bucket = "tasks"
}

func (slurmSvc *SlurmService) RemoveAll() error {
	err := slurmSvc.databaseSvc.RemoveAll(slurmSvc.Bucket)
	return err
}

func (slurmSvc *SlurmService) AddTask(id int, task slurm.Task) error {
	data, err := slurmSvc.marshallingSvc.BinEncode(task)
	if err != nil {
		return err
	}

	err = slurmSvc.databaseSvc.Set(slurmSvc.Bucket, []byte(strconv.Itoa(id)), data)
	return err
}

func (slurmSvc *SlurmService) GetTask(id int) (slurm.Task, error) {
	task := slurm.Task{}
	data, err := slurmSvc.databaseSvc.Get(slurmSvc.Bucket, []byte(strconv.Itoa(id)))
	if err != nil {
		return task, err
	}

	err = slurmSvc.marshallingSvc.BinDecode(data, &task)
	return task, err
}

func (slurmSvc *SlurmService) GenerateJob(numTasks, maxConcurrentTasks, nodes, cpuCores, memory, timeout int, command string) (string, error) {
	randomSvc := slurmSvc.randomSvc
	config := slurmSvc.configSvc.Config
	commandTmpl := "#!/bin/bash\n#SBATCH --nodes={{.Nodes}}\n#SBATCH --cpus-per-task={{.CpuCores}}\n#SBATCH --mem={{.Memory}}M\n#SBATCH --time=00:{{.GlobalTimeout}}\n#SBATCH --array=1-{{.NumTasks}}%{{.MaxConcurrentTasks}}\n\n{{.Command}}\n"

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
		"Nodes":              nodes,
		"CpuCores":           cpuCores,
		"Memory":             memory,
		"Timeout":            timeout,
		"GlobalTimeout":      timeout + 5, // 5 extra seconds to gracefully shutdown
		"BenchmarkBin":       config.Paths.Bin.Benchmark,
		"Command":            command,
		"NumTasks":           numTasks,
		"MaxConcurrentTasks": maxConcurrentTasks,
	}); err != nil {
		return "", err
	}

	return filePath, nil
}
