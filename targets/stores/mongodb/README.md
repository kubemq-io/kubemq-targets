# Kubemq Mongodb Target Connector

Kubemq mongodb target connector allows services using kubemq server to access mongodb database services.

## Prerequisites
The following are required to run the mongodb target connector:

- kubemq cluster
- mongodb server
- kubemq-targets deployment

## Configuration

Mongodb target connector configuration properties:

| Properties Key            | Required | Description                          | Example                   |
|:--------------------------|:---------|:-------------------------------------|:--------------------------|
| host                      | yes      | mongodb host address                 | "localhost:27017"         |
| username                  | no       | mongodb username                     | "admin"                   |
| password                  | no       | mongodb password                     | "password"                |
| database                  | no       | set database name                    | "admin"                   |
| collection                | no       | set database collection              | "test"                    |
| params                    | no       | set connection additional parameters | ""                        |
| write_concurrency         | no       | set write concurrency                | "","majority","1","2"     |
| read_concurrency          | no       | set read concurrency                 | "","local"                |
|                           |          |                                      | "","local"                |
|                           |          |                                      | "majority","available"    |
|                           |          |                                      | "linearizable","snapshot" |
| operation_timeout_seconds | no       | set operation timeout in seconds     | "30"                      |



Example:

```yaml
bindings:
  - name: kubemq-query-mongodb
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-mongodb-connector"
        auth_token: ""
        channel: "query.mongodb"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind:stores.mongodb
      name: target-mongodb
      properties:
        host: "localhost:27017"
        username: "admin"
        password: "password"
        database: "admin"
        collection: "test"
        write_concurrency: "majority"
        read_concurrency: ""
        params: ""
        operation_timeout_seconds: "2"
```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | mongodb key string | any string      |
| method       | yes      | get              | "get"           |

Example:

```json
{
  "metadata": {
    "key": "your-mongodb-key",
    "method": "get"
  },
  "data": null
}
```

### Set Request

Set request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | mongodb key string | any string      |
| method       | yes      | set              | "set"           |

Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the mongodb key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "key": "your-mongodb-key",
    "method": "set"
  },
  "data": "c29tZS1kYXRh" 
}
```
### Delete Request

Delete request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | mongodb key string | any string      |
| method       | yes      | delete           | "delete"        |


Example:

```json
{
  "metadata": {
    "key": "your-mongodb-key",
    "method": "delete"
  },
  "data": null
}
```
