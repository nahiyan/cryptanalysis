package services

import (
	cubeslurmtask "benchmark/internal/cube_slurm_task"
	"strconv"
)

type Properties struct {
	Bucket string
}

func (cubeSlurmTaskSvc *CubeSlurmTaskService) Init() {
	cubeSlurmTaskSvc.Bucket = "cube_slurm_tasks"
}

func (cubeSlurmTaskSvc *CubeSlurmTaskService) RemoveAll() error {
	err := cubeSlurmTaskSvc.databaseSvc.RemoveAll(cubeSlurmTaskSvc.Bucket)
	return err
}

func (cubeSlurmTaskSvc *CubeSlurmTaskService) AddTask(id int, task cubeslurmtask.Task) error {
	data, err := cubeSlurmTaskSvc.marshallingSvc.BinEncode(task)
	if err != nil {
		return err
	}

	err = cubeSlurmTaskSvc.databaseSvc.Set(cubeSlurmTaskSvc.Bucket, []byte(strconv.Itoa(id)), data)
	return err
}

func (cubeSlurmTaskSvc *CubeSlurmTaskService) GetTask(id int) (cubeslurmtask.Task, error) {
	task := cubeslurmtask.Task{}
	data, err := cubeSlurmTaskSvc.databaseSvc.Get(cubeSlurmTaskSvc.Bucket, []byte(strconv.Itoa(id)))
	if err != nil {
		return task, err
	}

	err = cubeSlurmTaskSvc.marshallingSvc.BinDecode(data, &task)
	return task, err
}
