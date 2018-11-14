package building

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type command struct {
	name    string
	dir     string
	env     []string
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

func (c command) WithDir(dir string) command {
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

func (c command) WithEnv(env ...string) command {
	c.env = append(c.env, env...)
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
	cmd.Dir = c.dir
	cmd.Env = append(os.Environ(), c.env...)
	if !c.success {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdout = c.output
	code, err := run(cmd, c.success)
	if err != nil {
		Fatalf("error running %s: %s", c.name, err)
	}
	return code
}
