package model

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	CompanyId   uint           `gorm:"primaryKey;column:company_id" json:"company_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Business    string         `gorm:"column:business" json:"business"`
	Address     string         `gorm:"column:address" json:"address"`
	Remark      string         `gorm:"column:remark" json:"remark"`
	WebsiteInfo string         `gorm:"column:websiteInfo" json:"websiteInfo"`
	EntryTime   time.Time      `gorm:"column:entryTime" json:"entryTime"`
	CreatedAt   time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deletedAt" json:"deletedAt"`
}
