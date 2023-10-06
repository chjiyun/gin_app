package bingService

import (
	"gin_app/app/common"
	"gin_app/app/common/myError"
	"gin_app/app/controller/bingController/bingVo"
	"gin_app/app/model"
	"gin_app/app/service"
	"gin_app/app/util"
	"gin_app/config"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func GetWallPaper(c *gin.Context, reqVo bingVo.WallPaperReqVo) (common.PageRes, error) {
	db := c.Value("DB").(*gorm.DB)

	var bing []model.Bing
	var count int64
	var thumbs []model.Thumb
	var respVo []bingVo.WallPaperRespVo
	var pageRes common.PageRes

	db.Model(&model.Bing{}).Count(&count)
	db.Where("pass = ?", reqVo.Pass).Limit(reqVo.PageSize).Offset((reqVo.Page - 1) * reqVo.PageSize).
		Order("created_at desc").Find(&bing)
	fileIds := make([]string, 0, len(bing))
	for _, item := range bing {
		fileIds = append(fileIds, item.FileId)
	}
	db.Where("file_id in ? and ext = ?", fileIds, "webp").Find(&thumbs)

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

// AddWallPaper 手动上传壁纸 审核通过后方可进入file
func AddWallPaper(c *gin.Context, reqVo bingVo.WallPaperCreateReqVo) (bool, error) {
	log := c.Value("Logger").(*zap.SugaredLogger)
	// 校验图片格式和质量  大小 分辨率
	f := reqVo.File
	ext := filepath.Ext(f.Filename)
	imgSuffix := regexp.MustCompile(`jpg|jpeg|png$`)
	if !imgSuffix.MatchString(ext) {
		return false, myError.NewET(common.InValidFile)
	}
	contentTypes := f.Header["Content-Type"]
	isImage := slices.ContainsFunc(contentTypes, func(s string) bool {
		return strings.Contains(s, "image")
	})
	if !isImage {
		return false, myError.NewET(common.InValidFile)
	}

	mFile, _ := f.Open()
	defer mFile.Close()
	width, height, err := service.GetImageXY(mFile)
	if err != nil {
		return false, myError.New("文件解码失败")
	}
	if width < 256 || height < 256 {
		return false, myError.New("图片尺寸过小，请上传高分辨率的图片")
	}

	// 保存临时文件
	sourcePath := filepath.Join(config.Cfg.Basedir, "files/temp", f.Filename)
	if err := c.SaveUploadedFile(f, sourcePath); err != nil {
		log.Error(err)
		return false, myError.NewET(common.UnknownError)
	}
	// 保存文件 禁止转缩略图  节约资源
	c.Set("noThumb", true)
	fileId, err := service.Upload(c, f)
	if err != nil {
		log.Error(err)
		return false, myError.NewET(common.UnknownError)
	}
	db := c.Value("DB").(*gorm.DB)
	bing := model.Bing{
		FileId:    fileId,
		Desc:      reqVo.Desc,
		ReleaseAt: reqVo.ReleaseAt,
		Pass:      false,
	}
	id := idgen.NextId()
	bing.ID = id
	if err := db.Create(&bing).Error; err != nil {
		return false, err
	}
	return true, nil
}
