# Kubemq athena target Connector

Kubemq athena target connector allows services using kubemq server to access aws athena service.

## Prerequisites
The following required to run the aws-athena target connector:

- kubemq cluster
- aws account with athena active service
- kubemq-source deployment

## Configuration

athena target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-athena
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-athena"
        auth_token: ""
        channel: "query.aws.athena"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.athena
      name: aws-athena
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "region"
        token: ""
```

## Usage

### List Catalogs 

list all catalogs

List Catalogs:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_data_catalogs"                     |


Example:

```json
{
  "metadata": {
    "method": "list_data_catalogs"
  },
  "data": null
}
```

### List Databases 

list all databases 

List Databases:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_databases"                     |
| catalog           | yes      | aws athena catalog                      | "string"                     |


Example:

```json
{
  "metadata": {
    "method": "list_databases",
    "catalog": "my_catalog"
  },
  "data": null
}
```


### Query 

create a query request return execution_id.

Query:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "query"                     |
| catalog           | yes      | aws athena catalog                      | "string"                     |
| db                | yes      | aws athena db name                      | "string"                     |
| output_location   | yes      | aws athena folder location              | "string"                     |
| query             | yes      | aws query to execute                    | "query"                     |


Example:

```json
{
  "metadata": {
    "method": "query",
    "catalog": "my_catalog",
    "db": "my_db",
    "output_location": "my_output_location/path",
    "query": "SELECT * FROM \"my_db\".\"new_table_name2\" limit 10;"
  },
  "data": null
}
```


### Get Query Result

get Query result by execution_id that return from Query result.

Get Query Result:

| Metadata Key      | Required | Description                                | Possible values                            |
|:------------------|:---------|:-------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                             | "get_query_result"                         |
| execution_id      | yes      | aws executionID that returned from query   | "string"                                   |


Example:

```json
{
  "metadata": {
    "method": "get_query_result",
    "execution_id": "123-456-789"
  },
  "data": null
}
```
