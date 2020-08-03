# Kubemq spanner target Connector

Kubemq gcp-spanner target connector allows services using kubemq server to access google spanner server.

## Prerequisites
The following are required to run the gcp-spanner target connector:

- kubemq cluster
- gcp-spanner set up
- kubemq-source deployment

## Configuration

spanner target connector configuration properties:

| Properties Key | Required | Description                                | Example                         |
|:---------------|:---------|:-------------------------------------------|:--------------------------------|
| db             | yes      | gcp spanner db name                        | "<googleurl>/mydb"              |
| credentials    | yes      | gcp credentials files                      | "<google json credentials"      |

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
        client_id: "kubemq-query-gcp-spanner-connector"
        auth_token: ""
        channel: "query.gcp.spanner"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.gcp.spanner
      name: target-gcp-spanner
      properties:
        db: "id"
        credentials: 'json'

```

## Usage

### Query Request

Query metadata setting:

| Metadata Key | Required | Description                  | Possible values       |
|:-------------|:---------|:-----------------------------|:----------------------|
| method       | yes      | type of method               | "query"               |
| query        | yes      | the query body               | "select * from table" |


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


### Read Table Request by columns

Read Table metadata setting:

| Metadata Key | Required | Description               | Possible values                         |
|:-------------|:---------|:--------------------------|:----------------------------------------|
| method       | yes      | type of method            | "read"                                  |
| table_name   | yes      | table name to read from   | "<your data set ID>"                    |


Example:

```json
{
  "metadata": {
    "method": "read",
    "table_name": "<myTable>"
  },
  "data": "W1wiaWRcIixcIm5hbWVcIl0="
}
```

### Insert Or Update Table

Insert Or Update metadata setting:

| Metadata Key | Required | Description                     | Possible values                         |
|:-------------|:---------|:--------------------------------|:----------------------------------------|
| method       | yes      | type of method                  | "insert","update","insert_or_update"    |


Example:

```json
{
  "metadata": {
    "method": "insert_or_update"
  },
  "data": "W3tcInRhYmxlX25hbWVcIjpcInRlc3QyXCIsXCJjb2x1bW5fbmFtZXNcIjpbXCJpZFwiLFwibmFtZVwiXSxcImNvbHVtbl92YWx1ZXNcIjpbMTcsXCJuYW1lMVwiXSxcImNvbHVtbl90eXBlXCI6W1wiSU5UNjRcIixcIlNUUklOR1wiXX0se1widGFibGVfbmFtZVwiOlwidGVzdDJcIixcImNvbHVtbl9uYW1lc1wiOltcImlkXCIsXCJuYW1lXCJdLFwiY29sdW1uX3ZhbHVlc1wiOlsxOCxcIm5hbWUyXCJdLFwiY29sdW1uX3R5cGVcIjpbXCJJTlQ2NFwiLFwiU1RSSU5HXCJdfV0="
}
```
