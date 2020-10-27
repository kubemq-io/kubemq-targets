# Kubemq Cassandra Target Connector

Kubemq cassandra target connector allows services using kubemq server to access cassandra database services.

## Prerequisites
The following are required to run the cassandra target connector:

- kubemq cluster
- cassandra server/cluster
- kubemq-targets deployment

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
| timeout_seconds      | no       |set default timeout seconds            | "60"                              |
| connect_timeout_seconds   | no       | set default connect timeout seconds         | "60"                              |




Example:

```yaml
bindings:
  - name: kubemq-query-cassandra
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-cassandra-connector"
        auth_token: ""
        channel: "query.cassandra"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind:stores.cassandra
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
        connect_timeout_seconds: "60"
        timeout_seconds: "60"
```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description          | Possible values       |
|:-------------|:---------|:---------------------|:----------------------|
| key          | yes      | cassandra key string | any string            |
| method       | yes      | get                  | "get"                 |
| consistency  | yes      | set consistency                   | "",strong","eventual" |
| table        | yes      | table name           | "table                |
| keyspace     | yes      | key space name       | "keyspace"            |

Example:

```json
{
  "metadata": {
    "key": "your-cassandra-key",
    "method": "get",
    "consistency": "",
    "table": "table",
    "keyspace": "keyspace"
  },
  "data": null
}
```

### Set Request

Set request metadata setting:

| Metadata Key   | Required | Description               | Possible values  |
|:---------------|:---------|:--------------------------|:-----------------|
| key          | yes      | cassandra key string | any string            |
| method       | yes      | method name set                  | "set"                 |
| consistency  | yes      | set consistency                  | "",strong","eventual" |
| table        | yes      | table name           | "table                |
| keyspace     | yes      | key space name       | "keyspace"            |

Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the cassandra key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "key": "your-cassandra-key",
    "method": "set",
    "consistency": "",
    "table": "table",
    "keyspace": "keyspace"
  },
  "data": "c29tZS1kYXRh" 
}
```
### Delete Request

Delete request metadata setting:

| Metadata Key | Required | Description          | Possible values |
|:-------------|:---------|:---------------------|:----------------|
| key          | yes      | cassandra key string | any string      |
| method       | yes      | method name delete   | "delete"        |
| table        | yes      | table name           | "table          |
| keyspace     | yes      | key space name       | "keyspace"      |



Example:

```json
{
  "metadata": {
    "key": "your-cassandra-key",
    "method": "set",
    "table": "table",
    "keyspace": "keyspace"
  },
  "data": null
}
```

### Query Request

Query request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | method name query   | "query"        |
| consistency  | yes      | set consistency                  | "",strong","eventual" |


Query request data setting:

| Data Key | Required | Description  | Possible values    |
|:---------|:---------|:-------------|:-------------------|
| data     | yes      | query string | base64 bytes array |

Example:

Query string: `SELECT value FROM test.test WHERE key = 'some-key`

```json
{
  "metadata": {
    "method": "query",
    "consistency": "strong"
  },
  "data": "U0VMRUNUIHZhbHVlIEZST00gdGVzdC50ZXN0IFdIRVJFIGtleSA9ICdzb21lLWtleQ=="
}
```

### Exec Request

Exec request metadata setting:

| Metadata Key    | Required | Description                            | Possible values    |
|:----------------|:---------|:---------------------------------------|:-------------------|
| method          | yes      | set type of request                    | "exec"             |
| consistency  | yes      | set consistency                  | "",strong","eventual" |


Exec request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | exec string | base64 bytes array |

Example:

Exec string:
```sql
INSERT INTO test.test (key, value) VALUES ('some-key',textAsBlob('some-data'))
```

```json
{
  "metadata": {
    "method": "exec",
    "consistency": "strong"
  },
  "data": "SU5TRVJUIElOVE8gdGVzdC50ZXN0IChrZXksIHZhbHVlKSBWQUxVRVMgKCdzb21lLWtleScsdGV4dEFzQmxvYignc29tZS1kYXRhJykp" 
}
```
