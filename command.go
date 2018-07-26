package b

import (
	"log"
	"os"
	"os/exec"
)

type command struct {
	name string
}

func Command(name string) command {
	return command{
		name: name,
	}
}

func (c command) Run(args ...string) command {
	if *verbose {
		log.Println("running", c.name, args)
	}
	cmd := exec.Command(c.name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("error running %s: %s", c.name, err)
	}
	return c
}
