package ginkgobase

import (
	"bytes"
	"encoding/json"
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Stream struct {
	isClosed bool
	stream   io.ReadCloser
}

func NewStream(stream io.ReadCloser) *Stream {
	s := Stream{
		stream:   stream,
		isClosed: false,
	}
	return &s
}

func (stream *Stream) ReadAsBuffer() (*bytes.Buffer, error) {
	defer func() {
		stream.isClosed = true
	}()
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, stream.stream)
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}

func (stream *Stream) ReadAsJSON(v any) error {
	buffer, err := stream.ReadAsBuffer()
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer.Bytes(), v)
}

func (stream *Stream) ReadAsProtobuffer(v any) error {
	protoMessage, ok := v.(protoreflect.ProtoMessage)
	if !ok {
		return INVALID_PROTOMESSAGE_TYPE
	}
	buffer, err := stream.ReadAsBuffer()
	if err != nil {
		return err
	}
	return protojson.Unmarshal(buffer.Bytes(), protoMessage)
}

func (stream *Stream) ReadAsMap() (map[string]any, error) {
	buffer, err := stream.ReadAsBuffer()
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	err = json.Unmarshal(buffer.Bytes(), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
