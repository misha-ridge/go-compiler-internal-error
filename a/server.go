package a

import (
	"github.com/misha-ridge/x/thttp"
)

func A() {
	thttp.Run(nil, func(spawn thttp.SpawnFn) error {
		spawn("server", thttp.Fail, thttp.NewServer(nil, nil).Run)
		return nil
	})
}
