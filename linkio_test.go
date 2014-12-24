// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linkio

import (
	"bytes"
	"io"
	"testing"
)

func TestOne(t *testing.T) {
	// a dummy buffer full of zeros to send over the link
	var y [1000]byte
	buf := bytes.NewBuffer(y[:])

	lr := NewLink(30 /* kbps */).NewLinkReader(buf)
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
