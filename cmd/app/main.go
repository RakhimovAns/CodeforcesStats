package main

import (
	"context"

	"github.com/RakhimovAns/CodeforcesStats/internal/app"
	"github.com/RakhimovAns/wrapper/pkg/closer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	closer.Add(func() error {
		cancel()

		return nil
	})

	a := app.New()
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}
