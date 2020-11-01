# Kubemq OpenFaas Target Connector

Kubemq OpenFaas target connector allows services using kubemq server to invoke OpensFaas's functions.

## Prerequisites
The following are required to run the OpenFaas target connector:

- kubemq cluster
- OpenFaas platform
- kubemq-targets deployment

## Configuration

OpenFaas target connector configuration properties:

| Properties Key | Required | Description               | Example                  |
|:---------------|:---------|:--------------------------|:-------------------------|
| gateway        | yes      | OpenFaas gateway address  | "http://localhost:31112" |
| username       | yes      | OpenFass gateway username | "admin"                  |
| password       | yes      | OpenFaas gateway password | "password"               |


Example:

```yaml
bindings:
  - name: kubemq-query-OpenFaas
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-OpenFaas-connector"
        auth_token: ""
        channel: "query.OpenFaas"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: serverless.openfaas
      name: target-serverless-openfaas
      properties:
        gateway: "http://localhost:31112"
        username: "admin"
        password: "password"
```

## Usage

### Request

Request metadata setting:

| Metadata Key | Required | Description             | Possible values          |
|:-------------|:---------|:------------------------|:-------------------------|
| topic        | yes      | OpenFaas function topic | "function/nslookup"      |


Request data setting:

| Data Key | Required | Description                          | Possible values     |
|:---------|:---------|:-------------------------------------|:--------------------|
| data     | yes      | data to set for the OpenFaas request | base64 bytes array |

Example:

```json
{
  "metadata": {
    "topic": "function/nslookup"
  },
 "data": "a3ViZW1xLmlv"
}
```
