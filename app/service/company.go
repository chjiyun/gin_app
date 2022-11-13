package service

import (
	"gin_app/app/model"
	"gin_app/app/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 查询企业
func GetCompanies(c *gin.Context) {
	db := c.Value("DB").(*gorm.DB)
	log := c.Value("Logger").(*logrus.Entry)

	var companies []model.Company
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	util.SetDefault(&page, 1)
	util.SetDefault(&pageSize, 10)
	offset := (page - 1) * pageSize

	res := db.Table("company as c").Limit(10).Offset(offset).Find(&companies)
	if res.Error != nil {
		log.Panicln(res.Error)
	}

	// 切换数据库
	// var list []model.Company
	// tx := db.Begin()
	// tx.Exec("use egg_test")
	// tx.Table("company as c").Limit(10).Offset(offset).Find(&list)
	// tx.Commit()

	log.Info("companies: ", companies)

	c.JSON(200, gin.H{"xc_dev": companies})
}
