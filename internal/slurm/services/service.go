package services

import (
	"cryptanalysis/internal/slurm"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"text/template"
)

func (slurmSvc *SlurmService) GenerateJob(command string, numTasks, nodes, cpuCores, memory, timeout int) (string, error) {
	randomSvc := slurmSvc.randomSvc
	config := slurmSvc.configSvc.Config
	filePath := path.Join(config.Paths.Tmp, randomSvc.RandString(10)+".sh")
	jobFile, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer jobFile.Close()

	// Prepare the template
	commandTmpl := "#!/bin/bash\n#SBATCH --nodes={{.Nodes}}\n#SBATCH --cpus-per-task={{.CpuCores}}\n#SBATCH --mem={{.Memory}}M\n#SBATCH --time=00:{{.Timeout}}\n#SBATCH --array=1-{{.NumTasks}}\n\n{{.Command}}\n"
	tmpl, err := template.New("tmpl").Parse(commandTmpl)
	if err != nil {
		return "", err
	}

	timeout_ := math.Round(float64(timeout) * config.Slurm.WorkerTimeMul)
	if err := tmpl.Execute(jobFile, map[string]interface{}{
		"Nodes":    nodes,
		"CpuCores": cpuCores,
		"Memory":   memory,
		"Timeout":  timeout_,
		"Command":  command,
		"NumTasks": numTasks,
	}); err != nil {
		return "", err
	}

	return filePath, nil
}

func (slurmSvc *SlurmService) ScheduleJob(jobPath string, dependencies []slurm.Job) (int, error) {
	args := []string{}

	// Dependencies
	if len(dependencies) > 0 {
		args = append(args, "-d", "afterok:")
	}
	for _, dependency := range dependencies {
		args = append(args, strconv.Itoa(dependency.Id))
	}

	// TODO: Make account dynamic
	args = append(args, "--account=def-cbright")

	// Job path
	args = append(args, jobPath)
	output_, err := exec.Command("sbatch", args...).Output()
	if err != nil {
		return 0, err
	}
	output := string(output_)

	// Job ID
	var jobId int
	_, err = fmt.Sscanf(output, "Submitted batch job %d", &jobId)
	if err != nil {
		return 0, err
	}
	log.Println("Submitted job", jobId)

	return jobId, nil
}
