package biz

import (
	"context"
	v1 "github.com/devhg/kratos-example/api/blog/v1"
	"go.uber.org/zap"
)

type CommentRepo interface {
	ListArticleComments(ctx context.Context, articleID int64) ([]*v1.Comment, error)
	CreateComment(ctx context.Context, articleID int64, comment *v1.Comment) error
	DeleteArticleComment(ctx context.Context, articleID, commentID int64) error
}

type CommentUsecase struct {
	repo CommentRepo
}

func NewCommentUsecase(repo CommentRepo, logger *zap.Logger) *CommentUsecase {
	return &CommentUsecase{repo: repo}
}

func (uc *CommentUsecase) Create(ctx context.Context, articleID int64, comment *v1.Comment) error {
	return uc.repo.CreateComment(ctx, articleID, comment)
}

func (uc *CommentUsecase) ListComments(ctx context.Context, articleID int64) ([]*v1.Comment, error) {
	comments, err := uc.repo.ListArticleComments(ctx, articleID)
	if err != nil {
		return nil, err
	}
	var ret []*v1.Comment
	for _, comment := range comments {
		ret = append(ret, &v1.Comment{
			Id:      comment.Id,
			Name:    comment.Name,
			Content: comment.Content,
			// UpdateAt:  timestamppb.New(comment.UpdateAt),
			UpdateAt:  comment.UpdateAt,
			ArticleId: articleID,
		})
	}
	return ret, nil
}
