# bpfps

[![Travis CI](https://travis-ci.org/genuinetools/bpfps.svg?branch=master)](https://travis-ci.org/genuinetools/bpfps)

A tool to list and diagnose bpf programs. (Who watches the watchers..? :)

Shoutout to [cilium's](https://github.com/cilium/cilium) 
[golang bpf package](https://godoc.org/github.com/cilium/cilium/pkg/bpf) for doing a lot of heavy lifting here.

## Installation

#### Binaries

- **linux** [amd64](https://github.com/genuinetools/bpfps/releases/download/v0.1.1/bpfps-linux-amd64) 

#### Via Go

```bash
$ go get github.com/genuinetools/bpfps
```

#### Using your package manager

- ArchLinux: [precompiled binary](https://aur.archlinux.org/packages/bpfps-bin), [compiled from source](https://aur.archlinux.org/packages/bpfps-git)

## Usage

```console
$ bpfps -h
 _            __
| |__  _ __  / _|_ __  ___
| '_ \| '_ \| |_| '_ \/ __|
| |_) | |_) |  _| |_) \__ \
|_.__/| .__/|_| | .__/|___/
      |_|       |_|

 A tool to list and diagnose bpf programs.  (Who watches the watchers..? :)
 Version: v0.1.1
 Build: be7363d

  -d    run in debug mode
  -v    print version and exit (shorthand)
  -version
        print version and exit
```

```console
$ sudo bpfps                                                                                                             
BID                 NAME                TYPE                UID                 MAPS                LOADTIME
21                                      Array               0                   []uint32(nil)       17h25m55.523229433s
22                                      Array               0                   []uint32(nil)       17h25m55.530713603s
```
