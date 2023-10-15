package adminService

import (
	"gin_app/app/common"
	"gin_app/app/common/myError"
	"gin_app/app/controller/adminController/adminVo"
	"gin_app/app/model"
	"gin_app/app/service"
	"gin_app/app/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slices"
)

func GetWallpaper(c *gin.Context, reqVo adminVo.WallPaperReqVo) (common.PageRes, error) {
	db := c.Value("DB").(*gorm.DB)

	var bing []model.Bing
	var count int64
	var thumbs []model.Thumb
	var respVo []adminVo.WallPaperRespVo
	var pageRes common.PageRes

	tx := db.Model(&model.Bing{})
	if reqVo.Status != "" {
		tx.Where("status = ?", reqVo.Status)
	}
	if reqVo.Keyword != "" {
		tx.Where("desc like ?", util.SqlLike(reqVo.Keyword))
	}
	tx.Count(&count)
	tx.Limit(reqVo.PageSize).Offset((reqVo.Page - 1) * reqVo.PageSize).
		Order("release_at desc").Order("created_at desc").Find(&bing)
	fileIds := make([]string, 0, len(bing))
	for _, item := range bing {
		fileIds = append(fileIds, item.FileId)
	}
	if len(fileIds) > 0 {
		db.Where("file_id in ? and ext = ?", fileIds, "webp").Find(&thumbs)
	}

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

func AuditWallPaper(c *gin.Context, reqVo adminVo.WallPaperAuditReqVo) (bool, error) {
	db := c.Value("DB").(*gorm.DB)
	log := c.Value("Logger").(*zap.SugaredLogger)

	var ins model.Bing
	err := db.Find(&ins, reqVo.ID).Error
	if err != nil {
		return false, err
	}
	if ins.Status == "1" {
		return false, myError.New("已经审核过了")
	}
	if reqVo.Status == ins.Status {
		return false, myError.New("不能重复审核")
	}
	var bing model.Bing
	_ = copier.Copy(&bing, &reqVo)
	// 使用事务
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Updates(&bing).Error; err != nil {
			return err
		}
		// 通过
		if reqVo.Status == "1" {
			// 审核通过 生成缩略图
			err = service.ToWebp(c, ins.FileId)
			if err != nil {
				log.Error(err)
				return myError.New("生成缩略图失败")
			}
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
