# Kubemq MSK Source Connector

Kubemq MSK target connector allows services using kubemq server to store messages on MSK specific topics.

## Prerequisites
The following are required to run the redis target connector:

- kubemq cluster
- MSK set up in aws
- kubemq-targets deployment

## Configuration

MSK source connector configuration properties:

| Properties Key | Required | Description                                | Example          |
|:---------------|:---------|:-------------------------------------------|:-----------------|
| brokers        | yes      | MSK brokers connection, comma separated  | "localhost:9092" |
| topic          | yes      | MSK stored topic                         | "TestTopic"      |
| consumerGroup  | yes      | MSK consumer group name                  | "Group1          |
| saslUsername   | no       | SASL based authentication with broker      | "user"           |
| saslPassword   | no       | SASL based authentication with broker      | "pass"           |

Example:

```yaml
bindings:
  - name: kubemq-query-msk
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-msk-connector"
        auth_token: ""
        channel: "query.MSK"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.msk
      name: aws-msk
      properties:
        brokers: "localhost:9092"
        topic: "TestTopic"
```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description                             | Possible values                         |
|:-------------|:---------|:----------------------------------------|:----------------------------------------|
| key          | yes      | MSK message key base64                | "a2V5"                                  |
| headers      | no       | MSK message headers Key Value base64 | `[{"Key": "ZG9n","Value": "bWV0YTE="}]` |


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
