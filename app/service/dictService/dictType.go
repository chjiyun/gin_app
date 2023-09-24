package dictService

import (
	"errors"
	"gin_app/app/controller/dictController/dictVo"
	"gin_app/app/model"
	"gin_app/app/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

// GetDictType 获取所有字典类型
func GetDictType(c *gin.Context, keyword string) *[]dictVo.DictTypeRespVo {
	db := c.Value("DB").(*gorm.DB)

	var data []model.DictType
	var respVo []dictVo.DictTypeRespVo

	tx := db.Model(&model.DictType{})
	if keyword != "" {
		str := util.WriteString("%", keyword, "%")
		tx.Where(db.Where("name like ?", str).Or("value like ?", str))
	}
	tx.Find(&data)
	_ = copier.Copy(&respVo, &data)
	return &respVo
}

func CreateDictType(c *gin.Context, reqVo dictVo.DictTypeCreateReqVo) (uint64, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictType
	err := copier.Copy(&data, &reqVo)
	if err != nil {
		return 0, err
	}
	id := idgen.NextId()
	data.ID = id
	db.Create(data)
	return id, nil
}

func UpdateDictType(c *gin.Context, reqVo dictVo.DictTypeUpdateReqVo) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictType
	if tx := db.First(&data, reqVo.ID); tx.RowsAffected == 0 {
		return false, errors.New("该字典类型不存在")
	}
	err := copier.Copy(&data, &reqVo)
	if err != nil {
		return false, err
	}
	db.Save(&data)
	return true, nil
}

func DeleteDictType(c *gin.Context, id string) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	err := db.Delete(&model.DictType{}, id).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
