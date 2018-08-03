package building

import (
	"fmt"
	"log"
	"os"
)

type failure struct {
}

func CatchFailure() {
	if e := recover(); e != nil {
		if _, ok := e.(failure); ok {
			fmt.Print("\nBuild failed\n")
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
