# Kubemq Elastic Search Target Connector

Kubemq elastic-search target connector allows services using kubemq server to access elastic-search database services.

## Prerequisites
The following are required to run the elastic-search target connector:

- kubemq cluster
- elastic-search server
- kubemq-targets deployment

## Configuration

Elastic Search target connector configuration properties:

| Properties Key            | Required | Description                          | Example                   |
|:--------------------------|:---------|:-------------------------------------|:--------------------------|
| urls                      | yes      | elastic-search list of urls separated by comma                 | "http://localhost:9200,http://localhost:9201"         |
| username                  | no       | elastic-search username                     | "admin"                   |
| password                  | no       | elastic-search password                     | "password"                |
| sniff                  | no       | set sniff opn connect                    | "true", "false"                   |



Example:

```yaml
bindings:
  - name: kubemq-query-elastic-search
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-elastic-search-connector"
        auth_token: ""
        channel: "query.elastic-search"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: stores.elastic
      name: target-elastic-search
      properties:
        urls: "http://localhost:9200"
        username: "admin"
        password: "password"
        sniff: "false"
```

## Usage

### Get Request

Get request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name get                        | "get"|
| index        | yes      | elastic-search index table | any string      |
| id           | yes      | document id                | any string      |


Example:

```json
{
  "metadata": {
    "method": "get",
    "index": "log",
    "id": "doc-id"
  },
  "data": null
}
```

### Set Request

Set request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name set                        | "set"           |
| index        | yes      | elastic-search index table | any string      |
| id           | yes      | document id                | any string      |


Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the elastic-search key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "method": "set",
    "index": "log",
    "id": "doc-id"
  },
  "data": "ewoJImlkIjogInNvbWUtaWQiLAoiZGF0YSI6InNvbWUtZGF0YSIKfQ==" 
}
```
### Delete Request

Delete request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name delete                        | "delete"           |
| index        | yes      | elastic-search index table | any string      |
| id           | yes      | document id                | any string      |

Example:

```json
{
  "metadata": {
    "method": "delete",
    "index": "log",
    "id": "doc-id"
  },
  "data": null
}
```
### Index Exists Request

Index exists request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name index.exists                        | "index.exists"|
| index        | yes      | elastic-search index table | any string      |


Example:

```json
{
  "metadata": {
    "method": "index.exists",
    "index": "log"
  },
  "data": null
}
```

### Index Create Request

Index create request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name index.create                        | "index.create"|
| index        | yes      | elastic-search index table | any string      |


Index create data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set index mapping| base64 bytes array |


Example:

Mapping log index
```json
{
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
		"properties": {
			"id": {
				"type": "keyword"
			},
			"data": {
				"type": "text"
			}
		}
	}
}

```

Request:

```json
{
  "metadata": {
    "method": "index.create",
    "index": "log"
  },
  "data": "ewoJInNldHRpbmdzIjogewoJCSJudW1iZXJfb2Zfc2hhcmRzIjogMSwKCQkibnVtYmVyX29mX3JlcGxpY2FzIjogMAoJfSwKCSJtYXBwaW5ncyI6IHsKCQkicHJvcGVydGllcyI6IHsKCQkJImlkIjogewoJCQkJInR5cGUiOiAia2V5d29yZCIKCQkJfSwKCQkJImRhdGEiOiB7CgkJCQkidHlwZSI6ICJ0ZXh0IgoJCQl9CgkJfQoJfQp9"
}
```
### Index Delete Request

Index Delete request metadata setting:

| Metadata Key | Required | Description                | Possible values |
|:-------------|:---------|:---------------------------|:----------------|
| method       | yes      | method name index.delete                        | "index.delete"           |
| index        | yes      | elastic-search index table | any string      |


Example:

```json
{
  "metadata": {
    "method": "index.delete",
    "index": "log"
  },
  "data": null
}
```
