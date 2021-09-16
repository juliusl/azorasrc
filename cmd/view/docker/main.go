package main

import (
	"context"
	"os"
	"time"

	"github.com/juliusl/azorasrc/pkg/auto"
)

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)

	defer cancel()

	manifest, err := auto.Resolve(ctx)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	written, err := auto.Fetch(ctx, manifest.Layers...)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	if written > 0 {
		err = auto.Commit(ctx)

		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}
}
