package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"
)

func main() {
	validBody, err := json.Marshal("valid body2")
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// sendRequest
	sendRequest := types.NewRequest().
		SetData(validBody)

	querySendResponse, err := client.SetQuery(sendRequest.ToQuery()).
		SetChannel("query.messaging.ibmmq").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	sendResponse, err := types.ParseResponse(querySendResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("send executed, response: %s", sendResponse.Data))
}
