package sources

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/middleware"
	"github.com/kubemq-hub/kubemq-targets/sources/command"
	"github.com/kubemq-hub/kubemq-targets/sources/events"
	events_store "github.com/kubemq-hub/kubemq-targets/sources/events-store"
	"github.com/kubemq-hub/kubemq-targets/sources/query"
	"github.com/kubemq-hub/kubemq-targets/sources/queue"
)

type Source interface {
	Init(ctx context.Context, cfg config.Spec) error
	Start(ctx context.Context, target middleware.Middleware) error
	Stop() error
	Name() string
}

func Init(ctx context.Context, cfg config.Spec) (Source, error) {

	switch cfg.Kind {

	case "source.command":
		source := command.New()
		if err := source.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return source, nil
	case "source.query":
		target := query.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "source.events":
		source := events.New()
		if err := source.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return source, nil
	case "source.events-store":
		source := events_store.New()
		if err := source.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return source, nil
	case "queue":
		source := queue.New()
		if err := source.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return source, nil

	default:
		return nil, fmt.Errorf("invalid kind %s for source %s", cfg.Kind, cfg.Name)
	}

}
