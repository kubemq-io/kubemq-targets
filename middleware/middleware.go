package middleware

import (
	"context"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/logger"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/ratelimit"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
)

type Middleware interface {
	Do(ctx context.Context, request *types.Request) (*types.Response, error)
}

type DoFunc func(ctx context.Context, request *types.Request) (*types.Response, error)

func (df DoFunc) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	return df(ctx, request)
}

type MiddlewareFunc func(Middleware) Middleware

func Log(log *logger.Logger) MiddlewareFunc {
	return func(df Middleware) Middleware {
		return DoFunc(func(ctx context.Context, request *types.Request) (*types.Response, error) {
			result, err := df.Do(ctx, request)
			if err != nil {
				log.Error(err.Error())
			}
			return result, err
		})
	}
}
func Limit(rl ratelimit.Limiter) MiddlewareFunc {
	return func(df Middleware) Middleware {
		return DoFunc(func(ctx context.Context, request *types.Request) (*types.Response, error) {
			if rl != nil {
				_ = rl.Take()
			}
			return df.Do(ctx, request)
		})
	}
}

func Chain(md Middleware, list ...MiddlewareFunc) Middleware {
	chain := md
	for _, middleware := range list {
		chain = middleware(chain)
	}
	return chain
}
