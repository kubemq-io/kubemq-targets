package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
)

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
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
		SetMetadataKeyValue("properties", `{"tag-1":"test","tag-2":"test2"}`).
		SetData(body)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	sendUploadResponse, err := client.SetQuery(sendRequest.ToQuery()).
		SetChannel("azure.eventhubs").
		SetTimeout(10 * time.Second).Send(ctx)
	if err != nil {
		log.Fatal(err)
	}
	sendResponse, err := types.ParseResponse(sendUploadResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("send request executed, error: %v", sendResponse.Error))
}
