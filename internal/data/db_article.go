package data

import (
	"context"
	"fmt"
	"github.com/devhg/kratos-example/internal/model"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

	"github.com/devhg/kratos-example/internal/biz"
)

type articleRepo struct {
	data *Data
	log  *zap.Logger
}

// NewArticleRepo .
func NewArticleRepo(data *Data, logger *zap.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  logger,
	}
}

func (ar *articleRepo) ListArticle(ctx context.Context) ([]*biz.Article, error) {
	dbConn := ar.data.db.WithContext(ctx)
	var articles []*biz.Article
	if result := dbConn.Debug().Find(&articles); result.Error != nil {
		return nil, result.Error
	}
	return articles, nil
}

func (ar *articleRepo) GetArticle(ctx context.Context, id int64) (*biz.Article, error) {
	dbConn := ar.data.db.WithContext(ctx)
	var article biz.Article
	result := dbConn.Debug().
		Where("id=?", id).
		Find(&article)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &article, nil
}

func (ar *articleRepo) CreateArticle(ctx context.Context, article *biz.Article) error {
	timeoutCtx, _ := context.WithTimeout(ctx, time.Second)
	tx := ar.data.db.WithContext(timeoutCtx)
	result := tx.Debug().
		Omit("id", "deleted_at").
		Create(&model.Article{
			Title:     article.Title,
			Content:   article.Content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}) // 通过数据的指针来创建
	return result.Error
}

func (ar *articleRepo) UpdateArticle(ctx context.Context, id int64, article *biz.Article) error {
	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Second)
	tx := ar.data.db.WithContext(timeoutCtx)
	// tx.Debug().Exec("update articles set title=?, content=?, updated_at=? where id=?",
	// 	article.Title, article.Content, time.Now(), id) // 带 timeoutCtx 的更新

	tx.Debug().Model(&model.Article{}).
		Where("id", id).
		Updates(map[string]interface{}{
			"title":      article.Title,
			"content":    article.Content,
			"updated_at": time.Now(),
		})
	return tx.Error
}

func (ar *articleRepo) DeleteArticle(ctx context.Context, id int64) error {
	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Second)
	tx := ar.data.db.WithContext(timeoutCtx)

	// sql := tx.ToSQL(func(tx *gorm.DB) *gorm.DB {
	// 	return tx.Debug().Where("id = ?", id).Delete(&model.Article{})
	// })
	// fmt.Println(sql)
	tx.Debug().
		Where("id = ?", id).
		Delete(&model.Article{})
	return tx.Error
}

func likeKey(id int64) string {
	return fmt.Sprintf("like:%d", id)
}

func (ar *articleRepo) GetArticleLike(ctx context.Context, id int64) (rv int64, err error) {
	get := ar.data.rdb.Get(ctx, likeKey(id))
	rv, err = get.Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return
}

func (ar *articleRepo) IncArticleLike(ctx context.Context, id int64) error {
	_, err := ar.data.rdb.Incr(ctx, likeKey(id)).Result()
	return err
}
