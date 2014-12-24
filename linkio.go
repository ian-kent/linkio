// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package provides an io.Reader that returns data,
// simulating a network connection of a certain speed.
package linkio

import (
	"io"
	"time"
)

// A LinkReader wraps an io.Reader, simulating reading from a
// shared access link with a fixed maximum speed.
type LinkReader struct {
	r    io.Reader
	link *Link
}

// A Link serializes requests to sleep, simulating the way data travels
// across a link which is running at a certain kbps (kilo = 1024).
// Multiple LinkReaders can share a link (simulating multiple apps
// sharing a link). The sharing behavior is approximately fair, as implemented
// by Go when scheduling reads from a contested blocking channel.
type Link struct {
	in    chan linkRequest
	speed int64 // nanosec per bit
}

// A linkRequest asks the link to simulate sending that much data
// and return a true on the channel when it has accomplished the request.
type linkRequest struct {
	bytes int
	done  chan bool
}

// NewLinkReader returns a LinkReader that returns bytes from r,
// simulating that they arrived from a shared link.
func (link *Link) NewLinkReader(r io.Reader) (s *LinkReader) {
	s = &LinkReader{r: r, link: link}
	return
}

// NewLink returns a new Link running at kbps.
func NewLink(kbps int) (l *Link) {
	// allow up to 100 outstanding requests
	l = &Link{in: make(chan linkRequest, 100)}
	_ = l.SetSpeed(kbps)

	// This goroutine serializes the requests. He could calculate
	// link utilization by comparing the time he sleeps waiting for
	// linkRequests to arrive and the time he spends sleeping to simulate
	// traffic flowing.

	go func() {
		for lr := range l.in {
			// bits * nanosec/bit = nano to wait
			delay := time.Duration(int64(lr.bytes*8) * l.speed)
			time.Sleep(delay)
			lr.done <- true
		}
	}()

	return
}

// SetSpeed set current speed of the bytes returned, returning the old speed.
// The speed is expressed in kilobits per second, where kilo = 1024.
func (l *Link) SetSpeed(kbps int) int {
	old := 1e9 * l.speed / 1024
	// l.speed is stored in ns/bit
	l.speed = 1e9 / int64(kbps*1024)
	return int(old)
}

// why isn't this in package math? hmm.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Satisfies interface io.Reader.
func (l *LinkReader) Read(buf []byte) (n int, err error) {
	// Read small chunks at a time, even if they ask for more,
	// preventing one LinkReader from saturating the simulated link.
	// 1500 is the MTU for Ethernet, i.e. a likely maximum packet
	// size.
	toRead := min(len(buf), 1500)
	n, err = l.r.Read(buf[0:toRead])
	if err != nil {
		return 0, err
	}

	// send in the request to sleep to the Link and sleep
	lr := linkRequest{bytes: n, done: make(chan bool)}
	l.link.in <- lr
	_ = <-lr.done

	return
}
