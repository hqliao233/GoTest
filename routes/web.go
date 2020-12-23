package routes

import (
	"goblog/app/http/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes 注册web路由
func RegisterWebRoutes(r *mux.Router) {
	pages := new(controllers.PagesController)
	// 静态页面
	r.HandleFunc("/about", pages.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pages.NotFound)

	// 文章路由
	articles := new(controllers.ArticlesController)
	r.HandleFunc("/", articles.Index).Methods("GET").Name("home")
	r.HandleFunc("/articles", articles.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles/{id:[0-9]+}", articles.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles/create", articles.Create).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles", articles.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", articles.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", articles.Update).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", articles.Delete).Methods("POST").Name("articles.delete")

	// 用戶認證
	auc := new(controllers.AuthController)
	r.HandleFunc("/auth/register", auc.Register).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", auc.DoRegister).Methods("POST").Name("auth.doregister")

	//静态资源路由
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// 中间件
	// r.Use(middlewares.ForceHTML)
}
