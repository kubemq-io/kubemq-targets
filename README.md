# kubemq-targets

## Concept

![concept](.github/assets/concept.jpeg)

## Supported Targets

### Standalone Services

| Category   | Target                                                              | Kind                         | Configuration                                  |
|:-----------|:--------------------------------------------------------------------|:-----------------------------|:-----------------------------------------------|
| Cache      |                                                                     |                              |                                                |
|            | [Redis](https://redis.io/)                                          | target.cache.redis           | [Usage](targets/cache/redis/README.md)         |
|            | [Memcached](https://memcached.org/)                                 | target.cache.memcached       | [Usage](targets/cache/memcached/README.md)     |
| Stores/db  |                                                                     |                              |                                                |
|            | [Postgres](https://www.postgresql.org/)                             | target.stores.postgres       | [Usage](targets/stores/postgres/README.md)     |
|            | [Mysql](https://www.mysql.com/)                                     | target.stores.mysql          | [Usage](targets/stores/mysql/README.md)        |
|            | [MSSql](https://www.microsoft.com/en-us/sql-server/sql-server-2019) | target.stores.mssql          | [Usage](targets/stores/mssql/README.md)        |
|            | [MongoDB](https://www.mongodb.com/)                                 | target.stores.mongodb        | [Usage](targets/stores/mongodb/README.md)      |
|            | [Elastic Search](https://www.elastic.co/)                           | target.stores.elastic-search | [Usage](targets/stores/elastic/README.md)      |
|            | [Cassandra](https://cassandra.apache.org/)                          | target.stores.cassandra      | [Usage](targets/stores/cassandra/README.md)    |
|            | [Couchbase](https://www.couchbase.com/)                             | target.stores.couchbase      | [Usage](targets/stores/couchbase/README.md)    |
| Messaging  |                                                                     |                              |                                                |
|            | [Kafka](https://kafka.apache.org/)                                  | target.messaging.kafka       | [Usage](targets/messaging/kafka/README.md)     |
|            | [RabbitMQ](https://www.rabbitmq.com/)                               | target.messaging.rabbitmq    | [Usage](targets/messaging/rabbitmq/README.md)  |
|            | [MQTT](http://mqtt.org/)                                            | target.messaging.mqtt        | [Usage](targets/messaging/mqtt/README.md)      |
|            | [ActiveMQ](http://activemq.apache.org/)                             | target.messaging.activemq    | [Usage](targets/messaging/postgres/README.md)  |
| Storage    |                                                                     |                              |                                                |
|            | [Minio/S3](https://min.io/)                                         | target.storage.minio         | [Usage](targets/storage/minio/README.md)       |
| Serverless |                                                                     |                              |                                                |
|            | [OpenFaas](https://www.openfaas.com/)                               | target.serverless.openfaas   | [Usage](targets/serverless/openfass/README.md) |
| Http       |                                                                     |                              |                                                |
|            | Http                                                                | target.http                  | [Usage](targets/http/README.md)                |



### Google Cloud Platform (GCP)

| Category   | Target                                                              | Kind                       | Configuration                                        |
|:-----------|:--------------------------------------------------------------------|:---------------------------|:-----------------------------------------------------|
| Cache      |                                                                     |                            |                                                      |
|            | [Redis](https://cloud.google.com/memorystore)                       | target.gcp.cache.redis     | [Usage](targets/gcp/memorystore/redis/README.md)     |
|            | [Memcached](https://cloud.google.com/memorystore)                   | target.gcp.cache.memcached | [Usage](targets/gcp/memorystore/memcached/README.md) |
| Stores/db  |                                                                     |                            |                                                      |
|            | [Postgres](https://cloud.google.com/sql)                            | target.gcp.stores.postgres | [Usage](targets/gcp/sql/postgres/README.md)          |
|            | [Mysql](https://cloud.google.com/sql)                               | target.gcp.stores.mysql    | [Usage](targets/gcp/sql/mysql/README.md)             |
|            | [BigQuery](https://cloud.google.com/bigquery)                       | target.gcp.bigquery        | [Usage](targets/gcp/bigquery/README.md)              |
|            | [BigTable](https://cloud.google.com/bigtable)                       | target.gcp.bigtable        | [Usage](targets/gcp/bigtable/README.md)              |
|            | [Firestore](https://cloud.google.com/firestore)                     | target.gcp.firestore       | [Usage](targets/gcp/firestore/README.md)             |
|            | [Spanner](https://cloud.google.com/spanner)                         | target.gcp.spanner         | [Usage](targets/gcp/spanner/README.md)               |
|            | [Firebase](https://firebase.google.com/products/realtime-database/) | target.gcp.firebase        | [Usage](targets/gcp/firebase/README.md)              |
| Messaging  |                                                                     |                            |                                                      |
|            | [Pub/Sub](https://cloud.google.com/pubsub)                          | target.gcp.pubsub          | [Usage](targets/gcp/pubsub/README.md)                |
| Storage    |                                                                     |                            |                                                      |
|            | [Storage](https://cloud.google.com/storage)                         | target.gcp.storage         | [Usage](targets/gcp/storage/README.md)               |
| Serverless |                                                                     |                            |                                                      |
|            | [Functions](https://cloud.google.com/functions)                     | target.gcp.cloudfunctions  | [Usage](targets/gcp/cloudfunctions/README.md)        |







## Installation


## Configuration


