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
      kind: kubemq.query
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
      kind:gcp.bigquery
      name: gcp-bigquery
      properties:
        project_id: "id"
        credentials: 'json'
```


## Usage

### Query Request

query request.

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

create a new table under data set

This method required a body of rows of string [bigquery.TableMetadata]



Example how to create the struct:
```go
    mySchema := bigquery.Schema{
    		{Name: "name", Type: bigquery.StringFieldType},
    		{Name: "age", Type: bigquery.IntegerFieldType},
    	}
    
    	metaData := &bigquery.TableMetadata{
    		Schema:         mySchema,
    		ExpirationTime: time.Now().AddDate(2, 1, 0), // Table will deleted in 2 years and 1 month.
    	}
    	bSchema, err := json.Marshal(metaData)
```

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
  "data": "eyJOYW1lIjoiIiwiTG9jYXRpb24iOiIiLCJEZXNjcmlwdGlvbiI6IiIsIlNjaGVtYSI6W3siTmFtZSI6Im5hbWUiLCJEZXNjcmlwdGlvbiI6IiIsIlJlcGVhdGVkIjpmYWxzZSwiUmVxdWlyZWQiOmZhbHNlLCJUeXBlIjoiU1RSSU5HIiwiUG9saWN5VGFncyI6bnVsbCwiU2NoZW1hIjpudWxsfSx7Ik5hbWUiOiJhZ2UiLCJEZXNjcmlwdGlvbiI6IiIsIlJlcGVhdGVkIjpmYWxzZSwiUmVxdWlyZWQiOmZhbHNlLCJUeXBlIjoiSU5URUdFUiIsIlBvbGljeVRhZ3MiOm51bGwsIlNjaGVtYSI6bnVsbH1dLCJNYXRlcmlhbGl6ZWRWaWV3IjpudWxsLCJWaWV3UXVlcnkiOiIiLCJVc2VMZWdhY3lTUUwiOmZhbHNlLCJVc2VTdGFuZGFyZFNRTCI6ZmFsc2UsIlRpbWVQYXJ0aXRpb25pbmciOm51bGwsIlJhbmdlUGFydGl0aW9uaW5nIjpudWxsLCJSZXF1aXJlUGFydGl0aW9uRmlsdGVyIjpmYWxzZSwiQ2x1c3RlcmluZyI6bnVsbCwiRXhwaXJhdGlvblRpbWUiOiIyMDIyLTExLTI1VDEwOjQxOjI1LjIyNjQ1NyswMjowMCIsIkxhYmVscyI6bnVsbCwiRXh0ZXJuYWxEYXRhQ29uZmlnIjpudWxsLCJFbmNyeXB0aW9uQ29uZmlnIjpudWxsLCJGdWxsSUQiOiIiLCJUeXBlIjoiIiwiQ3JlYXRpb25UaW1lIjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJMYXN0TW9kaWZpZWRUaW1lIjoiMDAwMS0wMS0wMVQwMDowMDowMFoiLCJOdW1CeXRlcyI6MCwiTnVtTG9uZ1Rlcm1CeXRlcyI6MCwiTnVtUm93cyI6MCwiU3RyZWFtaW5nQnVmZmVyIjpudWxsLCJFVGFnIjoiIn0="
}
```


### Delete Table Request

delete a new table under data set

Delete table metadata setting:

| Metadata Key | Required | Description                             | Possible values       |
|:-------------|:---------|:----------------------------------------|:----------------------|
| method       | yes      | type of method                          | "create_table"        |
| dataset_id   | yes      | dataset to assign the table to          | "your data set ID"  |
| table_name   | yes      | table name                              | "unique name" |


Example:

```json
{
  "metadata": {
    "method": "delete_table",
    "dataset_id": "<mySet>",
    "table_name": "<myTable>"
  },
  "data":null
}
```

### Create Data Set Request

Create a Data Set

Create Data Set metadata setting:

| Metadata Key | Required | Description                             | Possible values       |
|:-------------|:---------|:----------------------------------------|:----------------------|
| method       | yes      | type of method                          | "create_data_set"        |
| dataset_id   | yes      | dataset to assign the table to          | "your data set ID"  |
| location     | yes      | dataset location to set                 | "US"  See https://cloud.google.com/bigquery/docs/locations |


Example:

```json
{
  "metadata": {
    "method": "create_data_set",
    "dataset_id": "<mySet>",
    "location": "US"
  },
  "data":null
}
```


### Delete Data Set Request

delete a Data Set

Delete Data Set metadata setting:

| Metadata Key | Required | Description                             | Possible values       |
|:-------------|:---------|:----------------------------------------|:----------------------|
| method       | yes      | type of method                          | "delete_data_set"        |
| dataset_id   | yes      | dataset to assign the table to          | "your data set ID"  |


Example:

```json
{
  "metadata": {
    "method": "delete_data_set",
    "dataset_id": "<mySet>"
  },
  "data":null
}
```

### Get DataSets Request

get data sets.

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

get basic information on a table by name

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

insert rows to table

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
