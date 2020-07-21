package targets

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/targets/gcp/storage"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/targets/aws/sqs"
	"github.com/kubemq-hub/kubemq-targets/targets/cache/memcached"
	"github.com/kubemq-hub/kubemq-targets/targets/cache/redis"
	"github.com/kubemq-hub/kubemq-targets/targets/gcp/bigquery"
	"github.com/kubemq-hub/kubemq-targets/targets/gcp/bigtable"
	"github.com/kubemq-hub/kubemq-targets/targets/gcp/firestore"
	gcpmemcached "github.com/kubemq-hub/kubemq-targets/targets/gcp/memorystore/memcached"
	gcpredis "github.com/kubemq-hub/kubemq-targets/targets/gcp/memorystore/redis"
	"github.com/kubemq-hub/kubemq-targets/targets/gcp/pubsub"
	"github.com/kubemq-hub/kubemq-targets/targets/gcp/spanner"
	"github.com/kubemq-hub/kubemq-targets/targets/http"
	"github.com/kubemq-hub/kubemq-targets/targets/kubemq/command"
	"github.com/kubemq-hub/kubemq-targets/targets/kubemq/events"
	events_store "github.com/kubemq-hub/kubemq-targets/targets/kubemq/events-store"
	"github.com/kubemq-hub/kubemq-targets/targets/kubemq/query"
	"github.com/kubemq-hub/kubemq-targets/targets/kubemq/queue"
	"github.com/kubemq-hub/kubemq-targets/targets/logs/elastic"
	"github.com/kubemq-hub/kubemq-targets/targets/messaging/activemq"
	"github.com/kubemq-hub/kubemq-targets/targets/messaging/kafka"
	"github.com/kubemq-hub/kubemq-targets/targets/messaging/mqtt"
	"github.com/kubemq-hub/kubemq-targets/targets/messaging/rabbitmq"
	"github.com/kubemq-hub/kubemq-targets/targets/serverless/openfass"
	"github.com/kubemq-hub/kubemq-targets/targets/storage/minio"
	"github.com/kubemq-hub/kubemq-targets/targets/stores/cassandra"
	"github.com/kubemq-hub/kubemq-targets/targets/stores/couchbase"
	"github.com/kubemq-hub/kubemq-targets/targets/stores/mongodb"
	"github.com/kubemq-hub/kubemq-targets/targets/stores/mssql"
	"github.com/kubemq-hub/kubemq-targets/targets/stores/mysql"
	"github.com/kubemq-hub/kubemq-targets/targets/stores/postgres"
	"github.com/kubemq-hub/kubemq-targets/types"
)

var (
	errTargetNotImplemented = fmt.Errorf("target not implemented")
)

type Target interface {
	Init(ctx context.Context, cfg config.Spec) error
	Do(ctx context.Context, request *types.Request) (*types.Response, error)
	Name() string
}

func Init(ctx context.Context, cfg config.Spec) (Target, error) {

	switch cfg.Kind {
	case "target.aws.sqs":
		target := sqs.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.cache.redis":
		target := redis.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.cache.memcached":
		target := memcached.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.cache.memcached":
		target := gcpmemcached.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.cache.redis":
		target := gcpredis.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.bigquery":
		target := bigquery.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.bigtable":
		target := bigtable.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.firestore":
		target := firestore.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.stores.postgres":
		target := firestore.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.pubsub":
		target := pubsub.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.spanner":
		target := spanner.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.storage":
		target := storage.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.http":
		target := http.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.command":
		target := command.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.query":
		target := query.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.events":
		target := events.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.events-store":
		target := events_store.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.queue":
		target := queue.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.logs.elastic":
		target := elastic.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.messaging.activemq":
		target := activemq.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.messaging.kafka":
		target := kafka.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.messaging.mqtt":
		target := mqtt.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.messaging.rabbitmq":
		target := rabbitmq.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.cassandra":
		target := cassandra.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.couchbase":
		target := couchbase.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.mongodb":
		target := mongodb.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.mssql":
		target := mssql.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.mysql":
		target := mysql.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.postgres":
		target := postgres.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.serverless.openfaas":
		target := openfass.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.storage.minio":
		target := minio.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	default:
		return nil, fmt.Errorf("invalid kind %s for target %s", cfg.Kind, cfg.Name)
	}

}
