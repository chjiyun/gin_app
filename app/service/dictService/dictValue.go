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

func GetDictValue(c *gin.Context, reqVo dictVo.DictValueReqVo) *[]dictVo.DictValueRespVo {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	var respVo []dictVo.DictValueRespVo

	tx := db.Where("type_id = ?", reqVo.TypeId)
	if reqVo.Keyword != "" {
		str := util.WriteString("%", reqVo.Keyword, "%")
		tx.Where("label like ?", str)
	}
	order := "sort asc"
	if reqVo.SortField != "" && reqVo.SortOrder != "" {
		order = util.WriteString(reqVo.SortField, " ", reqVo.SortOrder)
	}
	tx.Order(order).Find(&data)
	_ = copier.Copy(&respVo, &data)

	return &respVo
}

func CreateDictValue(c *gin.Context, reqVo dictVo.DictValueCreateReqVo) (uint64, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	err := copier.Copy(&data, &reqVo)
	if err != nil {
		return 0, err
	}
	if ok, _ := checkValueUnique(c, data.Value, data.TypeId, data.ID); ok {
		return 0, errors.New("字典映射值不能重复")
	}
	if ok, _ := checkLabelUnique(c, data.Label, data.TypeId, data.ID); ok {
		return 0, errors.New("字典映射名称不能重复")
	}
	// 查询并写入最大sort
	maxSort := 0
	db.Select("max(sort)").Find(&maxSort)
	data.Sort = maxSort + 1

	id := idgen.NextId()
	data.ID = id
	err = db.Create(data).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateDictValue(c *gin.Context, reqVo dictVo.DictValueUpdateReqVo) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	if tx := db.First(&data, reqVo.ID); tx.RowsAffected == 0 {
		return false, errors.New("该字典不存在")
	}
	err := copier.Copy(&data, &reqVo)
	if err != nil {
		return false, err
	}
	// 零值也会保存
	err = db.Save(&data).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteDictValue(c *gin.Context, id string) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	err := db.Delete(&model.DictValue{}, id).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func SortDictValue(c *gin.Context, id uint64) {

}

func checkValueUnique(c *gin.Context, value int16, typeId uint64, id uint64) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	var count int64
	// 结构体参数会自动忽略零值
	w := &model.DictValue{
		Value:  value,
		TypeId: typeId,
	}
	w.ID = id
	err := db.Model(&data).Where(&w).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func checkLabelUnique(c *gin.Context, label string, typeId uint64, id uint64) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	var count int64
	// 结构体参数会自动忽略零值
	w := &model.DictValue{
		Label:  label,
		TypeId: typeId,
	}
	w.ID = id
	err := db.Model(&data).Where(&w).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
