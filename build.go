package building

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	containers = flag.Bool("c", false, "always build in containers")
	verbose    = flag.Bool("v", false, "verbose")
)

type B struct {
	root    string
	targets map[string]target
	tools   map[string]Tool
	mutex   sync.Mutex
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
	return t
}

func (b *B) Build(t target) {
	if *verbose {
		log.Println(">", t.name)
	}
	start := time.Now()
	t.f(b)
	if *verbose {
		delta := time.Now().Sub(start)
		log.Printf("< %s (took %s)", t.name, delta)
	}
}

// func (t target) Default() target {
// 	if defaultTarget.f != nil {
// 		Fatalf("%s cannot be set as default: %s already set", t.name, defaultTarget.name)
// 	}
// 	defaultTarget = t
// 	return t
// }

func init() {
	// manual flags parsing to enable verbose and containers before any actual work
	for _, arg := range os.Args {
		switch arg {
		case "-v":
			*verbose = true
		case "-c":
			*containers = true
		}
	}
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
	flag.Parse()

	var runs []target
	args := flag.Args()
	if len(args) == 0 {
		Fatal("no target specified")
		// if defaultTarget.f == nil {
		// 	Fatal("missing default target")
		// }
		// runs = append(runs, defaultTarget)
	}
	for _, a := range args {
		if t, ok := b.targets[a]; ok {
			runs = append(runs, t)
		} else {
			Fatalln("invalid target", a)
		}
	}

	if *verbose {
		log.Print("build started")
	}
	start := time.Now()
	for _, t := range runs {
		b.Build(t)
	}
	if *verbose {
		delta := time.Now().Sub(start)
		log.Printf("build finished (took %s)", delta)
	}
}

func (b *B) ExecExt(os string) string {
	if os == "windows" {
		return ".exe"
	}
	return ""
}
