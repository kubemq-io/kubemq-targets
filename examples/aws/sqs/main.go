package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/aws/sqs/queue.txt")
	if err != nil {
		log.Fatal(err)
	}
	queue := string(dat)
	validBody, err := json.Marshal("valid body2")
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// sendRequest
	sendRequest := types.NewRequest().
		SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
		SetMetadataKeyValue("delay", "0").
		SetMetadataKeyValue("queue", queue).
		SetData(validBody)

	querySendResponse, err := client.SetQuery(sendRequest.ToQuery()).
		SetChannel("query.aws.sqs").
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
