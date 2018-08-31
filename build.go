package building

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type B struct {
	root          string
	targets       map[string]target
	defaultTarget *target
	tools         map[string]Tool
	mutex         sync.Mutex
}

type target struct {
	name        string
	description string
	f           func(*B)
}

var b *B

func Init(pkgName string) *B {
	if b != nil {
		panic("build.Init(...) called twice")
	}
	b = &B{
		root:    pkgName,
		targets: make(map[string]target),
		tools:   make(map[string]Tool),
	}
	return b
}

func Builder() *B {
	if b == nil {
		panic("missing building.Init(...)")
	}
	return b
}

func (b *B) MakeTarget(name, description string, f func(*B)) target {
	t := target{
		name:        name,
		description: description,
		f:           f,
	}
	b.targets[name] = t
	if b.defaultTarget == nil {
		b.defaultTarget = &t
	}
	return t
}

func (b *B) Build(t target) {
	Println(">", t.name)
	start := time.Now()
	t.f(b)
	delta := time.Now().Sub(start)
	Printf("< %s (took %s)", t.name, delta)
}

func (b *B) printTargets() {
	fmt.Printf("\nTargets:\n")
	align := 6
	sorted := sort.StringSlice{}
	for _, t := range b.targets {
		sorted = append(sorted, t.name)
	}
	sorted.Sort()
	for _, s := range sorted {
		t := b.targets[s]
		lf := ""
		spaces := align - len(t.name)
		if spaces <= 0 {
			if len(t.description) > 0 {
				lf = "\n"
			}
			spaces = align + 2
		}
		fmt.Printf("  %s%s%s\n", t.name, lf+strings.Repeat(" ", spaces), t.description)
	}
}

func (b *B) Run() {
	flag.Usage = func() {
		fmt.Print(`Usage: build [OPTIONS] [TARGETS]

Options:
`)
		flag.PrintDefaults()
		b.printTargets()
	}
	*quiet = false
	flag.Parse()

	var runs []target
	args := flag.Args()
	if len(args) == 0 {
		if b.defaultTarget == nil {
			Fatal("no target defined")
		}
		runs = append(runs, *b.defaultTarget)
	}
	for _, a := range args {
		if t, ok := b.targets[a]; ok {
			runs = append(runs, t)
		} else {
			Fatalln("invalid target", a)
		}
	}

	Print("build started")
	start := time.Now()
	for _, t := range runs {
		b.Build(t)
	}
	delta := time.Now().Sub(start)
	Printf("build finished (took %s)", delta)
}

func (b *B) Exe(os string) string {
	if os == "windows" {
		return ".exe"
	}
	return ""
}
