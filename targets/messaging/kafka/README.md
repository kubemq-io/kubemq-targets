# Kubemq Kafka Source Connector

Kubemq kafka source connector allows services using kubemq server to access redis server. TODO

## Prerequisites
The following are required to run the redis target connector:

- kubemq cluster
- kafka TODO version
- kubemq-target-connectors deployment

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
  - name: kafka-store-kubemq
    source:
      kind: source.kafka
      name: kafka-stream
      properties:
     	brokers: "localhost:9092,localhost:9093",
		topics: "TestTopic",
		consumerGroup: "cg",
    target:
      kind: target.kubemq.event-store
      name: target-kubemq-event-store
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-redis-connector"
        auth_token: ""
        channel: "store.kafka"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
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
