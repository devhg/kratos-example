package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	Id        int64  `gorm:"primary_key;AUTO_INCREMENT"`      // 文章ID
	Title     string `gorm:"not null;index;type:varchar(36)"` // 文章标题
	Content   string `gorm:"not null;type:text"`              // 文章内容
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt // 逻辑删除 for Delete()
}

func (news *Article) TableName() string {
	return "articles"
}

func (news *Article) BeforeDelete(db *gorm.DB) error {
	fmt.Println("BeforeDelete")
	return nil
}

func (news *Article) AfterDelete(db *gorm.DB) error {
	fmt.Println("AfterDelete")
	return nil
}

func (news *Article) BeforeCreate(db *gorm.DB) error {
	fmt.Println("BeforeCreate")
	return nil
}

func (news *Article) AfterCreate(db *gorm.DB) error {
	fmt.Println("AfterCreate")
	return nil
}
