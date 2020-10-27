# Kubemq keyspaces Target Connector

Kubemq keyspaces target connector allows services using kubemq server to access keyspaces database services.

## Prerequisites
The following are required to run the keyspaces target connector:

- kubemq cluster
- IAM user keyspaces credentials 
- aws keyspaces server/cluster
- kubemq-targets deployment

## Configuration

keyspaces target connector configuration properties:

| Properties Key            | Required | Description                            | Example                             |
|:--------------------------|:---------|:---------------------------------------|:------------------------------------|
| hosts                     | yes      | aws end point                          | "localhost"                         |
| port                      | yes      | keyspaces port                         | "9142"                              |
| proto_version             | no       | keyspaces proto version                | "4"                                 |
| replication_factor        | no       | set replication factor                 | "1"                            |
| username                  | no       | set keyspaces username                 | "keyspaces"                         |
| password                  | no       | set keyspaces password                 | "keyspaces"                         |
| consistency               | no       | set keyspaces consistency              | "One","LocalOne","LocalQuorum"  see https://docs.aws.amazon.com/keyspaces/latest/devguide/consistency.html    |
| default_table             | no       | set table name                         | "test"                              |
| default_keyspace          | no       | set keyspace name                      | "test"                              |
| tls                       | yes      | aws keyspace certificate               | aws tls link see https://docs.aws.amazon.com/keyspaces/latest/devguide/using_go_driver.html                  |
| timeout_seconds           | no       | set default timeout seconds            | "60"                              |
| connect_timeout_seconds   | no       | set default connect timeout seconds    | "60"                              |




Example:

```yaml
bindings:
  - name: kubemq-query-keyspaces
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-keyspaces-connector"
        auth_token: ""
        channel: "aws.query.keyspaces"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.keyspaces
      name: aws-keyspaces
      properties:
        hosts: "cassandra.us-east-2.amazonaws.com"
        port: "9142"
        username: "keyspaces"
        password: "keyspaces"
        proto_version: "4"
        replication_factor: "1"
        consistency: "LocalQuorum"
        default_table: "test"
        default_keyspace: "test"
        tls: "https://www.amazontrust.com/repository/AmazonRootCA1.pem"
```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description          | Possible values       |
|:-------------|:---------|:---------------------|:----------------------|
| key          | yes      | keyspaces key string | any string            |
| method       | yes      | get                  | "get"                 |
| consistency  | yes      | set consistency      | "",strong","eventual" |
| table        | yes      | table name           | "table                |
| keyspace     | yes      | key space name       | "keyspace"            |

Example:

```json
{
  "metadata": {
    "key": "your-keyspaces-key",
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
| key          | yes      | keyspaces key string | any string            |
| method       | yes      | method name set                  | "set"                 |
| consistency  | yes      | set consistency                  | "",strong","eventual" |
| table        | yes      | table name           | "table                |
| keyspace     | yes      | key space name       | "keyspace"            |

Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the keyspaces key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "key": "your-keyspaces-key",
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
| key          | yes      | keyspaces key string | any string      |
| method       | yes      | method name delete   | "delete"        |
| table        | yes      | table name           | "table          |
| keyspace     | yes      | key space name       | "keyspace"      |



Example:

```json
{
  "metadata": {
    "key": "your-keyspaces-key",
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
