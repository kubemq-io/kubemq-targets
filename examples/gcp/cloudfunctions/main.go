package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
)

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	publishRequest := types.NewRequest().
		SetMetadataKeyValue("name", "test-kubemq").
		SetMetadataKeyValue("project", "pubsubdemo-281010").
		SetMetadataKeyValue("location", "us-central1").
		SetData([]byte(`{"message":"test"}`))
	queryPublishResponse, err := client.SetQuery(publishRequest.ToQuery()).
		SetChannel("query.gcp.functions").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	publishResponse, err := types.ParseResponse(queryPublishResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("publish message, response: %s", publishResponse.Metadata.String()))

}
