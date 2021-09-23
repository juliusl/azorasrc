package auto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/juliusl/azorasrc/pkg/auto/docker"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	artifactspec "github.com/oras-project/artifacts-spec/specs-go/v1"
	"oras.land/oras-go/pkg/remotes"
)

var (
	registry   *remotes.Registry
	fetch      func(ctx context.Context, desc ocispec.Descriptor) (io.ReadCloser, error)
	discover   func(ctx context.Context, desc ocispec.Descriptor, artifactType string) (*remotes.Artifacts, error)
	files      []string
	workingDir string
	outputDir  string
	reference  string
	host       string
	namespace  string
	loc        string
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

func SetupV1ManifestStore(desc ocispec.Descriptor, manifest ocispec.Manifest) error {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	var ok bool
	host, ok = manifest.Annotations["host"]
	if !ok {
		return errors.New("missing host annotation")
	}

	namespace, ok = manifest.Annotations["namespace"]
	if !ok {
		return errors.New("missing namespace annotation")
	}

	loc, ok = manifest.Annotations["loc"]
	if !ok {
		return errors.New("missing loc annotation")
	}

	storeDir := path.Join(cachedir, host, namespace, loc, desc.Digest.String())
	err = os.MkdirAll(storeDir, 0755)
	if err != nil {
		return err
	}

	err = SetDirectories("./work", storeDir)
	if err != nil {
		return err
	}

	pathToManifest := path.Join(storeDir, "manifest.json")
	m, err := os.Create(pathToManifest)
	if err != nil {
		return err
	}

	defer m.Close()

	err = json.NewEncoder(m).Encode(manifest)
	if err != nil {
		return err
	}

	return nil
}

func Discover(ctx context.Context, desc ocispec.Descriptor, artifactType string) ([]artifactspec.Manifest, error) {
	refs, err := discover(ctx, desc, artifactType)
	if err != nil {
		return nil, err
	}

	artifacts := make([]artifactspec.Manifest, len(refs.References))

	for i, r := range refs.References {
		// assemble ref
		aref := fmt.Sprintf("%s/%s@%s", host, namespace, r.Digest.String())

		_, manifest, err := registry.GetArtifactManifest(ctx, aref)
		if err != nil {
			return nil, err
		}

		artifacts[i] = *manifest
	}

	apath := path.Join(outputDir, "discovered-artifacts.json")
	f, err := os.Create(apath)
	if err != nil {
		return nil, err
	}

	err = json.NewEncoder(f).Encode(artifacts)
	if err != nil {
		return nil, err
	}

	return artifacts, nil
}

func Resolve(ctx context.Context) (*ocispec.Descriptor, *ocispec.Manifest, error) {
	return registry.GetManifest(ctx, reference)
}

func FetchArtifact(ctx context.Context, descs ...artifactspec.Descriptor) (bytes int, err error) {
	normalize := make([]ocispec.Descriptor, len(descs))

	for i, d := range descs {
		normalize[i] = ocispec.Descriptor{
			Size:        d.Size,
			MediaType:   d.MediaType,
			Digest:      d.Digest,
			Annotations: d.Annotations,
		}

		normalize[i].Annotations["artifactType"] = d.ArtifactType
	}

	return Fetch(ctx, normalize...)
}

func Fetch(ctx context.Context, descs ...ocispec.Descriptor) (bytes int, err error) {
	for _, desc := range descs {
		blob, err := fetch(ctx, desc)
		if err != nil {
			err = fmt.Errorf("error: %w", err)
			return bytes, err
		}
		defer blob.Close()

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

func tempFile(desc ocispec.Descriptor) (*os.File, error) {
	filename := fmt.Sprintf("%s-*", desc.Digest.String())

	path := path.Join(workingDir, filename)

	return os.Create(path)
}

func committedFile(tempFile string) (*os.File, error) {
	filename := tempFile[:strings.IndexByte(tempFile, '-')]

	path := path.Join(outputDir, filename)

	return os.Create(path)
}
