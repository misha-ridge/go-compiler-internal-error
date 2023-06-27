package a

import (
	"github.com/misha-ridge/x/b"
)

func A() {
	b.R(func(spawn func(func() error)) error {
		spawn(b.NewS().R)
		return nil
	})
}
