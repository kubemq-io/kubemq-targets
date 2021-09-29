// +build !container

package targets

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/targets/aws/amazonmq"
	"github.com/kubemq-io/kubemq-targets/targets/aws/athena"
	"github.com/kubemq-io/kubemq-targets/targets/aws/cloudwatch/events"
	"github.com/kubemq-io/kubemq-targets/targets/aws/cloudwatch/logs"
	"github.com/kubemq-io/kubemq-targets/targets/aws/cloudwatch/metrics"
	"github.com/kubemq-io/kubemq-targets/targets/aws/dynamodb"
	"github.com/kubemq-io/kubemq-targets/targets/aws/elasticsearch"
	"github.com/kubemq-io/kubemq-targets/targets/aws/keyspaces"
	"github.com/kubemq-io/kubemq-targets/targets/aws/kinesis"
	"github.com/kubemq-io/kubemq-targets/targets/aws/lambda"
	"github.com/kubemq-io/kubemq-targets/targets/aws/msk"
	"github.com/kubemq-io/kubemq-targets/targets/aws/s3"
	"github.com/kubemq-io/kubemq-targets/targets/aws/sns"
	"github.com/kubemq-io/kubemq-targets/targets/azure/eventhubs"
	"github.com/kubemq-io/kubemq-targets/targets/azure/servicebus"
	"github.com/kubemq-io/kubemq-targets/targets/azure/storage/blob"
	"github.com/kubemq-io/kubemq-targets/targets/azure/storage/files"
	"github.com/kubemq-io/kubemq-targets/targets/azure/storage/queue"
	"github.com/kubemq-io/kubemq-targets/targets/azure/stores/azuresql"
	azurmysql "github.com/kubemq-io/kubemq-targets/targets/azure/stores/mysql"
	azurpostgres "github.com/kubemq-io/kubemq-targets/targets/azure/stores/postgres"
	"github.com/kubemq-io/kubemq-targets/targets/cache/hazelcast"
	"github.com/kubemq-io/kubemq-targets/targets/echo"
	"github.com/kubemq-io/kubemq-targets/targets/gcp/firebase"
	"github.com/kubemq-io/kubemq-targets/targets/storage/filesystem"
	"github.com/kubemq-io/kubemq-targets/targets/stores/crate"

	//"github.com/kubemq-io/kubemq-targets/targets/messaging/ibmmq"
	"github.com/kubemq-io/kubemq-targets/targets/messaging/nats"
	"github.com/kubemq-io/kubemq-targets/targets/storage/hdfs"
	"github.com/kubemq-io/kubemq-targets/targets/stores/aerospike"
	"github.com/kubemq-io/kubemq-targets/targets/stores/cockroachdb"
	"github.com/kubemq-io/kubemq-targets/targets/stores/consulkv"
	"github.com/kubemq-io/kubemq-targets/targets/stores/elastic"
	"github.com/kubemq-io/kubemq-targets/targets/stores/percona"
	"github.com/kubemq-io/kubemq-targets/targets/stores/rethinkdb"
	"github.com/kubemq-io/kubemq-targets/targets/stores/singlestore"

	"github.com/kubemq-io/kubemq-targets/config"
	awsmariadb "github.com/kubemq-io/kubemq-targets/targets/aws/rds/mariadb"
	awsmssql "github.com/kubemq-io/kubemq-targets/targets/aws/rds/mssql"
	awsmysql "github.com/kubemq-io/kubemq-targets/targets/aws/rds/mysql"
	awspostgres "github.com/kubemq-io/kubemq-targets/targets/aws/rds/postgres"
	"github.com/kubemq-io/kubemq-targets/targets/aws/rds/redshift"
	redshiftsvc "github.com/kubemq-io/kubemq-targets/targets/aws/redshift"
	"github.com/kubemq-io/kubemq-targets/targets/aws/sqs"
	"github.com/kubemq-io/kubemq-targets/targets/cache/memcached"
	"github.com/kubemq-io/kubemq-targets/targets/cache/redis"
	"github.com/kubemq-io/kubemq-targets/targets/gcp/bigquery"
	"github.com/kubemq-io/kubemq-targets/targets/gcp/bigtable"
	"github.com/kubemq-io/kubemq-targets/targets/gcp/cloudfunctions"
	"github.com/kubemq-io/kubemq-targets/targets/gcp/firestore"
	gcpmemcached "github.com/kubemq-io/kubemq-targets/targets/gcp/memorystore/memcached"
	gcpredis "github.com/kubemq-io/kubemq-targets/targets/gcp/memorystore/redis"
	"github.com/kubemq-io/kubemq-targets/targets/gcp/pubsub"
	"github.com/kubemq-io/kubemq-targets/targets/gcp/spanner"
	gcpmysql "github.com/kubemq-io/kubemq-targets/targets/gcp/sql/mysql"
	gcppostgres "github.com/kubemq-io/kubemq-targets/targets/gcp/sql/postgres"
	"github.com/kubemq-io/kubemq-targets/targets/gcp/storage"
	"github.com/kubemq-io/kubemq-targets/targets/http"
	"github.com/kubemq-io/kubemq-targets/targets/messaging/activemq"
	"github.com/kubemq-io/kubemq-targets/targets/messaging/kafka"
	"github.com/kubemq-io/kubemq-targets/targets/messaging/mqtt"
	"github.com/kubemq-io/kubemq-targets/targets/messaging/rabbitmq"
	"github.com/kubemq-io/kubemq-targets/targets/serverless/openfaas"
	"github.com/kubemq-io/kubemq-targets/targets/storage/minio"
	"github.com/kubemq-io/kubemq-targets/targets/stores/cassandra"
	"github.com/kubemq-io/kubemq-targets/targets/stores/couchbase"
	"github.com/kubemq-io/kubemq-targets/targets/stores/mongodb"
	"github.com/kubemq-io/kubemq-targets/targets/stores/mssql"
	"github.com/kubemq-io/kubemq-targets/targets/stores/mysql"
	"github.com/kubemq-io/kubemq-targets/targets/stores/postgres"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Target interface {
	Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error
	Do(ctx context.Context, request *types.Request) (*types.Response, error)
	Stop() error
	Connector() *common.Connector
}

