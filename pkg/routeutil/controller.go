package routeutil

import (
	"github.com/clevergo/clevergo"
)

type Controllers []Controller

func (cs Controllers) Routes() (routes Routes) {
	for _, c := range cs {
		routes = append(routes, c.Routes()...)
	}
	return
}

type Controller interface {
	Routes() Routes
}

type RESTController interface {
	Path() string
	Index(*clevergo.Context)
	Create(*clevergo.Context)
	Delete(*clevergo.Context)
	Opitions(*clevergo.Context)
	Update(*clevergo.Context)
}
