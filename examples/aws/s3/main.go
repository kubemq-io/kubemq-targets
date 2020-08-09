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
	// listRequest
	listRequest := types.NewRequest().
		SetMetadataKeyValue("method", "list_buckets")
	queryListResponse, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.s3").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(queryListResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("list buckets executed, response: %s", listResponse.Data))

	// Create Bucket
	attributes := make(map[string]*string)
	BucketName := "testmykubemqbucketname"
	b, err := json.Marshal(attributes)
	if err != nil {
		log.Fatal(err)
	}
	createRequest := types.NewRequest().
		SetMetadataKeyValue("method", "create_bucket").
		SetMetadataKeyValue("bucket_name", BucketName).
		SetMetadataKeyValue("wait_for_completion", "true").
		SetData(b)

	getCreate, err := client.SetQuery(createRequest.ToQuery()).
		SetChannel("query.aws.s3").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	createResponse, err := types.ParseResponse(getCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("create bucket executed, error: %v", createResponse.IsError))
}
