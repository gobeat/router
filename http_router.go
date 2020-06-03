package router

import (
	"context"
	"net/http"

	"github.com/gobeat/interfaces"
	"github.com/gobeat/tools"
	"github.com/julienschmidt/httprouter"
)

type httpRouter struct {
	httprouter                *httprouter.Router
	errorHandler              interfaces.ErrorHandler
	responseHandler           interfaces.Middleware
	beforeDispatchMiddlewares []interfaces.Middleware
	afterDispatchMiddlewares  []interfaces.Middleware
}

// NewHTTPRouter returns an instance of Route using httprouter package
func NewHTTPRouter(errorHandler interfaces.ErrorHandler, responseHandler interfaces.Middleware, notFoundHandler http.Handler) Router {
	if errorHandler == nil {
		panic("errorHandler must be specified.")
	}
	httprouter := httprouter.New()
	httprouter.NotFound = notFoundHandler
	return &httpRouter{
		httprouter:                httprouter,
		errorHandler:              errorHandler,
		responseHandler:           responseHandler,
		beforeDispatchMiddlewares: make([]interfaces.Middleware, 0),
		afterDispatchMiddlewares:  make([]interfaces.Middleware, 0),
	}
}

func (hR *httpRouter) BeforeDispatch(middlewares ...interfaces.Middleware) Router {
	hR.beforeDispatchMiddlewares = append(hR.beforeDispatchMiddlewares, hR.removeNilMiddlewares(middlewares)...)
	return hR
}

func (hR *httpRouter) AfterDispatch(middlewares ...interfaces.Middleware) Router {
	hR.afterDispatchMiddlewares = append(hR.afterDispatchMiddlewares, hR.removeNilMiddlewares(middlewares)...)
	return hR
}

func (hR *httpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hR.httprouter.ServeHTTP(w, r)
}

func (hR *httpRouter) ROUTE(method string, path string, middlewares ...interfaces.Middleware) Router {
	mws := make([]interfaces.Middleware, 0)
	mws = append(mws, hR.beforeDispatchMiddlewares...)
	mws = append(mws, hR.removeNilMiddlewares(middlewares)...)
	mws = append(mws, hR.afterDispatchMiddlewares...)
	mws = append(mws, hR.responseHandler)
	hR.httprouter.Handle(method, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		rb := tools.NewBag()
		for _, param := range params {
			rb.Set(param.Key, param.Value)
		}
		r = r.WithContext(context.WithValue(r.Context(), requestBagKey{}, rb))
		for _, mw := range mws {
			ctx, err := mw(w, r)
			if err != nil {
				hR.errorHandler(w, r, err)
				return
			}
			if ctx != nil {
				r = r.WithContext(ctx)
			}
		}
	})

	return hR
}

func (hR *httpRouter) GET(path string, middlewares ...interfaces.Middleware) Router {
	return hR.ROUTE(http.MethodGet, path, middlewares...)
}

func (hR *httpRouter) POST(path string, middlewares ...interfaces.Middleware) Router {
	return hR.ROUTE(http.MethodPost, path, middlewares...)
}

func (hR *httpRouter) PUT(path string, middlewares ...interfaces.Middleware) Router {
	return hR.ROUTE(http.MethodPut, path, middlewares...)
}

func (hR *httpRouter) PATCH(path string, middlewares ...interfaces.Middleware) Router {
	return hR.ROUTE(http.MethodPatch, path, middlewares...)
}

func (hR *httpRouter) DELETE(path string, middlewares ...interfaces.Middleware) Router {
	return hR.ROUTE(http.MethodDelete, path, middlewares...)
}

func (hR *httpRouter) OPTIONS(path string, middlewares ...interfaces.Middleware) Router {
	return hR.ROUTE(http.MethodOptions, path, middlewares...)
}

func (hR *httpRouter) HEAD(path string, middlewares ...interfaces.Middleware) Router {
	return hR.ROUTE(http.MethodHead, path, middlewares...)
}

func (hR *httpRouter) removeNilMiddlewares(middlewares []interfaces.Middleware) []interfaces.Middleware {
	mws := make([]interfaces.Middleware, 0)
	for _, mw := range middlewares {
		if mw != nil {
			mws = append(mws, mw)
		}
	}
	return mws
}
