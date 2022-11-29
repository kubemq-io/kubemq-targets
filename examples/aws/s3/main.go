package main

import (
	"context"
	"fmt"
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
	BucketName := "testmykubemqbucketname"
	createRequest := types.NewRequest().
		SetMetadataKeyValue("method", "create_bucket").
		SetMetadataKeyValue("bucket_name", BucketName).
		SetMetadataKeyValue("wait_for_completion", "true")

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
