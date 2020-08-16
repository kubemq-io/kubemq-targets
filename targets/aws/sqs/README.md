# Kubemq sqs target Connector

Kubemq aws-sqs target connector allows services using kubemq server to access aws sqs service.

## Prerequisites
The following required to run the aws-sqs target connector:

- kubemq cluster
- aws account with sqs active service
- kubemq-source deployment

## Configuration

sqs target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-sqs
    source:
      kind: source.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-aws-sqs-connector"
        auth_token: ""
        channel: "query.aws.sqs"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.aws.sqs
      name: target-aws-sqs
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "instance"
        max_retries: "0"
        max_receive: '10'
        dead_letter:  "my_dead_letter_queue"
        max_retries_backoff_seconds: "0"
```

## Usage

### Do Topics

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

