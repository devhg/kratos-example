package migrate

import (
	"fmt"
	"log"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/jinzhu/gorm"

	"github.com/devhg/kratos-example/internal/conf"
	"github.com/devhg/kratos-example/internal/model"
	"github.com/devhg/kratos-example/internal/utils"
	"github.com/devhg/kratos-example/pkg/dotenv"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string

	Db *gorm.DB

	bconf *conf.Bootstrap
)

func loadBconf(flagConf string) {
	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	bconf = &bc
}

func Close() {
	if Db != nil {
		_ = Db.Close()
	}
}

func Migrate(conf string, db *gorm.DB) error {
	loadBconf(conf)

	defer Close()
	log.Println("正在连接数据库...")

	if db == nil {
		var err error
		db, err = gorm.Open(bconf.Data.Database.Driver, bconf.Data.Database.Source)

		if err != nil {
			return err
		}
	}

	db.LogMode(bconf.Common.Mode != "production")

	// Migrate the schema
	if err := db.AutoMigrate(
		new(model.Admin),           // 管理员表
		new(model.User),            // 用户表
		new(model.WalletCny),       // 钱包 - CNY
		new(model.WalletUsd),       // 钱包 - USD
		new(model.WalletCoin),      // 钱包 - COIN
		new(model.InviteHistory),   // 邀请表
		new(model.LoginLog),        // 登陆成功表
		new(model.TransferLogCny),  // 转账记录 - CNY
		new(model.TransferLogUsd),  // 转账记录 - USD
		new(model.TransferLogCoin), // 转账记录 - COIN
		new(model.FinanceLogCny),   // 流水列表 - CNY
		new(model.FinanceLogUsd),   // 流水列表 - USD
		new(model.FinanceLogCoin),  // 流水列表 - COIN
		new(model.Address),         // 收货地址
		new(model.Banner),          // Banner 表
		new(model.Help),            // 帮助中心
		new(model.WechatOpenID),    // 微信 open_id 外键表
		new(model.OAuth),           // oAuth2 表
	).Error; err != nil {
		return err
	}

	log.Println("数据库同步完成.")

	superAdminInfo := model.Admin{Username: "admin", IsSuper: true}

	// 确保超级管理员账号存在
	if err := db.Where(&superAdminInfo).First(&superAdminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&model.Admin{
				Username:  "admin",
				Name:      "admin",
				Password:  utils.GeneratePassword(dotenv.GetByDefault("ADMIN_DEFAULT_PASSWORD", "123456")),
				Accession: []string{},
				Status:    model.AdminStatusInit,
				IsSuper:   true,
			}).Error

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func Connect() {
	log.Println("正在连接数据库...")

	db, err := gorm.Open(bconf.Data.Database.Driver, bconf.Data.Database.Source)
	if err != nil {
		log.Fatalln(err)
	}
	db.LogMode(bconf.Common.Mode != "production")

	log.Println("连接数据库成功...")
	Db = db
}

// DeleteRowByTable WANING: 该操作会删除数据，并且不可恢复
// 通常只用于测试中
func DeleteRowByTable(tableName string, field string, value interface{}) {
	var (
		err error
		tx  *gorm.DB
	)

	defer func() {
		CommitOrRollback(tx, err)
	}()

	tx = Db.Begin()

	raw := fmt.Sprintf("DELETE FROM \"%s\" WHERE %s = '%s'", tableName, field, value)

	if err = tx.Exec(raw).Error; err != nil {
		return
	}
}

func CommitOrRollback(tx *gorm.DB, err error) {
	fmt.Println(tx == nil, err == nil)
	if tx != nil {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}
}
