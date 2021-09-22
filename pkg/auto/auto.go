package auto

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/juliusl/azorasrc/pkg/auto/docker"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/remotes"
)

var (
	registry   *remotes.Registry
	fetch      func(ctx context.Context, desc v1.Descriptor) (io.ReadCloser, error)
	discover   func(ctx context.Context, desc v1.Descriptor, artifactType string) (*remotes.Artifacts, error)
	files      []string
	workingDir string
	outputDir  string
	reference  string
)

func init() {
	registry = docker.Registry()

	if len(os.Args) < 1 {
		os.Stderr.WriteString("Usage: <reference>")
		os.Exit(1)
	}
	reference = os.Args[1]

	workingDir = "/temp"
	outputDir = ""
}

func init() {
	// Set the fetcher
	fetch = registry.AsFunctions().Fetcher()

	// Set the discoverer
	discover = registry.AsFunctions().Discoverer()
}

func SetDirectories(working, output string) error {
	for i, d := range []string{working, output} {
		if d == "" {
			continue
		}

		info, err := os.Stat(d)
		if err != nil {
			return err
		}

		if info.IsDir() {
			switch i {
			case 0:
				workingDir = d
			case 1:
				outputDir = d
			}
		} else {
			return errors.New("path is not a directory")
		}
	}

	if working != "" && workingDir != working {
		return errors.New("did not set working dir")
	}

	if output != "" && outputDir != output {
		return errors.New("did not set output dir")
	}

	return nil
}

func Discover(ctx context.Context, desc v1.Descriptor, artifactType string) (*remotes.Artifacts, error) {
	return discover(ctx, desc, artifactType)
}

func Resolve(ctx context.Context) (*v1.Manifest, error) {
	return registry.GetManifest(ctx, reference)
}

func Fetch(ctx context.Context, descs ...v1.Descriptor) (bytes int, err error) {
	for _, desc := range descs {
		blob, err := fetch(ctx, desc)
		if err != nil {
			err = fmt.Errorf("error: %w", err)
			return bytes, err
		}

		f, err := tempFile(desc)
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
	for _, work := range files {
		src, err := os.Open(work)
		if err != nil {
			return err
		}
		defer src.Close()

		dest, err := committedFile(path.Base(work))
		if err != nil {
			return err
		}

		defer dest.Close()

		written, err := io.Copy(dest, src)
		if err != nil {
			return err
		}

		if written > 0 {
			os.Stdout.WriteString(dest.Name())
			os.Stdout.WriteString("\n")
			os.Remove(work)
		} else {
			return errors.New("could not write to destination")
		}
	}

	return nil
}

func tempFile(desc v1.Descriptor) (*os.File, error) {
	filename := fmt.Sprintf("%s-*", desc.Digest.String())

	path := path.Join(workingDir, filename)

	return os.Create(path)
}

func committedFile(tempFile string) (*os.File, error) {
	filename := tempFile[:strings.IndexByte(tempFile, '-')]

	path := path.Join(outputDir, filename)

	return os.Create(path)
}
