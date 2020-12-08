# Kubemq hazelcast Target Connector

Kubemq hazelcast target connector allows services using kubemq server to access hazelcast server functions such `set`, `get` and `delete`.

## Prerequisites
The following are required to run the hazelcast target connector:

- kubemq cluster
- hazelcast
- kubemq-targets deployment

## Configuration

hazelcast target connector configuration properties:

| Properties Key            | Required| Description                  | Example          |
|:--------------------------|:--------|:-----------------------------|:-----------------|
| address                   | yes     | hazelcast connection string                     | "localhost:5701" |
| username                  | no      | hazelcast username                              | "admin" |
| password                  | no      | hazelcast password                              | "password" |
| connectionAttemptLimit    | no      | hazelcast connection attempts(default 1)        | 1                            |
| connectionAttemptPeriod   | no      | hazelcast attempt period seconds(default 5)     | 5 |
| connectionTimeout         | no      | hazelcast connection timeout seconds(default 5) | 5 |
| ssl                       | no      | hazelcast use ssl                               | false |
| sslcertificatefile        | no      | hazelcast certificate file                      | "" |
| sslcertificatekey         | no      | hazelcast certificate key                       | "" |
| serverName                | no      | hazelcast server name                           | "myserver" |

Example:

```yaml
bindings:
  - name: kubemq-hazelcast
    source:
      kind: kubemq.query
      properties:
        address: localhost:50000
        channel: query.hazelcast
    target:
      kind: cache.hazelcast
      properties:
        address: localhost:5701
        server_name: test
    properties: {}

```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description          | Possible values |
|:-------------|:---------|:---------------------|:----------------|
| key          | yes      | hazelcast key string | any string      |
| method       | yes      | get                  | "get"           |
| map_name     | yes      | hazelcast map name   | "my_map"        |


Example:

```json
{
  "metadata": {
    "key": "your-hazelcast-key",
    "map_name": "my_map",
    "method": "get"
  },
  "data": null
}
```

### Set Request

Set request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | hazelcast key    | any string      |
| method       | yes      | set              | "set"           |
| map_name     | yes      | hazelcast map name   | "my_map"        |

Set request data setting:

| Data Key | Required | Description                       | Possible values     |
|:---------|:---------|:----------------------------------|:--------------------|
| data     | yes      | data to set for the hazelcast key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "key": "your-hazelcast-key",
    "map_name": "my_map",
    "method": "set"
  },
  "data": "c29tZS1kYXRh" 
}
```
### Delete Request

Delete request metadata setting:

| Metadata Key | Required | Description          | Possible values |
|:-------------|:---------|:---------------------|:----------------|
| key          | yes      | hazelcast key string | any string      |
| method       | yes      | delete               | "delete"        |
| map_name     | yes      | hazelcast map name   | "my_map"        |


Example:

```json
{
  "metadata": {
    "key": "your-hazelcast-key",
    "map_name": "my_map",
    "method": "delete"
  },
  "data": null
}
```
