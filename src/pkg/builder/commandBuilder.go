package builder

import (
	"io"
	"os/exec"
)

type CommandBuilder struct {
	command exec.Cmd
}

func (b *CommandBuilder) SetPath(path string) *CommandBuilder {
	b.command.Path = path
	return b
}

func (b *CommandBuilder) SetArgs(args []string) *CommandBuilder {
	b.command.Args = args
	return b
}

func (b *CommandBuilder) SetStdout(out io.Writer) *CommandBuilder {
	b.command.Stdout = out
	return b
}

func (b *CommandBuilder) GetCommandInstance() exec.Cmd {
	return exec.Cmd{
		Path:   b.command.Path,
		Args:   b.command.Args,
		Stdout: b.command.Stdout,
	}
}
