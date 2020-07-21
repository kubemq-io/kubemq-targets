# Kubemq Events Store Source

Kubemq Events Store source provides rpc events subscriber for processing targets commands.

## Prerequisites
The following are required to run events source connector:

- kubemq cluster
- kubemq-targets deployment


## Configuration

Events Store target connector configuration properties:

| Properties Key             | Required | Description                           | Example            |
|:---------------------------|:---------|:--------------------------------------|:-------------------|
| host                       | yes      | kubemq server host address            | "localhost         |
| port                       | yes      | kubemq server port number             | "50000"            |
| client_id                  | no       | set client id                         | "client_id"        |
| auth_token                 | no       | set authentication token              | jwt token          |
| channel                    | yes      | set channel to subscribe              |                    |
| group                      | no       | set subscriber group                  |                    |
| concurrency                | yes      | set parallel subscribers count        | "10"               |
| response_channel             | no       | set send target response to channel   | "response.channel" |
| auto_reconnect             | no       | set auto reconnect on lost connection | "false", "true"    |
| reconnect_interval_seconds | no       | set reconnection seconds              | "5"                |
| max_reconnects             | no       | set how many time to reconnect        | "0"                |






Example:

```yaml
bindings:
  - name: kubemq-events-store-elastic-search
    source:
      kind: source.events-store
      name: kubemq-events
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-events-store-elastic-search-connector"
        auth_token: ""
        channel: "events-store.elastic-search"
        group:   ""
        concurrency: "1"
        response_channel: "events-store.response.elastic"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.stores.elastic-search
      name: target-elastic-search
      properties:
        urls: "http://localhost:9200"
        username: "admin"
        password: "password"
        sniff: "false"
```
