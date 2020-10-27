# Kubemq dynamodb target Connector

Kubemq dynamodb target connector allows services using kubemq server to access aws dynamodb service.

## Prerequisites
The following required to run the aws-dynamodb target connector:

- kubemq cluster
- aws account with dynamodb active service
- kubemq-source deployment

## Configuration

dynamodb target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-dynamodb
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-dynamodb"
        auth_token: ""
        channel: "query.aws.dynamodb"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.dynamodb
      name: aws-dynamodb
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "region"
        token: ""
```

## Usage

### List Tables 

list all tables under dynamodb.

List Tables:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_tables"                     |


Example:

```json
{
  "metadata": {
    "method": "list_tables"
  },
  "data": null
}
```

### Create Table 

create a new table under dynamodb.

Create Table:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "create_table"                     |
| data              | yes      | dynamodb.CreateTableInput as json       | "string"                     |


Example:

```json
{
  "metadata": {
    "method": "create_table"
  },
  "data": "ewoJCQkJCSJBdHRyaWJ1dGVEZWZpbml0aW9ucyI6IFsKCQkJCQkJewoJCQkJCQkJIkF0dHJpYnV0ZU5hbWUiOiAiWWVhciIsCgkJCQkJCQkiQXR0cmlidXRlVHlwZSI6ICJOIgoJCQkJCQl9LAoJCQkJCQl7CgkJCQkJCQkiQXR0cmlidXRlTmFtZSI6ICJUaXRsZSIsCgkJCQkJCQkiQXR0cmlidXRlVHlwZSI6ICJTIgoJCQkJCQl9CgkJCQkJXSwKCQkJCQkiQmlsbGluZ01vZGUiOiBudWxsLAoJCQkJCSJHbG9iYWxTZWNvbmRhcnlJbmRleGVzIjogbnVsbCwKCQkJCQkiS2V5U2NoZW1hIjogWwoJCQkJCQl7CgkJCQkJCQkiQXR0cmlidXRlTmFtZSI6ICJZZWFyIiwKCQkJCQkJCSJLZXlUeXBlIjogIkhBU0giCgkJCQkJCX0sCgkJCQkJCXsKCQkJCQkJCSJBdHRyaWJ1dGVOYW1lIjogIlRpdGxlIiwKCQkJCQkJCSJLZXlUeXBlIjogIlJBTkdFIgoJCQkJCQl9CgkJCQkJXSwKCQkJCQkiTG9jYWxTZWNvbmRhcnlJbmRleGVzIjogbnVsbCwKCQkJCQkiUHJvdmlzaW9uZWRUaHJvdWdocHV0IjogewoJCQkJCQkiUmVhZENhcGFjaXR5VW5pdHMiOiAxMCwKCQkJCQkJIldyaXRlQ2FwYWNpdHlVbml0cyI6IDEwCgkJCQkJfSwKCQkJCQkiU1NFU3BlY2lmaWNhdGlvbiI6IG51bGwsCgkJCQkJIlN0cmVhbVNwZWNpZmljYXRpb24iOiBudWxsLAoJCQkJCSJUYWJsZU5hbWUiOiAibXl0YWJsZW5hbWUiLAoJCQkJCSJUYWdzIjogbnVsbAoJCQkJfQ=="
}
```

### Delete Table 

delete a table under dynamodb.

Delete Table:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "delete_table"                     |
| table_name        | yes      | table name to delete                    | "string"                     |


Example:

```json
{
  "metadata": {
    "method": "delete_table",
    "table_name": "my_table_name"
  },
  "data": null
}
```


### Insert Item 

insert item to table under dynamodb.

Insert Item :

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "insert_item"                     |
| table_name        | yes      | table name to delete                    | "string"                     |
| data              | yes      | dynamodb.AttributeValue as json         | "string"                     |


Example:

```json
{
  "metadata": {
    "method": "insert_item",
    "table_name": "my_table_name"
  },
  "data": "ewoJCSJQbG90IjogewoJCQkiUyI6ICJzb21lIHBsb3QiCgkJfSwKCQkiUmF0aW5nIjogewoJCQkiTiI6ICIxMC4xIgoJCX0sCgkJIlRpdGxlIjogewoJCQkiUyI6ICJLdWJlTVEgdGVzdCBNb3ZpZSIKCQl9LAoJCSJZZWFyIjogewoJCQkiTiI6ICIyMDIwIgoJCX0KCX0="
}
```

### Get Item 

get an item from a table under dynamodb.

Get Item :

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "get_item"                     |
| data              | yes      | dynamodb.GetItemInput as json           | "string"                     |


Example:

```json
{
  "metadata": {
    "method": "get_item"
  },
  "data": "ewoJCQkiQXR0cmlidXRlc1RvR2V0IjogbnVsbCwKCQkJIkNvbnNpc3RlbnRSZWFkIjogbnVsbCwKCQkJIkV4cHJlc3Npb25BdHRyaWJ1dGVOYW1lcyI6IG51bGwsCgkJCSJLZXkiOiB7CgkJCQkiVGl0bGUiOiB7CgkJCQkJIlMiOiAiS3ViZU1RIHRlc3QgTW92aWUiCgkJCQl9LAoJCQkJIlllYXIiOiB7CgkJCQkJIk4iOiAiMjAyMCIKCQkJCX0KCQkJfSwKCQkJIlByb2plY3Rpb25FeHByZXNzaW9uIjogbnVsbCwKCQkJIlJldHVybkNvbnN1bWVkQ2FwYWNpdHkiOiBudWxsLAoJCQkiVGFibGVOYW1lIjogIm15dGFibGVuYW1lIgoJCX0="
}
```

### Update Item 

update an item from a table under dynamodb.


Update Item :

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "update_item"                     |
| data              | yes      | dynamodb.UpdateItemInput as json        | "string"                     |


Example:

```json
{
  "metadata": {
    "method": "update_item"
  },
  "data": "ewoJCQkiRXhwcmVzc2lvbkF0dHJpYnV0ZVZhbHVlcyI6IHsKCQkJCSI6ciI6IHsKCQkJCQkiTiI6ICIwLjkiCgkJCQl9CgkJCX0sCgkJCSJLZXkiOiB7CgkJCQkiVGl0bGUiOiB7CgkJCQkJIlMiOiAiS3ViZU1RIHRlc3QgTW92aWUiCgkJCQl9LAoJCQkJIlllYXIiOiB7CgkJCQkJIk4iOiAiMjAyMCIKCQkJCX0KCQkJfSwKCQkJIlJldHVyblZhbHVlcyI6ICJVUERBVEVEX05FVyIsCgkJCSJUYWJsZU5hbWUiOiAibXl0YWJsZW5hbWUiLAoJCQkiVXBkYXRlRXhwcmVzc2lvbiI6ICJzZXQgUmF0aW5nID0gOnIiCgkJfQ=="
}
```

### Delete Item 

delete an item from a table under dynamodb.

Delete Item :

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "delete_item"                     |
| data              | yes      | dynamodb.DeleteItemInput as json        | "string"                     |


Example:

```json
{
  "metadata": {
    "method": "delete_item"
  },
  "data": "ewoJCQkJCSJLZXkiOiB7CgkJCQkJCSJUaXRsZSI6IHsKCQkJCQkJCSJTIjogIkt1YmVNUSB0ZXN0IE1vdmllIgoJCQkJCQl9LAoJCQkJCQkiWWVhciI6IHsKCQkJCQkJCSJOIjogIjIwMjAiCgkJCQkJCX0KCQkJCQl9LAoJCQkJCSJUYWJsZU5hbWUiOiAibXl0YWJsZW5hbWUiCgkJCQl9"
}
```
