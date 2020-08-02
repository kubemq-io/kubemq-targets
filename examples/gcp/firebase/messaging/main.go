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
		SetMetadataKeyValue("method", "SendMessage").
		SetData([]byte(`{"Topic":"test",
								"data":{
									"key1":"value1"
										}}`))

	queryPublishResponse, err := client.SetQuery(publishRequest.ToQuery()).
		SetChannel("query.gcp.firebase").
		SetTimeout(90 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	publishResponse, err := types.ParseResponse(queryPublishResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("publish message, response: %s", publishResponse.Metadata.String()))

}
