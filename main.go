package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// ArticlesFormData 文章数据类型
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

// Article 文章模型
type Article struct {
	Title, Body string
	ID          int64
}

// Link 文章访问连接
func (a Article) Link() string {
	id := a.ID
	articleURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(id, 10))
	if err != nil {
		return ""
	}
	return articleURL.String()
}

// Delete 删除文章
func (a Article) Delete() (int64, error) {
	id := a.ID
	rs, err := db.Exec("delete from articles where id = ?", id)
	if err != nil {
		checkError(err)
		return 0, err
	}
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}

// RouteNameToURL 通过路由名称获取链接
func RouteNameToURL(routeName string, pairs ...string) string {
	url, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		checkError(err)
		return ""
	}
	return url.String()
}

// Int64ToString int转string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

var router = mux.NewRouter()
var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "gotest",
		AllowNativePasswords: true,
	}

	// 准备数据库连接池
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	// 设置最大连接数
	db.SetMaxOpenConns(25)
	// 设置最大空闲数
	db.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)
	// 尝试连接，失败报错
	err = db.Ping()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createTable() {
	createArticlesSQL := `
	CREATE TABLE IF NOT EXISTS articles(
		id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
		body longtext COLLATE utf8mb4_unicode_ci
	);
	`
	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 这里是goblog</h1>")
}

func aboutHander(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "这里是关于blog的内容12123123"+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "_我不知道你在说什么123")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	article := Article{}
	query := "select * from articles where id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		}
	} else {
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteNameToURL": RouteNameToURL,
				"Int64ToString":  Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		tmpl.Execute(w, article)
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	// 1.获取所有文章数据
	rows, err := db.Query("select * from articles")
	checkError(err)
	defer rows.Close()

	// 2.循环赋值
	var articles []Article

	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		checkError(err)
		articles = append(articles, article)
	}

	// 2.1检查遍历是否存在错误
	err = rows.Err()
	checkError(err)

	// 3.渲染模板
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	checkError(err)

	tmpl.Execute(w, articles)
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	// Form会比PostForm多一些URL中自带的参数
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errors := validateFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		storeURL, _ := router.Get("articles.store").URL()
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		tmpl.Execute(w, data)
	}
}

func saveArticleToDB(title string, body string) (int64, error) {
	// 变量处理话
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 1.获取一个prepare执行语句
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	if err != nil {
		return 0, err
	}

	// 2. 方法结束关闭请求，防止占用sql连接
	defer stmt.Close()

	// 3.执行请求，传参进入绑定内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	// 4.成功返回ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, data)
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)

	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		}
	} else {
		updateURL, _ := router.Get("articles.update").URL("id", id)
		articleData := ArticlesFormData{
			URL:    updateURL,
			Title:  article.Title,
			Body:   article.Body,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		checkError(err)
		tmpl.Execute(w, articleData)
	}
}

// ArticlesUpdateHandler 文章更新路由
func ArticlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 1 获取路由
	id := getRouteVariable("id", r)

	// 2 获取文章判断是否存在即可
	_, err := getArticleByID(id)

	// 3.检查是否有错误
	if err != nil {
		// 3.1 错误处理
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		}
	} else {
		// 3.2 获取POST参数
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")
		errors := validateFormData(title, body)

		// 3.3 不存在错误进行保存错误返回上一个页面继续进行编辑操作
		if len(errors) == 0 {
			query := "update articles set title = ?, body = ? where id = ?"
			// db.Exec方法为防sql注入，提供给create update delete方法使用
			rs, err := db.Exec(query, title, body, id)
			if err != nil {
				checkError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 系统错误")
			}
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改")
			}
		} else {
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")

			checkError(err)
			tmpl.Execute(w, data)
		}
	}
}

// ArticlesDeleteHandler 文章删除路由
func ArticlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1.获取文章
	id := getRouteVariable("id", r)
	article, err := getArticleByID(id)

	// 2.排查错误
	checkError(err)
	if err != nil {
		// 2.1存在错误进行提示
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		}
	} else {
		// 2.2不存在错误进行删除
		rowsAffected, err := article.Delete()

		// 2.3删除出现错误进行处理未发现错误调整首页
		if err != nil {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		} else {
			if rowsAffected > 0 {
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}

// 根据ID获取文章
func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "select * from articles where id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

// 获取路由参数
func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

// validateFormData POST参数校验
func validateFormData(title string, body string) map[string]string {
	errors := make(map[string]string)

	// 标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度必须在3到10之前"
	}

	// 内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度不能小于10"
	}

	return errors
}

// forceHTMLMiddleware 设置请求头
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1.设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//2.继续处理请求
		next.ServeHTTP(w, r)
	})
}

// removeTrailingSlash 清除路由斜杠
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 数据库连接
	initDB()
	createTable()

	// 起始路由
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHander).Methods("GET").Name("about")
	// 文章路由
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", ArticlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", ArticlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 强制未找到路由跳转
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 设置请求头
	router.Use(forceHTMLMiddleware)

	// 服务
	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
