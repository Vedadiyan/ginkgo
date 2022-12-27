package ginkgo

import (
	"fmt"
	"strings"
)

const ROOT string = ""

type IRoute interface {
	Extend(suffix string) *route
	GetRoute() string
}

type route struct {
	prefix string
}

func NewRoute(prefix string) IRoute {
	return &route{prefix}
}

func (_route route) Extend(suffix string) *route {
	prefix := _route.prefix
	prefix = strings.TrimLeft(prefix, "/")
	prefix = strings.TrimRight(prefix, "/")
	_suffix := suffix
	_suffix = strings.TrimLeft(_suffix, "/")
	_suffix = strings.TrimRight(_suffix, "/")
	return &route{fmt.Sprintf("%s/%s", prefix, _suffix)}
}

func (_route route) GetRoute() string {
	return _route.prefix
}
