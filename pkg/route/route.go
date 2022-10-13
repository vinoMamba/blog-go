package route

import (
	"github.com/gorilla/mux"
	"github.com/vinoMamba/goblog/routes"
)

var Router *mux.Router

func Initialize() {
	Router = mux.NewRouter()
	routes.RegisterWebRoutes(Router)
}

func RouteNameToUrl(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		return ""
	}
	return url.String()
}
