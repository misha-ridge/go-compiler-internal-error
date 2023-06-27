package a

import (
	"context"
	"github.com/misha-ridge/x/thttp"
)

func A() {
	thttp.Run(nil, func(spawn func(func(context.Context) error)) error {
		spawn(thttp.NewServer(nil, nil).Run)
		return nil
	})
}
