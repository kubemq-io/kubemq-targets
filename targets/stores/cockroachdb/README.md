# Kubemq Cockroach Target Connector

Kubemq Cockroach target connector allows services using kubemq server to access Cockroach database services.

## Prerequisites
The following are required to run the Cockroach target connector:

- kubemq cluster
- Cockroach server
- kubemq-targets deployment

## Configuration

Cockroach target connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| connection                      | yes      | Cockroach connection string address         | "postgres://root:postgres@localhost:26257/Cockroach?sslmode=disable" |
| max_idle_connections            | no       | set max idle connections                    | "10"                                                                   |
| max_open_connections            | no       | set max open connections                    | "100"                                                                  |
| connection_max_lifetime_seconds | no       | set max lifetime for connections in seconds | "3600"                                                                 |


Example:

```yaml
bindings:
- name: cockroachdb
  source:
    kind: kubemq.query
    properties:
      address: localhost:50000
      channel: query.cockroachdb
  target:
    kind: stores.cockroachdb
    properties:
      connection: postgres://root:postgres@localhost:26257/postgres?sslmode=disable
  properties: {}

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

### Transaction Request

Transaction request metadata setting:

| Metadata Key    | Required | Description                            | Possible values    |
|:----------------|:---------|:---------------------------------------|:-------------------|
| method          | yes      | set type of request                    | "transaction"             |
| isolation_level | no       | set isolation level for exec operation | ""                 |
|                 |          |                                        | "read_uncommitted" |
|                 |          |                                        | "read_committed"   |
|                 |          |                                        | "repeatable_read"  |
|                 |          |                                        | "serializable"     |


Transaction request data setting:

| Data Key | Required | Description                   | Possible values     |
|:---------|:---------|:------------------------------|:--------------------|
| data     | yes      | string string | base64 bytes array |

Example:

Transaction string:
```sql
DROP TABLE IF EXISTS post;
CREATE TABLE post (
         ID serial,
         TITLE varchar(40),
         CONTENT varchar(255),
         CONSTRAINT pk_post PRIMARY KEY(ID)
       );
INSERT INTO post(ID,TITLE,CONTENT) VALUES
                       (1,NULL,'Content One'),
                       (2,'Title Two','Content Two');
```
```json
{
  "metadata": {
    "key": "your-Cockroach-key",
    "method": "delete"
  },
  "data": "CURST1AgVEFCTEUgSUYgRVhJU1RTIHBvc3Q7CiAgICBDUkVBVEUgVEFCTEUgcG9zdCAoCgkgICAgICAgICBJRCBzZXJpYWwsCgkgICAgICAgICBUSVRMRSB2YXJjaGFyKDQwKSwKCSAgICAgICAgIENPTlRFTlQgdmFyY2hhcigyNTUpLAoJICAgICAgICAgQ09OU1RSQUlOVCBwa19wb3N0IFBSSU1BUlkgS0VZKElEKQoJICAgICAgICk7CiAgICBJTlNFUlQgSU5UTyBwb3N0KElELFRJVExFLENPTlRFTlQpIFZBTFVFUwoJICAgICAgICAgICAgICAgICAgICAgICAoMSxOVUxMLCdDb250ZW50IE9uZScpLAoJICAgICAgICAgICAgICAgICAgICAgICAoMiwnVGl0bGUgVHdvJywnQ29udGVudCBUd28nKTs="
}
```
