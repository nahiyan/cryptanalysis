package services

import (
	"benchmark/internal/command"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
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

func (commandSvc *CommandService) Create(command string) *exec.Cmd {
	segments := strings.Fields(command)
	if len(segments) == 0 {
		logrus.Fatal("Command: empty")
	}

	cmd := exec.Command(segments[0], segments[1:]...)
	return cmd
}
