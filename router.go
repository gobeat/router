package router

import (
	"net/http"

	"github.com/gobeat/interfaces"
)

// Router is a routing service
type Router interface {
	http.Handler

	ROUTE(method string, path string, middlewares ...interfaces.Middleware) Router
	GET(path string, middlewares ...interfaces.Middleware) Router
	POST(path string, middlewares ...interfaces.Middleware) Router
	PUT(path string, middlewares ...interfaces.Middleware) Router
	PATCH(path string, middlewares ...interfaces.Middleware) Router
	DELETE(path string, middlewares ...interfaces.Middleware) Router
	OPTIONS(path string, middlewares ...interfaces.Middleware) Router
	HEAD(path string, middlewares ...interfaces.Middleware) Router

	BeforeDispatch(...interfaces.Middleware) Router
	AfterDispatch(...interfaces.Middleware) Router
}
