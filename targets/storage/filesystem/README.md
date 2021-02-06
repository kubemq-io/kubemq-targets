# Kubemq Filesystem Target Connector

Kubemq Filesystem target connector allows services using kubemq server to perform filesystem operation such as save,load, delete and list.

## Prerequisites
The following are required to run the minio target connector:

- kubemq cluster
- kubemq-targets deployment

## Configuration

Filesystem target connector configuration properties:

| Properties Key    | Required | Description                              | Example          |
|:------------------|:---------|:-----------------------------------------|:-----------------|
| base_path          | yes      | base root for all functions                     | "./" |

Example:

```yaml
bindings:
  - name: kubemq-query-filesystem
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-fs-connector"
        auth_token: ""
        channel: "query.fs"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: stores.filesystem
      name: target-filesystem
      properties:
        base_path: "./"
 ```

## Usage

### Save File Request

Save file request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "save"   |
| path       | no      | set path for filename     | "path"        |
| filename       | yes       | set filename | "filename.txt"              |


Example:

```json
{
  "metadata": {
    "method": "save",
    "path": "path",
    "filename": "filename.txt"
  },
  "data": "c29tZS1kYXRh"
}
```

### Load File Request

Load file request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "load"   |
| path       | no      | set path for filename     | "path"        |
| filename       | yes       | set filename | "filename.txt"              |

Example:

```json
{
  "metadata": {
    "method": "load",
    "path": "path",
    "filename": "filename.txt"
  },
  "data": null
}
```
### Delete File Request

Delete file request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "delete"   |
| path       | no      | set path for filename     | "path"        |
| filename       | yes       | set filename | "filename.txt"              |

Example:

```json
{
  "metadata": {
    "method": "delete",
    "path": "path",
    "filename": "filename.txt"
  },
  "data": null
}
```


### List Request

List files in directory request metadata setting:

| Metadata Key | Required | Description         | Possible values |
|:-------------|:---------|:--------------------|:----------------|
| method       | yes      | method name         | "list"   |
| path       | no      | set path for filename     | "path"        |


Example:

```json
{
  "metadata": {
    "method": "list",
    "path": "path"
  },
  "data": null
}
```
