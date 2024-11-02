package io

import "bytes"

type Stream struct {
	buf *bytes.Buffer
}

func NewStream(buf *bytes.Buffer) *Stream {
	return &Stream{
		buf: buf,
	}
}
