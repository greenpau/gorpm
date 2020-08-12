# gorpm

Build RPM Packages with Go.

***Note***: This library was forked from [mh-cbon/go-bin-rpm](https://github.com/mh-cbon/go-bin-rpm).

<!-- begin-markdown-toc -->
## Table of Contents

* [Purpose](#purpose)
* [Packagin Examples](#packagin-examples)
* [Getting Started](#getting-started)

<!-- end-markdown-toc -->

## Purpose

* Generate RPM SPEC file
* Generate RPM package

## Packagin Examples

The `assets/projects` contains real-life examples of packaging:

* Prometheus Node Exporter
* Prometheus Blackbox Exporter

First, try building the `node_exporter`. If you experience failures,
please open Github Issue to get assistance.

Build `node_exporter` project in the following way:

```
cd assets/projects/node_exporter
make
```

The RPM will be placed into `dist` directory:

```
SCP:       scp ./dist/node-exporter-0.18.1-5.el7.x86_64.rpm root@remote:/tmp/
Install:   sudo yum -y localinstall ./dist/node-exporter-0.18.1-5.el7.x86_64.rpm
RPM File:  ./dist/node-exporter-0.18.1-5.el7.x86_64.rpm
```

## Getting Started

Build `gorpm` binary:

```
make
```

Next, review `gorpm` help:

```
$ gorpm -h

NAME:
   gorpm - RPM utilities in Go

USAGE:
   gorpm <cmd> <options>

VERSION:
   1.0.0

COMMANDS:
   generate-spec  Generate the SPEC file
   generate       Generate the package
   test           Test the package json file
   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

Next, review more specific `gorpm` arguments, e.g.
`generate` or `generate-spec`:

```
$ gorpm generate-spec --help
NAME:
   gorpm generate-spec - Generate the SPEC file

USAGE:
   gorpm generate-spec [command options] [arguments...]

OPTIONS:
   --file value     Path to the config.json file (default: "config.json")
   --arch value     Target CPU architecture of the build, e.g. amd64
   --version value  Target version of the build
   --release value  Target release of the build
   --distro value   Target distribution of the build
   --cpu value      Target CPU Instruction Set Architecture (ISA) of the build, e.g. x86_64
   --output value   File path to the resulting RPM .spec file
   --help, -h       show help (default: false)
```
