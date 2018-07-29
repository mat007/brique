package b

import (
	"log"
	"os"
)

func (b *B) Remove(paths ...string) {
	err := remove(func(path string) {
	}, paths...)
	if err != nil {
		log.Fatalf("remove failed: %s", err)
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
		} else if *verbose {
			log.Printf("removed %q", match)
		}
	}
	return nil
}
