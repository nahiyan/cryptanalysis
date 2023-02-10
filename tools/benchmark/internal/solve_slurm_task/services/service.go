package services

import (
	encoder "benchmark/internal/encoder/services"
	solveslurmtask "benchmark/internal/solve_slurm_task"
	"crypto/sha1"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

type Properties struct {
	Bucket string
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) Init() {
	solveSlurmTaskSvc.Bucket = "solve_slurm_tasks"
	gob.Register(encoder.EncoderService{})
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) RemoveAll() error {
	err := solveSlurmTaskSvc.databaseSvc.RemoveAll(solveSlurmTaskSvc.Bucket)
	return err
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) GenerateId(task solveslurmtask.Task) string {
	combination := task.Encoding.BasePath + string(task.Solver) + task.Timeout.Round(time.Second).String()

	checksum := sha1.Sum([]byte(combination))
	return fmt.Sprintf("%x", checksum)
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) Add(task solveslurmtask.Task) error {
	data, err := solveSlurmTaskSvc.marshallingSvc.BinEncode(task)
	if err != nil {
		return err
	}

	id := solveSlurmTaskSvc.GenerateId(task)
	err = solveSlurmTaskSvc.databaseSvc.Set(solveSlurmTaskSvc.Bucket, []byte(id), data)

	return err
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) AddMultiple(tasks []solveslurmtask.Task) error {
	startTime := time.Now()
	defer solveSlurmTaskSvc.filesystemSvc.LogInfo("Solve slurm task: add multiple took", time.Since(startTime).String())

	keys := lo.Map(tasks, func(task solveslurmtask.Task, _ int) []byte {
		id := solveSlurmTaskSvc.GenerateId(task)
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

func (solveSlurmTaskSvc *SolveSlurmTaskService) Book() (mo.Option[solveslurmtask.Task], string, error) {
	task := solveslurmtask.Task{}
	taskId := ""

	startTime := time.Now()
	defer solveSlurmTaskSvc.filesystemSvc.LogInfo("Solve slurk task: book took", time.Since(startTime).String())

	key, value, err := solveSlurmTaskSvc.databaseSvc.Consume(solveSlurmTaskSvc.Bucket)
	if err == os.ErrNotExist {
		return mo.None[solveslurmtask.Task](), taskId, nil
	}
	if err != nil {
		return mo.None[solveslurmtask.Task](), taskId, err
	}

	err = solveSlurmTaskSvc.marshallingSvc.BinDecode(value, &task)
	if err != nil {
		return mo.None[solveslurmtask.Task](), taskId, err
	}
	taskId = fmt.Sprintf("%x", key)

	return mo.Some(task), taskId, err
}

func (solveSlurmTaskSvc *SolveSlurmTaskService) Remove(id string) error {
	err := solveSlurmTaskSvc.databaseSvc.Remove(solveSlurmTaskSvc.Bucket, []byte(id))
	return err
}
