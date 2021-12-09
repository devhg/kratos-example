package model

import (
	"gorm.io/gorm"
	"time"
)

type AdminStatus int32

const (
	AdminStatusBanned      AdminStatus = -100 // 账号被禁用
	AdminStatusInactivated AdminStatus = -1   // 账号未激活
	AdminStatusInit        AdminStatus = 0    // 初始化状态
)

type Admin struct {
	Id        int64       `gorm:"primary_key" json:"id"`                                  // 用户ID
	Username  string      `gorm:"not null;unique;index;type:varchar(36)" json:"username"` // 用户名, 用于登陆
	Name      string      `gorm:"not null;index;type:varchar(36)" json:"name"`            // 管理员名
	Password  string      `gorm:"not null;type:varchar(36)" json:"password"`              // 登陆密码
	IsSuper   bool        `gorm:"not null;" json:"is_super"`                              // 是否是超级管理员, 超级管理员全站应该只有一个
	Status    AdminStatus `gorm:"not null;" json:"status"`                                // 状态
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (news *Admin) TableName() string {
	return "admin"
}

func (news *Admin) BeforeCreate(db *gorm.DB) error {
	// panic("implement me")
	return nil
}

// func (news *Admin) BeforeCreate(scope *gorm.DB) error {
// 	return scope.SetColumn("id", utils.GenerateId())
// }
