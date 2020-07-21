package middleware

import (
	"context"
	"github.com/kubemq-hub/kubemq-targets/pkg/retry"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Middleware interface {
	Do(ctx context.Context, request *types.Request) (*types.Response, error)
}

type DoFunc func(ctx context.Context, request *types.Request) (*types.Response, error)

func (df DoFunc) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	return df(ctx, request)
}

type MiddlewareFunc func(Middleware) Middleware

func Log(log *LogMiddleware) MiddlewareFunc {
	return func(df Middleware) Middleware {
		return DoFunc(func(ctx context.Context, request *types.Request) (*types.Response, error) {
			result, err := df.Do(ctx, request)
			switch log.minLevel {
			case "debug":
				reqStr := ""
				if request != nil {
					reqStr = request.String()
				}
				resStr := ""
				if result != nil {
					resStr = result.String()
				}
				log.Infof("request: %s, response: %s, error:%+v", reqStr, resStr, err)
			case "info":
				if err != nil {
					log.Errorf("error processing request: %s", err.Error())
				} else {
					log.Infof("request processed with successful response")
				}
			case "error":
				if err != nil {
					log.Errorf("error processing request: %s", err.Error())
				}
			}
			return result, err
		})
	}
}
func RateLimiter(rl *RateLimitMiddleware) MiddlewareFunc {
	return func(df Middleware) Middleware {
		return DoFunc(func(ctx context.Context, request *types.Request) (*types.Response, error) {
			rl.Take()
			return df.Do(ctx, request)
		})
	}
}

func Retry(r *RetryMiddleware) MiddlewareFunc {
	return func(df Middleware) Middleware {
		return DoFunc(func(ctx context.Context, request *types.Request) (*types.Response, error) {
			var resp *types.Response
			err := retry.Do(func() error {
				var doErr error
				resp, doErr = df.Do(ctx, request)
				if doErr != nil {
					return doErr
				}
				return nil
			}, r.opts...)
			return resp, err
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
