package route

import "net/http"

type Route struct {
	name     string
	httpVerb string
	handler  func(responseWriter http.ResponseWriter, request *http.Request)
}

func NewRoute(name, httpVerb string, handler func(responseWriter http.ResponseWriter, request *http.Request)) Route {
	return Route{
		name:     name,
		handler:  handler,
		httpVerb: httpVerb,
	}
}

func (r *Route) Name() string {
	return r.name
}

func (r *Route) Handler() func(responseWriter http.ResponseWriter, request *http.Request) {
	return r.handler
}
