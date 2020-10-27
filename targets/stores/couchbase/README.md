# Kubemq Couchbase Target Connector

Kubemq couchbase target connector allows services using kubemq server to access couchbase database services.

## Prerequisites
The following are required to run the couchbase target connector:

- kubemq cluster
- couchbase server
- kubemq-targets deployment

## Configuration

Couchbase target connector configuration properties:

| Properties Key   | Required | Description            | Example          |
|:-----------------|:---------|:-----------------------|:-----------------|
| url              | yes      | couchbase host address | "localhost"      |
| username         | no       | couchbase username     | "couchdb"        |
| password         | no       | couchbase password     | "couchdb"        |
| bucket           | no       | set bucket name        | "bucket"         |
| num_to_replicate | no       | set replication number | "1"              |
| num_to_persist   | no       | set persistence number | "1"              |
| collection       | no       | set collection name    | "collection"     |



Example:

```yaml
bindings:
  - name: kubemq-query-couchbase
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-couchbase-connector"
        auth_token: ""
        channel: "query.couchbase"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind:stores.couchbase
      name: target-couchbase
      properties:
        url: "localhost"
        username: "couchbase"
        password: "couchbase"
        bucket: "bucket"
        collection: "test"
        num_to_replicate: "1"
        num_to_persist: "1"
```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | couchbase key string | any string      |
| method       | yes      | get              | "get"           |


Example:

```json
{
  "metadata": {
    "key": "your-couchbase-key",
    "method": "get"
  },
  "data": null
}
```

### Set Request

Set request metadata setting:

| Metadata Key   | Required | Description               | Possible values  |
|:---------------|:---------|:--------------------------|:-----------------|
| key            | yes      | couchbase key string      | any string       |
| method         | yes      | set                       | "set"            |
| cas            | no       | set cas value             | "0"              |
| expiry_seconds | no       | set key expiry in seconds | "3600"           |

Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the couchbase key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "key": "your-couchbase-key",
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
| key            | yes      | couchbase key string      | any string       |
| method         | yes      | set                       | "delete"            |
| cas            | no       | set cas value             | "0"              |


Example:

```json
{
  "metadata": {
    "key": "your-couchbase-key",
    "method": "delete",
    "cas": "0"
  },
  "data": null
}
```
