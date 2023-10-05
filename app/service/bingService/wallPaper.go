package bingService

import (
	"gin_app/app/common"
	"gin_app/app/controller/bingController/bingVo"
	"gin_app/app/model"
	"gin_app/app/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"slices"
)

func GetWallPaper(c *gin.Context, reqVo bingVo.WallPaperReqVo) (common.PageRes, error) {
	db := c.Value("DB").(*gorm.DB)

	var bing []model.Bing
	var count int64
	var thumbs []model.Thumb
	var respVo []bingVo.WallPaperRespVo
	var pageRes common.PageRes

	db.Model(&model.Bing{}).Count(&count)
	db.Limit(reqVo.PageSize).Offset((reqVo.Page - 1) * reqVo.PageSize).
		Order("created_at desc").Find(&bing)
	fileIds := make([]string, 0, len(bing))
	for _, item := range bing {
		fileIds = append(fileIds, item.FileId)
	}
	db.Where("file_id in ?", fileIds).Where("ext = ?", "webp").Find(&thumbs)

	_ = copier.Copy(&respVo, &bing)
	// type 取file 中的ext
	for i := range respVo {
		index := slices.IndexFunc(thumbs, func(m model.Thumb) bool {
			return m.FileId == respVo[i].FileId
		})
		if index < 0 {
			continue
		}
		thumb := thumbs[index]
		respVo[i].Name = thumb.Name
		respVo[i].Width = thumb.Width
		respVo[i].Height = thumb.Height
		respVo[i].Ext = util.GetFileExt(respVo[i].FileId)
		respVo[i].ThumbId = util.WriteString(util.ToString(thumb.ID), ".", thumb.Ext)
	}
	pageRes.Count = count
	pageRes.Rows = respVo
	return pageRes, nil
}

// CreateWallPaper 手动上传壁纸
func CreateWallPaper() (bool, error) {

	return false, nil
}
