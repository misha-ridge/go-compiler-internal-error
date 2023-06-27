package a

import (
	"github.com/misha-ridge/x/thttp"
)

func A() {
	thttp.Run(func(spawn func(func() error)) error {
		spawn(thttp.NewServer().Run)
		return nil
	})
}
