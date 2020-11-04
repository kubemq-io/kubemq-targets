package main

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"log"
	"time"
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
		SetMetadataKeyValue("queue", "some-queue").
		SetMetadataKeyValue("exchange", "").
		SetMetadataKeyValue("delivery_mode", "2").
		SetData([]byte("some-data"))
	queryPublishResponse, err := client.SetQuery(publishRequest.ToQuery()).
		SetChannel("messaging.rabbitmq").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	publishResponse, err := types.ParseResponse(queryPublishResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("publish message, response: %s", publishResponse.Metadata.String()))
	time.Sleep(time.Second)

}
