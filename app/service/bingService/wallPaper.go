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
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func GetWallPaperPage(c *gin.Context, reqVo bingVo.WallPaperReqVo) (common.PageRes, error) {
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
		respVo[i].FileId = ""
	}
	pageRes.Count = count
	pageRes.Rows = respVo
	return pageRes, nil
}

func GetWallPaper(c *gin.Context, id string) (bingVo.WallPaperRespVo, error) {
	db := c.Value("DB").(*gorm.DB)
	var bing model.Bing
	var thumb model.Thumb
	var respVo bingVo.WallPaperRespVo

	if err := db.Take(&bing, "id = ?", id).Error; err != nil {
		return respVo, err
	}
	_ = copier.Copy(&respVo, &bing)
	err := db.Where("file_id = ? and ext = ?", bing.FileId, "webp").First(&thumb).Error
	if err != nil {
		return respVo, err
	}
	respVo.Name = thumb.Name
	respVo.Width = thumb.Width
	respVo.Height = thumb.Height
	respVo.Ext = util.GetFileExt(respVo.FileId)
	respVo.ThumbId = util.WriteString(util.ToString(thumb.ID), ".", thumb.Ext)
	return respVo, nil
}

// AddWallPaper 手动上传壁纸 审核通过后方可进入file
func AddWallPaper(c *gin.Context, reqVo bingVo.WallPaperCreateReqVo) (bool, error) {
	db := c.Value("DB").(*gorm.DB)
	var isExist int
	err := db.Select("1").Where("file_id = ? and status = ?", reqVo.FileId, "0").Scan(&isExist).Error
	if err != nil {
		return false, err
	}
	if isExist == 1 {
		return false, myError.New("存在待审核的图片")
	}
	// 校验图片格式和质量  大小 分辨率
	file, err := service.GetFile(c, reqVo.FileId)
	// 文件待审核去重
	if err != nil {
		return false, err
	}
	if err := validateLocalImage(file.Path); err != nil {
		return false, err
	}

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
	c.Set("genThumb", false)
	fileId, err := service.Upload(c, f)
	if err != nil {
		return "", err
	}
	return fileId, nil
}

// 识别已上传的图片文件
func validateLocalImage(sourcePath string) error {
	f, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer f.Close()
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
	if _, err = f.Seek(0, 0); err != nil {
		return err
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
		if err = validateLocalImage(file.Path); err != nil {
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
