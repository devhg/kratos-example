package model

import (
	"errors"
	"github.com/devhg/kratos-example/internal/utils"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type UserStatus int32

type Gender int

const (
	// 用户状态
	UserStatusBanned      UserStatus = -100 // 账号被禁用
	UserStatusInactivated UserStatus = -1   // 账号未激活
	UserStatusInit        UserStatus = 1    // 初始化状态

	// 用户性别
	GenderUnknown Gender = 0 // 未知性别
	GenderMale    Gender = 1 // 男
	GenderFemale  Gender = 2 // 女
)

type User struct {
	Id                      string         `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 用户ID
	Username                string         `gorm:"not null;type:varchar(36)unique;index" json:"username"`        // 用户名
	Password                string         `gorm:"not null;type:varchar(36);index" json:"password"`              // 登陆密码
	PayPassword             *string        `gorm:"null;type:varchar(36)" json:"pay_password"`                    // 支付密码
	Nickname                *string        `gorm:"null;type:varchar(36)" json:"nickname"`                        // 昵称
	Phone                   *string        `gorm:"null;unique;type:varchar(16);index" json:"phone"`              // 手机号
	Email                   *string        `gorm:"null;unique;type:varchar(36);index" json:"email"`              // 邮箱
	Status                  UserStatus     `gorm:"not null" json:"status"`                                       // 状态
	Role                    pq.StringArray `gorm:"not null;type:varchar(36)[]" json:"role"`                      // 角色, 用户可以拥有多个角色
	Avatar                  string         `gorm:"not null;type:varchar(128)" json:"avatar"`                     // 头像
	Level                   int32          `gorm:"default(1)" json:"level"`                                      // 用户等级
	Gender                  Gender         `gorm:"default(0)" json:"gender"`                                     // 性别
	EnableTOTP              bool           `gorm:"not null;" json:"enable_totp"`                                 // 是否启用双重身份认证
	Secret                  string         `gorm:"not null;type:varchar(32)" json:"secret"`                      // 用户自己的密钥
	InviteCode              string         `gorm:"not null;unique;type:varchar(8)" json:"invite_code"`           // 用户的邀请码，邀请码唯一
	UsernameRenameRemaining int            `gorm:"not null;" json:"username_rename_remaining"`                   // 用户名还有几次重新更改的机会， 主要是如果用第三方注册登陆，则用户名随机生成，这里给用户一个重新命名的机会

	// 外键关联
	WechatOpenID *string       `gorm:"null;unique;type:varchar(255);index" json:"wechat_open_id"` // 绑定的微信帐号 open_id
	Wechat       *WechatOpenID `gorm:"foreignkey:wechatOpenID" json:"wechat"`                     // **外键**
	OAuth        []OAuth       `gorm:"foreignkey:Uid" json:"oauth"`                               // **外键**

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	uid := utils.GenerateId()
	if err := scope.SetColumn("id", uid); err != nil {
		return err
	}

	// 生成邀请码
	if err := scope.SetColumn("invite_code", utils.GenerateInviteCode()); err != nil {
		return err
	}

	// 默认关闭启用谷歌验证码
	if err := scope.SetColumn("enable_totp", false); err != nil {
		return err
	}

	// 生成用户自己的密钥
	if secret, err := utils.Generate2FASecret(uid); err != nil {
		return err
	} else {
		if err := scope.SetColumn("secret", secret); err != nil {
			return err
		}
	}

	return nil
}

// 检查用户状态是否正常
func (u *User) CheckStatusValid() error {
	if u.Status != UserStatusInit {
		switch u.Status {
		case UserStatusInactivated:
			return errors.New("exception.UserIsInActive")
		case UserStatusBanned:
			return errors.New("exception.UserHaveBeenBan")
		}
	}

	return nil
}
