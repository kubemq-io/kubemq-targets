package main

import (
	"context"
	"encoding/json"
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
	body, err := json.Marshal("test")
	if err != nil {
		panic(err)
	}
	// send
	sendRequest := types.NewRequest().
		SetMetadataKeyValue("method", "send").
		SetMetadataKeyValue("label", `test`).
		SetData(body)
	sendUploadResponse, err := client.SetQuery(sendRequest.ToQuery()).
		SetChannel("azure.servicebus").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	sendResponse, err := types.ParseResponse(sendUploadResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("send request executed, error: %v , response:%s", sendResponse.Error,sendResponse.Metadata["result"]))

}
