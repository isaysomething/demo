package routeutil

import (
	"github.com/clevergo/clevergo"
)

type Routes []*Route

func (rs Routes) Register(router clevergo.IRouter) {
	for _, r := range rs {
		r.Register(router)
	}
}

type Route struct {
	method      string
	path        string
	handle      clevergo.Handle
	name        string
	middlewares []clevergo.MiddlewareFunc
}

func NewRoute(method, path string, handle clevergo.Handle) *Route {
	return &Route{
		method: method,
		path:   path,
		handle: handle,
	}
}
func (r *Route) Name(name string) *Route {
	r.name = name
	return r
}

func (r *Route) Middlewares(middlewares ...clevergo.MiddlewareFunc) *Route {
	r.middlewares = middlewares
	return r
}

func (r *Route) Register(router clevergo.IRouter) {
	router.Handle(r.method, r.path, r.handle, r.options()...)
}

func (r *Route) options() (opts []clevergo.RouteOption) {
	if r.name != "" {
		opts = append(opts, clevergo.RouteName(r.name))
	}
	return
}
