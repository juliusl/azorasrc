package docker

import (
	"context"
	"os"

	"oras.land/oras-go/pkg/auth"
	"oras.land/oras-go/pkg/auth/docker"
	"oras.land/oras-go/pkg/remotes"
)

var (
	reference string
	ctx       context.Context
	registry  *remotes.Registry
)

func init() {
	// Parse arguments
	if len(os.Args) < 1 {
		os.Stderr.WriteString("Usage: <reference>")
		os.Exit(1)
	}
	reference = os.Args[1]
	ctx = context.Background()
}

func init() {
	host, namespace, _, err := remotes.Parse(reference)
	if err != nil {
		os.Stderr.WriteString("could not parse reference")
		os.Exit(1)
	}

	registry, err = docker.NewRegistryWithAccessProvider(
		host,
		namespace,
		[]string{},
		auth.WithLoginContext(ctx))
	if err != nil {
		os.Stderr.WriteString("could not create an access provider")
		os.Exit(1)
	}
}

func Registry() *remotes.Registry {
	return registry
}
