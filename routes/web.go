package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vinoMamba/goblog/app/http/controllers"
)

func RegisterWebRoutes(router *mux.Router) {
	pc := new(controllers.PagesController)
	router.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	router.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	router.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	ac := new(controllers.ArticlesController)
	router.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
}
