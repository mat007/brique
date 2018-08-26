package building

import (
	"flag"
	"log"
	"os"
)

var (
	quiet   = flag.Bool("q", false, "quiet output")
	verbose = flag.Bool("v", false, "verbose output")
)

func init() {
	log.SetFlags(0)
	// Manual flags parsing to disable logging before calling the target
	// functions unless -v is passed.
	*quiet = true
	for _, arg := range os.Args {
		if arg == "-v" {
			*verbose = true
			*quiet = false
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
			log.Print("build failed")
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
	if isInfo() {
		log.Printf(format, v...)
	}
}

func Print(v ...interface{}) {
	if isInfo() {
		log.Print(v...)
	}
}

func Println(v ...interface{}) {
	if isInfo() {
		log.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if isDebug() {
		log.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if isDebug() {
		log.Print(v...)
	}
}

func Debugln(v ...interface{}) {
	if isDebug() {
		log.Println(v...)
	}
}

func isInfo() bool {
	return !*quiet
}

func isDebug() bool {
	return isInfo() && *verbose
}

func Quiet() {
	*verbose = false
	*quiet = true
}
