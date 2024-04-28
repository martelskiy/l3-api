package route

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/martelskiy/l3-api/api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router interface {
	WithAPIDocumentation() Router
	WithCORSMiddleware() Router
	WithRoute(route Route) Router
	GetRouter() *mux.Router
}

type WebRouter struct {
	muxRouter *mux.Router
}

func NewRouter() *WebRouter {
	muxRouter := mux.NewRouter().StrictSlash(true)
	return &WebRouter{
		muxRouter: muxRouter,
	}
}

func (r *WebRouter) WithCORSMiddleware() Router {
	r.muxRouter.Use(accessControlMiddleware)
	return r
}

func (r *WebRouter) WithAPIDocumentation() Router {
	r.muxRouter.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	return r
}

func (r *WebRouter) WithRoute(route Route) Router {
	r.muxRouter.HandleFunc(route.name, route.handler).Methods(route.httpVerb)
	return r
}

func (r *WebRouter) GetRouter() *mux.Router {
	return r.muxRouter
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		next.ServeHTTP(w, r)
	})
}
