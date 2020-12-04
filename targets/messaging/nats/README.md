# Kubemq Nats Target Connector

Kubemq Nats target connector allows services using kubemq server to access Nats messaging services.

## Prerequisites
The following are required to run the Nats target connector:

- kubemq cluster
- Nats server
- kubemq-targets deployment

## Configuration

Nats target connector configuration properties:

| Properties Key                  | Required | Description                                             | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------------------|:-----------------------------------------------------------------------|
| url                             | yes      | nats connection host                                    | "localhost:1883" |
| subject                         | yes      | set subject name                                        | any string |
| username                        | no       | set nats username                                       | "username" |
| password                        | no       | set nats password                                       | "password" |
| token                           | no       | set nats token                                          | "my_token" |
| tls                             | no       | set if tls is needed                                    | "false","true" |
| cert_file                       | no       | tls certificate file in string format                   | "my_file" |
| cert_key                        | no       | tls certificate key in string format                    | "my_key"  |
| timeout                         | no       | connection timeout in seconds                           | "130"  |


Example:

```yaml
bindings:
  - name: nats
    source:
      kind: kubemq.events
      properties:
        address: localhost:50000
        channel: event.messaging.nats
    target:
      kind: messaging.nats
      properties:
        cert_file: |-
          -----BEGIN CERTIFICATE-----
          -----END CERTIFICATE-----
        cert_key: |-
          -----BEGIN PRIVATE KEY-----
          -----END PRIVATE KEY-----
        dynamic_mapping: "false"
        url: nats://localhost:4222
    properties: {}

```

## Usage

### Request


Query request data setting:

| Data Key          | Required | Description                               | Possible values    |
|:------------------|:---------|:------------------------------------------|:-------------------|
| data              | yes      | data to publish                           | base64 bytes array |

Example:


```json
{
  "data": "U0VMRUNUIGlkLHRpdGxlLGNvbnRlbnQgRlJPTSBwb3N0Ow=="
}
```
