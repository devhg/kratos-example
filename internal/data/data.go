package data

import (
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/devhg/kratos-example/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewArticleRepo, NewCommentRepo)

// Data .
type Data struct {
	db  *gorm.DB
	rdb *redis.Client
}

// NewData .
func NewData(conf *conf.Data, logger *zap.Logger) (*Data, func(), error) {
	sqlConf := mysql.Config{
		DriverName: conf.Database.Driver,
		DSN:        conf.Database.Source, // Data Source Name，参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
	}
	db, err := gorm.Open(mysql.New(sqlConf), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: gormLogger.New(
			log.New(os.Stdout, "\r", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			gormLogger.Config{
				SlowThreshold:             time.Second,       // 慢 SQL 阈值
				LogLevel:                  gormLogger.Silent, // 日志级别
				IgnoreRecordNotFoundError: true,              // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,              // 禁用彩色打印
			},
		),
	})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	d := &Data{
		db:  db,
		rdb: rdb,
	}
	return d, func() {
		logger.Info("closing the data resources")
		sqlDB, _ := d.db.DB()
		if err := sqlDB.Close(); err != nil {
			logger.Error("db close err", zap.Error(err))
		}
		if err := d.rdb.Close(); err != nil {
			logger.Error("rdb close err", zap.Error(err))
		}
	}, nil
}
