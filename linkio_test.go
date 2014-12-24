// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linkio

import (
	"bytes"
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	// a dummy buffer full of zeros to send over the link
	var y [1000]byte
	buf := bytes.NewBuffer(y[:])

	lr := NewLink(Throughput(30) * KilobitPerSecond /* kbps */).NewLinkReader(buf)
	for {
		var x [1024]byte
		n, err := lr.Read(x[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error("err ", err)
		}
		t.Log("got", n, "bytes")
	}
}

type fw struct {
	buf []byte
}

func (f *fw) Write(buf []byte) (int, error) {
	f.buf = append(f.buf, buf...)
	return len(buf), nil
}

func TestWriter(t *testing.T) {
	// a fake writer
	w := &fw{}
	// a dummy buffer full of zeros to send over the link
	var y [1000]byte
	lw := NewLink(Throughput(30) * KilobitPerSecond /* kbps */).NewLinkWriter(w)
	n, err := lw.Write(y[:])
	if err != nil {
		t.Error("err ", err)
	}
	t.Log("wrote", n, "bytes")
	t.Log("got", len(w.buf), "bytes")
}
