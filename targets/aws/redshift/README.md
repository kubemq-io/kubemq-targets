# Kubemq redshift target Connector (Service admin)
Kubemq aws-redshift target connector allows services using kubemq server to access aws redshift service.

## Prerequisites
The following required to run the aws-redshift target connector:

- kubemq cluster
- aws account with redshift active service
- kubemq-source deployment

## Configuration

redshift target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |


Example:

```yaml
bindings:
  - name: kubemq-query-aws-redshift-service
    source:
      kind: source.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-aws-redshift-connector"
        auth_token: ""
        channel: "query.aws.redshift.service"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.aws.redshift.service
      name: target-aws-redshift-service
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "region"
```

## Usage

### Create Tags

Create Tags:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "create_tags"                     |
| resource_arn      | yes      | aws resource ARN                        | "arn:aws:redshift:region:account_id:cluster:cluster_name"                     |
| data              | yes      | key value of string string(tag-value)   | "eyJ0ZXN0MS1rZXkiOiJ0ZXN0MS12YWx1ZSJ9"                     |



Example:

```json
{
  "metadata": {
    "method": "create_tags",
    "resource_arn": "arn:aws:redshift:region:account_id:cluster:cluster_name"
  },
  "data": "eyJ0ZXN0MS1rZXkiOiJ0ZXN0MS12YWx1ZSJ9"
}
```

### Delete Tags

Delete Tags:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "delete_tags"                     |
| resource_arn      | yes      | aws resource ARN                        | "arn:aws:redshift:region:account_id:cluster:cluster_name"                     |
| data              | yes      | key slice of tags to remove(by keys)    | "WyJ0ZXN0MS1rZXkiXQ=="                     |



Example:

```json
{
  "metadata": {
    "method": "delete_tags",
    "resource_arn": "arn:aws:redshift:region:account_id:cluster:cluster_name"
  },
  "data": "WyJ0ZXN0MS1rZXkiXQ=="
}
```

### List Tags

List Tags:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_tags"                     |            |



Example:

```json
{
  "metadata": {
    "method": "list_tags"
  },
  "data": null
}
```

### List Snapshots

List Snapshots:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_snapshots"                     |            |



Example:

```json
{
  "metadata": {
    "method": "list_snapshots"
  },
  "data": null
}
```

### List Snapshots By Tag Keys

List Snapshots By Tag Keys:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_snapshots_by_tags_keys"                     |            |
| data              | yes      | key slice of tags to search by(by keys) | "WyJ0ZXN0MS1rZXkiXQ=="                     |


Example:

```json
{
  "metadata": {
    "method": "list_snapshots_by_tags_keys"
  },
  "data": "WyJ0ZXN0MS1rZXkiXQ=="
}
```

### List Snapshots By Tag Values

List Snapshots By Tag Values:

| Metadata Key      | Required | Description                               | Possible values                            |
|:------------------|:---------|:------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                            | "list_snapshots_by_tags_values"                     |            |
| data              | yes      | key slice of tags to search by(by values) | "WyJ0ZXN0MS1rZXkiXQ=="                     |


Example:

```json
{
  "metadata": {
    "method": "list_snapshots_by_tags_keys"
  },
  "data": "WyJ0ZXN0MS1rZXkiXQ=="
}
```

### Describe Clusters

Describe Clusters:

| Metadata Key      | Required | Description                               | Possible values                            |
|:------------------|:---------|:------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                            | "describe_cluster"                     |            |
| resource_name     | yes      | aws resource name                         | "my_cluster_name"                     |


Example:

```json
{
  "metadata": {
    "method": "list_snapshots_by_tags_keys",
    "resource_name": "my_cluster_name",
  },
  "data": null
}
```

### List Clusters

List Clusters:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_clusters"                     |            |



Example:

```json
{
  "metadata": {
    "method": "list_clusters"
  },
  "data": null
}
```

### List Clusters By Tag Keys

List Clusters By Tag Keys:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                          | "list_clusters_by_tags_keys"                     |            |
| data              | yes      | key slice of tags to search by(by keys) | "WyJ0ZXN0MS1rZXkiXQ=="                     |


Example:

```json
{
  "metadata": {
    "method": "list_clusters_by_tags_keys"
  },
  "data": "WyJ0ZXN0MS1rZXkiXQ=="
}
```

### List Clusters By Tag Values

List Clusters By Tag Values:

| Metadata Key      | Required | Description                               | Possible values                            |
|:------------------|:---------|:------------------------------------------|:-------------------------------------------|
| method            | yes      | type of method                            | "list_clusters_by_tags_values"                     |            |
| data              | yes      | key slice of tags to search by(by values) | "WyJ0ZXN0MS1rZXkiXQ=="                     |


Example:

```json
{
  "metadata": {
    "method": "list_clusters_by_tags_values"
  },
  "data": "WyJ0ZXN0MS1rZXkiXQ=="
}
```

