package codecio

import (
	"bytes"
	"gurms/internal/infra/io"
)

type CodecStream struct {
	stream io.Stream
}

func NewCodecStream(buf *bytes.Buffer) *CodecStream {
	return &CodecStream{
		stream: *io.NewStream(buf),
	}
}
