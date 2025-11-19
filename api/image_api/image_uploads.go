package image_api

import (
	"blogx_server/common/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	filePath := fmt.Sprintf("uploads/images/%s", fileHeader.Filename)
	c.SaveUploadedFile(fileHeader, filePath)
	res.Ok("/"+filePath, "图片上传成功", c)
}
