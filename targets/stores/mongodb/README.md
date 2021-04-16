# Kubemq Mongodb Target Connector

Kubemq mongodb target connector allows services using kubemq server to access mongodb database services.

## Prerequisites
The following are required to run the mongodb target connector:

- kubemq cluster
- mongodb server
- kubemq-targets deployment

## Configuration

Mongodb target connector configuration properties:

| Properties Key            | Required | Description                          | Example                   |
|:--------------------------|:---------|:-------------------------------------|:--------------------------|
| host                      | yes      | mongodb host address                 | "localhost:27017"         |
| username                  | no       | mongodb username                     | "admin"                   |
| password                  | no       | mongodb password                     | "password"                |
| database                  | no       | set database name                    | "admin"                   |
| collection                | no       | set database collection              | "test"                    |
| params                    | no       | set connection additional parameters | ""                        |
| write_concurrency         | no       | set write concurrency                | "","majority","1","2"     |
| read_concurrency          | no       | set read concurrency                 | "","local"                |
|                           |          |                                      | "","local"                |
|                           |          |                                      | "majority","available"    |
|                           |          |                                      | "linearizable","snapshot" |
| operation_timeout_seconds | no       | set operation timeout in seconds     | "30"                      |



Example:

```yaml
bindings:
  - name: kubemq-query-mongodb
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-mongodb-connector"
        auth_token: ""
        channel: "query.mongodb"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: stores.mongodb
      name: target-mongodb
      properties:
        host: "localhost:27017"
        username: "admin"
        password: "password"
        database: "admin"
        collection: "test"
        write_concurrency: "majority"
        read_concurrency: ""
        params: ""
        operation_timeout_seconds: "2"
```

## Usage



### Find Request

Find document by filter request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | find document by set filter   | "find"           |
| filter       | yes      | filter json object   | '{"_id":"some-id"}'       |
Example:

```json
{
  "metadata": {
     "method": "find",
     "filter": "{\"_id\":\"e7141dc8-d793-4f3b-8290-bc9d9a3f9fda\"}"
  },
  "data": null
}
```


### Find Many Request

Find many documents by filter request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | find document by set filter   | "find_many"           |
| filter       | yes      | filter json object   | '{"color":"white"}'       |

Example:

```json
{
  "metadata": {
     "method": "find_many",
     "filter": "{\"color\":\"white\"}"
  },
  "data": null
}
```


### Insert Request

Insert single document request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | insert document  | "insert"           |


insert request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | json document to insert | base64 bytes array |

Example:

```json
{
  "metadata": {
     "method": "insert"
  },
  "data":"eyJfaWQiOiI2ZjNhNGJlNy00OTk4LTRjYmYtYjc5ZS0xOWQ3Y2FiOGU4MDkiLCJzdHJpbmciOiJzIiwiaW50ZWdlciI6IjEifQ=="
}
```


### Insert Many Request

Insert multiple documents request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | insert document  | "insert_many"           |


Insert multiple request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | array of json documents to insert | base64 bytes array |

Example:

```json
{
  "metadata": {
     "method": "insert_many"
  },
  "data":"W3siX2lkIjoiMmNjZjkyYmEtNTdmOC00NDQxLWFmZmEtOTJmOGY5MTAxMGYzIiwic3RyaW5nIjoicyIsImludGVnZXIiOiIxIn1d"
}
```

### Update Request

Update single document request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | update document  | "update"           |
| filter       | yes      | filter json object   | '{"_id":"some-id"}'       |
| set_upsert       | no      | perform insert if not found   | true       |

Update request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | updated json document | base64 bytes array |

Example:

```json
{
  "metadata": {
     "method": "update",
      "filter": "{\"_id\":\"e7141dc8-d793-4f3b-8290-bc9d9a3f9fda\"}",
      "set_upsert": "true"
  },
  "data":"eyJfaWQiOiI2ZjNhNGJlNy00OTk4LTRjYmYtYjc5ZS0xOWQ3Y2FiOGU4MDkiLCJzdHJpbmciOiJzIiwiaW50ZWdlciI6IjEifQ=="
}
```


### Update Many Request

Update many documents request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | update document  | "update_many"           |
| filter       | yes      | filter json object   | '{"color":"blue"}'       |
| set_upsert       | no      | perform insert if not found   | true       |

Update many request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | updated json document | base64 bytes array |

Example:

```json
{
  "metadata": {
     "method": "update_many",
      "filter": "{\"color\":\"blue\"}",
      "set_upsert": "true"
  },
  "data":"eyJfaWQiOiI2ZjNhNGJlNy00OTk4LTRjYmYtYjc5ZS0xOWQ3Y2FiOGU4MDkiLCJzdHJpbmciOiJzIiwiaW50ZWdlciI6IjEifQ=="
}
```


### Delete Request

Delete document by filter request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | delete document by set filter   | "delete"           |
| filter       | yes      | filter json object   | '{"_id":"some-id"}'       |
Example:

```json
{
  "metadata": {
     "method": "delete",
     "filter": "{\"_id\":\"e7141dc8-d793-4f3b-8290-bc9d9a3f9fda\"}"
  },
  "data": null
}
```

### Delete Many Request

Delete many documents by filter request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | delete many documents by set filter   | "delete_many"           |
| filter       | yes      | filter json object   | '{"color":"white"}'       |
Example:

```json
{
  "metadata": {
     "method": "delete_many",
     "filter": "{\"color\":\"white\"}"
  },
  "data": null
}
```

### Aggregate Request

Aggregate documents by pipeline request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | aggregate documents by set pipeline   | "aggregate"           |


Aggregate pipeline request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | pipeline json data | base64 bytes array |

Example:

```json
{
  "metadata": {
     "method": "aggregate"
  },
  "data": "WwogICB7ICRtYXRjaDogeyBzdGF0dXM6ICJBIiB9IH0sCiAgIHsgJGdyb3VwOiB7IF9pZDogIiRjdXN0X2lkIiwgdG90YWw6IHsgJHN1bTogIiRhbW91bnQiIH0gfSB9Cl0="
}
```

### Distinct Request

Distinct documents by filter and field name request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method       | yes      | get distinct documents by set filter   | "distinct"           |
| filter       | yes      | filter json object   | '{"_id":"some-id"}'       |
| field_name       | yes      | distinct field name   | "id"       |

Example:

```json
{
  "metadata": {
     "method": "distinct",
     "filter": "{\"_id\":\"e7141dc8-d793-4f3b-8290-bc9d9a3f9fda\"}",
     "field_name": "id"
  },
  "data": null
}
```

### Get By Key Request

Get request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | mongodb key string | any string      |
| method       | yes      | get document by key   | "get_by_key"           |

Example:

```json
{
  "metadata": {
    "key": "key",
    "method": "get_by_key"
  },
  "data": null
}
```

### Set By Key Request

Set request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | mongodb key string | any string      |
| method       | yes      | set document by key             | "set_by_key"           |

Set request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | data to set for the mongodb key | base64 bytes array |

Example:

```json
{
  "metadata": {
    "key": "key",
    "method": "set_by_key"
  },
  "data": "c29tZS1kYXRh" 
}
```
### Delete By Key Request

Delete request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| key          | yes      | mongodb key string | any string      |
| method       | yes      | delete document by key           | "delete_by_key"        |

Example:

```json
{
  "metadata": {
    "key": "key",
    "method": "delete_by_key"
  },
  "data": null
}
```
