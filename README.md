cla (Under development :construction:)
===

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ryhszk/cla)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryhszk/cla)](https://goreportcard.com/report/github.com/ryhszk/cla)
[![Go](https://github.com/ryhszk/cll/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ryhszk/cll/actions/workflows/go.yml)
[![GitHub](https://img.shields.io/github/license/ryhszk/cll)](https://github.com/ryhszk/cll/blob/main/LICENSE)

# Description

cla is a Command line based LAuncher.

(Under development. I may not be able to achieve what I want to do with this project, so it will be on hiatus for a while)

# Support OS

- Linux

# Installation

```
$ go get github.com/ryhszk/cla/cmd/cla
```

# Usage

1. Start `cla` command (Launches a TUI APP that looks like next).
    ```
    $ cla
    +--------------+
    | MODE: Normal | 
    +--------------+
    |  0:   free -h
    |  1: > dstat -c -C 0,1,2
    +---------------------------------------------+
    | ctrl+c            | Exit.                   |
    | enter             | Execute selected line.  |
    | ctrl+s            | Save all lines.         |
    | ctrl+a            | Add a line at end.      |
    | ctrl+d            | Remove current line.    |
    | ↓ [tab]           | Move down.              |
    | ↑ [shift+tab]     | Move up.                |
    +---------------------------------------------+
    ```
3. Move the cursor to the command line you want to execute.
4. Press `enter` to execute the command (current line).

:memo: The data of the registration command and the configuration file are stored in `$HOME/.cla`.

# Demo

![demo](https://github.com/ryhszk/cla/blob/main/assets/cla.gif)

:warning: If the registered command is too long, the display may be lost.
