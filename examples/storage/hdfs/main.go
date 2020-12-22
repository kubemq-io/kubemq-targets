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
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// listRequest
	statRequest := types.NewRequest().
		SetMetadataKeyValue("file_path", "/test/foo.txt").
		SetMetadataKeyValue("method", "stat")
	queryStatResponse, err := client.SetQuery(statRequest.ToQuery()).
		SetChannel("query.hdfs").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	statResponse, err := types.ParseResponse(queryStatResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("stat executed, response: %s", statResponse.Data))

}
