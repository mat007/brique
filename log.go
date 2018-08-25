package building

import (
	"flag"
	"log"
	"os"
)

var verbose = flag.Bool("v", false, "verbose")

func init() {
	// manual flags parsing to enable verbose before any actual work
	for _, arg := range os.Args {
		if arg == "-v" {
			*verbose = true
			return
		}
	}
}

type failure struct {
}

func CatchFailure() {
	if e := recover(); e != nil {
		if _, ok := e.(failure); ok {
			// $$$$ MAT: print stack trace
			Print("build failed")
			os.Exit(1)
		}
		panic(e)
	}
}

func Fatalf(format string, v ...interface{}) {
	log.Printf(format, v...)
	panic(failure{})
}

func Fatal(v ...interface{}) {
	log.Print(v...)
	panic(failure{})
}

func Fatalln(v ...interface{}) {
	log.Println(v...)
	panic(failure{})
}

func Printf(format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v...)
	}
}

func Print(v ...interface{}) {
	if *verbose {
		log.Print(v...)
	}
}

func Println(v ...interface{}) {
	if *verbose {
		log.Println(v...)
	}
}
