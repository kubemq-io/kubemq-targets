# Kubemq Queue Source

Kubemq Queue source provides rpc queue subscriber for processing targets commands.

## Prerequisites
The following are required to run queue source connector:

- kubemq cluster
- kubemq-targets deployment


## Configuration

Queue source connector configuration properties:

| Properties Key | Required | Description                                            | Example     |
|:---------------|:---------|:-------------------------------------------------------|:------------|
| host           | yes      | kubemq server host address                             | "localhost  |
| port           | yes      | kubemq server port number                              | "50000"     |
| client_id      | no       | set client id                                          | "client_id" |
| auth_token     | no       | set authentication token                               | jwt token   |
| channel        | yes      | set channel to subscribe                               |             |
| concurrency    | yes      | set parallel subscribers count                         | "10"        |
| response_channel             | no       | set send target response to channel   | "response.channel" |
| batch_size     | yes      | set how many messages to pull from queue               | "1"         |
| wait_timeout   | yes      | set how long to wait for messages to arrive in seconds | "60"        |


Example:

```yaml
bindings:
  - name: kubemq-queue-elastic-search
    source:
      kind: queue
      name: kubemq-queue
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-queue-elastic-search-connector"
        auth_token: ""
        channel: "queue.elastic-search"
        response_channel: "queue.response.elastic"
        concurrency: "1"
        batch_size: "1"
        wait_timeout: "60"
    target:
      kind: target.stores.elastic-search
      name: target-elastic-search
      properties:
        urls: "http://localhost:9200"
        username: "admin"
        password: "password"
        sniff: "false"
```
