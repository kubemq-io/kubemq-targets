# Kubemq Kafka Source Connector

Kubemq kafka target connector allows services using kubemq server to store messages on kafka specific topics.

## Prerequisites
The following are required to run the redis target connector:

- kubemq cluster
- kafka server
- kubemq-targets deployment

## Configuration

Kafka source connector configuration properties:

| Properties Key | Required | Description                                | Example          |
|:---------------|:---------|:-------------------------------------------|:-----------------|
| brokers        | yes      | kafka brokers connection, comma separated  | "localhost:9092" |
| topic          | yes      | kafka stored topic                         | "TestTopic"      |
| sasl_username  | no       | SASL based authentication with broker      | "user"           |
| sasl_password  | no       | SASL based authentication with broker      | "pass"           |

Example:

```yaml
bindings:
  - name: kubemq-query-kafka
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-kafka-connector"
        auth_token: ""
        channel: "query.kafka"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: messaging.kafka
      name: kafka-stream
      properties:
        brokers: "localhost:9092"
        topic: "TestTopic"
        sasl_username: "test"
        sasl_password: "pass"
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
