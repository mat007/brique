package building

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
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

func CatchFailure(start time.Time) {
	if e := recover(); e != nil {
		if _, ok := e.(failure); ok {
			b.Debugf("build failed (took %s)", time.Since(start))
			os.Exit(1)
		}
		panic(e)
	}
	b.Debugf("build finished (took %s)", time.Since(start))
}

func (b *B) Helper() {
	var pc [2]uintptr
	n := runtime.Callers(2, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.helpers == nil {
		b.helpers = make(map[string]bool)
	}
	b.helpers[frame.Function] = true
}

func (b *B) location(suffix string) string {
	var pc [10]uintptr
	n := runtime.Callers(1, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		if frame.Function == "runtime.main" {
			return ""
		}
		if !b.skip(frame) {
			return fmt.Sprintf("%s:%d:%s", frame.File, frame.Line, suffix)
		}
		if !more {
			return ""
		}
	}
}

func (b *B) skip(frame runtime.Frame) bool {
	if strings.Contains(frame.Function, "github.com/mat007/brique") {
		return true
	}
	if b == nil {
		return false
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.helpers[frame.Function]
}

func (b *B) Fatal(v ...interface{}) {
	log.Print(append([]interface{}{b.location(" ")}, v...)...)
	panic(failure{})
}

func (b *B) Fatalf(format string, v ...interface{}) {
	log.Printf(b.location(" ")+format, v...)
	panic(failure{})
}

func (b *B) Fatalln(v ...interface{}) {
	log.Println(append([]interface{}{b.location("")}, v...)...)
	panic(failure{})
}

func (b *B) Print(v ...interface{}) {
	if isInfo() {
		log.Print(v...)
	}
}

func (b *B) Printf(format string, v ...interface{}) {
	if isInfo() {
		log.Printf(format, v...)
	}
}

func (b *B) Println(v ...interface{}) {
	if isInfo() {
		log.Println(v...)
	}
}

func (b *B) Debug(v ...interface{}) {
	if isDebug() {
		log.Print(v...)
	}
}

func (b *B) Debugf(format string, v ...interface{}) {
	if isDebug() {
		log.Printf(format, v...)
	}
}

func (b *B) Debugln(v ...interface{}) {
	if isDebug() {
		log.Println(v...)
	}
}

func (b *B) Assert(err error) {
	if err != nil {
		log.Println([]interface{}{b.location(""), err}...)
		panic(failure{})
	}
}

func (b *B) Check(err error) {
	if err != nil {
		log.Println([]interface{}{b.location(""), err}...)
	}
}

func (b *B) Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println([]interface{}{b.location(""), err}...)
	}
}

func isInfo() bool {
	return !*Quiet
}

func isDebug() bool {
	return isInfo() && *Verbose
}
