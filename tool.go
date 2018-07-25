package b

import (
	"archive/tar"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	// PackageName stores the go package name of the project, this must not be left empty.
	PackageName   string
	AlpineVersion = "3.7"
)

type Tool struct {
	name         string
	url          string
	env          []string
	instructions string
	names        string
	container    bool
	output       io.Writer
	input        io.Reader
	success      bool
}

func (t Tool) WithEnv(env ...string) Tool {
	t.env = append(t.env, env...)
	return t
}

func (t Tool) WithOutput(w io.Writer) Tool {
	t.output = w
	return t
}

func (t Tool) WithInput(r io.Reader) Tool {
	t.input = r
	return t
}

func (t Tool) WithSuccess() Tool {
	t.success = true
	return t
}

func (t Tool) WithTool(tool Tool) Tool {
	t.instructions += "\n" + tool.instructions
	if t.container || tool.container {
		t.container = true
		t.names += "-" + tool.name
		t.buildImage()
	}
	return t
}

func MakeTool(name, check, url, instructions string) Tool {
	t := Tool{
		name:         name,
		url:          url,
		instructions: instructions,
		names:        name,
		container:    noApplication(name, check),
	}
	if t.container && name != "" && url != "" {
		log.Print("missing " + name + ": consider installing it to speed up the build, see " + url)
	}
	if t.container {
		t.buildImage()
	}
	return t
}

func WithOS(f func(goos string)) {
	platforms := []string{runtime.GOOS}
	if *cross {
		platforms = []string{"linux", "darwin", "windows"}
	}
	wg := sync.WaitGroup{}
	for _, goos := range platforms {
		if *verbose {
			log.Println("building for", goos)
		}
		if *parallel {
			wg.Add(1)
			go func(goos string) {
				f(goos)
				wg.Done()
			}(goos)
		} else {
			f(goos)
		}
	}
	wg.Wait()
}

func noApplication(name, check string) bool {
	if *verbose {
		log.Println("checking for", name)
	}
	if len(check) == 0 {
		log.Fatalf("missing check for %s", name)
	}
	cmd := exec.Command(name, check)
	err := cmd.Run()
	if err == nil {
		return *containers
	}
	if _, ok := err.(*exec.ExitError); ok {
		log.Fatalf("error checking %s: %s", name, err)
	}
	return true
}

func (t Tool) buildImage() {
	var buf bytes.Buffer
	tarFile(t.instructions, "Dockerfile", &buf)
	cmd := exec.Command("docker", "build", "-t", t.image(), "-")
	cmd.Stderr = os.Stderr
	if *verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stdin = &buf
	if err := cmd.Run(); err != nil {
		log.Fatalf("error building image for %s: %s", t.name, err)
	}
}

func (t Tool) image() string {
	if PackageName == "" {
		log.Fatalf("error building image for %s: missing PackageName", t.name)
	}
	return strings.Replace(PackageName, "/", "-", -1) + "-build-" + t.names
}

func tarFile(content, filename string, writer io.Writer) {
	tw := tar.NewWriter(writer)
	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}
}

func (t Tool) Run(args ...string) {
	if t.output == nil {
		t.output = os.Stdout
	}
	if t.input == nil {
		t.input = os.Stdin
	}
	if t.container {
		t.runContainer(args)
	} else {
		t.runApplication(args)
	}
}

func (t Tool) runApplication(args []string) {
	if *verbose {
		log.Println("running", t.name, args)
	}
	cmd := exec.Command(t.name, args...)
	cmd.Env = append(os.Environ(), t.env...)
	if !t.success {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdout = t.output
	cmd.Stdin = t.input
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok || !t.success {
			log.Fatalf("error running %s: %s", t.name, err)
		}
	}
}

func (t Tool) runContainer(args []string) {
	if *verbose {
		log.Println("running (container)", t.name, args)
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error running container for %s: %s", t.name, err)
	}
	if PackageName == "" {
		log.Fatalf("error running container for %s: missing PackageName", t.name)
	}
	w := "/go/src/" + PackageName
	var envs []string
	for _, e := range t.env {
		envs = append(envs, "-e", e)
	}
	arg := append([]string{"run", "--rm", "-v", wd + ":" + w, "-w", w, "-t"}, envs...)
	arg = append(arg, t.image(), t.name)
	arg = append(arg, args...)
	if t.success {
		// $$$$ MAT ignore error for run cmd
	}
	cmd := exec.Command("docker", arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = t.output
	cmd.Stdin = t.input
	if err := cmd.Run(); err != nil {
		log.Fatalf("error running container for %s: %s", t.name, err)
	}
}
