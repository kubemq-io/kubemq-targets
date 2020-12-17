# Kubemq Consulkv Target Connector

Kubemq consulkv target connector allows services using kubemq server to access consulkv database services.

## Prerequisites
The following are required to run the consulkv target connector:

- kubemq cluster
- consulkv server
- kubemq-targets deployment

## Configuration

Consulkv target connector configuration properties:

| Properties Key            | Required | Description                                   | Example                   |
|:--------------------------|:---------|:----------------------------------------------|:--------------------------|
| address                   | yes      | consulkv host address                         | "localhost:8500"          |
| scheme                    | no       | consulkv scheme(empty will use default)       | "my_scheme"          |
| data_center               | no       | consulkv data center(empty will use default)  | "my_center"          |
| token                     | no       | consulkv token for all requests               | ""          |
| token_file                | no       | consulkv token file                           | ""          |
| insecure_skip_verify      | no       | if true will disable TLS host verification    | "false"          |
| wait_time                 | no       | WaitTime limits how long a Watch will block   | "36000"          |
| tls                       | no       | consulkv use tls                              | false |
| cert_file                 | no       | consulkv certificate file                     | "" |
| cert_key                  | no       | consulkv certificate key                      | "" |



Example:

```yaml
bindings:
- name: consulkv
  source:
    kind: kubemq.query
    properties:
      address: localhost:50000
      channel: query.consulkv
  target:
    kind: stores.consulkv
    properties:
      address: localhost:8500
  properties: {}

```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key          | Required | Description                                             | Possible values |
|:----------------------|:---------|:--------------------------------------------------------|:----------------|
| key                   | yes      | key to get                                              | any string      |
| method                | yes      | method name                                             | "get"           |
| allow_stale           | no       | allows any Consul server (non-leader) to service a read | "false"           |
| require_consistent    | no       | RequireConsistent forces the read to be fully consistent| "false"           |
| use_cache             | no       | UseCache requests that the agent cache results locally  | "false"           |
| max_age               | no       | MaxAge limits how old a cached value will be returned   | "36000"           |
| max_age               | no       | MaxAge limits how old a cached value will be returned   | "36000"           |
| max_age               | no       | MaxAge limits how old a cached value will be returned   | "36000"           |

Example:

```json
{
  "metadata": {
    "key": "your-consulkv-key",
    "method": "get"
  },
  "data": null
}
```

### Put Request

Put request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | method name      | "put"           |


Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | struct of consul.KVPair class | base64 bytes array |

Example:

```json
{
  "metadata": {
    "method": "set"
  },
  "data": "eyJLZXkiOiJzb21lLWtleSIsIkNyZWF0ZUluZGV4IjowLCJNb2RpZnlJbmRleCI6MCwiTG9ja0luZGV4IjowLCJGbGFncyI6MCwiVmFsdWUiOiJiWGtnZG1Gc2RXVT0iLCJTZXNzaW9uIjoiIn0=" 
}
```


### Delete Request

Delete request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      |  key string      | any string      |
| method       | yes      | method name      | "delete"        |


Example:

```json
{
  "metadata": {
    "key": "your-consulkv-key",
    "method": "delete"
  },
  "data": null
}
```


### List Request

List request metadata setting:

| Metadata Key          | Required | Description                                             | Possible values |
|:----------------------|:---------|:--------------------------------------------------------|:----------------|
| perfix                | no       | prefix name to look up                                  | any string      |
| method                | yes      | method name                                             | "List"          |
| allow_stale           | no       | allows any Consul server (non-leader) to service a read | "false"           |
| require_consistent    | no       | RequireConsistent forces the read to be fully consistent| "false"           |
| use_cache             | no       | UseCache requests that the agent cache results locally  | "false"           |
| max_age               | no       | MaxAge limits how old a cached value will be returned   | "36000"           |
| max_age               | no       | MaxAge limits how old a cached value will be returned   | "36000"           |
| max_age               | no       | MaxAge limits how old a cached value will be returned   | "36000"           |

Example:

```json
{
  "metadata": {
    "key": "your-consulkv-key",
    "method": "get"
  },
  "data": null
}
```