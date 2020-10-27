# Kubemq sns target Connector

Kubemq aws-sns target connector allows services using kubemq server to access aws sns service.

## Prerequisites
The following required to run the aws-sns target connector:

- kubemq cluster
- aws account with sns active service
- kubemq-source deployment

## Configuration

sns target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string)         | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-sns
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-sns-connector"
        auth_token: ""
        channel: "query.aws.sns"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.sns
      name: aws-sns
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "instance"
        token:   ""
```

## Usage

### List Topics

list all topics

List Topics:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_topics"                     |


Example:

```json
{
  "metadata": {
    "method": "list_topics"
  },
  "data": null
}
```



### List Subscriptions

list all subscriptions

List Subscriptions:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_subscriptions"                     |


Example:

```json
{
  "metadata": {
    "method": "list_subscriptions"
  },
  "data": null
}
```


### List Subscriptions By Topic

list all Subscriptions of the selected topic

List Subscriptions By Topic:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_subscriptions_by_topic"       |
| topic             | yes      | topic_name                              | "arn:aws-my-topic"                     |


Example:

```json
{
  "metadata": {
    "method": "list_subscriptions_by_topic",
    "topic": "arn:aws-my-topic"
  },
  "data": null
}
```


### Create Topic

create a new topic , topic name must by unique 

Create Topic:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "create_topic"       |
| topic             | yes      | topic_name                              | "arn:aws-my-topic"                     |
| data              | no       | create attributes as base64             | '{"DisplayName":"my-display-name"}' 


Example:

```json
{
  "metadata": {
    "method": "create_topic",
    "topic": "arn:aws-my-topic"
  },
  "data": "eyJEaXNwbGF5TmFtZSI6Im15LWRpc3BsYXktbmFtZSJ9"
}
```


### Send Message

send a message to topic.

Send Message:

| Metadata Key      | Required                            | Description                  | Possible values                            |
|:------------------|:------------------------------------|:-----------------------------|:-------------------------------------------|
| method            | yes                                 | type of method               | "send_message"                           |
| topic             | no(unless target_arn is missing)    | topic_name                   | "arn:aws-my-topic"                     |
| target_arn        | no(unless topic is missing)         | target_arn                   | "arn:aws-my-topic"                     |
| message           | yes                                 | message body as string       | 'some message in string format'        || 
| message           | yes                                 | message body as string       | 'some message in string format'        || 
| subject           | no                                  | sns subject name             | "string name of sns subject"           || 
| phone_number      | no                                  | valid phone number           | "valid phone number"                   || 
| data              | no                                  | message attributes as base64 | "[{"name":"store","data_type":"String","string_value":"my_store"},{"name":"event","data_type":"String","string_value":"my_event"}]"                  || 


Example:

```json
{
  "metadata": {
    "method": "send_message",
    "topic": "arn:aws-my-topic"
    "message": "my message to send"
  },
  "data": "W3sibmFtZSI6InN0b3JlIiwiZGF0YV90eXBlIjoiU3RyaW5nIiwic3RyaW5nX3ZhbHVlIjoibXlfc3RvcmUifSx7Im5hbWUiOiJldmVudCIsImRhdGFfdHlwZSI6IlN0cmluZyIsInN0cmluZ192YWx1ZSI6Im15X2V2ZW50In1d"
}
```


### Subscribe

Subscribe to topic

Subscribe:

| Metadata Key      | Required | Description                             | Possible values                                           |
|:--------------------|:------------------------------------|:---------------------------------|:--------------------------------------|
| method              | yes                                 | type of method                   | "subscribe"                           |
| topic               | yes    | topic_name                 | "arn:aws-my-topic"               | "arn:aws-my-topic"   |
| data                | no                                  | Subscribe attributes as base64   | '{"store": ["mystore"],"event": [{"anything-but": "my-event"}]}'       


Example:

```json
{
  "metadata": {
    "method": "subscribe",
    "topic": "arn:aws-my-topic"
  },
  "data": "eyJzdG9yZSI6IFsibXlzdG9yZSJdLCJldmVudCI6IFt7ImFueXRoaW5nLWJ1dCI6ICJteS1ldmVudCJ9XX0="
}
```

### Delete Topic

delete the selected topic

Delete Topic:

| Metadata Key      | Required | Description                             | Possible values                                           |
|:--------------------|:------------------------------------|:-------------------------------|:--------------------------------------|
| method              | yes                                 | type of method                 | "delete_topic"                           |
| topic               | yes                                 | topic_name                     | "arn:aws-my-topic"         | 


Example:

```json
{
  "metadata": {
    "method": "delete_topic",
    "topic": "arn:aws-my-topic"
  },
  "data": null
}
```
