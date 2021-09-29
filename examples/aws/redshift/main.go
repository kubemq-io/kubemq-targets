package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/aws/redshift-svc/resourceARN.txt")
	if err != nil {
		log.Fatal(err)
	}
	resourceARN := fmt.Sprintf("%s", dat)
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}

	// Create tag
	tags := make(map[string]string)
	tags["test-key"] = "test-value"
	b, err := json.Marshal(tags)
	if err != nil {
		log.Fatal(err)
	}
	createRequest := types.NewRequest().
		SetMetadataKeyValue("method", "create_tags").
		SetMetadataKeyValue("resource_arn", resourceARN).
		SetData(b)

	getCreate, err := client.SetQuery(createRequest.ToQuery()).
		SetChannel("query.aws.redshift.service").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	createResponse, err := types.ParseResponse(getCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("create tag executed, error: %v", createResponse.IsError))

	// listRequest
	listRequest := types.NewRequest().
		SetMetadataKeyValue("method", "list_tags")
	queryListResponse, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.redshift.service").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(queryListResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("list tags executed, response: %s", listResponse.Data))

}
