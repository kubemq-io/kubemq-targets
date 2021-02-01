# Kubemq Queue Stream Source

Kubemq Queue stream source provides a queue subscriber for processing messages with ack/reject way

## Prerequisites
The following are required to run queue source connector:

- kubemq cluster
- kubemq-targets deployment


## Configuration

Queue source connector configuration properties:

| Properties Key | Required | Description                                            | Example     |
|:---------------|:---------|:-------------------------------------------------------|:------------|
| address                    | yes      | kubemq server address (gRPC interface) | kubemq-cluster:50000 |
| client_id      | no       | set client id                                          | "client_id" |
| auth_token     | no       | set authentication token                               | jwt token   |
| channel        | yes      | set channel to subscribe                               |             |
| sources        | no      | set how many concurrent sources to subscribe                               |    1        |
| response_channel             | no       | set send target response to channel   | "response.channel" |
| visibility_timeout_seconds     | no      | set how long to hold the message for processing before reject | "3600"         |
| wait_timeout   | no      | set how long to wait for messages to arrive in seconds | "3600"        |


Example:

```yaml
bindings:
  - name: kubemq-queue-elastic-search
    source:
      kind: kubemq.queue-stream
      name: kubemq-queue
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-queue-elastic-search-connector"
        auth_token: ""
        channel: "queue.elastic-search"
        sources: "1"
        response_channel: "queue.response.elastic"
        visibility_timeout_seconds: "3600"
        wait_timeout: "3600"
    target:
      kind: stores.elastic-search
      name: target-elastic-search
      properties:
        urls: "http://localhost:9200"
        username: "admin"
        password: "password"
        sniff: "false"
```
