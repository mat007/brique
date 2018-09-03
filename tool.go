package building

import (
	"archive/tar"
	"bytes"
	"flag"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var (
	AlpineVersion = "3.7"

	containers = flag.Bool("containers", false, "always build in containers")
	cross      = flag.Bool("cross", false, "build for all platforms (linux, darwin, windows)")
	parallel   = flag.Bool("parallel", false, "build in parallel")
)

func init() {
	// Manual flags parsing to enable containers before any actual work
	for _, arg := range os.Args {
		if arg == "-containers" {
			*containers = true
			return
		}
	}
}

type Tool struct {
	root         string
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

func (b *B) MakeTool(name, check, url, instructions string, args ...string) Tool {
	t := b.makeTool(name, check, url, instructions)
	if len(args) > 0 {
		t.Run(args...)
	}
	return t
}

func (b *B) makeTool(name, check, url, instructions string) Tool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if t, ok := b.tools[name]; ok {
		return t
	}
	t := Tool{
		root:         b.root,
		name:         name,
		url:          url,
		instructions: instructions,
		names:        name,
		container:    noApplication(name, check),
	}
	if t.container && name != "" && url != "" && !*containers {
		Print("missing " + name + ": consider installing it to speed up the build, see " + url)
	}
	if t.container {
		t.buildImage()
	}
	b.tools[name] = t
	return t
}

func (b *B) WithOS(f func(goos string)) {
	platforms := []string{runtime.GOOS}
	if *cross {
		platforms = []string{"linux", "darwin", "windows"}
	}
	wg := sync.WaitGroup{}
	for _, goos := range platforms {
		Println("building for", goos)
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
	Debugln("checking for", name)
	if len(check) == 0 {
		Fatalf("missing check for %s", name)
	}
	cmd := exec.Command(name, check)
	err := cmd.Run()
	if err == nil {
		return *containers
	}
	_, ok := err.(*exec.ExitError)
	return !ok
}

func (t Tool) buildImage() {
	Println("preparing image for", t.name)
	buf := &bytes.Buffer{}
	tarFile(t.instructions, "Dockerfile", buf)
	cmd := exec.Command("docker", "build", "-t", t.image(), "-")
	cmd.Stderr = os.Stderr
	if isDebug() {
		cmd.Stdout = os.Stdout
	}
	cmd.Stdin = buf
	if err := cmd.Run(); err != nil {
		Fatalf("error building image for %s: %s", t.name, err)
	}
}

func (t Tool) image() string {
	if t.root == "" {
		Fatalf("error building image for %s: missing root", t.name)
	}
	return strings.Replace(t.root, "/", "-", -1) + "-build-" + t.names
}

func tarFile(content, filename string, writer io.Writer) {
	tw := tar.NewWriter(writer)
	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		Fatal(err)
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		Fatal(err)
	}
	if err := tw.Close(); err != nil {
		Fatal(err)
	}
}

func (t Tool) Run(args ...string) int {
	if t.output == nil {
		t.output = os.Stdout
	}
	if t.input == nil {
		t.input = os.Stdin
	}
	if t.container {
		return t.runContainer(args)
	}
	return t.runApplication(args)
}

func (t Tool) runApplication(args []string) int {
	Println("running", append([]string{t.name}, args...))
	cmd := exec.Command(t.name, args...)
	cmd.Env = append(os.Environ(), t.env...)
	if !t.success {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdout = t.output
	cmd.Stdin = t.input
	code, err := run(cmd, t.success)
	if err != nil {
		Fatalf("error running %s: %s", t.name, err)
	}
	return code
}

func (t Tool) runContainer(args []string) int {
	// $$$$ MAT error out if docker in windows containers mode
	Println("running (in container)", append([]string{t.name}, args...))
	wd, err := os.Getwd()
	if err != nil {
		Fatalf("error running container for %s: %s", t.name, err)
	}
	if t.root == "" {
		Fatalf("error running container for %s: missing root", t.name)
	}
	w := "/go/src/" + t.root
	var envs []string
	for _, e := range t.env {
		envs = append(envs, "-e", e)
	}
	// $$$$ MAT use --net=none by default and allow to customize by tool
	arg := append([]string{"run", "--rm", "-v", wd + ":" + w, "-w", w, "-i"}, envs...)
	arg = append(arg, t.image(), t.name)
	// $$$$ MAT try and replace wd in args with w ?
	// $$$$ do the same with TEMPDIR -> /tmp, and mount it ? any dir ?
	// $$$$ MAT if GOPATH set, mount it instead of wd ?
	arg = append(arg, args...)
	Debugln("running", append([]string{"docker"}, arg...))
	cmd := exec.Command("docker", arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = t.output
	cmd.Stdin = t.input
	code, err := run(cmd, t.success)
	if err != nil {
		Fatalf("error running container for %s: %s", t.name, err)
	}
	return code
}

func run(cmd *exec.Cmd, success bool) (int, error) {
	if err := cmd.Run(); err != nil {
		exit, ok := err.(*exec.ExitError)
		if ok {
			if success {
				err = nil
			}
			if status, ok := exit.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), err
			}
		}
		return 1, err
	}
	return 0, nil
}
