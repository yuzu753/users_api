package main

import (
	"users_api/src/runtime"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		runtime.Module,
		fx.Invoke(runtime.StartLambda), // Use Lambda start function
	).Run()
}