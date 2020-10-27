# Kubemq Postgres-aws Target Connector
Please see before starting:

- https://aws.amazon.com/premiumsupport/knowledge-center/users-connect-rds-iam/
- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.DBAccounts.html

Kubemq postgres-aws  target connector allows services using kubemq server to access postgres database services.

## Prerequisites
The following are required to run the postgres-aws target connector:

- kubemq cluster
- aws rds with IAM user postgres server
- db user with IAM Authentication privileges 
- kubemq-targets deployment

## Configuration

Postgres target connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| max_idle_connections            | no       | set max idle connections                    | "10"                                                                   |
| max_open_connections            | no       | set max open connections                    | "100"                                                                  |
| connection_max_lifetime_seconds | no       | set max lifetime for connections in seconds | "3600"     
| db_user                         | yes      | aws db user name                            | "<aws user"               |
| db_name                         | yes      | aws db name                                 | "<aws instance name"      |
| aws_key                         | yes      | aws key                                     | aws key supplied by aws         |
| aws_secret_key                  | yes      | aws secret key                              | aws secret key supplied by aws  |
| region                          | yes      | region                                      | aws region                      |
| token                           | no       | aws token ("default" empty string           | aws token                       |
| end_point                       | yes      | aws rds end point                           | "<aws rds end point"        |

                                                              |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-postgres
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-postgres-connector"
        auth_token: ""
        channel: "query.aws.rds.postgres"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.rds.postgres
      name: aws-rds-postgres
      properties:
              instance_connection_name: "test"
              db_user:  "test"
              db_name:  "test"
              end_point: "my-end-point"
              aws_key: "id"
              aws_secret_key: 'json'
              region:  "region"
              token: ""
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
    "key": "your-postgres-key",
    "method": "delete"
  },
  "data": "CURST1AgVEFCTEUgSUYgRVhJU1RTIHBvc3Q7CiAgICBDUkVBVEUgVEFCTEUgcG9zdCAoCgkgICAgICAgICBJRCBzZXJpYWwsCgkgICAgICAgICBUSVRMRSB2YXJjaGFyKDQwKSwKCSAgICAgICAgIENPTlRFTlQgdmFyY2hhcigyNTUpLAoJICAgICAgICAgQ09OU1RSQUlOVCBwa19wb3N0IFBSSU1BUlkgS0VZKElEKQoJICAgICAgICk7CiAgICBJTlNFUlQgSU5UTyBwb3N0KElELFRJVExFLENPTlRFTlQpIFZBTFVFUwoJICAgICAgICAgICAgICAgICAgICAgICAoMSxOVUxMLCdDb250ZW50IE9uZScpLAoJICAgICAgICAgICAgICAgICAgICAgICAoMiwnVGl0bGUgVHdvJywnQ29udGVudCBUd28nKTs="
}
```
