package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/middleware"
	"github.com/kubemq-hub/kubemq-target-connectors/sources"
	"github.com/kubemq-hub/kubemq-target-connectors/targets"
)

type Binder struct {
	source sources.Source
	target targets.Target
	md     middleware.Middleware
}

func New() *Binder {
	return &Binder{}
}
func (b *Binder) Init(ctx context.Context, cfg config.BindingConfig) error {
	var err error
	b.target, err = targets.Init(ctx, cfg.Target)
	if err != nil {
		return err
	}
	b.md = middleware.Chain(b.target)
	b.source, err = sources.Init(ctx, cfg.Source)
	if err != nil {
		return err
	}
	return nil
}

func (b *Binder) Start(ctx context.Context) error {
	if b.md == nil {
		return fmt.Errorf("no valid initialzed target middleware found")
	}
	if b.source == nil {
		return fmt.Errorf("no valid initialzed source found")
	}
	return b.source.Start(ctx, b.md)
}
func (b *Binder) Stop() error {
	return b.source.Stop()
}
