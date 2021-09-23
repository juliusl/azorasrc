module github.com/juliusl/azorasrc

go 1.16

replace oras.land/oras-go => ./deps/oras-go

require (
	github.com/opencontainers/image-spec v1.0.1
	github.com/oras-project/artifacts-spec v0.0.0-20210824220838-2a6a33ce09ff
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	oras.land/oras-go v0.4.0
)
