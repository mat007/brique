package building

import (
	"io"
	"os"
	"os/exec"
)

type command struct {
	name    string
	output  io.Writer
	success bool
}

func (b *B) Command(name string, args ...string) command {
	c := command{
		name: name,
	}
	if len(args) > 0 {
		c.Run(args...)
	}
	return c
}

func (c command) WithOutput(w io.Writer) command {
	c.output = w
	return c
}

func (c command) WithSuccess() command {
	c.success = true
	return c
}

func (c command) Run(args ...string) int {
	Println("running", append([]string{c.name}, args...))
	if c.output == nil {
		c.output = os.Stdout
	}
	cmd := exec.Command(c.name, args...)
	cmd.Stdout = c.output
	cmd.Stderr = os.Stderr
	code, err := run(cmd, c.success)
	if err != nil {
		Fatalf("error running %s: %s", c.name, err)
	}
	return code
}
