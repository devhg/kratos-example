package data

import (
	"context"
	"github.com/devhg/kratos-example/internal/biz"
	"go.uber.org/zap"
)

type commentRepo struct {
	data *Data
	log  *zap.Logger
}

func NewCommentRepo(data *Data, logger *zap.Logger) biz.CommentRepo {
	return &commentRepo{data: data, log: logger}
}

func (c *commentRepo) ListArticleComments(ctx context.Context, articleID int64) ([]*biz.Comment, error) {
	return nil, nil
}

func (c *commentRepo) CreateComment(ctx context.Context, articleID int64, comment *biz.Comment) error {
	_, err := c.data.db.Comment.
		Create().
		SetName(comment.Name).
		SetContent(comment.Content).
		SetPostID(articleID).
		Save(ctx)
	return err
}

func (c *commentRepo) DeleteArticleComment(ctx context.Context, articleID, commentID int64) error {
	panic("implement me")
}
