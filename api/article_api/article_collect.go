package article_api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/models/enum"
	"minibar_server/utils/jwts"
)

type ArticleCollectRequest struct {
	ArticleID uint `json:"articleID" binding:"required"`
	CollectID uint `json:"collectID"`
}

func (ArticleApi) ArticleCollectView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCollectRequest](c)

	var article models.ArticleModel
	err := global.DB.Take(&article, "status = ? and id = ?", enum.ArticleStatusPublished, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	var collectModel models.CollectModel
	claims := jwts.GetClaims(c)

	if cr.CollectID == 0 {
		// 使用默认收藏夹
		err = global.DB.Take(&collectModel, "user_id = ? and is_default = ?", claims.UserID, 1).Error
		if err != nil {
			// 创建一个默认收藏夹
			collectModel.Title = "默认收藏夹"
			collectModel.UserID = claims.UserID
			collectModel.IsDefault = true
			global.DB.Create(&collectModel)
			cr.CollectID = collectModel.ID
		}
		cr.CollectID = collectModel.ID
	} else {
		// 判断收藏夹是否存在，是否为自己创建的
		err = global.DB.Take(&collectModel, "user_id = ?", claims.UserID).Error
		if err != nil {
			res.FailWithMsg("收藏夹不存在", c)
			return
		}
	}

	// 判断是否收藏
	var articleCollect models.UserArticleCollectModel
	err = global.DB.Where(models.UserArticleCollectModel{
		UserID:    claims.UserID,
		ArticleID: cr.ArticleID,
		CollectID: cr.CollectID,
	}).Take(&articleCollect).Error
	if err != nil {
		// 收藏
		err = global.DB.Create(&models.UserArticleCollectModel{
			UserID:    claims.UserID,
			ArticleID: cr.ArticleID,
			CollectID: cr.CollectID,
		}).Error
		if err != nil {
			res.FailWithMsg("收藏失败", c)
			return
		}
		res.OkWithMsg("收藏成功", c)

		// 增加收藏计数
		global.DB.Model(&collectModel).Update("article_count", gorm.Expr("article_count + 1"))
		return
	}

	// 取消收藏
	err = global.DB.Where(models.UserArticleCollectModel{
		UserID:    claims.UserID,
		ArticleID: cr.ArticleID,
		CollectID: cr.CollectID,
	}).Delete(models.UserArticleCollectModel{}).Error
	if err != nil {
		res.FailWithMsg("取消收藏失败", c)
	}
	res.OkWithMsg("取消收藏成功", c)
	// TODO:收藏数同步到缓存
	global.DB.Model(&collectModel).Update("article_count", gorm.Expr("article_count - 1"))
}
