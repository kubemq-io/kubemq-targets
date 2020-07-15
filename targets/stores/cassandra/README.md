# Kubemq Cassandra Target Connector

Kubemq cassandra target connector allows services using kubemq server to access cassandra database services.

## Prerequisites
The following are required to run the cassandra target connector:

- kubemq cluster
- cassandra server
- kubemq-target-connectors deployment

## Configuration

Cassandra target connector configuration properties:

| Properties Key     | Required | Description               | Example                             |
|:-------------------|:---------|:--------------------------|:------------------------------------|
| hosts              | yes      | cassandra hosts addresses | "localhost"                         |
| port               | yes      | cassandra port            | "9042"                              |
| proto_version      | no       | cassandra proto version   | "4"                                 |
| replication_factor | no       | set replication factor           | "1"                            |
| username           | no       | set cassandra username    | "cassandra"                         |
| password           | no       | set cassandra password    | "cassandra"                         |
| consistency        | no       | set cassandra consistency | "", "All","One","Two"               |
|                    |          |                           | "Quorum","LocalQuorum","EachQuorum" |
|                    |          |                           | "LocalOne","Any"                    |
| default_table      | no       | set table name            | "test"                              |
| default_keyspace   | no       | set keyspace name         | "test"                              |




Example:

```yaml
bindings:
  - name: kubemq-query-cassandra
    source:
      kind: source.kubemq.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-cassandra-connector"
        auth_token: ""
        channel: "query.cassandra"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.stores.cassandra
      name: target-cassandra
      properties:
        hosts: "localhost"
        port: "9042"
        username: "cassandra"
        password: "cassandra"
        proto_version: "4"
        replication_factor: "1"
        consistency: "All"
        default_table: "test"
        default_keyspace: "test"
```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | cassandra key string | any string      |
| method       | yes      | get              | "get"           |
| consistency       | yes      | get              | "strong"           |


Example:

```json
{
  "metadata": {
    "key": "your-cassandra-key",
    "method": "get"
  },
  "data": null
}
```

### Set Request

Set request metadata setting:

| Metadata Key   | Required | Description               | Possible values  |
|:---------------|:---------|:--------------------------|:-----------------|
| key            | yes      | cassandra key string      | any string       |
| method         | yes      | set                       | "set"            |
| cas            | no       | set cas value             | "0"              |
| expiry_seconds | no       | set key expiry in seconds | "3600"           |

Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the cassandra key | base 64 bytes array |

Example:

```json
{
  "metadata": {
    "key": "your-cassandra-key",
    "method": "set",
    "cas": "0",
    "expiry_seconds": "3600"
  },
  "data": "c29tZS1kYXRh" 
}
```
### Delete Request

Delete request metadata setting:

| Metadata Key   | Required | Description               | Possible values  |
|:---------------|:---------|:--------------------------|:-----------------|
| key            | yes      | cassandra key string      | any string       |
| method         | yes      | set                       | "delete"            |
| cas            | no       | set cas value             | "0"              |


Example:

```json
{
  "metadata": {
    "key": "your-cassandra-key",
    "method": "delete",
    "cas": "0"
  },
  "data": null
}
```
