package server

import (
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"go.uber.org/zap"

	v1 "github.com/devhg/kratos-example/api/blog/v1"
	"github.com/devhg/kratos-example/internal/conf"
	"github.com/devhg/kratos-example/internal/service"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, logger *zap.Logger,
	blog *service.BlogService, comment *service.CommentService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			ratelimit.Server(),
			recovery.Recovery(),
			validate.Validator(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	v1.RegisterBlogServiceHTTPServer(srv, blog)
	v1.RegisterCommentServiceHTTPServer(srv, comment)

	openAPIhandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIhandler)
	return srv
}
