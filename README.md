linkio [![GoDoc](https://godoc.org/github.com/ian-kent/linkio?status.svg)](https://godoc.org/github.com/ian-kent/linkio) [![Build Status](https://travis-ci.org/ian-kent/linkio.svg?branch=master)](https://travis-ci.org/ian-kent/linkio)
======

linkio provides an io.Reader and io.Writer that simulate a network connection of a certain speed, e.g. to simulate a mobile connection.

### Quick start

You can use `linkio` to wrap existing io.Reader and io.Writer interfaces:

```go
// Create a new link at 512kbps
link = linkio.NewLink(512 * linkio.KilobitPerSecond)

// Open a connection
conn, err := net.Dial("tcp", "google.com:80")
if err != nil {
  // handle error
}

// Create a link reader/writer
linkReader := link.NewLinkReader(io.Reader(conn))
linkWriter := link.NewLinkWriter(io.Writer(conn))

// Use them as you would normally...
fmt.Fprintf(linkWriter, "GET / HTTP/1.0\r\n\r\n")
status, err := bufio.NewReader(linkReader).ReadString('\n')

```

### History and license

This repository is a fork of [Jeff R. Allen's
linkio](https://github.com/jeffallen/jra-go/tree/master/linkio). linkio
was brought into the world to help make a proxy to simulate slow
Internet links (see [this blog
posting](https://blog.nella.org/a-rate-limiting-http-proxy-in-go/)).

Jeff's linkio was licensed via the [BSD 3-clause
license](http://opensource.org/licenses/BSD-3-Clause).

Any modifications since the initial commit are Copyright ©‎ 2014, Ian
Kent (http://iankent.uk), and are released under the terms of the [MIT
License](http://opensource.org/licenses/MIT).
