# go-rpm-build-lib

Build RPM Packages with Go.

***Note***: This library was forked from [mh-cbon/go-bin-rpm](https://github.com/mh-cbon/go-bin-rpm).

## Why?

* Generate RPM SPEC file
* Generate RPM package

## Getting Started

Build `go-rpm-builder` binary:

```
make
```

Next, review `go-rpm-builder` help:

```
$ go-rpm-builder -h

NAME:
   go-rpm-builder - RPM utilities in Go

USAGE:
   go-rpm-builder <cmd> <options>

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

Next, review more specific `go-rpm-builder` arguments, e.g.
`generate` or `generate-spec`:

```
$ bin/go-rpm-builder generate -h

NAME:
   go-rpm-builder generate - Generate the package

USAGE:
   go-rpm-builder generate [command options] [arguments...]

OPTIONS:
   --file value     Path to the rpm_config.json file (default: "rpm_config.json")
   -b value         Path to the build area (default: "pkg-build")
   -a value         Target architecture of the build
   -o value         File path to the resulting rpm file
   --version value  Target version of the build
   --help, -h       show help (default: false)
```
