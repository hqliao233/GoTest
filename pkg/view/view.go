package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// D 通用数据模型
type D map[string]interface{}

// Render 页面渲染
func Render(w io.Writer, data interface{}, tplFiles ...string) {
	RenderTemplate(w, "app", data, tplFiles...)
}

// RenderSimple 简单页面渲染
func RenderSimple(w io.Writer, data interface{}, tplFiles ...string) {
	RenderTemplate(w, "simple", data, tplFiles...)
}

// RenderTemplate 根据名称渲染模板
func RenderTemplate(w io.Writer, name string, data interface{}, tplFiles ...string) {
	// 1.设置模板相对路径
	viewDir := "resources/views/"

	// 2.路由名称替换
	for i, tplFile := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(tplFile, ".", "/", -1) + ".gohtml"
	}

	// 3.所有布局模板文件
	files, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	// 4.新增目标文件
	newFiles := append(files, tplFiles...)

	// 5.解析所有模板文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteNameToURL": route.NameToURL,
		}).
		ParseFiles(newFiles...)
	logger.LogError(err)

	// 6.渲染模板
	tmpl.ExecuteTemplate(w, name, data)
}
