package auto

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/juliusl/azorasrc/pkg/auto/docker"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/remotes"
)

var (
	outputDir string
	reference string
	registry  *remotes.Registry
	fetch     func(ctx context.Context, desc v1.Descriptor) (io.ReadCloser, error)
	discover  func(ctx context.Context, desc v1.Descriptor, artifactType string) (*remotes.Artifacts, error)

	files []string
)

func init() {
	registry = docker.Registry()

	if len(os.Args) < 1 {
		os.Stderr.WriteString("Usage: <reference>")
		os.Exit(1)
	}
	reference = os.Args[1]
}

func init() {
	// Set the fetcher
	fetch = registry.AsFunctions().Fetcher()

	// Set the discoverer
	discover = registry.AsFunctions().Discoverer()
}

func Resolve(ctx context.Context) (*v1.Manifest, error) {
	return registry.GetManifest(ctx, reference)
}

func Fetch(ctx context.Context, descs []v1.Descriptor) (bytes int, err error) {
	for _, desc := range descs {
		blob, err := fetch(ctx, desc)
		if err != nil {
			err = fmt.Errorf("error: %w", err)
			return bytes, err
		}

		f, err := os.CreateTemp("", fmt.Sprintf("%s-*", desc.Digest.String()))
		if err != nil {
			return bytes, err
		}

		defer f.Close()

		files = append(files, f.Name())

		written, err := io.Copy(f, blob)
		if err != nil {
			err = fmt.Errorf("error: %w", err)
			return bytes, err
		}

		if written <= 0 {
			err = fmt.Errorf("error: %w", err)
			return bytes, err
		}

		bytes += int(written)
	}

	return bytes, nil
}

func Commit(ctx context.Context) error {
	for _, path := range files {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		f := path[:strings.IndexByte(path, '-')]
		err = os.WriteFile(f, bytes, 0644)
		if err != nil {
			return err
		}

		os.Stdout.WriteString(f)
		os.Stdout.WriteString("\n")
		os.Remove(path)
	}

	return nil
}
