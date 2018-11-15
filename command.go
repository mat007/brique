package building

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	name    string
	dir     string
	env     []string
	output  io.Writer
	success bool
}

func (b *B) MakeCommand(name string, args ...string) Command {
	c := Command{
		name: name,
	}
	if len(args) > 0 {
		c.Run(args...)
	}
	return c
}

func (c Command) WithDir(dir string) Command {
	dir = filepath.Clean(dir)
	if filepath.IsAbs(dir) {
		Fatalln("dir must be relative", dir)
	}
	if strings.Contains(dir, "..") {
		Fatalln("dir must be a folder under project root", dir)
	}
	c.dir = dir
	return c
}

func (c Command) WithEnv(env ...string) Command {
	c.env = append(c.env, env...)
	return c
}

func (c Command) WithOutput(w io.Writer) Command {
	c.output = w
	return c
}

func (c Command) WithSuccess() Command {
	c.success = true
	return c
}

func (c Command) Run(args ...string) int {
	Println("running", append([]string{c.name}, args...))
	if c.output == nil {
		c.output = os.Stdout
	}
	cmd := exec.Command(c.name, args...)
	cmd.Dir = c.dir
	cmd.Env = append(os.Environ(), c.env...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = c.output
	code, err := run(cmd, c.success)
	if err != nil {
		Fatalln(err)
	}
	return code
}
