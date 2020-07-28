# Kubemq firebase target Connector

Kubemq gcp-firebase target connector allows services using kubemq server to access google firebase server.

## Prerequisites
The following required to run the gcp-firebase target connector:

- kubemq cluster
- gcp-firebase set up
- kubemq-source deployment

## Configuration

firebase target connector configuration properties:

| Properties Key  | Required | Description                                | Example                    |
|:----------------|:---------|:-------------------------------------------|:---------------------------|
| project_id      | yes      | gcp firebase project_id                    | "<googleurl>/myproject"    |
| credentials     | yes      | gcp credentials files                      | "<google json credentials" |
| db_client       | no       | initialize db client if true               | true/false                 |
| db_url          | no       | gcp db full path                           | <google db url"            |
| auth_client     | no       | initialize auth client if true             | true/false                 |
| messagingClient | no       | initialize messaging client                | true/false                 |
| defaultmsg      | no       | default Firebase Cloud Messaging           | json                       |
| defaultmultimsg | no       | default Firebase Cloud MulticastMessage    | json                       |   


Example:

```yaml
bindings:
  - name: kubemq-query-gcp-firebase
    source:
      kind: source.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-gcp-firebase-connector"
        auth_token: ""
        channel: "query.gcp.firebase"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: target.gcp.firebase
      name: target-gcp-firebase
      properties:
        project_id: "project_id"
        credentials: 'json'
        db_client: "true"
        db_url: "db_url"
        auth_client: "true"

```

## Usage


##DB:

### Get DB

Get DB metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | get_db                |
| ref_path     | yes      | ref path for the data                  | valid string          |
| child_ref    | no       | path for child ref data                | valid string          |


Example:

```json
{
  "metadata": {
    "method": "get_db",
    "ref_path": "string"
  },
  "data": null
}
```


### Set DB

Set DB metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | set_db                |
| ref_path     | yes      | ref path for the data                  | valid string          |
| child_ref    | no       | path for child ref data                | valid string          |


Example:

```json
{
  "metadata": {
    "method": "get_db",
    "ref_path": "string"
  },
  "data": null
}
```



### Update DB

Update DB metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | update_db             |
| ref_path     | yes      | ref path for the data                  | valid string          |
| child_ref    | no       | path for child ref data                | valid string          |


Example:

```json
{
  "metadata": {
    "method": "delete_db",
    "ref_path": "string"
  },
  "data": "updated_string"
}
```





### Delete DB

Delete DB metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | delete_db             |
| ref_path     | yes      | ref path for the data                  | valid string          |
| child_ref    | no       | path for child ref data                | valid string          |


Example:

```json
{
  "metadata": {
    "method": "delete_db",
    "ref_path": "string"
  },
  "data": null
}
```


##User:


### Create User

Create User metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | create_user           |


Example:

```json
{
  "metadata": {
    "method": "create_user"
  },
  "data": "\"email\":\"John@due.com\", \"password\":30, \"display_name\":\"John\""
}
```



### Retrieve User

Retrieve User metadata setting:

| Metadata Key | Required | Description                            | Possible values           |
|:-------------|:---------|:---------------------------------------|:--------------------------|
| method       | yes      | method type                            | retrieve_user             |
| retrieve_by  | yes      | type of retrieval                      | by_email ,by_uid,by_phone |
| uid          | no       | valid unique string                    | string                    |
| phone        | no       | valid phone number                     | string                    |
| email        | no       | valid email                            | string                    |


Example:

```json
{
  "metadata": {
    "method": "retrieve_user",
    "retrieve_by": "uid",
    "uid": "1223131"
  },
  "data": null
}
```


### Delete User

Delete User metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | delete_user           |
| uid          | yes      | valid unique string                    | string                |


Example:

```json
{
  "metadata": {
    "method": "delete_user",
    "uid": "1223131"
  },
  "data": null
}
```

### Delete Multiple Users

 Delete Multiple Users metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | delete_multiple_users |


Example:

```json
{
  "metadata": {
    "method": "delete_multiple_users"
  },
  "data": [ "user1", "user2", "user3" ]
}
```


### Update User

Update User metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | update_user           |        
| uid          | yes      | valid unique string                    | string                |      


Example:

```json
{
  "metadata": {
    "method": "update_user",
    "uid": "1223131"
  },
  "data": "\"email\":\"newJohn@due.com\", \"password\":30, \"display_name\":\"newJohn\""
}
```


### List Users

List User metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | list_users            |


Example:

```json
{
  "metadata": {
    "method": "list_users"
  },
  "data": null
}
```


##Token:

### Custom Token

Custom Token metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | custom_token                |        |      |
| uid          | yes      | valid unique string                    | string                     |        |


Example:

```json
{
  "metadata": {
    "method": "custom_token",
    "token_id": "some-uid"
  },
  "data": null
}
```

### Verify Token

Verify Token metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | verify_token                |        |      |
| uid          | yes      | valid unique string                    | string                     |        |


Example:

```json
{
  "metadata": {
    "method": "verify_token",
    "token_id": "some-uid"
  },
  "data": null
}
```

##Messaging:



### Send Message

  Firebase FCM messaging information can be found in https://firebase.google.com/docs/cloud-messaging/concept-options

Create User metadata setting:

| Metadata Key | Required | Description                            | Possible values       |
|:-------------|:---------|:---------------------------------------|:----------------------|
| method       | yes      | method type                            | SendMessage/SendBatch |




Example:
SendMessage
```json 
{
  "metadata": {
    "method": "SendMessage"
  },
  "data": {
    "topic":"app topic",
     "data":{"key1","val1"}
  }
}
```

SendBatch
```json 
{
  "metadata": {
    "method": "SendMessage"
  },
  "data": {
    "topic":"app topic",
     "data":{"Tokens":["123","456"],"Data":{"key":"val"},"Notification":{"title":"title"}}
  }
}
```