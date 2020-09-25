package blob

import (
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/kubemq-hub/kubemq-targets/config"
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
	"retry_policy_exponential": "retry_policy_exponential",
	"retry_policy_fixed":       "retry_policy_fixed",
}

type options struct {
	storageAccessKey string
	storageAccount   string

	policy        azblob.RetryPolicy
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
		o.policy = azblob.RetryPolicyFixed
	} else if policy == "retry_policy_exponential" {
		o.policy = azblob.RetryPolicyExponential
	} else {
		o.policy = azblob.RetryPolicyExponential
	}
	o.maxTries = int32(cfg.Properties.ParseInt("max_tries", defaultMaxTries))
	o.tryTimeout = cfg.Properties.ParseTimeDuration("try_timeout", defaultTryTimeout)
	o.retryDelay = cfg.Properties.ParseTimeDuration("retry_delay", defaultRetryDelay)
	o.maxRetryDelay = cfg.Properties.ParseTimeDuration("max_retry_delay", defaultMaxRetryDelay)

	return o, nil
}
