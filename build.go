package b

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

var (
	containers = flag.Bool("c", false, "always build in containers")
	verbose    = flag.Bool("v", false, "verbose")
)

type target struct {
	name        string
	description string
	f           func()
}

var (
	targets       = make(map[string]target)
	defaultTarget target
)

func Target(name, description string, f func()) target {
	t := target{
		name:        name,
		description: description,
		f:           f,
	}
	targets[name] = t
	return t
}

func Build(t target) {
	if *verbose {
		log.Println(">", t.name)
	}
	start := time.Now()
	t.f()
	if *verbose {
		delta := time.Now().Sub(start)
		log.Printf("< %s (took %s)", t.name, delta)
	}
}

func (t target) Default() target {
	if defaultTarget.f != nil {
		log.Fatalf("%s cannot be set as default: %s already set", t.name, defaultTarget.name)
	}
	defaultTarget = t
	return t
}

func Run() {
	flag.Usage = func() {
		fmt.Print(`Usage: build [OPTIONS] [TARGETS]

Options:
`)
		flag.PrintDefaults()
		fmt.Printf("\nTargets:\n")
		align := 6
		sorted := sort.StringSlice{}
		for _, t := range targets {
			sorted = append(sorted, t.name)
		}
		sorted.Sort()
		for _, s := range sorted {
			t := targets[s]
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
	flag.Parse()

	var runs []target
	args := flag.Args()
	if len(args) == 0 {
		if defaultTarget.f == nil {
			log.Fatal("missing default target")
		}
		runs = append(runs, defaultTarget)
	}
	for _, a := range args {
		if t, ok := targets[a]; ok {
			runs = append(runs, t)
		} else {
			log.Fatalln("invalid target", a)
		}
	}

	if *verbose {
		log.Print("build started")
	}
	start := time.Now()
	for _, t := range runs {
		Build(t)
	}
	if *verbose {
		delta := time.Now().Sub(start)
		log.Printf("build finished (took %s)", delta)
	}
}

func ExecExt(os string) string {
	if os == "windows" {
		return ".exe"
	}
	return ""
}
