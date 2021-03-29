module github.com/kubemq-hub/kubemq-targets

go 1.14

require (
	cloud.google.com/go v0.76.0
	cloud.google.com/go/bigquery v1.15.0
	cloud.google.com/go/bigtable v1.7.1
	cloud.google.com/go/firestore v1.4.0
	cloud.google.com/go/pubsub v1.10.1
	cloud.google.com/go/spanner v1.13.0
	cloud.google.com/go/storage v1.13.0
	firebase.google.com/go/v4 v4.2.0
	github.com/Azure/azure-event-hubs-go/v3 v3.3.0
	github.com/Azure/azure-pipeline-go v0.2.3
	github.com/Azure/azure-sdk-for-go v46.1.0+incompatible // indirect
	github.com/Azure/azure-service-bus-go v0.10.3
	github.com/Azure/azure-storage-blob-go v0.10.0
	github.com/Azure/azure-storage-file-go v0.8.0
	github.com/Azure/azure-storage-queue-go v0.0.0-20191125232315-636801874cdd
	github.com/Azure/go-autorest/autorest v0.11.6 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.0 // indirect
	github.com/GoogleCloudPlatform/cloudsql-proxy v1.19.1
	github.com/Shopify/sarama v1.27.2
	github.com/aerospike/aerospike-client-go v4.0.0+incompatible
	github.com/apache/thrift v0.13.0 // indirect
	github.com/aws/aws-sdk-go v1.37.6
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/cockroachdb/cockroach-go v2.0.1+incompatible
	github.com/colinmarc/hdfs/v2 v2.2.0
	github.com/couchbase/gocb/v2 v2.2.0
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/eclipse/paho.mqtt.golang v1.3.2
	github.com/fsnotify/fsnotify v1.4.9
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-resty/resty/v2 v2.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-stomp/stomp v2.0.6+incompatible
	github.com/gocql/gocql v0.0.0-20200815110948-5378c8f664e9
	github.com/golang/protobuf v1.4.3
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/hashicorp/consul/api v1.1.0
	github.com/hazelcast/hazelcast-go-client v0.6.0
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/kubemq-hub/builder v0.7.2
	github.com/kubemq-hub/ibmmq-sdk v0.3.8
	github.com/kubemq-io/kubemq-go v1.4.7
	github.com/labstack/echo/v4 v4.1.17
	github.com/lib/pq v1.9.0
	github.com/minio/minio-go/v7 v7.0.8
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/nats-io/nats-server/v2 v2.1.9 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/olivere/elastic/v7 v7.0.22
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.mongodb.org/mongo-driver v1.4.1
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/oauth2 v0.0.0-20210113205817-d3ed898aa8a3
	google.golang.org/api v0.38.0
	google.golang.org/genproto v0.0.0-20210203152818-3206188e46ba
	google.golang.org/grpc v1.35.0
	gopkg.in/rethinkdb/rethinkdb-go.v6 v6.2.1
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/Azure/azure-service-bus-go => github.com/Azure/azure-service-bus-go v0.10.3

//replace github.com/kubemq-hub/builder => ../builder
