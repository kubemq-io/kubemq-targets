# Kubemq cloudwatch-events target Connector

Kubemq cloudwatch-events target connector allows services using kubemq server to access aws cloudwatch-events service.

## Prerequisites
The following required to run the aws-cloudwatch-events target connector:

- kubemq cluster
- aws account with cloudwatch-events active service (IAM Permission under EventBridge)
- kubemq-source deployment

## Configuration

cloudwatch-events target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-cloudwatch-events
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-cloudwatch-events"
        auth_token: ""
        channel: "query.aws.cloudwatch.events"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.cloudwatch.events
      name: aws-cloudwatch-events
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "region"
        token: ""
```

## Usage

### Put Target

Put Target:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "put_targets"                     |
| rule              | yes      | aws existing rule name                  | "string"                     |
| data              | yes      | Key value pair of target ARN and ID     |  `{"my_arn_id":"arn:aws:test:number:function:id"}`     |



Example:

```json
{
  "metadata": {
    "method": "put_targets",
    "rule": "my_rule"
  },
  "data": "eyJteV9hcm5faWQiOiJhcm46YXdzOnRlc3Q6bnVtYmVyOmZ1bmN0aW9uOmlkIn0"
}
```


### List Event Buses

List Event Buses:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_buses"                     |
| limit             | no       | limit of return buses                   | "int"                     |



Example:

```json
{
  "metadata": {
    "method": "list_buses",
    "limit": "1"
  },
  "data": null
}
```


### Send Event

Send Event:

| Metadata Key      | Required | Description                             | Possible values                            |
|:-----------------|:---------|:----------------------------------------|:-------------------------------------------|
| method           | yes      | type of method                          | "send_event"                     |
| detail           | no       | general details on the event            | "string"                     |
| detail_type      | no       | event type                              | "string"                     |
| source           | no       | aws source to assign the message to     | "string"                     |
| data             | yes      | aws resources to assign the event to    |  slice of strings ("arn:string")                     |



Example:

```json
{
  "metadata": {
    "method": "send_event",
    "detail": "{ some detail }",
    "detail_type": "appRequestSubmitted",
    "source": "kubemq_testing"
  },
  "data": "WyJhcm46YXdzOnNpdGU6cmVnaW9uOmlkOmZ1bmN0aW9uOm15LWZ1bmN0aXNvbnMiXQ=="
}
```
