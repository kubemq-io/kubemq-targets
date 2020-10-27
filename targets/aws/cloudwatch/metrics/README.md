# Kubemq cloudwatch-metrics target Connector

Kubemq cloudwatch-metrics target connector allows services using kubemq server to access aws cloudwatch-metrics service.

## Prerequisites
The following required to run the aws-cloudwatch-metrics target connector:

- kubemq cluster
- aws account with cloudwatch-metrics active service
- kubemq-source deployment

## Configuration

cloudwatch-metrics target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-cloudwatch-metrics
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-aws-cloudwatch-metrics"
        auth_token: ""
        channel: "query.aws.cloudwatch.metrics"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.cloudwatch.metrics
      name: aws-cloudwatch-metrics
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "region"
        token: ""
```

## Usage

### Put Metrics

Put Metrics:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "put_metrics"                     |
| namespace         | yes      | aws namespace name                      | "string"                     |
| data              | yes      | array of aws MetricDatum                |  `W3siQ291bnRzIjpudWxsLCJEaW1lbnNpb25zIjpudWxsLCJNZXRyaWNOYW1lIjoiTmV3IE1ldHJpYyIsIlN0YXRpc3RpY1ZhbHVlcyI6bnVsbCwiU3RvcmFnZVJlc29sdXRpb24iOm51bGwsIlRpbWVzdGFtcCI6IjIwMjAtMDgtMTJUMTc6MDk6NDguMzg5NTgyMiswMzowMCIsIlVuaXQiOiJDb3VudCIsIlZhbHVlIjoxMzEsIlZhbHVlcyI6bnVsbH1d`     |



Example:

```json
{
  "metadata": {
    "method": "put_metrics",
    "namespace": "Logs"
  },
  "data": "W3siQ291bnRzIjpudWxsLCJEaW1lbnNpb25zIjpudWxsLCJNZXRyaWNOYW1lIjoiTmV3IE1ldHJpYyIsIlN0YXRpc3RpY1ZhbHVlcyI6bnVsbCwiU3RvcmFnZVJlc29sdXRpb24iOm51bGwsIlRpbWVzdGFtcCI6IjIwMjAtMDgtMTJUMTc6MDk6NDguMzg5NTgyMiswMzowMCIsIlVuaXQiOiJDb3VudCIsIlZhbHVlIjoxMzEsIlZhbHVlcyI6bnVsbH1d"
}
```


### List Metrics 

List Metrics:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_metrics"                     |
| namespace         | no       | aws namespace name                      | "string"                     |


Example:

```json
{
  "metadata": {
    "method": "list_metrics"
  },
  "data": null
}
```
