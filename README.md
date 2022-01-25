# lognsq

This app takes lines from STDIN and forwards them as messages to an NSQ
messaging service.

It is useful for a reliable aggregation of log messages from several similar
services, for grabbing logs from Docker instances of Kubernetes pods.

<!-- vim-markdown-toc GFM -->

* [TLDR](#tldr)
* [Installation](#installation)
  * [Linux or OS X](#linux-or-os-x)
  * [Windows](#windows)
  * [Install with Go](#install-with-go)

<!-- vim-markdown-toc -->

## TLDR

If a service sends logs to STDIN:

```bash
myservice | lognsq --topic="mylogs" --nsqd-tcp-address="127.0.0.1:4150"
```

If a service sends logs to STDERR:

```bash
myservice 2>&1 | lognsq --topic="mylogs" --nsqd-tcp-address="127.0.0.1:4150"
```

## Installation

Download and uncompress the app from the [latest release].

### Linux or OS X

Move ``lognsq`` executable somewhere in your PATH
(for example ``/usr/local/bin``)

```bash
sudo mv path_to/lognsq /usr/local/bin
```

### Windows

One possible way would be to create a default folder for executables and place
``lognsq`` there.

Use ``Windows+R`` keys
combination and type "``cmd``". In the appeared terminal window type:

```cmd
mkdir C:\bin
copy path_to\lognsq.exe C:\bin
```

[Add ``C:\bin`` directory to your ``PATH``][winpath] environment variable.

It is also possible to install [Windows Subsystem for Linux][wsl] on Windows
10, and use ``lognsq`` as a Linux executable.

### Install with Go

If you have Go installed on your computer use

```bash
go get -u github.com/gnames/lognsq/lognsq
```

For development install gnu make and use the following:

```bash
git clone https://github.com/gnames/lognsq.git
cd lognsq
make tools
make install
```

You do need your ``PATH`` to include ``$HOME/go/bin``


[latest release]: https://github.com/sfgrp/lognsq/releases/latest
[winpath]: https://www.computerhope.com/issues/ch000549.htm
[wsl]: https://docs.microsoft.com/en-us/windows/wsl/
