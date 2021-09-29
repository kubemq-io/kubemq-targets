package queue

import (
	"fmt"
	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/kubemq-io/kubemq-targets/config"
	"time"
)

const (
	defaultPolicy        = "retry_policy_exponential"
	defaultMaxTries      = 1
	defaultTryTimeout    = 1000
	defaultRetryDelay    = 60
	defaultMaxRetryDelay = 180
)

var policyMap = map[string]string{
	"exponential": "retry_policy_exponential",
	"fixed":       "retry_policy_fixed",
}

type options struct {
	storageAccessKey string
	storageAccount   string

	policy        azqueue.RetryPolicy
	maxTries      int32
	tryTimeout    time.Duration
	retryDelay    time.Duration
	maxRetryDelay time.Duration
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.storageAccessKey, err = cfg.Properties.MustParseString("storage_access_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing storage_access_key , %w", err)
	}
	o.storageAccount, err = cfg.Properties.MustParseString("storage_account")
	if err != nil {
		return options{}, fmt.Errorf("error parsing storage_account , %w", err)
	}
	var policy string
	policy, err = cfg.Properties.ParseStringMap("policy", policyMap)
	if err != nil {
		policy = defaultPolicy
	}
	if policy == "retry_policy_fixed" {
		o.policy = azqueue.RetryPolicyFixed
	} else if policy == "retry_policy_exponential" {
		o.policy = azqueue.RetryPolicyExponential
	} else {
		o.policy = azqueue.RetryPolicyExponential
	}
	o.maxTries = int32(cfg.Properties.ParseInt("max_tries", defaultMaxTries))
	o.tryTimeout = cfg.Properties.ParseTimeDuration("try_timeout", defaultTryTimeout)
	o.retryDelay = cfg.Properties.ParseTimeDuration("retry_delay", defaultRetryDelay)
	o.maxRetryDelay = cfg.Properties.ParseTimeDuration("max_retry_delay", defaultMaxRetryDelay)

	return o, nil
}
