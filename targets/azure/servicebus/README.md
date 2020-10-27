# Kubemq servicebus Target Connector

Kubemq servicebus target connector allows services using kubemq server to access servicebus messaging services.

## Prerequisites
The following are required to run the servicebus target connector:

- kubemq cluster
- Azure active with servicebus enable 
- kubemq-targets deployment



## Configuration

servicebus target connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| end_point                       | yes      | event hubs target endpoint                  | "sb://my_account.net" |
| shared_access_key_name          | yes      | event hubs access key name                  | "keyname" |
| shared_access_key               | yes      | event hubs shared access key name           | "213ase123" |


Example:

```yaml
bindings:
  - name: kubemq-query-azure-servicebus
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-azure-servicebus-connector"
        auth_token: ""
        channel: "azure.servicebus"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind:azure.servicebus
      name: target-azure-servicebus
      properties:
        end_point: "sb://my_account.net"
        shared_access_key_name: "keyname"
        shared_access_key: "213ase123"
        queue_name: "mypath"
```

## Usage

### send

send metadata setting:

| Metadata Key      | Required | Description                                    | Possible values                                  |
|:------------------|:---------|:-----------------------------------------------|:-------------------------------------------------|
| method            | yes      | type of method                                 | "send"                                         |
| label             | no       | the message label                              | "my_label"                                     |
| content_type      | no       | message content type                           | "content_type"                               |
| time_to_live      | no       | message time to live                           | "1000000000"default(1000000000)          |



Example:

```json
{
  "metadata": {
    "method": "send",
    "label": "my_label"
  },
  "data": "bXktZmlsZS1kYXRh"
}
```

### send batch

send batch metadata setting:

| Metadata Key                   | Required | Description                                     | Possible values                            |
|:-------------------------------|:---------|:------------------------------------------------|:-------------------------------------------|
| method                         | yes      | type of method                                  | "send_batch"                                  |
| label                          | no       | the message label                              | "my_label"                                     |
| content_type                   | no       | message content type                           | "content_type"                               |
| time_to_live                   | no       | message time to live                           | "1000000000"default(1000000000)          |


Example:

```json
{
  "metadata": {
    "method": "send_batch",
        "label": "my_label"
  },
  "data": "WyJ0ZXN0MSIsInRlc3QyIiwidGVzdDMiLCJ0ZXN0NCJd"
}
``````
