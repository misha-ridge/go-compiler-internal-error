package a

import (
	"github.com/misha-ridge/x/b"
)

func A() {
	b.Run(func(spawn func(func() error)) error {
		spawn(b.NewS().Run)
		return nil
	})
}
