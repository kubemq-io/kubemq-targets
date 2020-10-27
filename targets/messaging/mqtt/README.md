# Kubemq MQTT Target Connector

Kubemq mqtt target connector allows services using kubemq server to access mqtt messaging services.

## Prerequisites
The following are required to run the mqtt target connector:

- kubemq cluster
- mqtt server
- kubemq-targets deployment

## Configuration

MQTT target connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| host                      | yes      | mqtt connection host          | "localhost:1883" |
| username                      | no      | set mqtt username          | "username" |
| password                      | no      | set mqtt password          | "password" |
| client_id                      | no      | mqtt connection string address          | "client_id" |

Example:

```yaml
bindings:
  - name: kubemq-query-mqtt
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-mqtt-connector"
        auth_token: ""
        channel: "query.mqtt"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind:messaging.mqtt
      name: target-mqtt
      properties:
        host: "localhost:1883"
        username: "username"
        password: "password"
        client_id: "client_id"
```

## Usage

### Request

Request metadata setting:

| Metadata Key   | Required | Description         | Possible values |
|:---------------|:---------|:--------------------|:----------------|
| topic          | yes      | set topic name | "topic"         |
| qos       | yes      | set qos level | "0","1","2"         |


Query request data setting:

| Data Key | Required | Description  | Possible values    |
|:---------|:---------|:-------------|:-------------------|
| data     | yes      | data to publish | base64 bytes array |

Example:


```json
{
  "metadata": {
    "topic": "topic",
    "qos": "0"
  },
  "data": "U0VMRUNUIGlkLHRpdGxlLGNvbnRlbnQgRlJPTSBwb3N0Ow=="
}
```
