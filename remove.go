package building

import (
	"os"
)

func (b *B) Remove(paths ...string) {
	err := remove(paths...)
	if err != nil {
		b.Fatalln(err)
	}
}

func remove(paths ...string) error {
	matches, err := glob("", paths, false)
	if err != nil {
		return err
	}
	for _, match := range matches {
		if err := os.RemoveAll(match); err != nil {
			return err
		}
		b.Debugf("removed %q", match)
	}
	return nil
}
