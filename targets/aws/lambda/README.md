# Kubemq lambda target Connector

Kubemq aws-lambda target connector allows services using kubemq server to access aws lambda service.

## Prerequisites
The following required to run the aws-lambda target connector:

- kubemq cluster
- aws account with lambda active service
- kubemq-source deployment

## Configuration

lambda target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-lambda
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-lambda-connector"
        auth_token: ""
        channel: "query.aws.lambda"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.lambda
      name: aws-lambda
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "region"
        token: ""
```

## Usage

### List Lambda

List all lambdas

List Lambda:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list"                     |



Example:

```json
{
  "metadata": {
    "method": "list"
  },
  "data": null
}
```

### Create Lambda

create a new lambda.

Create Lambda:

| Metadata Key      | Required | Description                                     | Possible values                            |
|:------------------|:---------|:------------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                                  | "create"                     |
| zip_file_name     | yes      | name of the zip file                            | "file.zip"                     |
| handler_name      | yes      | lambda handler name                             | "handler-path"                     |
| role              | yes      | aws role name                                   | "arn:aws:iam::0000000:myRole"                     |
| runtime           | yes      | lambda runtime version                          | see https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html|
| function_name     | yes      | lambda function name                            | string                |
| data              | yes      | the function code , in byte array               | byte array            |
| memory_size       | no       | memory_size needed default of 256               | int                   |
| timeout           | no       | timeout set for task default of 15 (seconds)    | int                   |
| description       | no       | function description default of ""              | string                |



Example:

```json
{
  "metadata": {
    "method": "create",
    "zip_file_name": "myfile.zip",
    "handler_name": "myhandler",
    "role": "arn:aws:iam::0000000:myRole",
    "runtime": "nodejs12.x",
    "function_name": "testfunction",
    "memory_size": "256",
    "timeout": "3",
    "description": "my awesome testing method"
  },
  "data": "ZXhwb3J0cy5oYW5kbGVyID0gYXN5bmMgKGV2ZW50KSA9PiB7CiAgICAvLyBUT0RPIGltcGxlbWVudAogICAgY29uc3QgcmVzcG9uc2UgPSB7CiAgICAgICAgc3RhdHVzQ29kZTogMjAwLAogICAgICAgIGJvZHk6IEpTT04uc3RyaW5naWZ5KCdIZWxsbyBmcm9tIExhbWJkYSEnKSwKICAgIH07CiAgICByZXR1cm4gcmVzcG9uc2U7Cn07Cg=="
}
```

### Run Lambda

run a specific lambda

Run Lambda:

| Metadata Key      | Required | Description                                     | Possible values                            |
|:------------------|:---------|:------------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                                  | "run"                     |
| function_name     | yes      | lambda function name                            | string                |
| data              | yes      | the run request code , in byte array            | byte array            |



Example:

```json
{
  "metadata": {
    "method": "run",
    "function_name": "testfunction",
  },
  "data": "bXkgb2JqZWN0"
}
```

### Delete Lambda

Delete Lambda:

| Metadata Key      | Required | Description                                     | Possible values                            |
|:------------------|:---------|:------------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                                  | "delete"                     |
| function_name     | yes      | lambda function name                            | string                |



Example:

```json
{
  "metadata": {
    "method": "delete",
    "function_name": "testfunction"
  },
  "data": null
}
```
