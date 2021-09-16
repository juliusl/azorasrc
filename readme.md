# AZ-ORAS-RC

# Why?
- Motivation, registry interface in a cloudshell
    - Cloudshell's are a good place for doing registry tasks because they provide basic access to a cloud
    - However, they are not VM's so they are limited in scope, for example you can't expect `apt install` to work in a cloudshell
    - But, there are enough shell utilities that can be combined to provide a better experience
    - So what if, we create a toolkit to supercharge the cloudshell env, and then interface that toolkit to a registry client

- Hence, azorasrc is the first implementation of this idea for the azure cloud shell

# What?

## `.sh`
```
src/
    az/             // deals with contextual az features (az-cli, azure.json, system-identity, etc)
    oras/           // deals with installing/getting oras releases
    rc/             // deals with shell environment settings and setting up functional aliases
    rc/encpass.sh   // dependency: provides encryption via shell script, uses openssl
```
- TODO Explain folder architecture

## `.go`
```
cmd/
    view/           // resolve a manifest from a reference and output it to stdout
```
- TODO Write man pages for commands

# How?

- TODO Diagram would be really helpful to explain all the moving parts 
- TODO Explain how new access providers can be implemented (system-identity, azure-cli, etc)
- TODO Explain how different cloud providers can be implemented
