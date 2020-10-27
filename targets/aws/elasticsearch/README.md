# Kubemq elasticsearch target Connector

Kubemq aws-elasticsearch target connector allows services using kubemq server to access aws elasticsearch service.

## Prerequisites
The following required to run the aws-elasticsearch target connector:

- kubemq cluster
- aws account with elasticsearch active service -elastic service with an active domain
- kubemq-source deployment

## Configuration

aws-elasticsearch target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| token          | no       | aws token ("default" empty string          | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-elasticsearch
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-elasticsearch-connector"
        auth_token: ""
        channel: "query.aws.elasticsearch"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.elasticsearch
      name: aws-elasticsearch
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        token: ""
```

## Usage

### Sign Message 

Sign Message :

| Metadata Key      | Required                 | Description                                                 | Possible values                            |
|:------------------|:-------------------------|:------------------------------------------------------------|:-------------------------------------------|
| method            | yes                      | type of HTTP method                                         | "GET", "POST","PUT","DELETE","OPTIONS"                 |
| region            | yes                      | aws region associated with domain                           | "region"                                                 |
| json              | yes (unless "GET")       | json body to send with the http request                     | "list"                                                 |
| domain            | yes                      | elastic domain to assign the request                        | "list"                                                 |
| index             | yes                      | name of the elastic index                                   | "list"                                                 |
| endpoint          | yes                      | aws domain end point                                        | "list"                                                 |
| service           | no(Default "es"          | type of service                                             | "list"                                                 |
| id                | yes                      | Message ID                                                  | "list"                                                 |



Example:

```json
{
  "metadata": {
    "method": "GET",
    "region": "us-west-2",
    "domain": "https://my-domain-12345asdfg.us-west-2.es.amazonaws.com",
    "index": "myindex",
    "endpoint": "https://my-domain-12345asdfg.us-west-2.es.amazonaws.com/my/end_point",
    "service": "es",
    "id": "123124"
  },
  "data": null
}
```
