package api

import "github.com/gin-gonic/gin"

type DomainAPI interface {
	BuildPath() string
	RegisterHandlers(group *gin.RouterGroup)
}

type Router struct {
	engine  *gin.Engine
	domains []DomainAPI
}

func NewRouter(engine *gin.Engine, domains ...DomainAPI) *Router {
	return &Router{
		engine:  engine,
		domains: domains,
	}
}

func (r *Router) Register() {
	for _, domain := range r.domains {
		group := r.engine.Group(domain.BuildPath())
		domain.RegisterHandlers(group)
	}
}
