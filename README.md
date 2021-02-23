CLL (Command Line Launcher)
============================

[![Go](https://github.com/ryhszk/cll/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ryhszk/cll/actions/workflows/go.yml)

This is a TUI command line launcher written in the Go lang.

# Support OS

- Linux
- Windows (TBD)

# Installation

```
$ git clone https://github.com/ryhszk/cll
$ cd cll
$ go install
```

# Usage

```
$ cll
Please select a command from next list.

[ ] ls -la
[>] free -h
[ ] top
[ ] ./count
[ ] dstat

Press q to quit.
              total        used        free      shared  buff/cache   available
Mem:          7.5Gi       3.6Gi       186Mi       229Mi       3.7Gi       3.5Gi
Swap:         7.8Gi       344Mi       7.4Gi
```
