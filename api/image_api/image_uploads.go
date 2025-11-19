package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

func (ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	// 文件大小判断
	s := global.Config.Upload.Size
	if fileHeader.Size > s*1024*1024 {
		res.FailWithMsg(fmt.Sprintf("文件大小大于 %dMB", s), c)
		return
	}
	// 后缀判断
	filename := fileHeader.Filename
	err = ImageSuffixJudge(filename)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	filePath := fmt.Sprintf("uploads/images/%s", fileHeader.Filename)
	c.SaveUploadedFile(fileHeader, filePath)
	res.Ok("/"+filePath, "图片上传成功", c)
}

func ImageSuffixJudge(filename string) (err error) {
	_list := strings.Split(filename, ".")
	if len(_list) == 1 {
		return errors.New("错误的文件名")
	}
	suffix := _list[len(_list)-1]
	if !utils.InList(suffix, global.Config.Upload.WhiteList) {
		return errors.New("文件非法")
	}
	return
}
