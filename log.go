package building

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	Quiet   = flag.Bool("q", false, "quiet output")
	Verbose = flag.Bool("v", false, "verbose output")
)

func init() {
	log.SetFlags(0)
	// Manual flags parsing to disable logging before calling the target
	// functions unless -v is passed.
	*Quiet = true
	for _, arg := range os.Args {
		if arg == "-v" {
			*Verbose = true
			*Quiet = false
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
	log.Printf(location()+" "+format, v...)
	panic(failure{})
}

func Fatal(v ...interface{}) {
	log.Print(append([]interface{}{location() + " "}, v...)...)
	panic(failure{})
}

func Fatalln(v ...interface{}) {
	log.Println(append([]interface{}{location()}, v...)...)
	panic(failure{})
}

func location() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf("%s:%d:", file, line)
	}
	return ""
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
	return !*Quiet
}

func isDebug() bool {
	return isInfo() && *Verbose
}

func ForceQuiet() {
	*Verbose = false
	*Quiet = true
}
