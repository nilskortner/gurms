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

func (s *Stream) WriteBoolean(value bool) *Stream {
	// protobuf marshalling value
	//s.buf.Write()
	return s
}

func (s *Stream) ReadBoolean() {

}
