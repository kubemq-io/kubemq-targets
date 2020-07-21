# Kubemq Kafka Source Connector

Kubemq kafka target connector allows services using kubemq server to store messages on kafka specific topics.

## Prerequisites
The following are required to run the redis target connector:

- kubemq cluster
- kafka TODO version
- kubemq-targets deployment

## Configuration

Kafka source connector configuration properties:

| Properties Key | Required | Description                                | Example          |
|:---------------|:---------|:-------------------------------------------|:-----------------|
| brokers        | yes      | kafka brokers connection, comma separated  | "localhost:9092" |
| topic          | yes      | kafka stored topic                         | "TestTopic"      |
| consumerGroup  | yes      | kafka consumer group name                  | "Group1          |
| saslUsername   | no       | SASL based authentication with broker      | "user"           |
| saslPassword   | no       | SASL based authentication with broker      | "pass"           |

Example:

```yaml
bindings:
  - name: kubemq-query-kafka
    source:
      kind: source.kubemq.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-kafka-connector"
        auth_token: ""
        channel: "query.kafka"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: targets.messaging.kafka
      name: kafka-stream
      properties:
        brokers: "localhost:9092",
        topic: "TestTopic"
```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description                             | Possible values                         |
|:-------------|:---------|:----------------------------------------|:----------------------------------------|
| key          | yes      | kafka message key base64                | "a2V5"                                  |
| headers      | no       | kafka message headers Key Value base64 | `[{"Key": "ZG9n","Value": "bWV0YTE="}]` |


Example:

```json
{
  "metadata": {
    "key": "a2V5",
    "headers": [{"Key": "ZG9n","Value": "bWV0YTE="}]
  },
  "data": null
}
```
