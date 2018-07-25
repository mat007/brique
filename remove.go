package b

import (
	"log"
	"os"
)

func Remove(paths ...string) {
	err := remove(func(path string) {
		if *verbose {
			log.Println("removed", path)
		}
	}, paths...)
	if err != nil {
		log.Fatalf("clean failed: %s", err)
	}
}

func remove(notify func(path string), paths ...string) error {
	matches, err := glob("", paths, false)
	if err != nil {
		return err
	}
	for _, match := range matches {
		if err := os.RemoveAll(match); err != nil {
			return err
		} else if notify != nil {
			notify(match)
		}
	}
	return nil
}
