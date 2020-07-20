# Kubemq Elastic Search Target Connector

Kubemq elastic-search target connector allows services using kubemq server to access elastic-search database services.

## Prerequisites
The following are required to run the elastic-search target connector:

- kubemq cluster
- elastic-search server
- kubemq-targets deployment

## Configuration

Elastic Search target connector configuration properties:

| Properties Key            | Required | Description                          | Example                   |
|:--------------------------|:---------|:-------------------------------------|:--------------------------|
| urls                      | yes      | elastic-search list of urls separated by comma                 | "http://localhost:9200,http://localhost:9201"         |
| username                  | no       | elastic-search username                     | "admin"                   |
| password                  | no       | elastic-search password                     | "password"                |
| sniff                  | no       | set sniff opn connect                    | "true", "false"                   |



Example:

```yaml
bindings:
  - name: kubemq-query-elastic-search
    source:
      kind: source.kubemq.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-elastic-search-connector"
        auth_token: ""
        channel: "query.elastic-search"
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

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name get                        | "get"           |
| index        | yes      | elastic-search index table | any string      |
| id           | yes      | document id                | any string      |


Example:

```json
{
  "metadata": {
    "method": "get",
    "index": "log",
    "id": "doc-id"
  },
  "data": null
}
```

### Set Request

Set request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name set                        | "set"           |
| index        | yes      | elastic-search index table | any string      |
| id           | yes      | document id                | any string      |


Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the elastic-search key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "method": "set",
    "index": "log",
    "id": "doc-id"
  },
  "data": "c29tZS1kYXRh" 
}
```
### Delete Request

Delete request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name delete                        | "delete"           |
| index        | yes      | elastic-search index table | any string      |
| id           | yes      | document id                | any string      |

Example:

```json
{
  "metadata": {
    "method": "delete",
    "index": "log",
    "id": "doc-id"
  },
  "data": null
}
```
