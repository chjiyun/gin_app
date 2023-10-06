package poetryService

import (
	"errors"
	"gin_app/app/controller/poetryController/poetryVo"
	"gin_app/app/model"
	"gin_app/app/util"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"os"
	"regexp"
	"strings"
)

// SearchPoetry 模糊搜索诗词
func SearchPoetry(c *gin.Context, reqVo poetryVo.PoetrySearchReqVo) ([]poetryVo.PoetryRespVo, error) {
	reqVo.Keyword = strings.TrimSpace(reqVo.Keyword)
	if reqVo.Keyword == "" {
		return nil, errors.New("invalid keyword")
	}

	db := c.Value("DB").(*gorm.DB)
	var data []model.Poetry
	var respVos []poetryVo.PoetryRespVo

	reg := regexp.MustCompile(`\s+`)
	keywords := reg.Split(reqVo.Keyword, -1)
	var str string
	var str1 string
	isMultiKeyword := len(keywords) > 1

	if isMultiKeyword {
		str = util.WriteString("%", keywords[0], "%")
		str1 = util.WriteString("%", keywords[1], "%")
	} else {
		str = util.WriteString("%", reqVo.Keyword, "%")
		str1 = str
	}
	//db.Select("b_poetry.id, b_poetry.author, b_poetry.title, b_poetry.tag, pc.content, pc.sort").
	//	Joins("left join b_poetry_content as pc on pc.poetry_id = b_poetry.id").
	//	Model(&model.Poetry{}).
	//	Where("title like ?", str).Or("author like ?", str).
	//	Or("exists (select 1 from b_poetry_content where poetry_id = b_poetry.id and content like ?)", str).
	//	Limit(10).Find(&data)

	tx := db.Select("id, author, title, tag").
		Preload("PoetryContent", func(db *gorm.DB) *gorm.DB {
			return db.Select("poetry_id, content, sort").Where("content like ?", str1).Order("sort")
		}).
		Where(db.Where("title like ?", str).Or("author like ?", str))
	if isMultiKeyword {
		tx.Where("exists (select 1 from b_poetry_content where poetry_id = b_poetry.id and content like ?)", str1)
	} else {
		tx.Or("exists (select 1 from b_poetry_content where poetry_id = b_poetry.id and content like ?)", str1)
	}
	tx.Limit(10).Find(&data)

	_ = copier.Copy(&respVos, &data)
	return respVos, nil
}

func GetPoetry(c *gin.Context, id string) poetryVo.PoetryRespVo {
	db := c.Value("DB").(*gorm.DB)
	var poetry model.Poetry
	var respVo poetryVo.PoetryRespVo

	err := db.Preload("PoetryContent").Find(&poetry, id).Error
	if err != nil {
		return respVo
	}
	_ = copier.Copy(&respVo, &poetry)
	return respVo
}

func CreatePoetry(c *gin.Context) (bool, error) {

	return true, nil
}

func PoetryImport(c *gin.Context) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data []model.Poetry
	var contents []model.PoetryContent
	var source []poetryVo.PoetryImportReqVo

	// 读取本地json
	bytes, err := os.ReadFile("../chinese-poetry/宋词/data2_x.json")
	if err != nil {
		return false, err
	}
	if err = sonic.Unmarshal(bytes, &source); err != nil {
		return false, err
	}
	for _, item := range source {
		var d model.Poetry
		_ = copier.Copy(&d, &item)
		id := idgen.NextId()
		d.ID = id
		d.Tag = 2
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
	//fmt.Println(data)
	//fmt.Println(contents)

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&model.Poetry{}).CreateInBatches(&data, 400).Error
		if err != nil {
			return err
		}
		err = tx.Model(&model.PoetryContent{}).CreateInBatches(&contents, 400).Error
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
