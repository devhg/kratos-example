package service

import (
	"context"
	v1 "github.com/devhg/kratos-example/api/blog/v1"
	"github.com/devhg/kratos-example/internal/biz"
	"go.uber.org/zap"
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
	comment := &v1.Comment{
		Name:     req.Name,
		Content:  req.Content,
	}
	err := c.comment.Create(ctx, req.ArticleId, comment)
	if err != nil {
		return nil, err
	}
	return &v1.CreateCommentReply{}, nil
}

func (c *CommentService) ListArticleComment(ctx context.Context, req *v1.ListCommentReq) (*v1.ListCommentReply, error) {
	c.log.Info("req", zap.Any("req", req))
	comments, err := c.comment.ListComments(ctx, req.ArticleId)
	if err != nil {
		return nil, err
	}
	return &v1.ListCommentReply{Comments: comments}, nil
}
