package article_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"minibar_server/common"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/models/enum"
	"minibar_server/service/redis_service/redis_article"
	"minibar_server/utils/jwts"
	"minibar_server/utils/sql"
)

type ArticleListRequest struct {
	common.PageInfo
	Type       int8               `form:"type" binding:"required,oneof=1 2 3"` // 1 用户查别人的 2 查自己的 3 管理员查
	UserID     uint               `form:"userID"`
	CategoryID *uint              `form:"categoryID"`
	Status     enum.ArticleStatus `form:"status"`
}

type ArticleListResponse struct {
	models.ArticleModel
	UserTop  bool `json:"userTop"`  // 是否是用户置顶文章
	AdminTop bool `json:"adminTop"` // 是否是管理员指定
}

func (ArticleApi) ArticleListView(c *gin.Context) {
	cr := middleware.GetBind[ArticleListRequest](c)

	var topArticleIDList []uint

	var orderColumMap = map[string]bool{
		"look_count desc":    true,
		"digg_count desc":    true,
		"comment_count desc": true,
		"collect_count desc": true,
		"look_count asc":     true,
		"digg_count asc":     true,
		"comment_count asc":  true,
	}

	switch cr.Type {
	case 1:
		// 查别人的文章，id 必须填写
		if cr.UserID == 0 {
			res.FailWithMsg("请填写用户 id ", c)
			return
		}
		if cr.Page > 2 || cr.Limit > 10 {
			res.FailWithMsg("请登陆后查询更多内容！", c)
			return
		}
		cr.Status = 0
		cr.Order = ""
	case 2:
		// 查自己的文章
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.UserID
	case 3:
		// 管理员
		claims, err := jwts.ParseTokenByGin(c)
		if !(err == nil && claims.Role == enum.AdminRole) {
			res.FailWithMsg("角色错误", c)
			return
		}
	}

	if cr.Order != "" {
		_, ok := orderColumMap[cr.Order]
		if !ok {
			res.FailWithMsg("不支持该排列方式！", c)
			return
		}
	}

	var userTopMap = map[uint]bool{}
	var adminTopMap = map[uint]bool{}
	if cr.UserID != 0 {
		var userTopArticleList []models.UserTopArticleModel
		global.DB.Preload("UserModel").Order("created_at desc").Find(&userTopArticleList, "user_id = ?", cr.UserID)

		for _, i2 := range userTopArticleList {
			topArticleIDList = append(topArticleIDList, i2.ArticleID)
			if i2.UserModel.Role == enum.AdminRole {
				adminTopMap[i2.ArticleID] = true
			}
			userTopMap[i2.ArticleID] = true
		}
	}

	var options = common.Options{
		Likes:        []string{"title"},
		PageInfo:     cr.PageInfo,
		DefaultOrder: "created_at desc",
	}
	if len(topArticleIDList) > 0 {
		options.DefaultOrder = fmt.Sprintf("%s, created_at desc", sql.ConvertSliceOrderSql(topArticleIDList))
	}
	_list, count, _ := common.ListQuery(models.ArticleModel{
		UserID:     cr.UserID,
		CategoryID: cr.CategoryID,
		Status:     cr.Status,
	}, options)
	fmt.Printf("%s\n", sql.ConvertSliceOrderSql(topArticleIDList))

	var list = make([]ArticleListResponse, 0)
	collectMap := redis_article.GetAllCollectCache()
	lookMap := redis_article.GetAllLookCache()
	diggMap := redis_article.GetAllDiggCache()

	for _, model := range _list {
		model.Content = ""
		model.CollectCount = model.CollectCount + collectMap[model.ID]
		model.LookCount = model.LookCount + lookMap[model.ID]
		model.DiggCount = model.DiggCount + diggMap[model.ID]

		list = append(list, ArticleListResponse{
			ArticleModel: model,
			UserTop:      userTopMap[model.ID],
			AdminTop:     adminTopMap[model.ID],
		})
	}
	res.OkWithList(list, count, c)
}
