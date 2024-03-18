# Kubemq AMQP Target Connector

Kubemq amqp target connector allows services using kubemq server to access amqp messaging services.

## Prerequisites
The following are required to run the amqp target connector:

- kubemq cluster
- amqp server
- kubemq-targets deployment

## Configuration

AMQP target connector configuration properties:

| Properties Key      | Required | Description                        | Example                            |
|:--------------------|:---------|:-----------------------------------|:-----------------------------------|
| url                 | yes      | amqp connection string address | "amqp://amqp:amqp@localhost:5672/" |
| ca_cert            | no       | SSL CA certificate                          | pem certificate value              |
| skip_insecure      | no       | skip insecure certificate | "false"                            |
| username            | no       | amqp username | "amqp"                             |

Example:

```yaml
bindings:
  - name: kubemq-query-amqp
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-amqp-connector"
        auth_token: ""
        channel: "query.amqp"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: messaging.amqp
      name: target-amqp
      properties:
        url: "amqp://localhost:5672/"
        username: "amqp"
        password: "amqp"
```

## Usage

### Request

Request metadata setting:

| Metadata Key | Required | Description                                                         | Possible values |
|:-------------|:---------|:--------------------------------------------------------------------|:----------------|
| address      | yes      | set address destination                                             | "some-address"  |
|durable       | no       | set durable message                                                 | "false"         |
|priority      | no       | set message priority                                                | 0               |
|message_id    | no       | set message id                                                      | "some-id"       |
|to           | no       | identifies the node that is the intended destination of the message | "some-node"     |
|subject      | no       | message subject                                                     | "some-subject"  |
|reply_to     | no       | The address of the node to send replies to                          | "some-reply"    |
|correlation_id | no     | The message-id of the message this message is a reply to            | "some-id"       |
|content_type  | no       | The content-type of the message data                                | "some-type"     |
|expiry_time   | no       | the time in unix time at which the message will expire              | 0               |
|group_id      | no       | The group-id of a message group                                     | "some-group"    |
|reply_to_group_id | no  | The group-id of the group the reply belongs to                      | "some-group"    |
|group_sequence | no    | The sequence number of this message in a group                      | 0               |




Query request data setting:

| Data Key | Required | Description  | Possible values    |
|:---------|:---------|:-------------|:-------------------|
| data     | yes      | data to publish | base64 bytes array |

Example:


```json
{
  "metadata": {
    "address": "some-address",
    "durable": true
  },
  "data": "U0VMRUNUIGlkLHRpdGxlLGNvbnRlbnQgRlJPTSBwb3N0Ow=="
}
```
