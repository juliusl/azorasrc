package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"oras.land/oras-go/pkg/auth"
	"oras.land/oras-go/pkg/auth/docker"
	"oras.land/oras-go/pkg/remotes"
)

//
// TITLE: AZORAS View Manifest
// DESCRIPTION: View a manifest in the shell
// USAGE: ./view <reference>
// ENV:
// 		ORAS_BEGIN_ENV="<path-to-root-directory>"
// 		ORAS_STORE_DIR="$ORAS_BEGIN_ENV/<path-to-store-directory>"
// 		ORAS_LOGIN_DIR="$ORAS_BEGIN_ENV/<path-to-login-directory>"
// ARGS:
// reference - this is a oci distribution v2 compliant reference string
//
func main() {
	// Parse arguments
	if len(os.Args) < 1 {
		os.Stderr.WriteString("Usage: view <reference>")
		os.Exit(1)
	}
	reference := os.Args[1]

	ctx := context.Background()

	host, namespace, _, err := remotes.Parse(reference)
	if err != nil {
		os.Stderr.WriteString("could not parse reference")
		os.Exit(1)
	}

	reg, err := docker.NewRegistryWithAccessProvider(host, namespace, []string{}, auth.WithLoginContext(ctx))
	if err != nil {
		os.Stderr.WriteString("could not create an access provider")
		os.Exit(1)
	}

	// Get a resolver
	resolve := reg.AsFunctions().Resolver()
	_, desc, err := resolve(ctx, reference)
	if err != nil {
		err = fmt.Errorf("error: %w", err)

		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	// Get a fetcher
	fetch := reg.AsFunctions().Fetcher()
	man, err := fetch(ctx, desc)
	if err != nil {
		err = fmt.Errorf("error: %w", err)

		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	defer man.Close()

	// Copy result to std out
	written, err := io.Copy(os.Stdout, man)
	if err != nil {
		err = fmt.Errorf("error: %w", err)

		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	if written <= 0 {
		err = fmt.Errorf("error: %w", err)

		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
	// Store the desc
	// writer, err := store.Writer(
	// 	ctx,
	// 	content.WithDescriptor(desc), // TODO - fix reference to containerd
	// 	content.WithRef(reference),
	// )
	// if err != nil {
	// 	os.Exit(1)
	// }

	// written, err := io.Copy(writer, man)
	// if err != nil {
	// 	os.Exit(1)
	// }

	// if written <= 0 {
	// 	os.Exit(1)
	// }

}
