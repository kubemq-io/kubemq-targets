package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"sync"
)

type Service struct {
	sync.Mutex
	bindings map[string]*Binder
	log      *logger.Logger
}

func New() *Service {
	return &Service{
		Mutex:    sync.Mutex{},
		bindings: make(map[string]*Binder),
		log:      logger.NewLogger("binding-service"),
	}
}
func (s *Service) Start(ctx context.Context, cfg *config.Config) error {
	wg := sync.WaitGroup{}
	wg.Add(len(cfg.Bindings))
	for _, bindingCfg := range cfg.Bindings {
		go func(cfg config.BindingConfig) {
			defer wg.Done()
			err := s.Add(ctx, cfg)
			if err != nil {
				s.log.Errorf("failed to initialized binding, %s", err.Error())
				return
			}
		}(bindingCfg)

	}
	wg.Wait()
	if len(s.bindings) == 0 {
		return fmt.Errorf("no valid bindings started")
	}
	return nil
}

func (s *Service) Stop() {
	for _, binder := range s.bindings {
		err := s.Remove(binder.name)
		if err != nil {
			s.log.Error(err)
		}
	}
}
func (s *Service) Add(ctx context.Context, cfg config.BindingConfig) error {
	s.Lock()
	defer s.Unlock()
	_, ok := s.bindings[cfg.Name]
	if ok {
		return fmt.Errorf("duplicate binding name")
	}
	binder := NewBinder()
	err := binder.Init(ctx, cfg)
	if err != nil {
		return err
	}
	err = binder.Start(ctx)
	if err != nil {
		return err
	}
	s.bindings[cfg.Name] = binder
	return nil
}

func (s *Service) Remove(name string) error {
	s.Lock()
	defer s.Unlock()
	binder, ok := s.bindings[name]
	if !ok {
		return fmt.Errorf("binding %s no found", name)
	}
	err := binder.Stop()
	if err != nil {
		return err
	}
	delete(s.bindings, name)
	return nil
}
