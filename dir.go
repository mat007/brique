package building

import (
	"log"
	"os"
)

func (b *B) Dir(dir string, f func()) {
	wd, err := os.Getwd()
	if err != nil {
		Fatalln("cannot get current directory:", err)
	}
	// $$$$ MAT: this fails if inside a parallel section...
	defer func() {
		if err = os.Chdir(wd); err != nil {
			Fatalln("cannot change back current directory:", err)
		}
	}()
	if err := os.MkdirAll(dir, 0755); err != nil {
		Fatalln("cannot create directory:", err)
	}
	// $$$$ MAT: does not work in container build because the working dir gets mounted as root
	if err := os.Chdir(dir); err != nil {
		Fatalln("cannot change current directory:", err)
	}
	if *verbose {
		log.Println("changed to directory", dir)
	}
	f()
}
