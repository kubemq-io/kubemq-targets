package main

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
	"github.com/kubemq-hub/kubemq-targets/types"
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
	publishRequest := types.NewRequest().
		SetMetadataKeyValue("topic", "some-queue").
		SetMetadataKeyValue("qos", "0").
		SetData([]byte("some-data"))
	queryPublishResponse, err := client.SetQuery(publishRequest.ToQuery()).
		SetChannel("query.mqtt").
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
