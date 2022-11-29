package binding

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/config"
)

type Status struct {
	Binding          string            `json:"binding"`
	Ready            bool              `json:"ready"`
	SourceType       string            `json:"source_type"`
	SourceConnection string            `json:"source_connection"`
	SourceConfig     map[string]string `json:"source_config"`
	TargetType       string            `json:"target_type"`
	TargetConfig     map[string]string `json:"target_config"`
}

func getSourceConnection(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", properties["address"], properties["channel"])
}

func newStatus(cfg config.BindingConfig) *Status {
	return &Status{
		Binding:          cfg.Name,
		Ready:            false,
		SourceType:       cfg.Source.Kind,
		SourceConnection: getSourceConnection(cfg.Source.Properties),
		SourceConfig:     cfg.Source.Properties,
		TargetType:       cfg.Target.Kind,
		TargetConfig:     cfg.Target.Properties,
	}
}
