package ginkgo

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	GET Methods = iota
	POST
	PUT
	PATCH
	DELETE
)

var webApiStorage []*IWebAPI

type Methods int

type Filter func(ctx *gin.Context) bool

type Engine struct {
	instance *gin.Engine
}

type IWebAPI interface {
	Register(engine *gin.Engine)
}

type webApi struct {
	route       string
	method      Methods
	onExecuting []Filter
	handler     func(ctx *gin.Context)
	onExecuted  []Filter
	initFn      func()
}

func (webApi webApi) Method() Methods {
	return webApi.method
}

func (webApi webApi) Route() string {
	return webApi.route
}

func (webApi webApi) Handler() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if webApi.onExecuting != nil {
			for _, filter := range webApi.onExecuting {
				filter(ctx)
				if !filter(ctx) {
					return
				}
			}
		}
		webApi.handler(ctx)
		if webApi.onExecuted != nil {
			for _, filter := range webApi.onExecuted {
				if !filter(ctx) {
					return
				}
			}
		}
	}
}

func (webApi webApi) InitFn() func() {
	return webApi.initFn
}

func NewAPI(route IRoute, method Methods, preFilters []Filter, handler func(ctx *gin.Context), postFilters []Filter, initFn func()) *webApi {
	webApi := webApi{route.GetRoute(), method, preFilters, handler, postFilters, initFn}
	var router IWebAPI = webApi
	registerApi(&router)
	return &webApi
}

func (webApi webApi) Register(engine *gin.Engine) {
	if webApi.initFn != nil {
		webApi.initFn()
	}
	route := webApi.Route()
	route = strings.TrimRight(route, "/")
	route = strings.TrimLeft(route, "/")
	handler := webApi.Handler()
	switch webApi.Method() {
	case GET:
		{
			engine.GET(route, handler)
			break
		}
	case POST:
		{
			engine.POST(route, handler)
			break
		}
	case PATCH:
		{
			engine.PATCH(route, handler)
			break
		}
	case PUT:
		{
			engine.PUT(route, handler)
			break
		}
	case DELETE:
		{
			engine.DELETE(route, handler)
			break
		}
	}
}

func NewEngine() *Engine {
	var engine Engine
	engine.instance = gin.Default()
	return &engine
}

func (engine *Engine) DisableCORS() {
	engine.instance.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders: []string{"*"},
		MaxAge:       12 * time.Hour,
	}))
}

func (engine *Engine) ConfigureCORS(origins []string, headers []string, methods []string) {
	engine.instance.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: methods,
		AllowHeaders: headers,
		MaxAge:       12 * time.Hour,
	}))
}

func (engine *Engine) Use(middlewares ...gin.HandlerFunc) {
	engine.instance.Use(middlewares...)
}

func (engin *Engine) GetEngine() *gin.Engine {
	return engin.instance
}

func (engine *Engine) Start(addr ...string) {
	for _, webApi := range getAllApis() {
		(*webApi).Register(engine.instance)
	}
	engine.instance.Run(addr...)
}

func registerApi(webApi *IWebAPI) {
	webApiStorage = append(webApiStorage, webApi)
}

func getAllApis() []*IWebAPI {
	return webApiStorage
}
