# Kubemq sqs target Connector

Kubemq aws-sqs target connector allows services using kubemq server to access aws sqs service.

## Prerequisites
The following required to run the aws-sqs target connector:

- kubemq cluster
- aws account with sqs active service
- kubemq-source deployment

## Configuration

sqs target connector configuration properties:

| Properties Key | Required | Description                                                       | Example                     |
|:---------------|:---------|:------------------------------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                                           | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                                                    | aws secret key supplied by aws  |
| region         | yes      | region                                                            | aws region                      |
| retries        | no       | number of retries on send                                         | 1 (default 0)                   |
| token          | no       | aws token ("default" empty string                                 | "my token"                      |
| dead_letter    | no       | dead letter queue name (only relevant to SetQueueAttributes)      | "my_dead_letter_queue"          |
| max_receive    | no       | max receive of queue (only relevant to SetQueueAttributes)        | "0"                              |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-sqs
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-sqs-connector"
        auth_token: ""
        channel: "query.aws.sqs"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.sqs
      name: aws-sqs
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "instance"
        token:  ""
        retries: "1"
```

## Usage

### Send Message

send message to sqs.

Send Message:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| queue             | yes      | name of queue to send                   | "my_queue"                           |
| delay             | yes      | message delay                           | "0"                                  |
| tags              | no       | message tags (key value string string)  | "{"tag-1":"test","tag-2":"test2"}"   |
| data              | yes      | type of method                          | "dmFsaWQgYm9keQ=="                        |


Example:

```json
{
  "metadata": {
    "queue": "my_queue",
    "delay": "0",
    "tags": "{\"tag-1\":\"test\",\"tag-2\":\"test2\"}"
  },
  "data": "dmFsaWQgYm9keQ=="
}
```



