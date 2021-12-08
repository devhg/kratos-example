package biz

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type CommentRepo interface {
	ListArticleComments(ctx context.Context, articleID int64) ([]*Comment, error)
	CreateComment(ctx context.Context, articleID int64, comment *Comment) error
	DeleteArticleComment(ctx context.Context, articleID, commentID int64) error
}

type CommentUsecase struct {
	repo CommentRepo
}

type Comment struct {
	Name    string
	Content string
	Date    time.Time
}

func NewCommentUsecase(repo CommentRepo, logger *zap.Logger) *CommentUsecase {
	return &CommentUsecase{repo: repo}
}

func (uc *CommentUsecase) Create(ctx context.Context, articleID int64, comment *Comment) error {
	return uc.repo.CreateComment(ctx, articleID, comment)
}
