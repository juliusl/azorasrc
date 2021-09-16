# get-oras.sh

## Summary
This script is a download/package manager built on top of the github api for oras cli github releases. 
It has features like specifying a tag and repository to query for releases.
The script will download the the .tar.gz to the $HOME directory and unpack it to $HOME/oras-releases/$TAG.
Once the script has finished unpacking, it will output an install-script.sh which you can call with `. ./install-script.sh`. This script 
will set oras as an alias, and then call the alias to output the usage documentation. 

```opinion
You may add this to your sh profile or $PATH if you wish, but in my opinion that is overkill. I find sourcing the install script good
enough.
```

Note that, each subsequent execution of this script will overwrite the artifacts that are written to the home directory (including `install-script.sh`). 
This is so that the last release downloaded will always be the current one in $HOME.

## Details
At the moment, this script is scoped, to be used by developers interactively. So unmonitored scenarios are 
**NOT** taken into account. That being said, the script itself is quite short so it will probably take an 
experienced developer an hour or two to swap out the bits that aren't needed.

All of the directory paths can be configured may be configured via Environment Variables. 

After sourcing the install scrip, you can do:

```sh
set | grep ORAS_
```

Which will give you all the directories and files that were just created/downloaded, etc. This is useful, 
if you wanted to integrate this into an automation script that picks up the location and calls oras to do something. 
So basically you can think of it as a simple event hook.

Although this script is scoped to oras-project/oras, technically it can query any repo you pass through to REPO= when searching.
So feel free to copy/paste/reuse/recycle, it's good for the environment.

## Environment Variables Descriptions

### Search Parameters
```sh
    # Search Parameters
    REPO=${REPO:-'oras-project/oras'} # This is used to query a specific github repo 
    TAG=${TAG:-'v0.11.0'}  # This is used to query a specific git tag
    OS=${OS:-'linux'} # This is used to filter which .tar.gz to download
    ARCH=${ARCH:-'amd64'} # ""

    # Github API Auth Paramters
    # If you are developing on this script, you should set this to your github username. Password may be left blank. Otherwise you will be quickly throttled.
    # If unset, will display a warning in the output
    GITHUB_USERNAME=${GITHUB_USERNAME:-''} # Set this variable if planning on calling frequently
    GITHUB_PASSWORD=${GITHUB_PASSWORD:-''} # Password may be left blank.

    # WGET/CURL Config Parameters
    # These set up the options for curl and wget
    # Note:
    # I use both out of convenience even though curl would've been enough.
    # That is because I enjoy the `wget`'s semantics for downloading, 
    # and I enjoy `curl` semantics for calling external api's
    BASIC_AUTH=${BASIC_AUTH:-"-u $GITHUB_USERNAME:"}
    BASIC_AUTH_WGET=${BASIC_AUTH_WGET:-"--http-user=$GITHUB_USERNAME --http-password=$GITHUB_PASSWORD"}
```

```sh
    # Script Paramters
    ARCHIVE_NAME=${ARCHIVE_NAME:-"oras_$TAG_$OS_$ARCH.tar.gz"} # the release will be downloaded to this file
    ORAS_RELEASE_DIR=${ORAS_RELEASE_DIR:-"$HOME/oras-releases/$TAG"} # this is the release directory for the tag
    ORAS_INSTALL_DIR=${ORAS_INSTALL_DIR:-"/usr/local/bin"} # if you would like to use this from $PATH this is the correct install directory
    ORAS_FILE_NAME=${ORAS_FILE_NAME:-"oras"} # The name of the binary that the archive unpacks, if you wanted to use this with another repo you might want to set this as well
    ORAS_RELEASE_LOCATION=${ORAS_RELEASE_LOCATION:-"$ORAS_RELEASE_DIR/$ORAS_FILE_NAME"} # This is the full path to the current set release
    ORAS_INSTALLER_NAME=${ORAS_INSTALLER_NAME:-'install-oras.sh'} # The name of the install script that will be called after this script
```

Lastly, you can directly set the download link this script uses, and it will skip querying github altogether.

```sh
$ORAS_DOWNLOAD_LINK= # If unset, will search github for the release's download link
```

## Usages:

Specifying repository and tag - 
```sh
.//get-oras.sh | 
    REPO='juliusl/oras' TAG='v0.11.5-alpha' sh;. ./install-oras.sh
```

Pulling from github -
```sh
curl https://raw.githubusercontent.com/juliusl/azorasrc/main/scripts/get-oras.sh | 
sh;. ./install-oras.sh
```
