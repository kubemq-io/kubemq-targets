package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/middleware"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/logger"
	"github.com/kubemq-hub/kubemq-target-connectors/sources"
	"github.com/kubemq-hub/kubemq-target-connectors/targets"
)

type Binder struct {
	name   string
	log    *logger.Logger
	source sources.Source
	target targets.Target
	md     middleware.Middleware
}

func New() *Binder {
	return &Binder{}
}
func (b *Binder) Init(ctx context.Context, cfg config.BindingConfig) error {
	var err error
	b.name = cfg.Name
	b.log = logger.NewLogger(b.name)
	b.target, err = targets.Init(ctx, cfg.Target)
	if err != nil {
		return fmt.Errorf("error loading target conntector %s on binding %s, %w", cfg.Target.Name, b.name, err)
	}
	b.md = middleware.Chain(b.target)
	b.source, err = sources.Init(ctx, cfg.Source)
	if err != nil {
		return fmt.Errorf("error loading source conntector %s on binding %s, %w", cfg.Source.Name, b.name, err)
	}
	b.log.Infof("binding %s initialized successfully", b.name)
	return nil
}

func (b *Binder) Start(ctx context.Context) error {
	if b.md == nil {
		return fmt.Errorf("error starting binding connector %s,no valid initialzed target middleware found", b.name)
	}
	if b.source == nil {
		return fmt.Errorf("error starting binding connector %s,no valid initialzed source found", b.name)
	}
	err := b.source.Start(ctx, b.md)
	if err != nil {
		return err
	}
	b.log.Infof("binding %s started successfully", b.name)
	return nil
}
func (b *Binder) Stop() error {
	err := b.source.Stop()
	if err != nil {
		return err
	}
	b.log.Infof("binding %s stopped successfully", b.name)
	return nil
}
