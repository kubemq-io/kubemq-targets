# Kubemq Crate Target Connector

Kubemq crate target connector allows services using kubemq server to access crate database services.

## Prerequisites
The following are required to run the crate target connector:

- kubemq cluster
- crate server
- kubemq-targets deployment

## Configuration

Crate target connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| connection                      | yes      | crate connection string address          | "postgresql://crate@localhost:5432/doc?sslmode=disable" |
| max_idle_connections            | no       | set max idle connections                    | "10"                                                                   |
| max_open_connections            | no       | set max open connections                    | "100"                                                                  |
| connection_max_lifetime_seconds | no       | set max lifetime for connections in seconds | "3600"                                                                 |


Example:

```yaml
bindings:
  - name: kubemq-query-crate
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-crate-connector"
        auth_token: ""
        channel: "query.crate"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: stores.crate
      name: target-crate
      properties:
        connection: ""postgresql://crate@localhost:5432/doc?sslmode=disable""
        max_idle_connections: "10"
        max_open_connections: "100"
        connection_max_lifetime_seconds: "3600"
```

## Usage

### Query Request

Query request metadata setting:

| Metadata Key | Required | Description      | Possible values |
|:-------------|:---------|:-----------------|:----------------|
| method          | yes      | set type of request | "query"      |

Query request data setting:

| Data Key | Required | Description  | Possible values    |
|:---------|:---------|:-------------|:-------------------|
| data     | yes      | query string | base64 bytes array |

Example:

Query string: `SELECT id,title,content FROM post;`

```json
{
  "metadata": {
    "method": "query"
  },
  "data": "U0VMRUNUIGlkLHRpdGxlLGNvbnRlbnQgRlJPTSBwb3N0Ow=="
}
```

### Exec Request

Exec request metadata setting:

| Metadata Key    | Required | Description                            | Possible values    |
|:----------------|:---------|:---------------------------------------|:-------------------|
| method          | yes      | set type of request                    | "exec"             |
| isolation_level | no       | set isolation level for exec operation | ""                 |
|                 |          |                                        | "read_uncommitted" |
|                 |          |                                        | "read_committed"   |
|                 |          |                                        | "repeatable_read"  |
|                 |          |                                        | "serializable"     |
|                 |          |                                        |                    |


Exec request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | exec string | base64 bytes array |

Example:

Exec string:
```sql
INSERT INTO post(ID,TITLE,CONTENT) VALUES (1,NULL,'Content One'),(2,'Title Two','Content Two');
```

```json
{
  "metadata": {
    "method": "exec",
    "isolation_level": "read_uncommitted"
  },
  "data": "SU5TRVJUIElOVE8gcG9zdChJRCxUSVRMRSxDT05URU5UKSBWQUxVRVMKCSAgICAgICAgICAgICAgICAgICAgICAgKDEsTlVMTCwnQ29udGVudCBPbmUnKSwKCSAgICAgICAgICAgICAgICAgICAgICAgKDIsJ1RpdGxlIFR3bycsJ0NvbnRlbnQgVHdvJyk7" 
}
```
