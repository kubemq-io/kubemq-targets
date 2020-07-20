# Kubemq Command Source

Kubemq Command source provides rpc command subscriber for processing targets commands.

## Prerequisites
The following are required to run command source connector:

- kubemq cluster
- kubemq-targets deployment


## Configuration

Command target connector configuration properties:

| Properties Key             | Required | Description                                                           | Example         |
|:---------------------------|:---------|:----------------------------------------------------------------------|:----------------|
| host                       | yes      | kubemq server host address                                            | "localhost      |
| port                       | yes      | kubemq server port number                                             | "50000"         |
| client_id                  | no       | set client id                                                         | "client_id"     |
| auth_token                 | no       | set authentication token                                              | jwt token       |
| channel                    | yes      | set channel to subscribe                                              |                 |
| group                      | no       | set subscriber group                                                  |                 |
| concurrency                | yes      | set parallel subscribers count       | "10"            |
| auto_reconnect             | no       | set auto reconnect on lost connection                                 | "false", "true" |
| reconnect_interval_seconds | no       | set reconnection seconds                                              | "5"             |
| max_reconnects             | no       | set how many time to reconnect                                        | "0"             |





Example:

```yaml
bindings:
  - name: kubemq-command-elastic-search
    source:
      kind: source.command
      name: kubemq-command
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-command-elastic-search-connector"
        auth_token: ""
        channel: "command.elastic-search"
        group:   ""
        concurrency: "1"
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
