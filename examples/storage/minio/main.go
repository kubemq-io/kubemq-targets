package main

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"
)

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// make bucket

	makeRequest := types.NewRequest().
		SetMetadataKeyValue("method", "make_bucket").
		SetMetadataKeyValue("param1", "testbucket")

	response, err := client.SetQuery(makeRequest.ToQuery()).
		SetChannel("query.minio").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	makeResponse, err := types.ParseResponse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("make bucket request response: %s", makeResponse.Metadata.String()))

	// put object
	putRequest := types.NewRequest().
		SetMetadataKeyValue("method", "put").
		SetMetadataKeyValue("param1", "testbucket").
		SetMetadataKeyValue("param2", "object").
		SetData([]byte("object"))

	response, err = client.SetQuery(putRequest.ToQuery()).
		SetChannel("query.minio").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	putResponse, err := types.ParseResponse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("put object request response: %s", putResponse.Metadata.String()))

	// get list of object
	getListRequest := types.NewRequest().
		SetMetadataKeyValue("method", "list_objects").
		SetMetadataKeyValue("param1", "testbucket")

	response, err = client.SetQuery(getListRequest.ToQuery()).
		SetChannel("query.minio").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getListResponse, err := types.ParseResponse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get list object request response: %s", string(getListResponse.Data)))

	// get object
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("param1", "testbucket").
		SetMetadataKeyValue("param2", "object")

	response, err = client.SetQuery(getRequest.ToQuery()).
		SetChannel("query.minio").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get object request response: %s", string(getResponse.Data)))

	// remove object
	removeRequest := types.NewRequest().
		SetMetadataKeyValue("method", "remove").
		SetMetadataKeyValue("param1", "testbucket").
		SetMetadataKeyValue("param2", "object")

	response, err = client.SetQuery(removeRequest.ToQuery()).
		SetChannel("query.minio").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	removeResponse, err := types.ParseResponse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("remove object request response: %s", removeResponse.Metadata.String()))

	// remove bucket
	removeBucketRequest := types.NewRequest().
		SetMetadataKeyValue("method", "remove_bucket").
		SetMetadataKeyValue("param1", "testbucket")

	response, err = client.SetQuery(removeBucketRequest.ToQuery()).
		SetChannel("query.minio").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	removeBucketResponse, err := types.ParseResponse(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("remove bucket request response: %s", removeBucketResponse.Metadata.String()))
}
