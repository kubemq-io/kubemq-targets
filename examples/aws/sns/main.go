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
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// listRequest
	listRequest := types.NewRequest().
		SetMetadataKeyValue("method", "list_topics")
	queryListResponse, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.sns").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(queryListResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("list topics executed, response: %s", listResponse.Data))

	// Create Topic
	attributes := make(map[string]*string)
	DisplayName := "my-display-name"
	attributes["DisplayName"] = &DisplayName
	b, err := json.Marshal(attributes)
	if err != nil {
		log.Fatal(err)
	}
	dat, err := ioutil.ReadFile("./credentials/aws/sns/topic.txt")
	if err != nil {
		panic(err)
	}
	topic := string(dat)
	createRequest := types.NewRequest().
		SetMetadataKeyValue("method", "create_topic").
		SetMetadataKeyValue("topic", topic).
		SetData(b)

	getCreate, err := client.SetQuery(createRequest.ToQuery()).
		SetChannel("query.aws.sns").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	createResponse, err := types.ParseResponse(getCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("create topic executed, error: %v", createResponse.IsError))
}
