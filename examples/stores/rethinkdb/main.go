package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"
)

func main() {
	args := make(map[string]interface{})
	args["id"] = "test_user"
	args["password"] = 1
	insert, err := json.Marshal(args)
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	randomKey := uuid.New().String()
	// set request
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "insert").
		SetMetadataKeyValue("db_name", "test").
		SetMetadataKeyValue("table", "test").
		SetData(insert)
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.rethinkdb").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	setResponse, err := types.ParseResponse(querySetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("set request for key: %s executed, response: %s", randomKey, setResponse.Metadata.String()))

	// get request
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("db_name", "test").
		SetMetadataKeyValue("key", "test_user").
		SetMetadataKeyValue("table", "test")

	queryGetResponse, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("query.rethinkdb").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(queryGetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get request for key: %s executed, response: %s, data: %s", randomKey, getResponse.Metadata.String(), string(getResponse.Data)))
}
