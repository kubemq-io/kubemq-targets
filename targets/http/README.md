# Kubemq Https Target Connector

Kubemq Https target connector allows services using kubemq server to invoke http rest function to any destination.

## Prerequisites
The following are required to run the http target connector:

- kubemq cluster
- kubemq-targets deployment

## Configuration

http target connector configuration properties:

| Properties Key     | Required | Description                                        | Example                          |
|:-------------------|:---------|:---------------------------------------------------|:---------------------------------|
| auth_type          | no       | http authentication type                           | "","no_auth","basic","auth_token |
| username           | no       | set username in auth_type=basic mode               | "admin"                          |
| password           | no       | set password in auth_type=basic mode               | "password"                       |
| token              | no       | set auth token in auth_type=auth_token mode        | valid JWT token                  |
| proxy              | no       | set proxy url                                      | "http://localhost:8080"          |
| root_certificate   | no       | set root ca certificate for mTLS handshake         | any x509 pem                     |
| client_private_key | no       | set private key for mTLS handshake                 | any x509 pem                     |
| client_public_key  | no       | set public key for mTLS handshake                  | any x509 pem                     |
| default_headers    | no       | set any default headers to be add for each call    | map of headers                   |


Example:

```yaml
bindings:
  - name: kubemq-query-https
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-http-connector"
        auth_token: ""
        channel: "query.http"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: http
      name: target-http
      properties:
        auth_type: "no-auth"
        default_headers: '{"Content-Type":"application/json"}'
```

## Usage

### Request

Request metadata setting:

| Metadata Key | Required | Description                     | Possible values                       |
|:-------------|:---------|:--------------------------------|:--------------------------------------|
| method       | yes      | http method to invoke           | "get","post","head","put"             |
|              |          |                                 | "delete","patch","options"            |
| url          | yes      | http url                        | "https://httpbin.org/get"             |
| headers      | no       | any headers required for method | '{"Content-Type":"application/json"}' |


Request data setting:

| Data Key | Required | Description                          | Possible values     |
|:---------|:---------|:-------------------------------------|:--------------------|
| data     | yes      | data to set for the http request | base64 bytes array |

Example:

```json
{
  "metadata": {
    "method": "get",
    "url": "https://httpbin.org/get" 
  },
 "data": null
}
```
