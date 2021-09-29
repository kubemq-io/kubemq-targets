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
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// listRequest
	listRequest := types.NewRequest().
		SetMetadataKeyValue("method", "list_streams")
	queryListResponse, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.kinesis").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(queryListResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("list_streams executed, response: %s", listResponse.Data))

	dat, err := ioutil.ReadFile("./credentials/aws/kinesis/streamName.txt")
	if err != nil {
		panic(err)
	}
	streamName := string(dat)

	dat, err = ioutil.ReadFile("./credentials/aws/kinesis/partitionKey.txt")
	if err != nil {
		panic(err)
	}
	partitionKey := string(dat)
	// putRecord
	putRequest := types.NewRequest().
		SetMetadataKeyValue("method", "put_record").
		SetMetadataKeyValue("partition_key", partitionKey).
		SetMetadataKeyValue("stream_name", streamName).
		SetData([]byte("{\"my_result\":\"ok\"})"))
	queryPutResponse, err := client.SetQuery(putRequest.ToQuery()).
		SetChannel("query.aws.kinesis").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	putResponse, err := types.ParseResponse(queryPutResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("put query executed, response: %s", putResponse.Data))

}
