package files

import (
	"fmt"
	"github.com/Azure/azure-storage-file-go/azfile"
	"github.com/kubemq-hub/kubemq-targets/config"
	"time"
)

const (
	defaultPolicy        = "retry_policy_exponential"
	defaultMaxTries      = 1
	defaultTryTimeout    = 10000
	defaultRetryDelay    = 600
	defaultMaxRetryDelay = 1800
)

var policyMap = map[string]string{
	"retry_policy_exponential": "retry_policy_exponential",
	"retry_policy_fixed":       "retry_policy_fixed",
}

type options struct {
	storageAccessKey string
	storageAccount   string

	policy        azfile.RetryPolicy
	maxTries      int32
	tryTimeout    time.Duration
	retryDelay    time.Duration
	maxRetryDelay time.Duration
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.storageAccessKey, err = cfg.MustParseString("storage_access_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing storage_access_key , %w", err)
	}
	o.storageAccount, err = cfg.MustParseString("storage_account")
	if err != nil {
		return options{}, fmt.Errorf("error parsing storage_account , %w", err)
	}

	var policy string
	policy, err = cfg.ParseStringMap("policy", policyMap)
	if err != nil {
		policy = defaultPolicy
	}
	if policy == "retry_policy_fixed" {
		o.policy = azfile.RetryPolicyFixed
	} else if policy == "retry_policy_exponential" {
		o.policy = azfile.RetryPolicyExponential
	}else{
		o.policy = azfile.RetryPolicyExponential
	}
	o.maxTries = int32(cfg.ParseInt("max_tries", defaultMaxTries))
	o.tryTimeout = cfg.ParseTimeDuration("try_timeout", defaultTryTimeout)
	o.retryDelay = cfg.ParseTimeDuration("retry_delay", defaultRetryDelay)
	o.maxRetryDelay = cfg.ParseTimeDuration ("max_retry_delay", defaultMaxRetryDelay)
	return o, nil
}
