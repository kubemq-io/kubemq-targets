package filesystem

import (
	"fmt"
	"strings"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	basePath string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.basePath, err = cfg.Properties.MustParseString("base_path")
	if err != nil {
		return options{}, fmt.Errorf("error parsing base_path, %w", err)
	}
	o.basePath = unixNormalize(o.basePath)
	return o, nil
}

func unixNormalize(in string) string {
	return strings.Replace(in, `\`, "/", -1)
}
