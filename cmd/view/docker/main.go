package main

import (
	"context"
	"os"
	"time"

	"github.com/juliusl/azorasrc/pkg/auto"
)

var (
	total int
)

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)

	defer cancel()

	desc, manifest, err := auto.Resolve(ctx)
	if err != nil {
		exitWithError(err)
	}

	// Setup a store for the manifest
	err = auto.SetupV1ManifestStore(*desc, *manifest)
	if err != nil {
		exitWithError(err)
	}

	written, err := auto.Fetch(ctx, manifest.Config)
	if err != nil {
		exitWithError(err)
	}
	total += written

	written, err = auto.Fetch(ctx, manifest.Layers...)
	if err != nil {
		exitWithError(err)
	}
	total += written

	commitIfContent(ctx, total)
	os.Exit(0)
}

func commitIfContent(ctx context.Context, written int) {
	if written > 0 {
		err := auto.Commit(ctx)

		if err != nil {
			exitWithError(err)
		}
	}
}

func exitWithError(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}
