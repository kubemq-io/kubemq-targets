# Kubemq ActiveMQ Target Connector

Kubemq activemq target connector allows services using kubemq server to access activemq messaging services.

## Prerequisites
The following are required to run the activemq target connector:

- kubemq cluster
- activemq server
- kubemq-targets deployment

## Configuration

ActiveMQ target connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| host                      | yes      | activemq connection host          | "localhost:1883" |
| username                      | no      | set activemq username          | "username" |
| password                      | no      | set activemq password          | "password" |


Example:

```yaml
bindings:
  - name: kubemq-query-activemq
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-activemq-connector"
        auth_token: ""
        channel: "query.activemq"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind:messaging.activemq
      name: target-activemq
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
| destination          | yes      | set destination name | "destination"         |



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
