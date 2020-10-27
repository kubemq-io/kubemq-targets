# Kubemq s3 target Connector

Kubemq aws-s3 target connector allows services using kubemq server to access aws s3 service.

## Prerequisites
The following required to run the aws-s3 target connector:

- kubemq cluster
- aws account with s3 active service
- kubemq-source deployment

## Configuration

s3 target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |
| downloader     | no       | if needed to create downloader instance    | true                      |
| uploader       | no       | if needed to create uploader instance      | false                     |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-s3
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-s3-connector"
        auth_token: ""
        channel: "query.aws.s3"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.s3
      name: aws-s3
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "region"
        token: ""
        downloader:  "true"
        uploader:  "true"
```

## Usage

### List Buckets

list all buckets.

List Buckets:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_buckets"                     |



Example:

```json
{
  "metadata": {
    "method": "list_buckets"
  },
  "data": null
}
```


### List Bucket Items

list all items in the selected bucket

List Bucket Items:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_bucket_items"                     |
| bucket_name       | yes      | s3 bucket name                          | "my_bucket_name"                     |


Example:

```json
{
  "metadata": {
    "method": "list_bucket_items",
    "bucket_name": "my_bucket_name"
  },
  "data": null
}
```



### Create Bucket 

create a new bucket, name must be unique

Create Bucket :

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method              | yes      | type of method                          | "create_bucket"                     |
| bucket_name         | yes      | s3 bucket name                          | "my_bucket_name"                    |
| wait_for_completion | no       | wait for operation to end               | "true","false" default of false     |


Example:

```json
{
  "metadata": {
    "method": "create_bucket",
    "bucket_name": "my_bucket_name"
  },
  "data": null
}
```

### Upload Item 

upload item to bucket.

Upload Bucket Items:

| Metadata Key        | Required | Description                             | Possible values                      |
|:--------------------|:---------|:----------------------------------------|:-------------------------------------|
| method              | yes      | type of method                          | "upload_item"                        |
| bucket_name         | yes      | s3 bucket name                          | "my_bucket_name"                     |
| wait_for_completion | no       | wait for operation to end               | "true","false" (default of false )   |
| item_name           | yes      | the name of the item                    | "valid-string"                       |
| data                | yes      | the object data in byte array           | "valid-string"                       |


Example:

```json
{
  "metadata": {
    "method": "create_bucket",
    "bucket_name": "my_bucket_name",
    "item_name": "my_item_name"
  },
  "data": "bXkgaXRlbSBoZXJl"
}
```

### Get Item 

Get item by item name from bucket

Get Bucket Items:

| Metadata Key        | Required | Description                             | Possible values                      |
|:--------------------|:---------|:----------------------------------------|:-------------------------------------|
| method              | yes      | type of method                          | "get_item"                        |
| bucket_name         | yes      | s3 bucket name                          | "my_bucket_name"                     |
| item_name           | yes      | the name of the item                    | "valid-string"   |


Example:

```json
{
  "metadata": {
    "method": "get_item",
    "bucket_name": "my_bucket_name",
    "item_name": "my_item_name"
  },
  "data": null
}
```

### Delete Item 

delete item by item name from bucket

Delete Item:

| Metadata Key        | Required | Description                             | Possible values                      |
|:--------------------|:---------|:----------------------------------------|:-------------------------------------|
| method              | yes      | type of method                          | "delete_item_from_bucket"                        |
| bucket_name         | yes      | s3 bucket name                          | "my_bucket_name"                     |
| wait_for_completion | no       | wait for operation to end               | "true","false" (default of false )   |
| item_name           | yes      | the name of the item                    | "valid-string"   |


Example:

```json
{
  "metadata": {
    "method": "delete_item_from_bucket",
    "bucket_name": "my_bucket_name",
    "item_name": "my_item_name"
  },
  "data": null
}


```

### Delete All Items

delete all items from a bucket

Delete All Items:

| Metadata Key        | Required | Description                             | Possible values                      |
|:--------------------|:---------|:----------------------------------------|:-------------------------------------|
| method              | yes      | type of method                          | "delete_all_items_from_bucket"                        |
| bucket_name         | yes      | s3 bucket name                          | "my_bucket_name"                     |
| wait_for_completion | no       | wait for operation to end               | "true","false" (default of false )   |


Example:

```json
{
  "metadata": {
    "method": "delete_item_from_bucket",
    "bucket_name": "my_bucket_name"
  },
  "data": null
}


```

### Copy Item

copy an item from one bucket to another

Copy Items:

| Metadata Key        | Required | Description                             | Possible values                      |
|:--------------------|:---------|:----------------------------------------|:-------------------------------------|
| method              | yes      | type of method                          | "copy_item"                          |
| bucket_name         | yes      | s3 bucket name                          | "my_bucket_name"                     |
| copy_source         | yes      | s3 bucket name source name              | "my_bucket_source_name"              |
| item_name           | yes      | the name of the item                    | "valid-string"                       |
| wait_for_completion | no       | wait for operation to end               | "true","false" (default of false )   |


Example:

```json
{
  "metadata": {
    "method": "copy_item",
    "bucket_name": "my_bucket_name",
    "copy_source": "my_bucket_source_name",
    "item_name": "my_item_name"
  },
  "data": null
}


```

### Delete Bucket

delete a bucket by name.

Delete Bucket:

| Metadata Key        | Required | Description                             | Possible values                      |
|:--------------------|:---------|:----------------------------------------|:-------------------------------------|
| method              | yes      | type of method                          | "delete_bucket"                      |
| bucket_name         | yes      | s3 bucket name                          | "my_bucket_name"                     |
| wait_for_completion | no       | wait for operation to end               | "true","false" (default of false )   |


Example:

```json
{
  "metadata": {
    "method": "delete_bucket",
    "bucket_name": "my_bucket_name"
  },
  "data": null
}


```
