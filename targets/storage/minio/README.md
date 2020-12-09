# Kubemq Minio/S3 Target Connector

Kubemq Minio/S3 target connector allows services using kubemq server to access minio/s3 storage functions.

## Prerequisites
The following are required to run the minio target connector:

- kubemq cluster
- minio cluster / AWS s3 service
- kubemq-targets deployment

## Configuration

Mongodb target connector configuration properties:

| Properties Key    | Required | Description                              | Example          |
|:------------------|:---------|:-----------------------------------------|:-----------------|
| endpoint          | yes      | minio host address                     | "localhost:9001" |
| use_ssl           | no       | set connection ssl                       | "true"           |
| access_key_id     | yes      | set access key id                        | "minio"          |
| secret_access_key | yes      | set secret access key                    | "minio123"       |

Example:

```yaml
bindings:
  - name: kubemq-query-minio
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-minio-connector"
        auth_token: ""
        channel: "query.minio"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: stores.minio
      name: target-minio
      properties:
        endpoint: "localhost:9001"
        use_ssl: "true"
        access_key_id: "minio"
        secret_access_key: "minio123"
```

## Usage

### Make Bucket Request

Make bucket request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "make_bucket"   |
| param1       | yes      | set bucket name     | "bucket"        |
| param2       | no       | set bucket location | ""              |


Example:

```json
{
  "metadata": {
    "method": "make_bucket",
    "param1": "bucket",
    "param2": ""
  },
  "data": null
}
```

### List Buckets Request

List Buckets request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | method name              | "list_buckets"           |


Example:

```json
{
  "metadata": {
    "method": "list_buckets"
  },
  "data": null 
}
```

### Bucket Exists Request

Bucket exists request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "bucket_exists"   |
| param1       | yes      | set bucket name     | "bucket"        |

Example:

```json
{
  "metadata": {
    "method": "make_bucket",
    "param1": "bucket"
  },
  "data": null
}
```
### Remove Bucket Request

Remove bucket request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "remove_bucket"   |
| param1       | yes      | set bucket name     | "bucket"        |

Example:

```json
{
  "metadata": {
    "method": "remove_bucket",
    "param1": "bucket"
  },
  "data": null
}
```
### List Objects Request

List objects request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "list_objects"   |
| param1       | yes      | set bucket name     | "bucket"        |

Example:

```json
{
  "metadata": {
    "method": "list_objects",
    "param1": "bucket"
  },
  "data": null
}
```

### Put Object Request

Put object request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "put"   |
| param1       | yes      | set bucket name     | "bucket"        |
| param2       | yes      | set object name     | "object"        |

Put request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to put for object | base64 bytes array |

Example:

```json
{
  "metadata": {
    "method": "remove_bucket",
    "param1": "bucket"
  },
  "data": "c29tZS1kYXRh"
}
```

### Get Object Request

Get object request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "get"   |
| param1       | yes      | set bucket name     | "bucket"        |
| param2       | yes      | set object name     | "object"        |

Example:

```json
{
  "metadata": {
    "method": "get",
    "param1": "bucket",
    "param2": "object"
  },
  "data": null
}
```


### Remove Object Request

Remove object request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "remove"   |
| param1       | yes      | set bucket name     | "bucket"        |
| param2       | yes      | set object name     | "object"        |

Example:

```json
{
  "metadata": {
    "method": "remove",
    "param1": "bucket",
    "param2": "object"
  },
  "data": null
}
```
