package router

import (
	"github.com/gin-gonic/gin"
	"minibar_server/api"
	"minibar_server/api/article_api"
	"minibar_server/middleware"
	"minibar_server/models"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi
	r.POST("article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCreateRequest], app.ArticleCreateView)
	r.POST("article/review", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleReviewRequest], app.ArticleReviewView)
	r.GET("article/digg/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[models.IDRequest], app.ArticleDiggView)
	r.PUT("article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleUpdateRequest], app.ArticleUpdateView)
	r.GET("article", middleware.BindQueryMiddleware[article_api.ArticleListRequest], app.ArticleListView)
	r.GET("article/:id", middleware.BindUriMiddleware[models.IDRequest], app.ArticleDetailView)
	r.POST("article/collect", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCollectRequest], app.ArticleCollectView)
	r.DELETE("article/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[models.IDRequest], app.ArticleRemoveUserView)
	r.DELETE("article", middleware.AdminMiddleware, middleware.BindJsonMiddleware[models.RemoveRequest], app.ArticleRemoveAdminView)

}
