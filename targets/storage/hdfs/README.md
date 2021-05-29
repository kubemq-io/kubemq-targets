# Kubemq hadoop target Connector

Kubemq -hadoop target connector allows services using kubemq server to access hadoop service.

## Prerequisites
The following required to run the -hadoop target connector:

- kubemq cluster
- hadoop active server
- kubemq-targets deployment

## Configuration

hadoop target connector configuration properties:

| Properties Key | Required | Description                                | Example                     |
|:---------------|:---------|:-------------------------------------------|:----------------------------|
| address        | yes      | hadoop address                             |  "localhost:9000"         |
| user           | no       | hadoop user                                |  "my_user"   |


Example:

```yaml
bindings:
  - name: kubemq-query--hadoop
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query--hadoop-connector"
        auth_token: ""
        channel: "query..hadoop"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: storage.hadoop
      name: hadoop
      properties:
        _key: "id"
        _secret_key: 'json'
        region:  "region"
        token: ""
        downloader:  "true"
        uploader:  "true"
```

## Usage

### Read File

Read File:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| file_path         | yes      | path to file                            | "/test/foo2.txt"                     |
| method            | yes      | type of method                          | "read_file"                     |




Example:

```json
{
  "metadata": {
    "method": "read_file",
    "file_path": "/test/foo2.txt"
  },
  "data": null
}
```


### Write File

Write File:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| file_path         | yes      | path to file                            | "/test/foo2.txt"                     |
| method            | yes      | type of method                          | "write_file"                     |
| file_mode         | no       | os permission mode default(0777)        | "0777"                     |
| data              | yes      | file as byte array                      | "TXkgZXhhbXBsZSBmaWxlIHRvIHVwbG9hZA=="                     |




Example:

```json
{
  "metadata": {
    "method": "write_file",
    "file_path": "/test/foo2.txt"
  },
  "data": "TXkgZXhhbXBsZSBmaWxlIHRvIHVwbG9hZA=="
}
```

### Remove File

Remove File:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| file_path         | yes      | path to file                            | "/test/foo2.txt"                     |
| method            | yes      | type of method                          | "remove_file"                     |




Example:

```json
{
  "metadata": {
    "method": "remove_file",
    "file_path": "/test/foo2.txt"
  },
  "data": null
}
```

### Rename File

Rename File:

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| file_path         | yes      | new path to file                        | "/test/foo3.txt"                     |
| old_file_path     | yes      | new path to file                        | "/test/foo2.txt"                     |
| method            | yes      | type of method                          | "rename_file"                     |




Example:

```json
{
  "metadata": {
    "method": "rename_file",
    "file_path": "/test/foo3.txt",
    "old_file_path": "/test/foo2.txt"
  },
  "data": null
}
```

### Make Dir

Make Dir :

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| file_path         | yes      | new path to file                        | "/test_folder"                     |
| file_mode         | no       | os permission mode default(0777)        | "0777"                     |
| method            | yes      | type of method                          | "mkdir"                     |




Example:

```json
{
  "metadata": {
    "method": "mkdir",
    "file_path": "/test_folder"
  },
  "data": null
}
```

### Stat

Stat :

| Metadata Key      | Required | Description                             | Possible values                            |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------|
| file_path         | yes      | new path to file                        | "/test/foo3.txt"                     |
| method            | yes      | type of method                          | "stat"                     |




Example:

```json
{
  "metadata": {
    "method": "stat",
    "file_path": "/test/foo2.txt"
  },
  "data": null
}
```
