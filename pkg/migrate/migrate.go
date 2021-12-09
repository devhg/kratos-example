package migrate

import (
	"fmt"
	"gorm.io/driver/mysql"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gorm.io/gorm"

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

func Migrate(conf string, db *gorm.DB) error {
	loadBconf(conf)

	log.Println("正在连接数据库...")

	if db == nil {
		var err error
		sqlConf := mysql.Config{
			DriverName: bconf.Data.Database.Driver,
			DSN:        bconf.Data.Database.Source, // Data Source Name，参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
		}
		db, err = gorm.Open(mysql.New(sqlConf), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger: gormLogger.New(
				log.New(os.Stdout, "\r", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
				gormLogger.Config{
					SlowThreshold:             time.Second,      // 慢 SQL 阈值
					LogLevel:                  gormLogger.Error, // 日志级别
					IgnoreRecordNotFoundError: true,             // 忽略ErrRecordNotFound（记录未找到）错误
					Colorful:                  true,             // 禁用彩色打印
				},
			),
		})
		if err != nil {
			return err
		}
	}

	// Migrate the schema
	if err := db.AutoMigrate(
		new(model.Article),
		new(model.Admin),   // 管理员表
		new(model.Address), // 收货地址
	); err != nil {
		return err
	}

	log.Println("数据库同步完成.")

	superAdminInfo := model.Admin{Username: "admin", IsSuper: true}

	// 确保超级管理员账号存在
	if err := db.Where(&superAdminInfo).First(&superAdminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&model.Admin{
				Username: "admin",
				Name:     "admin",
				Password: utils.GeneratePassword(dotenv.GetByDefault("ADMIN_DEFAULT_PASSWORD", "123456")),
				Status:   model.AdminStatusInit,
				IsSuper:  true,
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

	var err error
	sqlConf := mysql.Config{
		DriverName: bconf.Data.Database.Driver,
		DSN:        bconf.Data.Database.Source, // Data Source Name，参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
	}
	Db, err = gorm.Open(mysql.New(sqlConf), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: gormLogger.New(
			log.New(os.Stdout, "\r", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			gormLogger.Config{
				SlowThreshold:             time.Second,     // 慢 SQL 阈值
				LogLevel:                  gormLogger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,            // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,            // 禁用彩色打印
			},
		),
	})
	if err != nil {
		panic(err)
	}

	log.Println("连接数据库成功...")
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
