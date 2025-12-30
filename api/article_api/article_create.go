package article_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/models/ctype"
	"minibar_server/models/enum"
	"minibar_server/utils/jwts"
)

type ArticleCreateRequest struct {
	Title       string             `json:"title"`
	Abstract    string             `json:"abstract"`
	Content     string             `json:"content"`
	CategoryID  *uint              `json:"categoryID"`
	TagList     ctype.List         `json:"tagList"`
	Cover       string             `json:"cover"`
	OpenComment bool               `json:"openComment"`
	Status      enum.ArticleStatus `json:"status" binding:"required,oneof=1 2"`
}

func (ArticleApi) ArticleCreateView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCreateRequest](c)

	user, err := jwts.GetClaims(c).GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// TODO:判断分类 id 是不是自己创建的

	// TODO:文章正文防 xss 注入

	// TODO:正文内容图片转存

	var article = models.ArticleModel{
		Title:       cr.Title,
		Abstract:    cr.Abstract,
		Content:     cr.Content,
		UserID:      user.ID,
		TagList:     cr.TagList,
		Cover:       cr.Cover,
		OpenComment: cr.OpenComment,
		CategoryID:  cr.CategoryID,
		Status:      cr.Status,
	}
	if global.Config.Site.Article.IsFreeReview {
		article.Status = enum.ArticleStatusPublished
	}

	err = global.DB.Create(&article).Error
	if err != nil {
		res.FailWithMsg("文章创建失败", c)
		return
	}

	res.OkWithMsg("文章创建成功", c)

}
