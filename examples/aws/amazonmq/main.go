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
		SetMetadataKeyValue("destination", "some-destination").
		SetData([]byte("some-data"))
	queryPublishResponse, err := client.SetQuery(publishRequest.ToQuery()).
		SetChannel("query.amazonmq").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v \n", queryPublishResponse)
	publishResponse, err := types.ParseResponse(queryPublishResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("publish message, response: %+v", publishResponse))

}
