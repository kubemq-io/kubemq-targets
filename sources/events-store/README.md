# Kubemq Events Store Source

Kubemq Events Store source provides an events-store subscriber for processing source events-store.

## Prerequisites
The following are required to run events source connector:

- kubemq cluster
- kubemq-targets deployment


## Configuration

Events Store source connector configuration properties:

| Properties Key             | Required | Description                           | Example            |
|:---------------------------|:---------|:--------------------------------------|:-------------------|
| address                    | yes      | kubemq server address (gRPC interface) | kubemq-cluster:50000 |
| client_id                  | no       | set client id                         | "client_id"        |
| auth_token                 | no       | set authentication token              | jwt token          |
| channel                    | yes      | set channel to subscribe              |                    |
| group                      | no       | set subscriber group                  |                    |
| response_channel             | no       | set send target response to channel   | "response.channel" |
| auto_reconnect             | no       | set auto reconnect on lost connection | "false", "true"    |
| reconnect_interval_seconds | no       | set reconnection seconds              | "5"                |
| max_reconnects             | no       | set how many time to reconnect        | "0"                |






Example:

```yaml
bindings:
  - name: kubemq-events-store-elastic-search
    source:
      kind: kubemq.events-store
      name: kubemq-events
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-events-store-elastic-search-connector"
        auth_token: ""
        channel: "events-store.elastic-search"
        group:   ""
        response_channel: "events-store.response.elastic"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: stores.elastic-search
      name: target-elastic-search
      properties:
        urls: "http://localhost:9200"
        username: "admin"
        password: "password"
        sniff: "false"
```
