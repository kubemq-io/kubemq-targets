package binding

import (
	"context"
	"fmt"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/middleware"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/pkg/metrics"
	"github.com/kubemq-io/kubemq-targets/sources"
	"github.com/kubemq-io/kubemq-targets/targets"
)

type Binder struct {
	name   string
	log    *logger.Logger
	source sources.Source
	target targets.Target
	md     middleware.Middleware
}

func NewBinder() *Binder {
	return &Binder{}
}

func (b *Binder) buildMiddleware(cfg config.BindingConfig, exporter *metrics.Exporter) (middleware.Middleware, error) {
	retry, err := middleware.NewRetryMiddleware(cfg.Properties, b.log)
	if err != nil {
		return nil, err
	}
	rateLimiter, err := middleware.NewRateLimitMiddleware(cfg.Properties)
	if err != nil {
		return nil, err
	}
	met, err := middleware.NewMetricsMiddleware(cfg, exporter)
	if err != nil {
		return nil, err
	}
	meta, err := middleware.NewMetadataMiddleware(cfg.Properties)
	if err != nil {
		return nil, err
	}
	md := middleware.Chain(b.target, middleware.RateLimiter(rateLimiter), middleware.Retry(retry), middleware.Metric(met), middleware.Metadata(meta))
	return md, nil
}

func (b *Binder) Init(ctx context.Context, cfg config.BindingConfig, exporter *metrics.Exporter, logLevel string) error {
	b.name = cfg.Name
	var err error
	b.log = logger.NewLogger(fmt.Sprintf("binding-%s", cfg.Name), logLevel)

	b.target, err = targets.Init(ctx, cfg.Target, b.log)
	if err != nil {
		return fmt.Errorf("error loading target conntector on binding %s, %w", b.name, err)
	}
	b.log.Infof("binding: %s target: initialized successfully", b.name)
	b.md, err = b.buildMiddleware(cfg, exporter)
	if err != nil {
		return fmt.Errorf("error loading middlewares on binding %s, %w", b.name, err)
	}
	b.source, err = sources.Init(ctx, cfg.Source, cfg.Name, b.log)
	if err != nil {
		return fmt.Errorf("error loading source conntector on binding %s, %w", b.name, err)
	}
	b.log.Infof("binding: %s source initialized successfully", b.name)
	b.log.Infof("binding: %s, initialized successfully", b.name)
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
	b.log.Infof("binding: %s, started successfully", b.name)
	return nil
}

func (b *Binder) Stop() error {
	err := b.source.Stop()
	if err != nil {
		return err
	}

	err = b.target.Stop()
	if err != nil {
		return err
	}
	b.log.Infof("binding: %s, stopped successfully", b.name)

	return nil
}
