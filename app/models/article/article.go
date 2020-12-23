package article

import (
	"goblog/pkg/model"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"strconv"
)

// Article 文章模型
type Article struct {
	model.BaseModel
	Title string
	Body  string
}

// Get 根据ID获取单篇文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}

// GetAll 獲取所有文章
func GetAll() ([]Article, error) {
	var articles []Article

	if err := model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

// Link 文章访问连接
func (a Article) Link() string {
	return route.NameToURL("articles.show", "id", strconv.FormatUint(a.ID, 10))
}
