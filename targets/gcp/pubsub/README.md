# Kubemq pubsub target Connector

Kubemq gcp-pubsub target connector allows services using kubemq server to access google pubsub server.

## Prerequisites
The following required to run the gcp-pubsub target connector:

- kubemq cluster
- gcp-pubsub set up
- kubemq-source deployment

## Configuration

pubsub target connector configuration properties:

| Properties Key | Required | Description                                | Example                    |
|:---------------|:---------|:-------------------------------------------|:---------------------------|
| project_id     | yes      | gcp firestore project_id                   | "<googleurl>/myproject"    |
| credentials    | yes      | gcp credentials files                      | "<google json credentials" |
| retries        | no       | number of sending retires                  | retries number             |


Example:

```yaml
bindings:
  - name: kubemq-query-gcp-pubsub
    source:
      kind: source.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-gcp-pubsub-connector"
        auth_token: ""
        channel: "query.gcp.pubsub"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.gcp.pubsub
      name: target-gcp-pubsub
      properties:
        project_id: "projectID"
        retries:    "0"
        credentials: 'json'

```

## Usage

### Send Message 

Send Message metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| topicID      | yes      | the name of the topicID to sent to     | "<valid topicID>"     |
| tags         | no       | type of method                         | string tags           |


Example:

```json
{
  "metadata": {
    "topic_id": "my_topic",
    "tags": "{\"tag-1\":\"test\",\"tag-2\":\"test2\"}"
  },
  "data": "Hello"
}
```
