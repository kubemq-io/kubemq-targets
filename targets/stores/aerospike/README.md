# Kubemq Aerospike Target Connector

Kubemq Aerospike target connector allows services using kubemq server to access Aerospike database services.

## Prerequisites
The following are required to run the Aerospike target connector:

- kubemq cluster
- Aerospike server
- kubemq-targets deployment

## Configuration

Aerospike target connector configuration properties:

| Properties Key            | Required | Description                          | Example                   |
|:--------------------------|:---------|:-------------------------------------|:--------------------------|
| host                      | yes      | Aerospike host address               | "localhost"         |
| port                      | yes      | Aerospike host port                  | "3000"         |
| username                  | no       | Aerospike username                   | "admin"                   |
| password                  | no       | Aerospike password                   | "password"                |
| timeout                   | no       | set  timeout in seconds              | "30"                      |



Example:

```yaml
bindings:
- name: aerospike
  source:
    kind: kubemq.query
    properties:
      address: localhost:50000
      channel: query.aerospike
  target:
    kind: stores.aerospike
    properties:
      host: localhost
      port: "3000"
  properties: {}

```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      |  key string      | any string      |
| user_key     | yes      |  user key string | any string      |
| namespace    | yes      |  namespace name  | any string      |
| method       | yes      | get              | "get"           |

Example:

```json
{
  "metadata": {
    "key": "your-Aerospike-key",
    "user_key": "your-user_key",
    "namespace": "test",
    "method": "get"
  },
  "data": null
}
```

### Set Request

Set request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | Aerospike key string | any string      |
| method       | yes      | set              | "set"           |

Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the Aerospike key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "method": "set"
  },
  "data": "eyJiaW5fbWFwIjp7ImJpbjEiOjQyLCJiaW4yIjoiQW4gZWxlcGhhbnQgaXMgYSBtb3VzZSB3aXRoIGFuIG9wZXJhdGluZyBzeXN0ZW0iLCJiaW4zIjpbIkdvIiwyMDA5XX0sImtleV9uYW1lIjoic29tZS1rZXkiLCJuYW1lc3BhY2UiOiJ0ZXN0IiwidXNlcl9rZXkiOiJ1c2VyX2tleTEifQ==" 
}
```
### Delete Request

Delete request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      |  key string      | any string      |
| user_key     | yes      |  user key string | any string      |
| namespace    | yes      |  namespace name  | any string      |
| method       | yes      |  delete          | "delete"        |


Example:

```json
{
  "metadata": {
    "key": "your-Aerospike-key",
    "user_key": "your-user_key",
    "namespace": "test",
    "method": "delete"
  },
  "data": null
}
```