func Init(ctx context.Context, cfg config.Spec, log *logger.Logger) (Target, error) {

	switch cfg.Kind {
	case "storage.filesystem":
		target := filesystem.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "echo":
		target := echo.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.sqs":
		target := sqs.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.sns":
		target := sns.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.s3":
		target := s3.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.lambda":
		target := lambda.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.dynamodb":
		target := dynamodb.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.athena":
		target := athena.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.kinesis":
		target := kinesis.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.elasticsearch":
		target := elasticsearch.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.cloudwatch.logs":
		target := logs.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.cloudwatch.events":
		target := events.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.cloudwatch.metrics":
		target := metrics.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.rds.mysql":
		target := awsmysql.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.rds.postgres":
		target := awspostgres.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.rds.mariadb":
		target := awsmariadb.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.rds.mssql":
		target := awsmssql.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.rds.redshift":
		target := redshift.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.redshift.service":
		target := redshiftsvc.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.keyspaces":
		target := keyspaces.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.msk":
		target := msk.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "aws.amazonmq":
		target := amazonmq.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "cache.redis":
		target := redis.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "cache.memcached":
		target := memcached.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "cache.hazelcast":
		target := hazelcast.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.cache.memcached":
		target := gcpmemcached.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.cache.redis":
		target := gcpredis.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.bigquery":
		target := bigquery.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.bigtable":
		target := bigtable.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.cloudfunctions":
		target := cloudfunctions.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.firestore":
		target := firestore.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.firebase":
		target := firebase.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.stores.postgres":
		target := gcppostgres.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.stores.mysql":
		target := gcpmysql.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.pubsub":
		target := pubsub.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.spanner":
		target := spanner.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "gcp.storage":
		target := storage.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "http":
		target := http.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "messaging.activemq":
		target := activemq.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "messaging.kafka":
		target := kafka.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "messaging.mqtt":
		target := mqtt.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "messaging.rabbitmq":
		target := rabbitmq.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	//case "messaging.ibmmq":
	//	target := ibmmq.New()
	//	if err := target.Init(ctx, cfg); err != nil {
	//		return nil, err
	//	}
	//	return target, nil
	case "messaging.nats":
		target := nats.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.cassandra":
		target := cassandra.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.couchbase":
		target := couchbase.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.mongodb":
		target := mongodb.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.mssql":
		target := mssql.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.mysql":
		target := mysql.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.elasticsearch":
		target := elastic.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.postgres":
		target := postgres.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.cockroachdb":
		target := cockroachdb.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.percona":
		target := percona.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.crate":
		target := crate.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.aerospike":
		target := aerospike.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.rethinkdb":
		target := rethinkdb.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.singlestore":
		target := singlestore.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "stores.consulkv":
		target := consulkv.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "serverless.openfaas":
		target := openfaas.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "storage.minio":
		target := minio.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "storage.hdfs":
		target := hdfs.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "azure.storage.blob":
		target := blob.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "azure.storage.queue":
		target := queue.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "azure.storage.files":
		target := files.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "azure.eventhubs":
		target := eventhubs.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "azure.servicebus":
		target := servicebus.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "azure.stores.azuresql":
		target := azuresql.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "azure.stores.postgres":
		target := azurpostgres.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "azure.stores.mysql":
		target := azurmysql.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil

	default:
		return nil, fmt.Errorf("invalid kind %s for target %s", cfg.Kind, cfg.Name)
	}

}

