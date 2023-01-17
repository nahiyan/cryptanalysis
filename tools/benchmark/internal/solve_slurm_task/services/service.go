package services

import (
	solveslurmtask "benchmark/internal/solve_slurm_task"
	"crypto/sha1"
	"fmt"
	"log"
	"time"

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

func GenerateId(task solveslurmtask.Task) string {
	combination := task.Encoding + string(task.Solver) + task.Timeout.Round(time.Second).String()

	checksum := sha1.Sum([]byte(combination))
	return fmt.Sprintf("%x", checksum)
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) Add(task solveslurmtask.Task) error {
	data, err := solveSlurmTaskSvc.marshallingSvc.BinEncode(task)
	if err != nil {
		return err
	}

	id := GenerateId(task)
	err = solveSlurmTaskSvc.databaseSvc.Set(solveSlurmTaskSvc.Bucket, []byte(id), data)

	return err
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) AddMultiple(tasks []solveslurmtask.Task) error {
	keys := lo.Map(tasks, func(task solveslurmtask.Task, _ int) []byte {
		id := GenerateId(task)
		// fmt.Println(id)
		return []byte(id)
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

func (solveSlurmTaskSvc *SolveSlurmTaskService) Get(id string) (solveslurmtask.Task, error) {
	task := solveslurmtask.Task{}
	data, err := solveSlurmTaskSvc.databaseSvc.Get(solveSlurmTaskSvc.Bucket, []byte(id))
	if err != nil {
		return task, err
	}

	err = solveSlurmTaskSvc.marshallingSvc.BinDecode(data, &task)
	return task, err
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) Book() (*solveslurmtask.Task, string, error) {
	var (
		task   *solveslurmtask.Task
		taskId string
	)

	err := solveSlurmTaskSvc.databaseSvc.FindAndReplace(solveSlurmTaskSvc.Bucket, func(key, value []byte) []byte {
		var task_ solveslurmtask.Task
		err := solveSlurmTaskSvc.marshallingSvc.BinDecode(value, &task_)
		if err != nil {
			return nil
		}

		if task_.Booked {
			return nil
		}

		task_.Booked = true
		task_.PingTime = time.Now()
		encodedTask, err := solveSlurmTaskSvc.marshallingSvc.BinEncode(task_)
		if err != nil {
			return nil
		}

		task = new(solveslurmtask.Task)
		*task = solveslurmtask.Task(task_)
		taskId = string(key)

		return encodedTask
	})

	return task, taskId, err
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) Remove(id string) error {
	err := solveSlurmTaskSvc.databaseSvc.Remove(solveSlurmTaskSvc.Bucket, []byte(id))
	return err
}
