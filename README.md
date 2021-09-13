# xk6-tcp
- TCP extension components used in K6
- Link: https://k6.io/blog/extending-k6-with-xk6/#creating-the-k6-extension

## Requirements
- [Go Installed](https://golang.org/doc/install)

## Install xk6 
- Link: https://github.com/grafana/xk6

You Can [download binaries](https://github.com/grafana/xk6/releases) that are already compiled for your platform, or build xk6 from source:

```shell
$ go install go.k6.io/xk6/cmd/xk6@latest
```

## Command usage
The xk6 command has two primary uses:

1. Compile custom k6 binaries
2. A replacement for go run while developing k6 extensions

The xk6 command will use the latest version of k6 by default. You can customize this for all invocations by setting the K6_VERSION environment variable.

As usual with go command, the xk6 command will pass the GOOS, GOARCH, and GOARM environment variables through for cross-compilation.

### **Custom builds**
Syntax:
```shell
$ xk6 build [<k6_version>]
    [--output <file>]
    [--with <module[@version][=replacement]>...]
```
- <k6_version> is the core k6 version to build; defaults to K6_VERSION env variable or latest.
- --output changes the output file.
- --with can be used multiple times to add extensions by specifying the Go module name and optionally its version, similar to go get. Module name is required, but specific version and/or local replacement are optional.

Examples:
```shell
$ xk6 build \
    --with github.com/k6io/xk6-sql

$ xk6 build v0.29.0 \
    --with github.com/k6io/xk6-sql@v0.0.1

$ xk6 build \
    --with github.com/k6io/xk6-sql=../../my-fork

$ xk6 build \
    --with github.com/k6io/xk6-sql=.

$ xk6 build \
    --with github.com/k6io/xk6-sql@v0.0.1=../../my-fork

# Build using a k6 fork repository. Note that a version is required if
# XK6_K6_REPO is a URI.
$ XK6_K6_REPO=github.com/example/k6 xk6 build master \
    --with github.com/k6io/xk6-sql

# Build using a k6 fork repository from a local path. The version must be omitted
# and the path must be absolute.
$ XK6_K6_REPO="$PWD/../../k6" xk6 build \
    --with github.com/k6io/xk6-sql
```
```shell
Next, let's build the k6 binary. To use the published version of the extension run:
$ xk6 build v0.32.0 --with github.com/shanyongsy/xk6-tcp

Or if you're working with a local directory run the following, replacing the path as needed:
$ xk6 build v0.32.0 --with github.com/shanyongsy/xk6-tcp="/home/sy/xk6-tcp"
```


## K6 Script
[k6-tcp script example](https://github.com/shanyongsy/xk6-tcp/blob/main/example/loadtest/test_tcp.js)

## Tcp Server
[Test server](https://github.com/shanyongsy/tcp-server-client-go)

## From
[xk6 - Custom k6 Builder](https://github.com/grafana/xk6/blob/master/README.md)


