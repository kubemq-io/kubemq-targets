package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/middleware"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/sources"
	"github.com/kubemq-hub/kubemq-targets/targets"
)

type Binder struct {
	name   string
	log    *logger.Logger
	source sources.Source
	target targets.Target
	md     middleware.Middleware
	opts   options
}

func NewBinder() *Binder {
	return &Binder{}
}
func (b *Binder) Init(ctx context.Context, cfg config.BindingConfig) error {
	var err error
	b.name = cfg.Name
	b.log = logger.NewLogger(b.name)
	b.opts, err = parseOptions(cfg.Properties)
	if err != nil {
		return fmt.Errorf("error parsing binding %s options", cfg.Name)
	}
	b.target, err = targets.Init(ctx, cfg.Target)
	if err != nil {
		return fmt.Errorf("error loading target conntector %s on binding %s, %w", cfg.Target.Name, b.name, err)
	}
	lm := middleware.NewLogMiddleware(cfg.Name, b.opts.logLevel)
	b.md = middleware.Chain(b.target, middleware.Log(lm))
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
