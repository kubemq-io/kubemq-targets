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
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// list buses
	listRequest := types.NewRequest().
		SetMetadataKeyValue("method", "list_buses").
		SetMetadataKeyValue("limit", "1")

	queryListResponse, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.cloudwatch.events").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(queryListResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("list executed, response: %s", listResponse.Data))

	dat, err := ioutil.ReadFile("./credentials/aws/cloudwatch/events/rule.txt")
	if err != nil {
		panic(err)
	}
	rule := string(dat)

	dat, err = ioutil.ReadFile("./credentials/aws/cloudwatch/events/resourceARN.txt")
	if err != nil {
		panic(err)
	}
	resourceARN := string(dat)

	m := make(map[string]string)
	m["my_arn_id"] = resourceARN
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	//put targets
	putRequest := types.NewRequest().
		SetMetadataKeyValue("method", "put_targets").
		SetMetadataKeyValue("rule", rule).
		SetData(b)

	putReq, err := client.SetQuery(putRequest.ToQuery()).
		SetChannel("query.aws.cloudwatch.events").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	putResponse, err := types.ParseResponse(putReq.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("put request executed error: %v",  putResponse.IsError))
}
