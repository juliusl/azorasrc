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

	// Fetch image config
	written, err := auto.Fetch(ctx, manifest.Config)
	if err != nil {
		exitWithError(err)
	}
	total += written

	// Fetch layers
	written, err = auto.Fetch(ctx, manifest.Layers...)
	if err != nil {
		exitWithError(err)
	}
	total += written

	// Discover artifacts
	artifacts, err := auto.Discover(ctx, *desc, "")
	if err != nil {
		if err.Error() != "HTTP 404" {
			exitWithError(err)
		}
	} else {
		// Fetch artifacts
		for _, a := range artifacts {
			written, err = auto.FetchArtifact(ctx, a.Blobs...)
			if err != nil {
				exitWithError(err)
			}
			total += written
		}
	}

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
