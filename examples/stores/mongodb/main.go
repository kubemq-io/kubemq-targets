package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"
)

func main() {
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
		SetMetadataKeyValue("method", "set").
		SetMetadataKeyValue("key", randomKey).
		SetData([]byte("some-data"))
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.mongodb").
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
		SetMetadataKeyValue("key", randomKey)

	queryGetResponse, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("query.mongodb").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(queryGetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get request for key: %s executed, response: %s, data: %s", randomKey, getResponse.Metadata.String(), string(getResponse.Data)))

	// delete request

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("key", randomKey)

	queryDelResponse, err := client.SetQuery(delRequest.ToQuery()).
		SetChannel("query.mongodb").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	delResponse, err := types.ParseResponse(queryDelResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("delete request for key: %s executed, response: %s", randomKey, delResponse.Metadata.String()))
}
