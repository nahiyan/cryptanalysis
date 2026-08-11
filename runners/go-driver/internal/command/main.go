package command

import "os/exec"

type Group struct {
	Commands []*exec.Cmd
}
