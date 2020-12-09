# Kubemq blob Target Connector

Kubemq blob target connector allows services using kubemq server to access blob messaging services.

## Prerequisites
The following are required to run the blob target connector:

- kubemq cluster
- Azure active storage account
- Azure active with storage enable - with access account
- kubemq-targets deployment


## Configuration

blob target connector configuration properties:

| Properties Key                  | Required | Description                                          | Example                                                                |
|:--------------------------------|:---------|:-----------------------------------------------------|:-----------------------------------------------------------------------|
| storage_account                 | yes     | azure storage account name                            | "my_account"                                                              |
| storage_access_key              | yes     | azure storage access key                              | "abcd1234"                                                                |
| policy                          | no      | azure blob retry policy                               | "retry_policy_exponential",retry_policy_fixed(default retry_policy_fixed) |
| max_tries                       | no      | try at most x times to perform the operation          | "3" default (1)                                                           |
| try_timeout                     | no      | Maximum time allowed for any single try (Millisecond) | "600"default (1000)                                                       |
| retry_delay                     | no      | Backoff amount for each retry (Millisecond)           | "60" default (60)                                                         |
| max_retry_delay                 | no      | Max delay between retries (Millisecond)               | "180"default (180                                                         |
 

Example:

```yaml
bindings:
  - name: kubemq-query-azure-blob
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-azure-blob-connector"
        auth_token: ""
        channel: "azure.storage.blob"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: azure.storage.blob
      name: azure-storage-blob
      properties:
        storage_account: "id"
        storage_access_key: "key"
```

## Usage

### Upload

Upload metadata setting:

| Metadata Key      | Required | Description                                    | Possible values                                  |
|:------------------|:---------|:-----------------------------------------------|:-------------------------------------------------|
| method            | yes      | type of method                                 | "upload"                                         |
| file_name         | yes      | the name to upload the file under              | "myfile.txt"                                     |
| service_url       | yes      | service url path                               | "https://test.blob.core.windows.net/test"        |
| data              | yes      | file data (byte array)                         | "bXktZmlsZS1kYXRh"                               |
| block_size        | no       | specifies the block size to use                | "0" ,default(azblob.BlockBlobMaxStageBlockBytes) |
| parallelism       | no       | maximum number of blocks to upload in parallel | "upload",default(0)                              |
| blob_metadata     | no       | Key value string string of blob_metadata       | "{"tag":"test","name":"myname"}"                              |


Example:

```json
{
  "metadata": {
    "method": "upload",
    "file_name": "myfile.txt",
    "service_url": "https://test.end.point.test.net/test"
  },
  "data": "bXktZmlsZS1kYXRh"
}
```

### Delete

Delete metadata setting:

| Metadata Key                   | Required | Description                             | Possible values                            |
|:-------------------------------|:---------|:----------------------------------------|:-------------------------------------------|
| method                         | yes      | type of method                          | "delete"                                  |
| file_name                      | yes      | the name of the file to delete          | "myfile.txt"                              |
| service_url                    | yes      | service url path                        | "https://test.blob.core.windows.net/test" |
| delete_snapshots_option_type   | no       | type of method                          | "only","include","" (default "")          |


Example:

```json
{
  "metadata": {
    "method": "delete",
    "file_name": "myfile.txt",
    "service_url": "https://test.end.point.test.net/test"
  },
  "data": null
}
```

### get

For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob
get metadata setting:

| Metadata Key      | Required | Description                             | Possible values                                        |
|:------------------|:---------|:----------------------------------------|:-------------------------------------------------------|
| method            | yes      | type of method                          | "get"                                                  |
| file_name         | yes      | the name of the file to get             | "myfile.txt"                                           |
| service_url       | yes      | service url path                        | "https://test.blob.core.windows.net/test"              |
| max_retry_request | no       | type of method                          | "20" (default "1")                                     |
| count             | no       | number of files to get                  | "20" (will get all from offset)                        |
| offset            | no       | start reading blob from offset          | "20" (will start from the first byte in blob)          |


Example:

```json
{
  "metadata": {
    "method": "get",
    "file_name": "myfile.txt",
    "service_url": "https://test.end.point.test.net/test"
  },
  "data": null
}
```
