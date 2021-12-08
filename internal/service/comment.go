package service

import (
	"context"
	v1 "github.com/devhg/kratos-example/api/blog/v1"
	"github.com/devhg/kratos-example/internal/biz"
	"go.uber.org/zap"
	"time"
)

type CommentService struct {
	v1.UnimplementedCommentServiceServer

	log     *zap.Logger
	comment *biz.CommentUsecase
}

func NewCommentService(log *zap.Logger, comment *biz.CommentUsecase) *CommentService {
	return &CommentService{log: log, comment: comment}
}

func (c *CommentService) CreateComment(ctx context.Context, req *v1.CreateCommentRequest) (*v1.CreateCommentReply, error) {
	c.log.Info("CreateComment", zap.Any("req", req))
	comment := &biz.Comment{
		Name:    req.Name,
		Content: req.Content,
		Date:    time.Now(),
	}
	err := c.comment.Create(ctx, req.ArticleId, comment)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCommentReply{}, nil
}
