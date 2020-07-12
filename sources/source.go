package sources

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/middleware"
	"github.com/kubemq-hub/kubemq-target-connectors/sources/command"
	"github.com/kubemq-hub/kubemq-target-connectors/sources/events"
	events_store "github.com/kubemq-hub/kubemq-target-connectors/sources/events-store"
	"github.com/kubemq-hub/kubemq-target-connectors/sources/query"
	"github.com/kubemq-hub/kubemq-target-connectors/sources/queue"
)

type Source interface {
	Init(ctx context.Context, cfg config.Metadata) error
	Start(ctx context.Context, target middleware.Middleware) error
	Stop() error
	Name() string
}

func Init(ctx context.Context, cfg config.Metadata) (Source, error) {

	switch cfg.Kind {

	case "source.kubemq.command":
		source := command.New()
		if err := source.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return source, nil
	case "source.kubemq.query":
		target := query.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "source.kubemq.events":
		source := events.New()
		if err := source.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return source, nil
	case "source.kubemq.events-store":
		source := events_store.New()
		if err := source.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return source, nil
	case "source.kubemq.queue":
		source := queue.New()
		if err := source.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return source, nil

	default:
		return nil, fmt.Errorf("invalid kind %s for source %s", cfg.Kind, cfg.Name)
	}

}
