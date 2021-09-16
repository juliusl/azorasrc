module github.com/juliusl/azorasrc

go 1.16

replace oras.land/oras-go => ./deps/oras-go

require (
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	oras.land/oras-go v0.4.0
	rsc.io/letsencrypt v0.0.3 // indirect
)
