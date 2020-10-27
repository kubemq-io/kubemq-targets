# Kubemq cloudwatch-logs target Connector

Kubemq cloudwatch-logs target connector allows services using kubemq server to access aws cloudwatch-logs service.

## Prerequisites
The following required to run the aws-cloudwatch-logs target connector:

- kubemq cluster
- aws account with cloudwatch-logs active service
- some action will need cloudwatch-logs permission (IAM User)
- kubemq-source deployment

## Configuration

cloudwatch-logs target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-cloudwatch-logs
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-cloudwatch-logs"
        auth_token: ""
        channel: "query.aws.cloudwatch.logs"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.cloudwatch.logs
      name: aws-cloudwatch-logs
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "region"
        token: ""
```

## Usage

### Create log Stream 

create a new log stream

Create log Stream:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "create_log_event_stream"                     |
| log_stream_name   | yes      | aws log stream name                     | "string"                     |
| log_group_name    | yes      | aws log group name                      | "string"                     |



Example:

```json
{
  "metadata": {
    "method": "create_log_event_stream",
    "log_stream_name": "my_stream_name",
    "log_group_name": "my_group_name"
  },
  "data": null
}
```


### Describe log Stream 

describe a selected log stream by group_name

Describe log Stream:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "describe_log_event_stream"                     |
| log_group_name    | yes      | aws log group name                      | "string"                     |



Example:

```json
{
  "metadata": {
    "method": "describe_log_event_stream",
    "log_group_name": "my_group_name"
  },
  "data": null
}
```

### Delete log Stream 

delete log stream

Delete log Stream:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "delete_log_event_stream"                     |
| log_stream_name   | yes      | aws log stream name                     | "string"                     |
| log_group_name    | yes      | aws log group name                      | "string"                     |



Example:

```json
{
  "metadata": {
    "method": "delete_log_event_stream",
    "log_stream_name": "my_stream_name",
    "log_group_name": "my_group_name"
  },
  "data": null
}
```

### Get log Event

get log event 

Get log Stream:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:-----------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                                 | "get_log_event"                     |
| log_stream_name   | yes      | aws log stream name                            | "string"                     |
| log_group_name    | yes      | aws log group name                             | "string"                                                 |



Example:

```json
{
  "metadata": {
    "method": "get_log_event",
    "log_stream_name": "my_stream_name",
    "log_group_name": "my_group_name"
  },
  "data": null
}
```

### Create Log Event Group

create a new log event group

Create Log Event Group:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:-----------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                                 | "create_log_group"                     |
| log_group_name    | yes      | aws log group name                             | "string"                                                 |
| data              | no       | aws tags                                       | key value pair string string                                                 |



Example:

```json
{
  "metadata": {
    "method": "create_log_group",
    "log_group_name": "my_group_name"
  },
  "data": null
}
```

### Put Log 

put a log in log stream

Put Log Event:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:---------------------------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                                                 | "put_log_event"                     |
| log_stream_name   | yes      | aws log stream name                                            | "string"                     |
| log_group_name    | yes      | aws log group name                                             | "string"                     |
| sequence_token    | yes      | aws stream sequence token                                      | "string"                     |
| data              | yes      | key value pair of int-string int-time - string-Message         | "string"                     |



Example:

```json
{
  "metadata": {
    "method": "put_log_event",
    "log_group_name": "my_group_name",
    "sequence_token": "my_token_from_aws"
  },
  "data": "eyIxNTk3MjM1NTU4NTEyIjoibXkgZmlyc3QgbWVzc2FnZSB0byBzZW5kIiwiMTU5NzIzNTU1ODUyNyI6Im15IHNlY29uZCBtZXNzYWdlIHRvIHNlbmQifQ=="
}
```

### Describe Log Event Group

describe log event group

Describe Log Event Group:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:-----------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                                 | "describe_log_group"                     |
| log_group_prefix  | yes      | aws log group prefix                           | "string"                                                 |



Example:

```json
{
  "metadata": {
    "method": "describe_log_group",
    "log_group_name": "my_group_name"
  },
  "data": null
}
```


### Delete Log Event Group

delete log event group

Delete Log Event Group:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:-----------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                                 | "delete_log_group"                     |
| log_group_name    | yes      | aws log group name                              | "string"                                                 |



Example:

```json
{
  "metadata": {
    "method": "delete_log_group",
    "log_group_name": "my_group_name"
  },
  "data": null
}
```
