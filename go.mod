module github.com/kubemq-hub/kubemq-targets

go 1.15

require (
	cloud.google.com/go v0.66.0
	cloud.google.com/go/bigquery v1.11.0
	cloud.google.com/go/bigtable v1.6.0
	cloud.google.com/go/firestore v1.3.0
	cloud.google.com/go/pubsub v1.6.2
	cloud.google.com/go/spanner v1.10.0
	cloud.google.com/go/storage v1.12.0
	firebase.google.com/go/v4 v4.0.0
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
	github.com/GoogleCloudPlatform/cloudsql-proxy v1.18.0
	github.com/Shopify/sarama v1.27.0
	github.com/aws/aws-sdk-go v1.34.31
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/couchbase/gocb/v2 v2.1.6
	github.com/denisenkom/go-mssqldb v0.0.0-20200910202707-1e08a3fab204
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-resty/resty/v2 v2.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-stomp/stomp v2.0.6+incompatible
	github.com/gocql/gocql v0.0.0-20200815110948-5378c8f664e9
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/kubemq-hub/builder v0.3.0
	github.com/kubemq-io/kubemq-go v1.4.3
	github.com/labstack/echo/v4 v4.1.17
	github.com/lib/pq v1.8.0
	github.com/minio/minio-go/v7 v7.0.0-20200623213017-e5cd2d129325
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/nats-io/nuid v1.0.1
	github.com/olivere/elastic/v7 v7.0.20
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	go.mongodb.org/mongo-driver v1.4.1
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/api v0.32.0
	google.golang.org/genproto v0.0.0-20200925023002-c2d885f95484
	google.golang.org/grpc v1.32.0
)

replace github.com/Azure/azure-service-bus-go => github.com/Azure/azure-service-bus-go v0.10.3
replace github.com/kubemq-hub/builder => ../builder
