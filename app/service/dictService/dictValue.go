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

func GetDictValue(c *gin.Context, reqVo dictVo.DictValueReqVo) []dictVo.DictValueRespVo {
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

	return respVo
}

func GetDictValueByType(c *gin.Context, value []string) ([]dictVo.DictValueRespVo, error) {
	if len(value) == 0 {
		return nil, errors.New("empty value")
	}

	db := c.Value("DB").(*gorm.DB)
	var data []model.DictValue
	var dictType []model.DictType
	var respVo []dictVo.DictValueRespVo

	err := db.Where("value in ?", value).Find(&dictType).Error
	if err != nil {
		return nil, err
	}
	if len(dictType) == 0 {
		return respVo, nil
	}
	typeIds := make([]uint64, 0, len(dictType))
	for _, item := range dictType {
		typeIds = append(typeIds, item.ID)
	}
	err = db.Where("type_id in ?", typeIds).Order("sort asc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	_ = copier.Copy(&respVo, &data)
	return respVo, nil
}

func CreateDictValue(c *gin.Context, reqVo dictVo.DictValueCreateReqVo) (uint64, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	err := copier.Copy(&data, &reqVo)
	if err != nil {
		return 0, err
	}
	if ok, err := checkDictValueValue(c, data.Value, data.TypeId, data.ID); !ok {
		return 0, err
	}
	if ok, err := checkDictValueLabel(c, data.Label, data.TypeId, data.ID); !ok {
		return 0, err
	}
	// 查询并写入最大sort
	maxSort := 0
	db.Model(&model.DictValue{}).Select("max(sort)").Find(&maxSort)
	data.Sort = maxSort + 1

	id := idgen.NextId()
	data.ID = id
	err = db.Create(&data).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateDictValue(c *gin.Context, reqVo dictVo.DictValueUpdateReqVo) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	if err := db.First(&data, reqVo.ID).Error; err != nil {
		return false, errors.New("该字典映射不存在")
	}
	err := copier.Copy(&data, &reqVo)
	if err != nil {
		return false, err
	}
	if ok, err := checkDictValueValue(c, data.Value, data.TypeId, data.ID); !ok {
		return false, err
	}
	if ok, err := checkDictValueLabel(c, data.Label, data.TypeId, data.ID); !ok {
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

func checkDictValueValue(c *gin.Context, value string, typeId uint64, id uint64) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	var count int64
	// 结构体参数会自动忽略零值
	w := model.DictValue{
		Value:  value,
		TypeId: typeId,
	}
	w.ID = id
	err := db.Model(&data).Where(&w).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, errors.New("字典对应值不能重复")
	}
	return true, nil
}

func checkDictValueLabel(c *gin.Context, label string, typeId uint64, id uint64) (bool, error) {
	db := c.Value("DB").(*gorm.DB)

	var data model.DictValue
	var count int64
	// 结构体参数会自动忽略零值
	w := model.DictValue{
		Label:  label,
		TypeId: typeId,
	}
	w.ID = id
	err := db.Model(&data).Where(&w).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, errors.New("字典名不能重复")
	}
	return true, nil
}
