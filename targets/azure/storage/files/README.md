# Kubemq files Target Connector

Kubemq files target connector allows services using kubemq server to access files messaging services.

## Prerequisites
The following are required to run the files target connector:

- kubemq cluster
- Azure active storage account
- Azure active with storage enable - with access account
- kubemq-targets deployment


## Configuration

files target connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| storage_account                 | yes      | azure storage account name                  | "my_account"                                                 |
| storage_access_key              | yes      | azure storage access key                    | "abcd1234"                                                   |
| policy                          | no       | exponential or linear                       | "retry_policy_exponential",default(retry_policy_exponential) |
| max_tries                       | no       | request max tries (1 disable)               | "1",default(1)                                                          |
| try_timeout                     | no       | Maximum time allowed for any single try     | "3000",default(10000) milliseconds                           |
| retry_delay                     | no       | Backoff amount for each retry               | "1000",default(600)   milliseconds                           |
| max_retry_delay                 | no       | delay between retries                       | "1000",default(1800)  milliseconds                           |


Example:

```yaml
bindings:
  - name: kubemq-query-azure-files
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-azure-files-connector"
        auth_token: ""
        channel: "azure.storage.files"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind:azure.storage.files
      name: azure-storage-files
      properties:
        storage_account: "id"
        storage_access_key: "key"
```

## Usage

### Create

Create metadata setting:

| Metadata Key      | Required | Description                                    | Possible values                                  |
|:------------------|:---------|:-----------------------------------------------|:-------------------------------------------------|
| method            | yes      | type of method                                 | "create"                                         |
| service_url       | yes      | service url path and filename                  | "https://test.files.core.windows.net/test/test.txt"   |                         |
| size              | no       | max file size                                  | "2000",default(1000000)   |   

```json
{
  "metadata": {
    "method": "create",
    "service_url": "https://test.end.point.test.net/test/test.txt"
  },
  "data": "bXktZmlsZS1kYXRh"
}
```

### Upload

Upload metadata setting:

| Metadata Key      | Required | Description                                    | Possible values                                  |
|:------------------|:---------|:-----------------------------------------------|:-------------------------------------------------|
| method            | yes      | type of method                                 | "upload"                                         |
| service_url       | yes      | service url path and filename                  | "https://test.files.core.windows.net/test/test.txt"   |
| data              | yes      | file data (byte array)                         | "bXktZmlsZS1kYXRh"                               |
| range_size        | no       | specifies the range size to use in bytes       | "0" ,default(4194304)                            | 
| parallelism       | no       | maximum number of ranges to upload in parallel | "upload",default(0)                              |
| file_metadata     | no       | key value string string file Metadata          | "{"tag":"test","name":"myname"}",default(none)   |                           |


Example:

```json
{
  "metadata": {
    "method": "upload",
    "service_url": "https://test.end.point.test.net/test/test.txt"
  },
  "data": "bXktZmlsZS1kYXRh"
}
```


### Get

Get metadata setting:

| Metadata Key      | Required | Description                             | Possible values                                        |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------------------|
| method            | yes      | type of method                          | "get"                                                  |                      |
| service_url       | yes      | service url path and filename           | "https://test.files.core.windows.net/test"              |
| max_retry_request | no       | max retry count                         | "20" (default "1")                                     |
| count             | no       | number of files to get                  | "20" (will get all from offset)                        |
| offset            | no       | start reading files from offset         | "20" (will start from the first byte in files)          |


Example:

```json
{
  "metadata": {
    "method": "get",
    "service_url": "https://test.end.point.test.net/test/test.txt"
  },
  "data": null
}
```

### Delete

Delete metadata setting:

| Metadata Key                   | Required | Description                             | Possible values                            |
|:-------------------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method                         | yes      | type of method                          | "delete"                                  |
| file_name                      | yes      | the name of the file to delete          | "myfile.txt"                              |
| service_url                    | yes      | service url path and filename           | "https://test.files.core.windows.net/test/test.txt" |


Example:

```json
{
  "metadata": {
    "method": "delete",
    "service_url": "https://test.end.point.test.net/test"
  },
  "data": null
}
```
