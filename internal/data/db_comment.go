package data

import (
	"context"

	"go.uber.org/zap"

	v1 "github.com/devhg/kratos-example/api/blog/v1"
	"github.com/devhg/kratos-example/internal/biz"
)

type commentRepo struct {
	data *Data
	log  *zap.Logger
}

func NewCommentRepo(data *Data, logger *zap.Logger) biz.CommentRepo {
	return &commentRepo{data: data, log: logger}
}

func (c *commentRepo) ListArticleComments(ctx context.Context, articleID int64) ([]*v1.Comment, error) {
	return nil, nil
}

func (c *commentRepo) CreateComment(ctx context.Context, articleID int64, comment *v1.Comment) error {
	return nil
}

func (c *commentRepo) DeleteArticleComment(ctx context.Context, articleID, commentID int64) error {
	panic("implement me")
}
