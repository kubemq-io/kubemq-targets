package main

import (
	"context"
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
	contents, err := ioutil.ReadFile("./credentials/azure/storage/files/contents.txt")
	if err != nil {
		panic(err)
	}
	dat, err := ioutil.ReadFile("./credentials/azure/storage/files/serviceURL.txt")
	if err != nil {
		panic(err)
	}
	serviceURL := string(dat)
	// upload
	uploadRequest := types.NewRequest().
		SetMetadataKeyValue("method", "upload").
		SetMetadataKeyValue("service_url", serviceURL).
		SetData(contents)
	queryUploadResponse, err := client.SetQuery(uploadRequest.ToQuery()).
		SetChannel("azure.storage.files").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	uploadResponse, err := types.ParseResponse(queryUploadResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("upload request executed, error: %v", uploadResponse.Error))

	// get request
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("service_url", serviceURL)

	getQueryResponse, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("azure.storage.files").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(getQueryResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get request done, response : %s", getResponse.Data))
}
