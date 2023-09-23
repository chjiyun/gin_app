package poetryService

import (
	"fmt"
	"gin_app/app/controller/poetryController/poetryVo"
	"gin_app/app/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"os"
)

// SearchPoetry 模糊搜索诗词
func SearchPoetry(c *gin.Context, keyword string) *[]poetryVo.PoetryRespVo {
	db := c.Value("DB").(*gorm.DB)
	var respVos []poetryVo.PoetryRespVo

	db.Find(&respVos)
	return &respVos
}

func GetPoetry(c *gin.Context, id string) *poetryVo.PoetryRespVo {
	db := c.Value("DB").(*gorm.DB)
	var poetry model.Poetry
	var contents []model.PoetryContent
	var respVo poetryVo.PoetryRespVo

	err := db.Find(&poetry, id).Error
	if err != nil {
		return nil
	}
	db.Where("poetry_id = ?", id).Find(&contents)
	_ = copier.Copy(&respVo, &poetry)
	var content []string
	for _, item := range contents {
		content = append(content, item.Content)
	}
	respVo.Content = content
	return &respVo
}

func CreatePoetry(c *gin.Context) (bool, error) {

	return true, nil
}

func ImportPoetry(c *gin.Context) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data []model.Poetry
	var contents []model.PoetryContent
	var source []poetryVo.PoetryImportReqVo

	// 读取本地json
	bytes, err := os.ReadFile("../chinese-poetry/唐诗/data5.json")
	if err != nil {
		return false, err
	}
	if err = jsoniter.Unmarshal(bytes, &source); err != nil {
		return false, err
	}
	for _, item := range source {
		var d model.Poetry
		_ = copier.Copy(&d, &item)
		id := idgen.NextId()
		d.ID = id
		d.Tag = 1
		data = append(data, d)
		for i, content := range item.Content {
			var pc model.PoetryContent
			pc.Content = content
			pc.Sort = i
			pc.PoetryId = id
			pc.ID = idgen.NextId()
			contents = append(contents, pc)
		}
	}
	fmt.Println(data)
	fmt.Println(contents)

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&model.Poetry{}).CreateInBatches(&data, 200).Error
		if err != nil {
			return err
		}
		err = tx.Model(&model.PoetryContent{}).CreateInBatches(&contents, 200).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
