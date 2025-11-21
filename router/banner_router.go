package router

import (
	"blogx_server/api"
	"blogx_server/middleware"
	"github.com/gin-gonic/gin"
)

func BannerRouter(r *gin.RouterGroup) {
	app := api.App.BannerApi
	r.GET("banner", app.BannerListView)
	r.POST("banner", middleware.AdminMiddleware, app.BannerCreateView)
	r.DELETE("banner", middleware.AdminMiddleware, app.BannerRemoveView)
	r.PUT("banner/:id", middleware.AdminMiddleware, app.BannerUpdateView)

}
