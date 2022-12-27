package ginkgobase

import (
	"github.com/gin-gonic/gin"
)

type ReadMethod int

const (
	NONE   ReadMethod = iota
	QUERY             = 1 << iota
	PARAMS            = 2 << iota
	BODY              = 3 << iota
)

type Reader struct {
	ctx *gin.Context
}

func NewReader(ctx *gin.Context) *Reader {
	r := Reader{
		ctx: ctx,
	}
	return &r
}

func (r Reader) ReadParams() map[string]string {
	m := make(map[string]string)
	for _, p := range r.ctx.Params {
		m[p.Key] = p.Value
	}
	return m
}

func (r Reader) ReadQuery() map[string]any {
	m := make(map[string]any)
	for key, value := range r.ctx.Request.URL.Query() {
		len := len(value)
		if len == 0 {
			continue
		}
		if len > 1 {
			m[key] = value
			continue
		}
		m[key] = value[0]
	}
	return m
}

func (r Reader) ReadBody() (map[string]any, error) {
	stream := NewStream(r.ctx.Request.Body)
	return stream.ReadAsMap()
}

func (r Reader) ReadMixed(rt ReadMethod) (map[string]any, error) {
	m := make(map[string]any)
	if rt&QUERY == QUERY {
		for key, value := range r.ReadQuery() {
			m[key] = value
		}
	}
	if rt&PARAMS == PARAMS {
		for key, value := range r.ReadParams() {
			m[key] = value
		}
	}
	if rt&BODY == BODY {
		b, err := r.ReadBody()
		if err != nil {
			return nil, err
		}
		for key, value := range b {
			m[key] = value
		}
	}
	return m, nil
}
