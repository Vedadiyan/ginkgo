package ginkgohelpers

import "github.com/gin-gonic/gin"

type GinkgoContext gin.Context

func (ctx GinkgoContext) BadRequest() {
	ginCtx := gin.Context(ctx)
	ginCtx.Status(400)
}

func (ctx GinkgoContext) BadGateway() {
	ginCtx := gin.Context(ctx)
	ginCtx.Status(502)
}

func (ctx GinkgoContext) InternalServerError() {
	ginCtx := gin.Context(ctx)
	ginCtx.Status(500)
}

func (ctx GinkgoContext) Unauthorized() {
	ginCtx := gin.Context(ctx)
	ginCtx.Status(401)
}

func (ctx GinkgoContext) Forbidden() {
	ginCtx := gin.Context(ctx)
	ginCtx.Status(403)
}
