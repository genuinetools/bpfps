# bpfps

[![Travis CI](https://img.shields.io/travis/genuinetools/bpfps.svg?style=for-the-badge)](https://travis-ci.org/genuinetools/bpfps)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/bpfps)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/bpfps/total.svg?style=for-the-badge)](https://github.com/genuinetools/bpfps/releases)

A tool to list and diagnose bpf programs. (Who watches the watchers..? :)

Shoutout to [cilium's](https://github.com/cilium/cilium) 
[golang bpf package](https://godoc.org/github.com/cilium/cilium/pkg/bpf) for doing a lot of heavy lifting here.

<!-- toc -->

- [Installation](#installation)
    + [Binaries](#binaries)
    + [Via Go](#via-go)
    + [Using your package manager](#using-your-package-manager)
- [Usage](#usage)

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/genuinetools/bpfps/releases).

#### Via Go

```console
$ go get github.com/genuinetools/bpfps
```

#### Using your package manager

- ArchLinux: [precompiled binary](https://aur.archlinux.org/packages/bpfps-bin), [compiled from source](https://aur.archlinux.org/packages/bpfps-git)

## Usage

```console
$ bpfps -h
bpfps -  A tool to list and diagnose bpf programs. (Who watches the watchers..? :).

Usage: bpfps <command>

Flags:

  -d  enable debug logging (default: false)

Commands:

  version  Show the version information.
```

```console
$ sudo bpfps                                                                                                             
BID                 NAME                TYPE                UID                 MAPS                LOADTIME
21                                      Array               0                   []uint32(nil)       17h25m55.523229433s
22                                      Array               0                   []uint32(nil)       17h25m55.530713603s
```