func Connectors() common.Connectors {
	return []*common.Connector{
		echo.Connector(),
		// cache
		redis.Connector(),
		memcached.Connector(),

		//stores
		postgres.Connector(),
		crate.Connector(),
		mysql.Connector(),
		mssql.Connector(),
		mongodb.Connector(),
		elastic.Connector(),
		cassandra.Connector(),
		couchbase.Connector(),
		cockroachdb.Connector(),
		percona.Connector(),
		aerospike.Connector(),
		rethinkdb.Connector(),
		singlestore.Connector(),
		consulkv.Connector(),

		// http
		http.Connector(),

		// messaging
		mqtt.Connector(),
		rabbitmq.Connector(),
		kafka.Connector(),
		activemq.Connector(),
		//ibmmq.Connector(),
		nats.Connector(),
		hazelcast.Connector(),

		//storage
		minio.Connector(),
		hdfs.Connector(),
		filesystem.Connector(),

		// serverless
		openfaas.Connector(),

		// aws
		sqs.Connector(),
		sns.Connector(),
		s3.Connector(),
		amazonmq.Connector(),

		awspostgres.Connector(),
		awsmysql.Connector(),
		awsmariadb.Connector(),
		awsmssql.Connector(),
		dynamodb.Connector(),
		redshift.Connector(),
		redshiftsvc.Connector(),
		athena.Connector(),

		msk.Connector(),
		lambda.Connector(),
		kinesis.Connector(),
		keyspaces.Connector(),
		elasticsearch.Connector(),

		events.Connector(),
		logs.Connector(),
		metrics.Connector(),

		// gcp

		pubsub.Connector(),
		gcpredis.Connector(),
		gcpmemcached.Connector(),
		gcppostgres.Connector(),
		gcpmysql.Connector(),
		spanner.Connector(),
		bigtable.Connector(),
		bigquery.Connector(),
		cloudfunctions.Connector(),
		firebase.Connector(),
		firestore.Connector(),
		storage.Connector(),

		// azure

		azurpostgres.Connector(),
		azurmysql.Connector(),
		azuresql.Connector(),
		queue.Connector(),
		files.Connector(),
		blob.Connector(),
		servicebus.Connector(),
		eventhubs.Connector(),
	}
}
