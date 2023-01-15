package services

import (
	solveslurmtask "benchmark/internal/solve_slurm_task"
	"log"
	"strconv"

	"github.com/samber/lo"
)

type Properties struct {
	Bucket string
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) Init() {
	solveSlurmTaskSvc.Bucket = "solve_slurm_tasks"
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) RemoveAll() error {
	err := solveSlurmTaskSvc.databaseSvc.RemoveAll(solveSlurmTaskSvc.Bucket)
	return err
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) AddTask(id int, task solveslurmtask.Task) error {
	data, err := solveSlurmTaskSvc.marshallingSvc.BinEncode(task)
	if err != nil {
		return err
	}

	err = solveSlurmTaskSvc.databaseSvc.Set(solveSlurmTaskSvc.Bucket, []byte(strconv.Itoa(id)), data)

	return err
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) AddTasks(ids []int, tasks []solveslurmtask.Task) error {
	keys := lo.Map(ids, func(id, _ int) []byte {
		id_ := strconv.Itoa(id)
		return []byte(id_)
	})

	values := lo.Map(tasks, func(task solveslurmtask.Task, _ int) []byte {
		value, err := solveSlurmTaskSvc.marshallingSvc.BinEncode(task)
		if err != nil {
			log.Fatal(err)
		}

		return value
	})

	err := solveSlurmTaskSvc.databaseSvc.BatchSet(solveSlurmTaskSvc.Bucket, keys, values)
	if err != nil {
		return err
	}

	return nil
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) GetTask(id int) (solveslurmtask.Task, error) {
	task := solveslurmtask.Task{}
	data, err := solveSlurmTaskSvc.databaseSvc.Get(solveSlurmTaskSvc.Bucket, []byte(strconv.Itoa(id)))
	if err != nil {
		return task, err
	}

	err = solveSlurmTaskSvc.marshallingSvc.BinDecode(data, &task)
	return task, err
}
