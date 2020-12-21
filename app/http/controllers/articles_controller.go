package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
	"strconv"
	"unicode/utf8"

	"gorm.io/gorm"
)

// ArticlesController 文章控制器
type ArticlesController struct {
}

// ArticlesFormData 文章数据类型
type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

// Index 文章列表首頁
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 系统错误")
	} else {
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, articles)
	}
}

// Show 显示文章
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		}
	} else {
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteNameToURL": route.NameToURL,
				"Int64ToString":  types.Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w, article)
	}
}

// Create 新建文章
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.NameToURL("articles.store")
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

// Store 保存文章
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	// Form会比PostForm多一些URL中自带的参数
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errors := validateFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		storeURL := route.NameToURL("articles.store")
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}

// Edit 编辑文章
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		}
	} else {
		updateURL := route.NameToURL("articles.update", "id", id)
		articleData := ArticlesFormData{
			URL:    updateURL,
			Title:  _article.Title,
			Body:   _article.Body,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)
		tmpl.Execute(w, articleData)
	}
}

// Update 文章更新路由
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 1 获取路由
	id := route.GetRouteVariable("id", r)

	// 2 获取文章判断是否存在即可
	_article, err := article.Get(id)

	// 3.检查是否有错误
	if err != nil {
		// 3.1 错误处理
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
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
			_article.Title = title
			_article.Body = body
			if rowsAffected, _ := _article.Update(); rowsAffected > 0 {
				showURL := route.NameToURL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改")
			}
		} else {
			updateURL := route.NameToURL("articles.update", "id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")

			logger.LogError(err)
			tmpl.Execute(w, data)
		}
	}
}

// Delete 文章删除路由
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	// 1.获取文章
	id := route.GetRouteVariable("id", r)
	article, err := article.Get(id)

	// 2.排查错误
	if err != nil {
		logger.LogError(err)
		// 2.1存在错误进行提示
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		}
	} else {
		// 2.2不存在错误进行删除
		rowsAffected, err := article.Delete()
		// 2.3删除出现错误进行处理未发现错误调整首页
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 系统错误")
		} else {
			if rowsAffected > 0 {
				indexURL := route.NameToURL("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}
