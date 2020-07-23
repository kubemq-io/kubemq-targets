# Kubemq bigquery Source Connector

Kubemq gcp-bigquery source connector allows services using kubemq server to access google bigquery server.

## Prerequisites
The following are required to run the gcp-bigquery target connector:

- kubemq cluster
- gcp-bigquery set up
- kubemq-targets deployment

## Configuration

bigquery source connector configuration properties:

| Properties Key | Required | Description                                | Example          |
|:---------------|:---------|:-------------------------------------------|:-----------------|
| project_id        | yes      | gcp bigquery project_id                    | "<googleurl>/myproject" |
| credentials       | yes      | gcp credentials files                      | "<google json credentials"      |

Example:

```yaml
bindings:
  - name: kubemq-query-gcp-bigquery
    source:
      kind: source.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-gcp-bigquery-connector"
        auth_token: ""
        channel: "query.gcp.bigquery"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.gcp.bigquery
      name: target-gcp-bigquery
      properties:
        project_id: "id"
        credentials: 'json'
```

## Usage

### Query Request

query metadata setting:

| Metadata Key | Required | Description                             | Possible values                         |
|:-------------|:---------|:----------------------------------------|:----------------------------------------|
| method          | yes      | type of method               | "query"                                  |
| query           | yes      | the query body               | "select * from table" |


Example:

```json
{
  "metadata": {
    "method": "query",
    "query": "select * from table"
  },
  "data": null
}
```
