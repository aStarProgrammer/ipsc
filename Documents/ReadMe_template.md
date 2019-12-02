# go-ipfs

[![banner](https://camo.githubusercontent.com/140c9a1fac2bafdc2f2a34208b5639edb6fe08e8/68747470733a2f2f697066732e696f2f697066732f516d566b37737272776168584c4e6d634459767955454a7074796f78706e646e52613537594a31314c346a5632362f697066732e676f2e706e67)](https://camo.githubusercontent.com/140c9a1fac2bafdc2f2a34208b5639edb6fe08e8/68747470733a2f2f697066732e696f2f697066732f516d566b37737272776168584c4e6d634459767955454a7074796f78706e646e52613537594a31314c346a5632362f697066732e676f2e706e67)

[![img](https://camo.githubusercontent.com/6d45f10dbc623fab8ed429bc1b4792649b798dda/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f6d61646525323062792d50726f746f636f6c2532304c6162732d626c75652e7376673f7374796c653d666c61742d737175617265)](http://ipn.io) [![img](https://camo.githubusercontent.com/64237ecf8e7bc07d91920a76ab70f10daa1c4585/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f70726f6a6563742d495046532d626c75652e7376673f7374796c653d666c61742d737175617265)](http://ipfs.io/) [![Matrix](https://camo.githubusercontent.com/2e6197e2a8b85d0456e8f2d19eda8fb6fe0b3c5c/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f6d61747269782d253233697066732533416d61747269782e6f72672d626c75652e7376673f7374796c653d666c61742d737175617265)](https://matrix.to/#/room/#ipfs:matrix.org) [![IRC](https://camo.githubusercontent.com/a45f0b8ee35c72ca687082b5e7e8dabaf3ad7d5b/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f667265656e6f64652d253233697066732d626c75652e7376673f7374796c653d666c61742d737175617265)](http://webchat.freenode.net/?channels=%23ipfs) [![Discord](https://camo.githubusercontent.com/b45eaf7f4ecf6793fbb6b524c25e9c35d023d6a4/68747470733a2f2f696d672e736869656c64732e696f2f646973636f72642f3437353738393333303338303438383730373f636f6c6f723d626c756576696f6c6574266c6162656c3d646973636f7264267374796c653d666c61742d737175617265)](https://discord.gg/24fmuwR) [![standard-readme compliant](https://camo.githubusercontent.com/a7e665f337914171fa0b60a110690af78fc5d943/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f7374616e646172642d2d726561646d652d4f4b2d677265656e2e7376673f7374796c653d666c61742d737175617265)](https://github.com/RichardLitt/standard-readme) [![GoDoc](https://camo.githubusercontent.com/2dde6c441cc79e2cf16685b1b1fd828d2b502bb5/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f697066732f676f2d697066733f7374617475732e737667)](https://godoc.org/github.com/ipfs/go-ipfs) [![Build Status](https://camo.githubusercontent.com/c6c57202211e4699c34b6f06a4a662d84e69fe40/68747470733a2f2f636972636c6563692e636f6d2f67682f697066732f676f2d697066732e7376673f7374796c653d737667)](https://circleci.com/gh/ipfs/go-ipfs)

## 

## What is IPFS?

IPFS is a global, versioned, peer-to-peer filesystem. It combines  good ideas from previous systems such Git, BitTorrent, Kademlia, SFS,  and the Web. It is like a single bittorrent swarm, exchanging git  objects. IPFS provides an interface as simple as the HTTP web, but with  permanence built in. You can also mount the world at /ipfs.

For more info see: https://docs.ipfs.io/introduction/overview/

Before opening an issue, consider using one of the following locations to ensure you are opening your thread in the right place:

- go-ipfs *implementation* bugs in [this repo](https://github.com/ipfs/go-ipfs/issues).
- Documentation issues in [ipfs/docs issues](https://github.com/ipfs/docs/issues).
- IPFS *design* in [ipfs/specs issues](https://github.com/ipfs/specs/issues).
- Exploration of new ideas in [ipfs/notes issues](https://github.com/ipfs/notes/issues).
- Ask questions and meet the rest of the community at the [IPFS Forum](https://discuss.ipfs.io).

## 

## Table of Contents

- [Security Issues](https://github.com/ipfs/go-ipfs#security-issues)
- Install
  - [System Requirements](https://github.com/ipfs/go-ipfs#system-requirements)
  - [Install prebuilt packages](https://github.com/ipfs/go-ipfs#install-prebuilt-packages)
  - [From Linux package managers](https://github.com/ipfs/go-ipfs#from-linux-package-managers)
  - Build from Source
    - [Install Go](https://github.com/ipfs/go-ipfs#install-go)
    - [Download and Compile IPFS](https://github.com/ipfs/go-ipfs#download-and-compile-ipfs)
    - [Troubleshooting](https://github.com/ipfs/go-ipfs#troubleshooting)
  - [Updating go-ipfs](https://github.com/ipfs/go-ipfs#updating-go-ipfs)
- Getting Started
  - [Some things to try](https://github.com/ipfs/go-ipfs#some-things-to-try)
  - [Usage](https://github.com/ipfs/go-ipfs#usage)
  - [Running IPFS inside Docker](https://github.com/ipfs/go-ipfs#running-ipfs-inside-docker)
  - [Troubleshooting](https://github.com/ipfs/go-ipfs#troubleshooting-1)
- [Packages](https://github.com/ipfs/go-ipfs#packages)
- Development
  - [CLI, HTTP-API, Architecture Diagram](https://github.com/ipfs/go-ipfs#cli-http-api-architecture-diagram)
  - [Testing](https://github.com/ipfs/go-ipfs#testing)
  - [Development Dependencies](https://github.com/ipfs/go-ipfs#development-dependencies)
- [Contributing](https://github.com/ipfs/go-ipfs#contributing)
- [License](https://github.com/ipfs/go-ipfs#license)

## 

## Security Issues

The IPFS protocol and its implementations are still in heavy  development. This means that there may be problems in our protocols, or  there may be mistakes in our implementations. And -- though IPFS is not  production-ready yet -- many people are already running nodes in their  machines. So we take security vulnerabilities very seriously. If you  discover a security issue, please bring it to our attention right away!

If you find a vulnerability that may affect live deployments -- for  example, by exposing a remote execution exploit -- please send your  report privately to [security@ipfs.io](mailto:security@ipfs.io). Please DO NOT file a public issue. The GPG key for [security@ipfs.io](mailto:security@ipfs.io) is [4B9665FB 92636D17 7C7A86D3 50AAE8A9 59B13AF3](https://pgp.mit.edu/pks/lookup?op=get&search=0x50AAE8A959B13AF3).

If the issue is a protocol weakness that cannot be immediately exploited or something not yet deployed, just discuss it openly.

## 

## Install

The canonical download instructions for IPFS are over at: https://docs.ipfs.io/guides/guides/install/. It is **highly suggested** you follow those instructions if you are not interested in working on IPFS development.

### 

### System Requirements

IPFS can run on most Linux, macOS, and Windows systems. We recommend  running it on a machine with at least 2 GB of RAM and 2 CPU cores  (go-ipfs is highly parallel). On systems with less memory, it may not be completely stable.

If your system is resource constrained, we recommend:

1. Installing OpenSSL and rebuilding go-ipfs manually with `make build GOFLAGS=-tags=openssl`. See the [download and compile](https://github.com/ipfs/go-ipfs#download-and-compile-ipfs) section for more information on compiling go-ipfs.
2. Initializing your daemon with `ipfs init --profile=lowpower`

### 

### Install prebuilt packages

We host prebuilt binaries over at our [distributions page](https://ipfs.io/ipns/dist.ipfs.io#go-ipfs).

From there:

- Click the blue "Download go-ipfs" on the right side of the page.
- Open/extract the archive.
- Move `ipfs` to your path (`install.sh` can do it for you).

You can also download go-ipfs from this project's GitHub releases page if you are unable to access ipfs.io.

### 

### From Linux package managers

- [Arch Linux](https://github.com/ipfs/go-ipfs#arch-linux)
- [Nix](https://github.com/ipfs/go-ipfs#nix)
- [Snap](https://github.com/ipfs/go-ipfs#snap)

#### 

#### Arch Linux

In Arch Linux go-ipfs is available as [go-ipfs](https://www.archlinux.org/packages/community/x86_64/go-ipfs/) package.

```
$ sudo pacman -S go-ipfs
```

Development version of go-ipfs is also on AUR under [go-ipfs-git](https://aur.archlinux.org/packages/go-ipfs-git/). You can install it using your favourite AUR Helper or manually from AUR.

#### 

#### Nix

For Linux and MacOSX you can use the purely functional package manager [Nix](https://nixos.org/nix/):

```
$ nix-env -i ipfs
```

You can also install the Package by using it's attribute name, which is also `ipfs`.

#### 

#### Guix

GNU's functional package manager, [Guix](https://www.gnu.org/software/guix/), also provides a go-ipfs package:

```
$ guix package -i go-ipfs
```

#### 

#### Snap

With snap, in any of the [supported Linux distributions](https://snapcraft.io/docs/core/install):

```
$ sudo snap install ipfs
```

### 

### From Windows package managers

- [Chocolatey](https://github.com/ipfs/go-ipfs#chocolatey)
- [Scoop](https://github.com/ipfs/go-ipfs#scoop)

#### 

#### Chocolatey

The package [ipfs](https://chocolatey.org/packages/ipfs) currently points to go-ipfs and is being maintained.

```
PS> choco install ipfs
```

#### 

#### Scoop

Scoop provides `go-ipfs` in its 'extras' bucket.

```
PS> scoop bucket add extras
PS> scoop install go-ipfs
```

### 

### Build from Source

go-ipfs's build system requires Go 1.13 and some standard POSIX build tools:

- GNU make
- Git
- GCC (or some other go compatible C Compiler) (optional)

To build without GCC, build with `CGO_ENABLED=0` (e.g., `make build CGO_ENABLED=0`).

#### 

#### Install Go

The build process for ipfs requires Go 1.12 or higher. If you don't have it: [Download Go 1.12+](https://golang.org/dl/).

You'll need to add Go's bin directories to your `$PATH` environment variable e.g., by adding these lines to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile`:

```
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
```

(If you run into trouble, see the [Go install instructions](https://golang.org/doc/install)).

#### 

#### Download and Compile IPFS

```
$ git clone https://github.com/ipfs/go-ipfs.git

$ cd go-ipfs
$ make install
```

Alternatively, you can run `make build` to build the go-ipfs binary (storing it in `cmd/ipfs/ipfs`) without installing it.

**NOTE:** If you get an error along the lines of "fatal  error: stdlib.h: No such file or directory", you're missing a C  compiler. Either re-run `make` with `CGO_ENABLED=0` or install GCC.

##### 

##### Cross Compiling

Compiling for a different platform is as simple as running:

```
make build GOOS=myTargetOS GOARCH=myTargetArchitecture
```

##### 

##### OpenSSL

To build go-ipfs with OpenSSL support, append `GOFLAGS=-tags=openssl` to your `make` invocation. Building with OpenSSL should significantly reduce the  background CPU usage on nodes that frequently make or receive new  connections.

Note: OpenSSL requires CGO support and, by default, CGO is disabled  when cross compiling. To cross compile with OpenSSL support, you must:

1. Install a compiler toolchain for the target platform.
2. Set the `CGO_ENABLED=1` environment variable.

#### 

#### Troubleshooting

- Separate [instructions are available for building on Windows](https://github.com/ipfs/go-ipfs/blob/master/docs/windows.md).
- `git` is required in order for `go get` to fetch all dependencies.
- Package managers often contain out-of-date `golang` packages. Ensure that `go version` reports at least 1.10. See above for how to install go.
- If you are interested in development, please install the development dependencies as well.
- *WARNING*: Older versions of OSX FUSE (for Mac OS X) can cause kernel panics when mounting!- We strongly recommend you use the [latest version of OSX FUSE](http://osxfuse.github.io/). (See https://github.com/ipfs/go-ipfs/issues/177)
- For more details on setting up FUSE (so that you can mount the filesystem), see the docs folder.
- Shell command completion is available in `misc/completion/ipfs-completion.bash`. Read [docs/command-completion.md](https://github.com/ipfs/go-ipfs/blob/master/docs/command-completion.md) to learn how to install it.
- See the [init examples](https://github.com/ipfs/website/tree/master/static/docs/examples/init) for how to connect IPFS to systemd or whatever init system your distro uses.

### 

### Updating go-ipfs

#### 

#### Using ipfs-update

IPFS has an updating tool that can be accessed through `ipfs update`. The tool is not installed alongside IPFS in order to keep that logic independent of the main codebase. To install `ipfs update`, [download it here](https://ipfs.io/ipns/dist.ipfs.io/#ipfs-update).

#### 

#### Downloading IPFS builds using IPFS

List the available versions of go-ipfs:

```
$ ipfs cat /ipns/dist.ipfs.io/go-ipfs/versions
```

Then, to view available builds for a version from the previous command ($VERSION):

```
$ ipfs ls /ipns/dist.ipfs.io/go-ipfs/$VERSION
```

To download a given build of a version:

```
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_darwin-386.tar.gz # darwin 32-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_darwin-amd64.tar.gz # darwin 64-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_freebsd-amd64.tar.gz # freebsd 64-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_linux-386.tar.gz # linux 32-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_linux-amd64.tar.gz # linux 64-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_linux-arm.tar.gz # linux arm build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_windows-amd64.zip # windows 64-bit build
```

## 

## Getting Started

See also: https://docs.ipfs.io/introduction/usage/

To start using IPFS, you must first initialize IPFS's config files on your system, this is done with `ipfs init`. See `ipfs init --help` for information on the optional arguments it takes. After initialization is complete, you can use `ipfs mount`, `ipfs add` and any of the other commands to explore!

### 

### Some things to try

Basic proof of 'ipfs working' locally:

```
echo "hello world" > hello
ipfs add hello
# This should output a hash string that looks something like:
# QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o
ipfs cat <that hash>
```

### 

### Usage

```
  ipfs - Global p2p merkle-dag filesystem.

  ipfs [<flags>] <command> [<arg>] ...

SUBCOMMANDS
  BASIC COMMANDS
    init          Initialize ipfs local configuration
    add <path>    Add a file to ipfs
    cat <ref>     Show ipfs object data
    get <ref>     Download ipfs objects
    ls <ref>      List links from an object
    refs <ref>    List hashes of links from an object

  DATA STRUCTURE COMMANDS
    block         Interact with raw blocks in the datastore
    object        Interact with raw dag nodes
    files         Interact with objects as if they were a unix filesystem

  ADVANCED COMMANDS
    daemon        Start a long-running daemon process
    mount         Mount an ipfs read-only mountpoint
    resolve       Resolve any type of name
    name          Publish or resolve IPNS names
    dns           Resolve DNS links
    pin           Pin objects to local storage
    repo          Manipulate an IPFS repository

  NETWORK COMMANDS
    id            Show info about ipfs peers
    bootstrap     Add or remove bootstrap peers
    swarm         Manage connections to the p2p network
    dht           Query the DHT for values or peers
    ping          Measure the latency of a connection
    diag          Print diagnostics

  TOOL COMMANDS
    config        Manage configuration
    version       Show ipfs version information
    update        Download and apply go-ipfs updates
    commands      List all available commands

  Use 'ipfs <command> --help' to learn more about each command.

  ipfs uses a repository in the local file system. By default, the repo is located
  at ~/.ipfs. To change the repo location, set the $IPFS_PATH environment variable:

    export IPFS_PATH=/path/to/ipfsrepo
```

### 

### Running IPFS inside Docker

An IPFS docker image is hosted at [hub.docker.com/r/ipfs/go-ipfs](https://hub.docker.com/r/ipfs/go-ipfs/). To make files visible inside the container you need to mount a host directory with the `-v` option to docker. Choose a directory that you want to use to import/export files from IPFS. You should also choose a directory to store IPFS files that will persist when you restart the container.

```
export ipfs_staging=</absolute/path/to/somewhere/>
export ipfs_data=</absolute/path/to/somewhere_else/>
```

Start a container running ipfs and expose ports 4001, 5001 and 8080:

```
docker run -d --name ipfs_host -v $ipfs_staging:/export -v $ipfs_data:/data/ipfs -p 4001:4001 -p 127.0.0.1:8080:8080 -p 127.0.0.1:5001:5001 ipfs/go-ipfs:latest
```

Watch the ipfs log:

```
docker logs -f ipfs_host
```

Wait for ipfs to start. ipfs is running when you see:

```
Gateway (readonly) server
listening on /ip4/0.0.0.0/tcp/8080
```

You can now stop watching the log.

Run ipfs commands:

```
docker exec ipfs_host ipfs <args...>
```

For example: connect to peers

```
docker exec ipfs_host ipfs swarm peers
```

Add files:

```
cp -r <something> $ipfs_staging
docker exec ipfs_host ipfs add -r /export/<something>
```

Stop the running container:

```
docker stop ipfs_host
```

When starting a container running ipfs for the first time with an empty data directory, it will call `ipfs init` to initialize configuration files and generate a new keypair. At this time, you can choose which profile to apply using the `IPFS_PROFILE` environment variable:

```
docker run -d --name ipfs_host -e IPFS_PROFILE=server -v $ipfs_staging:/export -v $ipfs_data:/data/ipfs -p 4001:4001 -p 127.0.0.1:8080:8080 -p 127.0.0.1:5001:5001 ipfs/go-ipfs:latest
```

It is possible to initialize the container with a swarm key file (`/data/ipfs/swarm.key`) using the variables `IPFS_SWARM_KEY` and `IPFS_SWARM_KEY_FILE`. The `IPFS_SWARM_KEY` creates `swarm.key` with the contents of the variable itself, whilst `IPFS_SWARM_KEY_FILE` copies the key from a path stored in the variable. The `IPFS_SWARM_KEY_FILE` **overwrites** the key generated by `IPFS_SWARM_KEY`.

```
docker run -d --name ipfs_host -e IPFS_SWARM_KEY=<your swarm key> -v $ipfs_staging:/export -v $ipfs_data:/data/ipfs -p 4001:4001 -p 127.0.0.1:8080:8080 -p 127.0.0.1:5001:5001 ipfs/go-ipfs:latest
```

The swarm key initialization can also be done using docker secrets **(requires docker swarm or docker-compose)**:

```
cat your_swarm.key | docker secret create swarm_key_secret -
docker run -d --name ipfs_host --secret swarm_key_secret -e IPFS_SWARM_KEY_FILE=/run/secrets/swarm_key_secret -v $ipfs_staging:/export -v $ipfs_data:/data/ipfs -p 4001:4001 -p 127.0.0.1:8080:8080 -p 127.0.0.1:5001:5001 ipfs/go-ipfs:latest
```

### 

### Troubleshooting

If you have previously installed IPFS before and you are running into problems getting a newer version to work, try deleting (or backing up  somewhere else) your IPFS config directory (~/.ipfs by default) and  rerunning `ipfs init`. This will reinitialize the config file to its defaults and clear out the local datastore of any bad entries.

Please direct general questions and help requests to our [forum](https://discuss.ipfs.io) or our IRC channel (freenode #ipfs).

If you believe you've found a bug, check the [issues list](https://github.com/ipfs/go-ipfs/issues) and, if you don't see your problem there, either come talk to us on IRC (freenode #ipfs) or file an issue of your own!

## 

## Packages

> This table is generated using the module [`package-table`](https://github.com/ipfs-shipyard/package-table) with `package-table --data=package-list.json`.

Listing of the main packages used in the IPFS ecosystem. There are also three specifications worth linking here:

| Name                                                         | CI/Travis                                                    | Coverage                                                     | Description                                               |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | --------------------------------------------------------- |
| **Files**                                                    |                                                              |                                                              |                                                           |
| [`go-unixfs`](https://github.com/ipfs/go-unixfs)             | [![Travis CI](https://camo.githubusercontent.com/944cf6d6cf011f1c4c6825c784a7c8408e6a2e0a/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d756e697866732e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-unixfs) | [![codecov](https://camo.githubusercontent.com/1486f136242f730501d071b872bd237552fc18c4/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d756e697866732f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-unixfs) | the core 'filesystem' logic                               |
| [`go-mfs`](https://github.com/ipfs/go-mfs)                   | [![Travis CI](https://camo.githubusercontent.com/931123f574bb3806aa661802eae4942433a32b61/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d6d66732e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-mfs) | [![codecov](https://camo.githubusercontent.com/1877d34814c02c409fed4fb81ed15068939b71da/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d6d66732f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-mfs) | a mutable filesystem editor for unixfs                    |
| [`go-ipfs-posinfo`](https://github.com/ipfs/go-ipfs-posinfo) | [![Travis CI](https://camo.githubusercontent.com/a23586b0295effb14ea83494c96b385a9319ebda/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d706f73696e666f2e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-posinfo) | [![codecov](https://camo.githubusercontent.com/14a17d349d9aa918f746381159e314380a2f9ab9/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d706f73696e666f2f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-posinfo) | helper datatypes for the filestore                        |
| [`go-ipfs-chunker`](https://github.com/ipfs/go-ipfs-chunker) | [![Travis CI](https://camo.githubusercontent.com/82555eeb5cf9f0b2fbca5b13d14d69139408b2e7/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d6368756e6b65722e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-chunker) | [![codecov](https://camo.githubusercontent.com/d58a2e2023c30ba8a9c38022e0033acbca877c7d/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d6368756e6b65722f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-chunker) | file chunkers                                             |
| **Exchange**                                                 |                                                              |                                                              |                                                           |
| [`go-ipfs-exchange-interface`](https://github.com/ipfs/go-ipfs-exchange-interface) | [![Travis CI](https://camo.githubusercontent.com/650cc8bcc9c483fc2f9a594f3d5b5125f4e384d6/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d65786368616e67652d696e746572666163652e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-exchange-interface) | [![codecov](https://camo.githubusercontent.com/3ee2a0751cea14dce5a3a6d33316e2b19bf770ab/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d65786368616e67652d696e746572666163652f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-exchange-interface) | exchange service interface                                |
| [`go-ipfs-exchange-offline`](https://github.com/ipfs/go-ipfs-exchange-offline) | [![Travis CI](https://camo.githubusercontent.com/f79015ab4aa01f42106bf05e79261fe35ae99acd/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d65786368616e67652d6f66666c696e652e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-exchange-offline) | [![codecov](https://camo.githubusercontent.com/183e63d758ccd05aec7fa680c1dd8748ec31a744/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d65786368616e67652d6f66666c696e652f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-exchange-offline) | (dummy) offline implementation of the exchange service    |
| [`go-bitswap`](https://github.com/ipfs/go-bitswap)           | [![Travis CI](https://camo.githubusercontent.com/58bb7e3ba872fa604d8e9e0c9f0daea344f6441f/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d626974737761702e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-bitswap) | [![codecov](https://camo.githubusercontent.com/7562e5c5960ccbf45da08f010c57b85facea0bd4/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d626974737761702f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-bitswap) | bitswap protocol implementation                           |
| [`go-blockservice`](https://github.com/ipfs/go-blockservice) | [![Travis CI](https://camo.githubusercontent.com/f37fac1a059909be8930cb042018d10b350375e1/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d626c6f636b736572766963652e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-blockservice) | [![codecov](https://camo.githubusercontent.com/d68223e17f6a7bd805d9fb1eaf4b06227dd93ccb/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d626c6f636b736572766963652f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-blockservice) | service that plugs a blockstore and an exchange together  |
| **Datastores**                                               |                                                              |                                                              |                                                           |
| [`go-datastore`](https://github.com/ipfs/go-datastore)       | [![Travis CI](https://camo.githubusercontent.com/50ad7dedc02603d8dfb7e62f3d83825efa01a562/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d6461746173746f72652e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-datastore) | [![codecov](https://camo.githubusercontent.com/5f4bf24f09b66f8e163ffa8c23ec9086e9f0ed70/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d6461746173746f72652f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-datastore) | datastore interfaces, adapters, and basic implementations |
| [`go-ipfs-ds-help`](https://github.com/ipfs/go-ipfs-ds-help) | [![Travis CI](https://camo.githubusercontent.com/fc5bd95a819501f31c61e5b0af576e22757031a6/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d64732d68656c702e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-ds-help) | [![codecov](https://camo.githubusercontent.com/7e5064ff975a28bbec098e1f6368ae8a5fb46821/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d64732d68656c702f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-ds-help) | datastore utility functions                               |
| [`go-ds-flatfs`](https://github.com/ipfs/go-ds-flatfs)       | [![Travis CI](https://camo.githubusercontent.com/078c48b4c76f466fdcba49763a66d94eb02aca15/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d64732d666c617466732e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ds-flatfs) | [![codecov](https://camo.githubusercontent.com/f66210dc121702db832b33cbb06357990e1f6b96/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d64732d666c617466732f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ds-flatfs) | a filesystem-based datastore                              |
| [`go-ds-measure`](https://github.com/ipfs/go-ds-measure)     | [![Travis CI](https://camo.githubusercontent.com/fee47bc9a434330ed6626d75011640204f5f4d35/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d64732d6d6561737572652e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ds-measure) | [![codecov](https://camo.githubusercontent.com/386a9410164754952564c2a2e34a7cc28184fb0e/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d64732d6d6561737572652f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ds-measure) | a metric-collecting database adapter                      |
| [`go-ds-leveldb`](https://github.com/ipfs/go-ds-leveldb)     | [![Travis CI](https://camo.githubusercontent.com/b1a9ed2c649a932dbce0a50d100c3d3639d6a59d/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d64732d6c6576656c64622e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ds-leveldb) | [![codecov](https://camo.githubusercontent.com/040ceb7e264a9a5683a654091d7bf6be580b885c/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d64732d6c6576656c64622f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ds-leveldb) | a leveldb based datastore                                 |
| [`go-ds-badger`](https://github.com/ipfs/go-ds-badger)       | [![Travis CI](https://camo.githubusercontent.com/e3d06361321d61f88e18d8f3edebcbac33d333ee/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d64732d6261646765722e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ds-badger) | [![codecov](https://camo.githubusercontent.com/34d773902a783259e8120c4213e551366ae056a9/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d64732d6261646765722f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ds-badger) | a badgerdb based datastore                                |
| **Namesys**                                                  |                                                              |                                                              |                                                           |
| [`go-ipns`](https://github.com/ipfs/go-ipns)                 | [![Travis CI](https://camo.githubusercontent.com/8b577a681cf383b4c44f7b436872b47c991196c6/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d69706e732e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipns) | [![codecov](https://camo.githubusercontent.com/84c1347d3faa1b675648ca9ab0cee0475abf862e/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d69706e732f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipns) | IPNS datastructures and validation logic                  |
| **Repo**                                                     |                                                              |                                                              |                                                           |
| [`go-ipfs-config`](https://github.com/ipfs/go-ipfs-config)   | [![Travis CI](https://camo.githubusercontent.com/878b8a0f6fd1aaada015e2ad70a139aed23832df/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d636f6e6669672e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-config) | [![codecov](https://camo.githubusercontent.com/c792f6698000fcfa6404eb094ce82c73c4495dae/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d636f6e6669672f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-config) | go-ipfs config file definitions                           |
| [`go-fs-lock`](https://github.com/ipfs/go-fs-lock)           | [![Travis CI](https://camo.githubusercontent.com/6bc7d404a9cfcc6433bfd6fbc68c2c0b166c606f/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d66732d6c6f636b2e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-fs-lock) | [![codecov](https://camo.githubusercontent.com/494de31b326fd0fd236bac3c4b7eed0dc943e9a3/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d66732d6c6f636b2f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-fs-lock) | lockfile management functions                             |
| [`fs-repo-migrations`](https://github.com/ipfs/fs-repo-migrations) | [![Travis CI](https://camo.githubusercontent.com/4393b4a03e5faa821934c8bc67a65c48b28100d7/68747470733a2f2f7472617669732d63692e636f6d2f697066732f66732d7265706f2d6d6967726174696f6e732e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/fs-repo-migrations) | [![codecov](https://camo.githubusercontent.com/f5e41e2a6edd14cb851fbc46e5d0c2681828410c/68747470733a2f2f636f6465636f762e696f2f67682f697066732f66732d7265706f2d6d6967726174696f6e732f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/fs-repo-migrations) | repo migrations                                           |
| **Blocks**                                                   |                                                              |                                                              |                                                           |
| [`go-block-format`](https://github.com/ipfs/go-block-format) | [![Travis CI](https://camo.githubusercontent.com/e5667880955744753c418732da6cdec69ecb3280/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d626c6f636b2d666f726d61742e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-block-format) | [![codecov](https://camo.githubusercontent.com/bb63729536cbea6076f830c0500f44c2952617f8/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d626c6f636b2d666f726d61742f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-block-format) | block interfaces and implementations                      |
| [`go-ipfs-blockstore`](https://github.com/ipfs/go-ipfs-blockstore) | [![Travis CI](https://camo.githubusercontent.com/d83387d0049e3fcb0c2dc9ba2cb614d0d7ec6283/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d626c6f636b73746f72652e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-blockstore) | [![codecov](https://camo.githubusercontent.com/15b77c5ad23445707fbdeb19263e13b5620834d1/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d626c6f636b73746f72652f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-blockstore) | blockstore interfaces and implementations                 |
| **Commands**                                                 |                                                              |                                                              |                                                           |
| [`go-ipfs-cmds`](https://github.com/ipfs/go-ipfs-cmds)       | [![Travis CI](https://camo.githubusercontent.com/65b349d04aaa4d685764816a75c6f28f869d2ad0/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d636d64732e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-cmds) | [![codecov](https://camo.githubusercontent.com/afee9c54558e0d22b2ca6326bb0ad139e276e2b9/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d636d64732f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-cmds) | CLI & HTTP commands library                               |
| [`go-ipfs-api`](https://github.com/ipfs/go-ipfs-api)         | [![Travis CI](https://camo.githubusercontent.com/df4da5b26743efd7c0a8e82df35be10b4bcfe440/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d6170692e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-api) | [![codecov](https://camo.githubusercontent.com/6dda58a220a376b379ab6c14070fba2c673afc27/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d6170692f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-api) | a shell for the IPFS HTTP API                             |
| **Metrics & Logging**                                        |                                                              |                                                              |                                                           |
| [`go-metrics-interface`](https://github.com/ipfs/go-metrics-interface) | [![Travis CI](https://camo.githubusercontent.com/cd51c4e49db406a2838e3caaeab7ab7dc539f8b1/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d6d6574726963732d696e746572666163652e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-metrics-interface) | [![codecov](https://camo.githubusercontent.com/935a975255481784c97550a9b5c4717f9c35c25f/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d6d6574726963732d696e746572666163652f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-metrics-interface) | metrics collection interfaces                             |
| [`go-metrics-prometheus`](https://github.com/ipfs/go-metrics-prometheus) | [![Travis CI](https://camo.githubusercontent.com/3e8cbfb26defb102ce0017ce05fc81ee972530d7/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d6d6574726963732d70726f6d6574686575732e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-metrics-prometheus) | [![codecov](https://camo.githubusercontent.com/f8228e7311451ecd348ca21e3f196dde4454cbb2/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d6d6574726963732d70726f6d6574686575732f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-metrics-prometheus) | prometheus-backed metrics collector                       |
| [`go-log`](https://github.com/ipfs/go-log)                   | [![Travis CI](https://camo.githubusercontent.com/2b8b8a30c6bbc88b15ce237becbf49ec7b607c03/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d6c6f672e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-log) | [![codecov](https://camo.githubusercontent.com/f3b4b7c3e62f076ad2f750ccfc8eeef34d31cd45/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d6c6f672f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-log) | logging framework                                         |
| **Generics/Utils**                                           |                                                              |                                                              |                                                           |
| [`go-ipfs-routing`](https://github.com/ipfs/go-ipfs-routing) | [![Travis CI](https://camo.githubusercontent.com/1114c4f839fec3f10179e8bec17d4bc88204edb5/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d726f7574696e672e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-routing) | [![codecov](https://camo.githubusercontent.com/65ddf9a0ea38f7b8e817d9012082ff964a007562/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d726f7574696e672f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-routing) | routing (content, peer, value) helpers                    |
| [`go-ipfs-util`](https://github.com/ipfs/go-ipfs-util)       | [![Travis CI](https://camo.githubusercontent.com/17faec5504b4450e22731de08e70449ced0d1bc6/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d7574696c2e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-util) | [![codecov](https://camo.githubusercontent.com/de2b4b7800d16130b07a5770903b6e753f49e8f3/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d7574696c2f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-util) | the kitchen sink                                          |
| [`go-ipfs-addr`](https://github.com/ipfs/go-ipfs-addr)       | [![Travis CI](https://camo.githubusercontent.com/2a72a69972958a66bc13450bebb03c99113f781b/68747470733a2f2f7472617669732d63692e636f6d2f697066732f676f2d697066732d616464722e7376673f6272616e63683d6d6173746572)](https://travis-ci.com/ipfs/go-ipfs-addr) | [![codecov](https://camo.githubusercontent.com/161d4b0c21b522c859f25ab43c11c860119fd921/68747470733a2f2f636f6465636f762e696f2f67682f697066732f676f2d697066732d616464722f6272616e63682f6d61737465722f67726170682f62616467652e737667)](https://codecov.io/gh/ipfs/go-ipfs-addr) | utility functions for parsing IPFS multiaddrs             |

For brevity, we've omitted go-libp2p and go-ipld packages. These  package tables can be found in their respective project's READMEs:

- [go-libp2p](https://github.com/libp2p/go-libp2p#packages)
- [go-ipld](https://github.com/ipld/go-ipld#packages)

## 

## Development

Some places to get you started on the codebase:

- Main file: [./cmd/ipfs/main.go](https://github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/main.go)
- CLI Commands: [./core/commands/](https://github.com/ipfs/go-ipfs/tree/master/core/commands)
- Bitswap (the data trading engine): [go-bitswap](https://github.com/ipfs/go-bitswap)
- libp2p 
  - libp2p: https://github.com/libp2p/go-libp2p
  - DHT: https://github.com/libp2p/go-libp2p-kad-dht
  - PubSub: https://github.com/libp2p/go-libp2p-pubsub
- [IPFS : The `Add` command demystified](https://github.com/ipfs/go-ipfs/tree/master/docs/add-code-flow.md)

### 

### Map of go-ipfs Subsystems

**WIP**: This is a high-level architecture diagram of  the various sub-systems of go-ipfs. To be updated with how they  interact. Anyone who has suggestions is welcome to comment [here](https://docs.google.com/drawings/d/1OVpBT2q-NtSJqlPX3buvjYhOnWfdzb85YEsM_njesME/edit) on how we can improve this! [![img](https://camo.githubusercontent.com/3759ceff46d10f763d6b288a574fe8edd6fd2d59/68747470733a2f2f646f63732e676f6f676c652e636f6d2f64726177696e67732f642f652f32504143582d3176535f6e3146765375366d646d5369726b427249494569623267716867746174443961776150325f576472474e347a544e65673632305851643950393557542d49766f676e5378494964434d3575452f7075623f773d3134343626683d31303336)](https://camo.githubusercontent.com/3759ceff46d10f763d6b288a574fe8edd6fd2d59/68747470733a2f2f646f63732e676f6f676c652e636f6d2f64726177696e67732f642f652f32504143582d3176535f6e3146765375366d646d5369726b427249494569623267716867746174443961776150325f576472474e347a544e65673632305851643950393557542d49766f676e5378494964434d3575452f7075623f773d3134343626683d31303336)

### 

### CLI, HTTP-API, Architecture Diagram

[![img](https://github.com/ipfs/go-ipfs/raw/master/docs/cli-http-api-core-diagram.png)](https://github.com/ipfs/go-ipfs/blob/master/docs/cli-http-api-core-diagram.png)

> [Origin](https://github.com/ipfs/pm/pull/678#discussion_r210410924)

Description: Dotted means "likely going away". The "Legacy" parts are thin wrappers around some commands to translate between the new system  and the old system. The grayed-out parts on the "daemon" diagram are  there to show that the code is all the same, it's just that we turn some pieces on and some pieces off depending on whether we're running on the client or the server.

### 

### Testing

```
make test
```

### 

### Development Dependencies

If you make changes to the protocol buffers, you will need to install the [protoc compiler](https://github.com/google/protobuf).

### 

### Developer Notes

Find more documentation for developers on [docs](https://github.com/ipfs/go-ipfs/blob/master/docs)

## 

## Contributing

[![img](https://camo.githubusercontent.com/2820cc493393fa993bef64b044c6d3ce1d4b56a4/68747470733a2f2f63646e2e7261776769742e636f6d2f6a62656e65742f636f6e747269627574652d697066732d6769662f6d61737465722f696d672f636f6e747269627574652e676966)](https://github.com/ipfs/community/blob/master/CONTRIBUTING.md)

We ❤️ all [our contributors](https://github.com/ipfs/go-ipfs/blob/master/docs/AUTHORS); this project wouldn’t be what it is without you! If you want to help out, please see [CONTRIBUTING.md](https://github.com/ipfs/go-ipfs/blob/master/CONTRIBUTING.md).

This repository falls under the IPFS [Code of Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

You can contact us on the freenode #ipfs-dev channel or attend one of our [weekly calls](https://github.com/ipfs/team-mgmt/issues/674).

## 

## License

[MIT](https://github.com/ipfs/go-ipfs/blob/master/LICENSE)