package routeutil

import "github.com/clevergo/clevergo"

type Groups []*Group

func (gs Groups) Register(router clevergo.IRouter) {
	for _, g := range gs {
		g.Register(router)
	}
}

type Group struct {
	path        string
	middlewares []clevergo.MiddlewareFunc
	routes      Routes
	groups      Groups
}

func NewGroup(path string, routes Routes) *Group {
	return &Group{
		path:   path,
		routes: routes,
	}
}

func (g *Group) Groups(groups ...*Group) *Group {
	g.groups = groups
	return g
}

func (g *Group) Middlewares(middlewares ...clevergo.MiddlewareFunc) *Group {
	g.middlewares = middlewares
	return g
}

func (g *Group) Register(router clevergo.IRouter) {
	router = router.Group(g.path, g.options()...)
	if len(g.routes) > 0 {
		g.routes.Register(router)
	}
	if len(g.groups) > 0 {
		g.groups.Register(router)
	}
}

func (g *Group) options() (opts []clevergo.RouteGroupOption) {
	if len(g.middlewares) > 0 {
		opts = append(opts, clevergo.RouteGroupMiddleware(g.middlewares...))
	}
	return
}
