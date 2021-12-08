package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

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
	ps, err := ar.data.db.Article.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.Article, 0)
	for _, p := range ps {
		rv = append(rv, &biz.Article{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}
	return rv, nil
}

func (ar *articleRepo) GetArticle(ctx context.Context, id int64) (*biz.Article, error) {
	p, err := ar.data.db.Article.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.Article{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}

func (ar *articleRepo) CreateArticle(ctx context.Context, article *biz.Article) error {
	_, err := ar.data.db.Article.
		Create().
		SetTitle(article.Title).
		SetContent(article.Content).
		Save(ctx)
	return err
}

func (ar *articleRepo) UpdateArticle(ctx context.Context, id int64, article *biz.Article) error {
	p, err := ar.data.db.Article.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = p.Update().
		SetTitle(article.Title).
		SetContent(article.Content).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	return err
}

func (ar *articleRepo) DeleteArticle(ctx context.Context, id int64) error {
	return ar.data.db.Article.DeleteOneID(id).Exec(ctx)
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
