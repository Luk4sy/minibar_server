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
	"minibar_server/utils/markdown"
	"minibar_server/utils/xss"
)

type ArticleUpdateRequest struct {
	ID          uint       `json:"id" binding:"required"`
	Title       string     `json:"title"`
	Abstract    string     `json:"abstract"`
	Content     string     `json:"content"`
	CategoryID  *uint      `json:"categoryID"`
	TagList     ctype.List `json:"tagList"`
	Cover       string     `json:"cover"`
	OpenComment bool       `json:"openComment"`
}

func (ArticleApi) ArticleUpdateView(c *gin.Context) {
	cr := middleware.GetBind[ArticleUpdateRequest](c)

	user, err := jwts.GetClaims(c).GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	if global.Config.Site.SiteInfo.Mode == 2 {
		if user.Role != enum.AdminRole {
			res.FailWithMsg("目前为博客模式，普通用户无法更新文章~", c)
		}
	}

	var article models.ArticleModel
	err = global.DB.Take(&article, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// 更新的文章必须是自己的
	if article.UserID != user.ID {
		res.FailWithMsg("无权修改他人文章", c)
		return
	}

	// 判断分类 id 是不是自己创建的
	var category models.CategoryModel
	if cr.CategoryID != nil {
		err = global.DB.Take(&category, "id = ? and user_id = ?", *cr.CategoryID, user.ID).Error
		if err != nil {
			res.FailWithMsg("文章分类不存在", c)
			return
		}
	}

	// 文章正文防 xss 注入
	cr.Content = xss.FilterSanitize(cr.Content)
	// 如果清洗完发现内容没了（说明用户发的全是脚本），报错
	if cr.Content == "" {
		res.FailWithMsg("文章内容包含非法字符或为空", c)
		return
	}

	// 如果未传如简介，自动从正文中取前 50 个字符
	if cr.Abstract == "" {
		// 把 markdown 转成 html，再取文本
		cr.Abstract = markdown.GetAbstract(cr.Content, 50)
	}

	mps := map[string]any{
		"title":        cr.Title,
		"abstract":     cr.Abstract,
		"content":      cr.Content,
		"category_id":  cr.CategoryID,
		"tag_list":     cr.TagList,
		"cover":        cr.Cover,
		"open_comment": cr.OpenComment,
	}

	// 如果文章是已经发布的文章，进行编辑后，需要把状态改回待审核
	if article.Status == enum.ArticleStatusPublished && !global.Config.Site.Article.IsFreeReview {
		mps["status"] = enum.ArticleStatusUnderReview
	}

	err = global.DB.Model(&article).Updates(mps).Error
	if err != nil {
		res.FailWithMsg("更新失败", c)
	}

	res.OkWithMsg("文章更新成功", c)

}
