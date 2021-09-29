package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/topicID.txt")
	if err != nil {
		panic(err)
	}
	topicID := string(dat)
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// add request
	setRequest := types.NewRequest().
		SetMetadataKeyValue("tags", `{"tag-1":"test","tag-2":"test2"}`).
		SetMetadataKeyValue("topic_id", topicID)
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.gcp.pubsub").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	setResponse, err := types.ParseResponse(querySetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("set request with tags for topic : %s executed, is_error: %v", topicID, setResponse.IsError))
}
