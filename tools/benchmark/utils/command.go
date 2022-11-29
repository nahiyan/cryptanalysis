package utils

import (
	"strings"
)

const (
	Pipe uint = iota
	Command
	Placeholder
)

type Segment struct {
	type_      uint
	components []string
}

type CommandStructure struct {
	segments []Segment
}

func NewCommand() *CommandStructure {
	return &CommandStructure{}
}

func (c *CommandStructure) AddPipe() *CommandStructure {
	c.segments = append(c.segments, Segment{
		type_: Pipe,
	})

	return c
}

func (c *CommandStructure) AddCommand(component string) *CommandStructure {
	return c.AddCommands([]string{component})
}

func (c *CommandStructure) AddCommands(components []string) *CommandStructure {
	c.segments = append(c.segments, Segment{
		type_:      Command,
		components: components,
	})

	return c
}

func (c *CommandStructure) AddPlaceholder() *CommandStructure {
	c.segments = append(c.segments, Segment{
		type_: Placeholder,
	})

	return c
}

func (c *CommandStructure) FillWithCommand(command string) *CommandStructure {
	return c.FillPlaceholder(Segment{
		type_:      Command,
		components: []string{command},
	})
}

func (c *CommandStructure) FillWithCommands(commands []string) *CommandStructure {
	return c.FillPlaceholder(Segment{
		type_:      Command,
		components: commands,
	})
}

func (c *CommandStructure) FillWithPipe() *CommandStructure {
	return c.FillPlaceholder(Segment{
		type_: Pipe,
	})
}

func (c *CommandStructure) FillPlaceholder(segment Segment) *CommandStructure {
	for i, segment_ := range c.segments {
		if segment_.type_ == Placeholder {
			c.segments[i] = segment

			break
		}
	}

	return c
}

func (c *CommandStructure) String() string {
	string := ""
	for _, segment := range c.segments {
		switch segment.type_ {
		case Pipe:
			string += " | "
		case Command:
			string += strings.Join(segment.components, " ")
		}
	}

	return string
}
