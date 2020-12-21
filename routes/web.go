package routes

import (
	"goblog/app/http/controllers"
	"goblog/app/http/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes 注册web路由
func RegisterWebRoutes(r *mux.Router) {
	pages := new(controllers.PagesController)
	// 静态页面
	r.HandleFunc("/", pages.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pages.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pages.NotFound)

	// 文章路由
	articles := new(controllers.ArticlesController)
	r.HandleFunc("/articles", articles.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles/{id:[0-9]+}", articles.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles/create", articles.Create).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles", articles.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", articles.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", articles.Update).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", articles.Delete).Methods("POST").Name("articles.delete")

	// 中间件
	r.Use(middlewares.ForceHTML)
}
