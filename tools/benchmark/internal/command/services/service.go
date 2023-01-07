package services

import (
	"benchmark/internal/command"
	"os/exec"
)

func (commandSvc *CommandService) CreateGroup() *command.Group {
	return &command.Group{
		Commands: []*exec.Cmd{},
	}
}

func (commandSvc *CommandService) AddToGroup(group *command.Group, cmd *exec.Cmd) {
	group.Commands = append(group.Commands, cmd)
}

func (commandSvc *CommandService) StopGroup(group *command.Group) error {
	var err error
	for _, cmd := range group.Commands {
		err = cmd.Process.Kill()
	}

	return err
}
