# Kubemq amazonMQ Target Connector

Kubemq AmazonMQ target connector allows services using kubemq server to access AmazonMQ messaging services.

## Prerequisites
The following are required to run the AmazonMQ target connector:

- kubemq cluster
- AmazonMQ server - with access 
- kubemq-targets deployment


- Please note the connector uses connection with stomp+ssl, when finishing handling messages need to call Close().

## Configuration

AmazonMQ target connector configuration properties:

| Properties Key                  | Required | Description                                  | Example                                                                |
|:--------------------------------|:---------|:---------------------------------------------|:-----------------------------------------------------------------------|
| host                            | yes      | AmazonMQ connection host (stomp+ssl endpoint)| "localhost:1883" |
| username                        | no       | set AmazonMQ username                        | "username" |
| password                        | no       | set AmazonMQ password                        | "password" |


Example:

```yaml
bindings:
  - name: kubemq-query-amazonmq
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-amazonmq-connector"
        auth_token: ""
        channel: "query.amazonmq"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.amazonmq
      name: aws-amazonmq
      properties:
        host: "localhost:61613"
        username: "admin"
        password: "admin"
```

## Usage

### Request

Request metadata setting:

| Metadata Key   | Required | Description         | Possible values |
|:---------------|:---------|:--------------------|:----------------|
| destination    | yes      | set destination name| "destination"         |



Query request data setting:

| Data Key | Required | Description  | Possible values    |
|:---------|:---------|:-------------|:-------------------|
| data     | yes      | data to publish | base64 bytes array |

Example:


```json
{
  "metadata": {
    "destination": "destination"
  },
  "data": "U0VMRUNUIGlkLHRpdGxlLGNvbnRlbnQgRlJPTSBwb3N0Ow=="
}
```
