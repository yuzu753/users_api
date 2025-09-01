package main

import (
	"users_api/src/runtime"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		runtime.Module,
		fx.Invoke(runtime.Start), // Use server start function
	).Run()
}