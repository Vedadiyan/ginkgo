package ginkgoproxies

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	ginkgo_base "github.com/vedadiyan/ginkgo/pkg/ginkgobase"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	_ERROR_TEMPLATE = `"%s": { "error": "%s" }`
)

type Proxy[TReqType protoreflect.ProtoMessage, TResType protoreflect.ProtoMessage, TFuncType ~func(req TReqType) (TResType, error)] struct {
	sc map[string]TFuncType
}

func NewProxy[TReqType protoreflect.ProtoMessage, TResType protoreflect.ProtoMessage, TFuncType ~func(req TReqType) (TResType, error)](sc TFuncType) *Proxy[TReqType, TResType, TFuncType] {
	p := Proxy[TReqType, TResType, TFuncType]{
		sc: map[string]TFuncType{
			"": sc,
		},
	}
	return &p
}

func (proxy Proxy[req, res, fn]) Make(rt ginkgo_base.ReadMethod) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		switch len(proxy.sc) {
		case 1:
			{
				proxy.one(rt)(ctx)
			}
		default:
			{
				proxy.many(rt)(ctx)
			}
		}
	}
}

func (proxy Proxy[req, res, fn]) request(rt ginkgo_base.ReadMethod) func(ctx *gin.Context) (*req, error) {
	return func(ctx *gin.Context) (*req, error) {
		reader := ginkgo_base.NewReader(ctx)
		data, err := reader.ReadMixed(rt)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		json := ginkgo_base.JSON(jsonBytes)
		var req req
		err = json.ToProtobuffer(&req)
		if err != nil {
			ctx.Error(err)
			return nil, err
		}
		return &req, nil
	}
}

func (proxy Proxy[req, res, fn]) many(rt ginkgo_base.ReadMethod) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if proxy.sc == nil {
			ctx.Error(fmt.Errorf("no proxies could be found"))
			return
		}
		req, err := proxy.request(rt)(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}
		ar := make([]string, 0)
		var wg sync.WaitGroup
		for key, sc := range proxy.sc {
			wg.Add(1)
			go func(key string, sc fn) {
				defer wg.Done()
				res, err := sc(*req)
				if err != nil {
					ar = append(ar, fmt.Sprintf(_ERROR_TEMPLATE, key, err.Error()))
					return
				}
				jsonBytes, err := protojson.Marshal(res)
				if err != nil {
					ar = append(ar, fmt.Sprintf(_ERROR_TEMPLATE, key, err.Error()))
					return
				}
				json := ginkgo_base.JSON(jsonBytes)
				ar = append(ar, fmt.Sprintf(`"%s": %s`, key, json))
			}(key, sc)
		}
		wg.Wait()
		ctx.Data(200, "application/json", []byte(fmt.Sprintf("{%s}", strings.Join(ar, ","))))
	}
}

func (proxy Proxy[req, res, fn]) one(rt ginkgo_base.ReadMethod) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		req, err := proxy.request(rt)(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}
		var key string
		for k := range proxy.sc {
			key = k
			break
		}
		res, err := proxy.sc[key](*req)
		if err != nil {
			ctx.Error(err)
			return
		}
		jsonBytes, err := protojson.Marshal(res)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Data(200, "application/json", jsonBytes)
	}
}
