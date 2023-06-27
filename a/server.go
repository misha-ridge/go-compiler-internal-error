package a

import (
	"github.com/misha-ridge/x/thttp"
)

func A() {
	thttp.Run(nil, func(spawn thttp.SpawnFn) error {
		spawn(thttp.NewServer(nil, nil).Run)
		return nil
	})
}
