package bingService

import (
	"errors"
	"gin_app/app/common"
	"gin_app/app/common/myError"
	"gin_app/app/controller/bingController/bingVo"
	"gin_app/app/model"
	"gin_app/app/service"
	"gin_app/app/util"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
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

	tx := db.Where("status = ?", "1")
	if reqVo.Keyword != "" {
		tx.Where("desc like ?", util.SqlLike(reqVo.Keyword))
	}
	tx.Model(&model.Bing{}).Count(&count)
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

// AddWallPaper 手动上传壁纸 审核通过后方可进入file
func AddWallPaper(c *gin.Context, reqVo bingVo.WallPaperCreateReqVo) (bool, error) {
	// 校验图片格式和质量  大小 分辨率
	file, err := service.GetFile(c, reqVo.FileId)
	if err != nil {
		return false, err
	}
	f, err := os.Open(file.Path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	if err := validateLocalImage(f); err != nil {
		return false, err
	}

	db := c.Value("DB").(*gorm.DB)
	var bing model.Bing
	_ = copier.Copy(&bing, &reqVo)
	bing.Status = "0"
	bing.IsBing = false
	bing.ID = idgen.NextId()
	if err := db.Create(&bing).Error; err != nil {
		return false, err
	}
	return true, nil
}

// 校验用户上传的壁纸
func validateWallPaper(f *multipart.FileHeader) error {
	ext := filepath.Ext(f.Filename)
	imgSuffix := regexp.MustCompile(`jpg|jpeg|png$`)
	if !imgSuffix.MatchString(ext) {
		return myError.NewET(common.InValidFile)
	}
	contentTypes := f.Header["Content-Type"]
	isImage := slices.ContainsFunc(contentTypes, func(s string) bool {
		return strings.Contains(s, "image")
	})
	if !isImage {
		return myError.NewET(common.InValidFile)
	}
	if f.Size > 1024*1024*10 {
		return myError.New("图片大小不能超过10M")
	}

	mFile, _ := f.Open()
	defer mFile.Close()
	width, height, err := service.GetImageXY(mFile)
	if err != nil {
		return myError.New("文件解码失败")
	}
	if width < 256 || height < 256 {
		return myError.New("图片尺寸过小，请上传高分辨率的图片")
	}
	return nil
}

func ValidateWallPaper(c *gin.Context, f *multipart.FileHeader) (string, error) {
	if err := validateWallPaper(f); err != nil {
		return "", err
	}
	// 审核前禁止转webp
	c.Set("noThumb", "1")
	fileId, err := service.Upload(c, f)
	if err != nil {
		return "", err
	}
	return fileId, nil
}

// 识别已上传的图片文件
func validateLocalImage(f *os.File) error {
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	ext := filepath.Ext(fileInfo.Name())
	imgSuffix := regexp.MustCompile(`jpg|jpeg|png$`)
	if !imgSuffix.MatchString(ext) {
		return myError.NewET(common.InValidFile)
	}
	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return myError.New("无效文件：mime识别失败")
	}
	mType := mime.String()
	index := strings.Index(mType, "/")
	if index < 0 || !strings.Contains(mType[:index], "image") {
		return myError.NewET(common.InValidFile)
	}
	if fileInfo.Size() > 1024*1024*10 {
		return myError.New("图片大小不能超过10M")
	}

	width, height, err := service.GetImageXY(f)
	if err != nil {
		return myError.New("文件解码失败")
	}
	if width < 256 || height < 256 {
		return myError.New("图片尺寸过小，请上传高分辨率的图片")
	}
	return nil
}

func UpdateWallPaper(c *gin.Context, reqVo bingVo.WallPaperUpdateReqVo) (bool, error) {
	db := c.Value("DB").(*gorm.DB)
	var ins model.Bing

	if err := db.Take(&ins, "id = ?", reqVo.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, myError.New("数据不存在")
		}
		return false, err
	}
	if ins.IsBing {
		return false, myError.NewET(common.IllegalVisit)
	}
	// 校验图片格式和质量  大小 分辨率
	if reqVo.FileId != "" && reqVo.FileId != ins.FileId {
		var file model.File
		file, err := service.GetFile(c, reqVo.FileId)
		if err != nil {
			return false, err
		}
		f, err := os.Open(file.Path)
		if err != nil {
			return false, err
		}
		defer f.Close()
		if err = validateLocalImage(f); err != nil {
			return false, err
		}
	}
	var bing model.Bing
	_ = copier.Copy(&bing, &reqVo)
	bing.Status = "0"
	if err := db.Updates(&bing).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteWallPaper(c *gin.Context, id string) (bool, error) {
	db := c.Value("DB").(*gorm.DB)
	var ins model.Bing

	if err := db.Take(&ins, "id = ?", id).Error; err != nil {
		return false, err
	}
	if ins.IsBing {
		return false, myError.NewET(common.IllegalVisit)
	}
	if ins.Status != "0" {
		return false, myError.New("已审核过的图片禁止删除")
	}
	if err := db.Delete(&ins).Error; err != nil {
		return false, err
	}
	return true, nil
}

func AuditWallPaper(c *gin.Context, reqVo bingVo.WallPaperAuditReqVo) (bool, error) {
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
