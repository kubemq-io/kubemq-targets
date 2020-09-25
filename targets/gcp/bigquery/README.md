# Kubemq bigquery target Connector

Kubemq gcp-bigquery target connector allows services using kubemq server to access google bigquery server.

## Prerequisites
The following are required to run the gcp-bigquery target connector:

- kubemq cluster
- gcp-bigquery set up
- kubemq-source deployment

## Configuration

bigquery target connector configuration properties:

| Properties Key | Required | Description                                | Example                    |
|:---------------|:---------|:-------------------------------------------|:---------------------------|
| project_id     | yes      | gcp bigquery project_id                    | "<googleurl>/myproject"    |
| credentials    | yes      | gcp credentials files                      | "<google json credentials" |

Example:

```yaml
bindings:
  - name: kubemq-query-gcp-bigquery
    source:
      kind: source.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-gcp-bigquery-connector"
        auth_token: ""
        channel: "query.gcp.bigquery"
        group:   ""
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


### Create Table Request

Create table metadata setting:

| Metadata Key | Required | Description                             | Possible values       |
|:-------------|:---------|:----------------------------------------|:----------------------|
| method       | yes      | type of method                          | "create_table"        |
| dataset_id   | yes      | dataset to assign the table to          | "your data set ID"  |
| table_name   | yes      | table name                              | "unique name" |


Example:

```json
{
  "metadata": {
    "method": "create_table",
    "dataset_id": "<mySet>",
    "table_name": "<myTable>"
  },
  "data": null
}
```



### Get DataSets Request

Get DataSets setting:

| Metadata Key | Required | Description                             | Possible values   |
|:-------------|:---------|:----------------------------------------|:------------------|
| method       | yes      | type of method                          | "get_data_sets"   |


Example:

```json
{
  "metadata": {
    "method": "get_data_sets"
  },
  "data": null
}
```

### Get Table Info

Get table Info

| Metadata Key | Required | Description                             | Possible values       |
|:-------------|:---------|:----------------------------------------|:----------------------|
| method       | yes      | type of method                          | "get_table_info"      |
| dataset_id   | yes      | dataset to assign the table to          | "your data set ID"  |
| table_name   | yes      | table name                              | "unique table name" |


Example:

```json
{
  "metadata": {
    "method": "get_table_info",
    "dataset_id": "<mySet>",
    "table_name": "<myTable>"
  },
  "data": null
}
```


### Insert To Table

Insert To Table this method required a body of rows of string [bigquery.value]



Example how to create the struct:
```go
    var rows = []map[string]bigquery.Value
	firstRow := make(map[string]bigquery.Value)
	firstRow["name"] = "myName4"
	firstRow["age"] = 25
	rows = append(rows, firstRow)
	rows = append(rows, secondRow)
	bRows, err := json.Marshal(&rows)
```

| Metadata Key | Required | Description                             | Possible values       |
|:-------------|:---------|:----------------------------------------|:----------------------|
| method       | yes      | type of method                          | "insert"              |
| dataset_id   | yes      | dataset to assign the table to          | "your data set ID"  |
| table_name   | yes      | table name                              | "unique table name" |


Example:

```json
{
  "metadata": {
    "method": "insert",
    "dataset_id": "<mySet>",
    "table_name": "<myTable>"
  },
  "data": "W3sgIm5hbWUiOiJteU5hbWU0IiwgImFnZSI6MjV9XQ=="
}
```
