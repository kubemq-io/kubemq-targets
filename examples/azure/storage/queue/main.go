package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
)

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	dat, err := ioutil.ReadFile("./credentials/azure/storage/queue/queueName.txt")
	if err != nil {
		panic(err)
	}
	queueName := string(dat)
	dat, err = ioutil.ReadFile("./credentials/azure/storage/queue/serviceURL.txt")
	if err != nil {
		panic(err)
	}
	serviceURL := string(dat)
	myMessage := "my message to send to queue"
	b, err := json.Marshal(myMessage)
	if err != nil {
		panic(err)
	}
	// upload
	uploadRequest := types.NewRequest().
		SetMetadataKeyValue("method", "push").
		SetMetadataKeyValue("queue_name", queueName).
		SetMetadataKeyValue("service_url", serviceURL).
		SetData(b)
	queryUploadResponse, err := client.SetQuery(uploadRequest.ToQuery()).
		SetChannel("azure.storage.queue").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	uploadResponse, err := types.ParseResponse(queryUploadResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("push request executed, error: %v", uploadResponse.Error))

	// get count
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get_messages_count").
		SetMetadataKeyValue("queue_name", queueName).
		SetMetadataKeyValue("service_url", serviceURL)

	getQueryResponse, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("azure.storage.queue").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(getQueryResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get request done, response : %s and count : %s", getResponse.Data, getResponse.Metadata["count"]))
}
