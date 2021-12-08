package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	v1 "github.com/devhg/kratos-example/api/blog/v1"
	"github.com/devhg/kratos-example/internal/biz"
)

type BlogService struct {
	v1.UnimplementedBlogServiceServer

	log     *zap.Logger
	article *biz.ArticleUsecase
}

func NewBlogService(article *biz.ArticleUsecase, logger *zap.Logger) *BlogService {
	return &BlogService{article: article, log: logger}
}

func (s *BlogService) CreateArticle(ctx context.Context, req *v1.CreateArticleRequest) (*v1.CreateArticleReply, error) {
	err := s.article.Create(ctx, &biz.Article{
		Title:   req.Title,
		Content: req.Content,
	})

	return &v1.CreateArticleReply{Article: &v1.Article{
		Id:      111,
		Title:   req.Title,
		Content: req.Content,
		Like:    1000,
	}}, err
}

func (s *BlogService) UpdateArticle(ctx context.Context, req *v1.UpdateArticleRequest) (*v1.UpdateArticleReply, error) {
	s.log.Info("input data %v", zap.Any("req", req))
	err := s.article.Update(ctx, req.Id, &biz.Article{
		Title:   req.Title,
		Content: req.Content,
	})
	return &v1.UpdateArticleReply{}, err
}

func (s *BlogService) DeleteArticle(ctx context.Context, req *v1.DeleteArticleRequest) (*v1.DeleteArticleReply, error) {
	s.log.Info("input data %v", zap.Any("req", req))
	err := s.article.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteArticleReply{Msg: "删除成功"}, nil
}

func (s *BlogService) GetArticle(ctx context.Context, req *v1.GetArticleRequest) (*v1.GetArticleReply, error) {
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetArticle")
	defer span.End()

	p, err := s.article.Get(ctx, req.Id)
	s.log.Error("", zap.Error(err))
	if err != nil {
		// err := v1.ErrorArticleNotFound("%v", err)
		err := errors.FromError(err)
		err = err.WithMetadata(map[string]string{
			"foo": "bar",
		})
		return nil, err
	}

	v1a := &v1.Article{
		Id:      p.ID,
		Title:   p.Title,
		Content: p.Content,
		Like:    p.Like,
	}
	return &v1.GetArticleReply{Article: v1a}, nil
}

func (s *BlogService) ListArticle(ctx context.Context, req *v1.ListArticleRequest) (*v1.ListArticleReply, error) {
	ps, err := s.article.List(ctx)
	reply := &v1.ListArticleReply{}
	for _, p := range ps {
		reply.Results = append(reply.Results, &v1.Article{
			Id:      p.ID,
			Title:   p.Title,
			Content: p.Content,
		})
	}
	return reply, err
}
