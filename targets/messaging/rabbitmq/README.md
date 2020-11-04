# Kubemq RabbitMQ Target Connector

Kubemq rabbitmq target connector allows services using kubemq server to access rabbitmq messaging services.

## Prerequisites
The following are required to run the rabbitmq target connector:

- kubemq cluster
- rabbitmq server
- kubemq-targets deployment

## Configuration

RabbitMQ target connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| url                             | yes      | rabbitmq connection string address          | "amqp://rabbitmq:rabbitmq@localhost:5672/" |

Example:

```yaml
bindings:
  - name: kubemq-query-rabbitmq
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-rabbitmq-connector"
        auth_token: ""
        channel: "query.rabbitmq"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: messaging.rabbitmq
      name: target-rabbitmq
      properties:
        url: "amqp://rabbitmq:rabbitmq@localhost:5672/"
```

## Usage

### Request

Request metadata setting:

| Metadata Key   | Required | Description         | Possible values |
|:---------------|:---------|:--------------------|:----------------|
| queue          | yes      | set queue name | "queue"         |
| exchange       | no       | set exchange name | "exchange"         |
| mandatory      | no       | set mandatory | "true","false"         |
| immediate      | no       | set immediate | "true","false"         |
| delivery_mode  | no       | set delivery mode | "1","2"         |
| priority       | no       | set priority | "0"-"9"         |
| correlation_id | no       | set correlation id | "some id"         |
| reply_to       | no       | set set reply to | ""         |
| expiry_seconds | no       | set message expiry in seconds| "3600"         |


Query request data setting:

| Data Key | Required | Description  | Possible values    |
|:---------|:---------|:-------------|:-------------------|
| data     | yes      | data to publish | base64 bytes array |

Example:


```json
{
  "metadata": {
    "queue": "queue",
    "exchange": "",
    "confirm": "true",
    "mandatory": "false",
    "immediate": "false",
    "delivery_mode": "1",
    "priority": "0",
    "correlation_id": "",
    "reply_to": "",
    "expiry_seconds": "3600"
    
  },
  "data": "U0VMRUNUIGlkLHRpdGxlLGNvbnRlbnQgRlJPTSBwb3N0Ow=="
}
```
