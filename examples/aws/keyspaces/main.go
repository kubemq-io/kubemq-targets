package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
)

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
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
		SetMetadataKeyValue("consistency", "LocalQuorum").
		SetData([]byte("some-data"))
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("aws.query.keyspaces").
		SetTimeout(60 * time.Second).Send(context.Background())
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
		SetMetadataKeyValue("key", randomKey).
		SetMetadataKeyValue("consistency", "LocalQuorum")

	queryGetResponse, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("aws.query.keyspaces").
		SetTimeout(60 * time.Second).Send(context.Background())
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
		SetMetadataKeyValue("consistency", "LocalQuorum").
		SetMetadataKeyValue("key", randomKey)

	queryDelResponse, err := client.SetQuery(delRequest.ToQuery()).
		SetChannel("aws.query.keyspaces").
		SetTimeout(60 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	delResponse, err := types.ParseResponse(queryDelResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("delete request for key: %s executed, response: %s", randomKey, delResponse.Metadata.String()))
}
