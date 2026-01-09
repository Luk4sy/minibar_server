package article_api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/models/enum"
	"minibar_server/service/redis_service/redis_article"
	"minibar_server/utils/jwts"
	"time"
)

type ArticleLookRequest struct {
	ArticleID  uint `json:"articleID" binding:"required"`
	TimeSecond int  `json:"timeSecond"` // 读文章一共用了多久
}

func (ArticleApi) ArticleLookView(c *gin.Context) {
	cr := middleware.GetBind[ArticleLookRequest](c)

	// TODO:未登录用户，浏览量如何算？
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.OkWithMsg("未登录", c)
		return
	}

	// 引入缓存
	// 当天用户请求该文章之后，将用户 id 和文章作为 key 存入缓存，进行判断，如果存在就直接返回，不走下面的逻辑
	if redis_article.GetUserArticleHistoryCache(cr.ArticleID, claims.UserID) {
		logrus.Infof("在缓存中")
		res.OkWithMsg("成功", c)
		return
	}

	var article models.ArticleModel
	err = global.DB.Take(&article, "status = ? and id = ?", enum.ArticleStatusPublished, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// 查这个文章当日是否在 “足迹” 里面
	var history models.UserArticleLookHistoryModel
	err = global.DB.Take(&history, "user_id = ? and article_id = ? and created_at < ? and created_at > ?",
		claims.UserID, cr.ArticleID,
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02")+" 00:00:00",
	).Error
	if err == nil {
		res.OkWithMsg("成功", c)
		return
	}

	err = global.DB.Create(&models.UserArticleLookHistoryModel{
		UserID:    claims.UserID,
		ArticleID: cr.ArticleID,
	}).Error
	if err != nil {
		res.FailWithMsg("失败", c)
		return
	}

	redis_article.SetLookCache(cr.ArticleID, true)
	redis_article.SetUserArticleHistoryCache(cr.ArticleID, claims.UserID)
	res.OkWithMsg("成功", c)
}
