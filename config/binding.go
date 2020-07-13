package config

import "fmt"

type BindingConfig struct {
	Name   string   `json:"name"`
	Source Metadata `json:"source"`
	Target Metadata `json:"target"`
}

func (b BindingConfig) Validate() error {
	if b.Name == "" {
		return fmt.Errorf("binding must have name")
	}
	if err := b.Source.Validate(); err != nil {
		return fmt.Errorf("binding source error, %w", err)
	}
	if err := b.Target.Validate(); err != nil {
		return fmt.Errorf("binding target error, %w", err)
	}
	return nil
}
