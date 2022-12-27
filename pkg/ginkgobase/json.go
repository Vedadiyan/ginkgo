package ginkgobase

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type JSON string

func (json JSON) ToProtobuffer(v any) error {
	protoMessage, ok := v.(protoreflect.ProtoMessage)
	if !ok {
		return INVALID_PROTOMESSAGE_TYPE
	}
	return protojson.Unmarshal([]byte(json), protoMessage)
}

func (json JSON) Minify() (string, error) {
	buffer := bytes.NewBufferString("")
	len := len(json)
	hold := false
	jump := false
	for i := 0; i < len; i++ {
		r := rune(json[i])
		switch r {
		case '"':
			{
				if !jump {
					hold = !hold
				}
				jump = false
			}
		case '\\':
			{
				if i < len-1 {
					if rune(json[i+1]) == '"' {
						jump = true
					}
				} else {
					return "", fmt.Errorf("invalid end of json")
				}
			}
		case ' ':
			fallthrough
		case '\r':
			fallthrough
		case '\n':
			fallthrough
		case '\t':
			{
				if !hold {
					continue
				}
			}
		}
		buffer.WriteRune(r)
	}
	return buffer.String(), nil
}
