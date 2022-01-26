# lognsq

[LogNSQ] performs log redirection and aggregation. It reads logs redirected by
a pipe from a server and uses NSQ messaging service to collect, filter and
store the logs in a centralized (or distributed) location.

## Introduction

Logs contain an important information about running services, and are great
for analysis, or for figuring out what caused a problem with some functionality.
Practically all services are able to generate logs.
Sometimes it is beneficial to collect logs into a central repository.

- a central repository simplifies logs handling. It is easier to read them,
  back them up, run them through such programs like [Logstash] for example.
- quite often applications run inside of virtual machines, Docker instances,
  Kubernetes pods, where logs are isolated from the host and its file system.
  Sending logs to central repository makes them much more accessible.
- Quite often there are several instances of a service running simultaneously,
  orchestrated by a load balancer. Aggregation of logs to a central
  repository allows to see whole picture at once.
- It is often beneficial to clean up logs from unimportant information,
  partition bot accesses from "real" use, separate different parts of a service
  into separate "bins" (for example web UI calls from API calls).

[LogNSQ] provides means for logs aggregation, partition, filtering. It saves
logs as messages to [NSQ] service. Independent application, such as a
`nsq_to_file` or a custom script can collect the logs and prepare them for
debugging, analysis, or backup.

<!-- vim-markdown-toc GFM -->

* [TLDR](#tldr)
* [Installation](#installation)
  * [Linux or OS X](#linux-or-os-x)
  * [Windows](#windows)
  * [Install with Go](#install-with-go)
* [Usage](#usage)
* [Flags](#flags)

<!-- vim-markdown-toc -->

## TLDR

If a service sends logs to STDIN:

```bash
myservice | lognsq --topic='mylogs' --nsqd-tcp-address='127.0.0.1:4150'
```

If a service sends logs to STDERR:

```bash
myservice 2>&1 | lognsq --topic='mylogs' --nsqd-tcp-address='127.0.0.1:4150'
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

## Usage

[LogNSQ] is used by consuming log lines from STDIN and redirecting them to
NSQ service and (optionally) to STDERR.

**IMPORTANT** use single quotes where possible to avoid shell
interpolation of its special characters like '$', '\', '!' etc.

Usually logs are coming from STDERR and need to be redirected to STDIN:

```bash
myservice 2>&1 | lognsq -t 'mylogs' -a 'localhost:4150'
```

To print logs to STDERR as well as sending them to an nsqd service:

```bash
myservice 2>&1 | lognsq -t 'mylogs' -a 'localhost:4150' -p
```

To filter bots, split logs from the same service to different topics

```bash
myservice 2>&1 | grep -v 'bot' | \
  lognsq -t 'web' -a 'localhost:4159' -p -c '!/api/v1' 2>&1 | \
  lognsq -t 'api' -a 'localhost:4159' -p -c '/api/v1'
```

## Flags

`--help -h`
: displays help message.

```bash
lognsq -h
```

`--topic -t`
: sets the `topic` to which messages will be sent to nsqd server (**required**).

```bash
myapp 2>&1 | lognsq --topic='web' --nsqd-tcp-address='localhost:4150'
myapp 2>&1 | lognsq -t 'web' -a 'localhost:4150'
```

`--nsqd-tcp-address -a`
: the address and port of nsqd TCP service (**required**).

```bash
myapp 2>&1 | lognsq --topic='web' --nsqd-tcp-address='localhost:4150'
myapp 2>&1 | lognsq -t 'web' -a 'localhost:4150'
```

`--contains-filter' -c`
: filters log lines by matching positive or negative patterns. Negative
  patters have '!' as the first character (e.g. '/api/v1', '!/api/v1').
  If more complex pattern matching is required, use `--regex-filter'.
  If both filters are given, they add their effect.

`--regex-filter -r`
: filters log lines by matching them to the regular expression. If both
  `--contains-filter` and `--regex-filter` are given, they effect is
  cumulative. Negative lookahead expressions are not supported, use
  ``--contains-filter` negation like '!api' instead.

```bash
myapp 2>&1 | lognsq -t 'api' -a 'localhost:4150' --regex-filter='api/v1'
```

`--print-log -p`
: outputs all logs to STDERR again. The logs are unfiltered. This allows to
apply `lognsq` again with different filters and topics.

```bash
myservice 2>&1 | grep -v 'bot' | \
  lognsq -t 'web' -a 'localhost:4159' -p -c '!/api/v1' 2>&1 | \
  lognsq -t 'api' -a 'localhost:4159' -p -c '/api/v1'
```

`--debug -d`
: prints out filtered out logs that are sent to NSQ service and shows all logs
of NSQ interaction. Without `--debug` flag NSQ interaction log is suppressed.

[LogNSQ]: https://github.com/sfgrp/lognsq
[Logstash]: https://www.elastic.co/downloads/logstash
[latest release]: https://github.com/sfgrp/lognsq/releases/latest
[winpath]: https://www.computerhope.com/issues/ch000549.htm
[wsl]: https://docs.microsoft.com/en-us/windows/wsl/
