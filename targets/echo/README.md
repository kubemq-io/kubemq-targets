# Kubemq Echo Target Connector

Kubemq Echo target connector allows to echo back any request for testing purposes.

## Prerequisites
The following are required to run the echo target connector:

- kubemq cluster
- kubemq-targets deployment

## Configuration

echo target does not need any configuration

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
      kind: echo
      name: echo
      properties: {}
        default_headers: '{"Content-Type":"application/json"}'
```

## Usage
Any request will be send back as a response with the host name data embed within the metadata 

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
Will response back with:
```json
{
  "metadata": {
    "host": "some-host",
    "method": "get",
    "url": "https://httpbin.org/get" 
  },
 "data": null
}
```

