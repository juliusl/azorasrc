package shell

import (
	"context"
	"fmt"

	remotessh "github.com/juliusl/azorasrc/pkg/remotes/shell"
	"oras.land/oras-go/pkg/auth"
	"oras.land/oras-go/pkg/remotes"
)

func New(ctx context.Context, reference string) (*remotes.Registry, error) {
	// Parse and validate reference
	_, host, ns, _, err := remotes.Parse(reference)
	if err != nil {
		err = fmt.Errorf("parse-error: %w", err)
		return nil, err
	}

	// Begin authn
	_, sh, err := Begin(ns)
	if err != nil {
		err = fmt.Errorf("begin-error: %w", err)
		return nil, err
	}

	// Login, behind the scenes this just checks the status
	err = sh.LoginWithOpts(
		auth.WithLoginContext(ctx),
		auth.WithLoginHostname(host),
	)
	if err != nil {
		err = fmt.Errorf("login-error: %w", err)
		return nil, err
	}

	// Configure the access provider directory passed back from the login process
	ap, err := remotessh.ConfigureAccessProvider(
		sh.LoginDir,
		sh.AccessProviderDir)
	if err != nil {
		err = fmt.Errorf("access-provider-error: %w", err)
		return nil, err
	}

	// Create new registry interface
	return remotes.NewRegistry(host, ns, ap), nil
}
