# Kubemq firestore target Connector

Kubemq gcp-firestore target connector allows services using kubemq server to access google firestore server.

## Prerequisites
The following required to run the gcp-firestore target connector:

- kubemq cluster
- gcp-firestore set up
- kubemq-source deployment

## Configuration

firestore target connector configuration properties:

| Properties Key | Required | Description                                | Example                    |
|:---------------|:---------|:-------------------------------------------|:---------------------------|
| project_id     | yes      | gcp firestore project_id                   | "<googleurl>/myproject"    |
| credentials    | yes      | gcp credentials files                      | "<google json credentials" |


Example:

```yaml
bindings:
  - name: kubemq-query-gcp-firestore
    source:
      kind: source.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-gcp-firestore-connector"
        auth_token: ""
        channel: "query.gcp.firestore"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.gcp.firestore
      name: target-gcp-firestore
      properties:
        project_id: "id"
        credentials: 'json'
        instance:  "instance"
```

## Usage

### Set Key 

Set Key  metadata setting:

| Metadata Key | Required | Description                            | Possible values      |
|:-------------|:---------|:---------------------------------------|:---------------------|
| method       | yes      | type of method                         | "add"                |
| collection   | yes      | the name of the collection to sent to  | "collection name"    |


Example:

```json
{
  "metadata": {
    "method": "add",
    "collection": "my_collection"
  },
  "data": "Hello"
}
```


### get Values by document key

Get Key  metadata setting:

| Metadata Key | Required | Description                            | Possible values            |
|:-------------|:---------|:---------------------------------------|:---------------------------|
| method       | yes      | type of method                         | "document_key"             |
| collection   | yes      | the name of the collection to sent to  | "collection name"          |
| document_key | yes      | the name of the key to get his value   | "valid existing key"       |


Example:

```json
{
  "metadata": {
    "method": "documents_all",
    "collection": "my_collection",
    "item": "<valid existing key>"
  },
  "data": "Hello"
}
```

### get all Values

Get all metadata setting:
| Metadata Key | Required | Description                             | Possible values        |
|:-------------|:---------|:----------------------------------------|:-----------------------|
| method       | yes      | type of method                         | "documents_all"         |
| collection   | yes      | the name of the collection to sent to  | "collection name"       |


Example:

```json
{
  "metadata": {
    "method": "documents_all",
    "collection": "my_collection"
  },
  "data": "Hello"
}
```


### delete key

Delete key metadata setting:

| Metadata Key | Required | Description                             | Possible values         |
|:-------------|:---------|:----------------------------------------|:------------------------|
| method       | yes      | type of method                          | "delete_document_key"   |
| collection   | yes      | the name of the collection to sent to   | "collection name"     |
| document_key | yes      | the name of the key to delete his value | "valid existing key"  |


Example:

```json
{
  "metadata": {
    "method": "delete_document_key",
    "collection": "my_collection",
    "item": "valid existing key"
  },
  "data": "Hello"
}
```
