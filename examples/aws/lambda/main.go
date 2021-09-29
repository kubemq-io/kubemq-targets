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
		SetMetadataKeyValue("method", "list")
	queryListResponse, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.lambda").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(queryListResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("list functions executed, response: %s", listResponse.Data))

	// Delete Lambda
	dat, err := ioutil.ReadFile("./credentials/aws/lambda/functionName.txt")
	if err != nil {
		panic(err)
	}
	functionName := string(dat)
	deleteRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("function_name", functionName)

	getDelete, err := client.SetQuery(deleteRequest.ToQuery()).
		SetChannel("query.aws.lambda").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	deleteResponse, err := types.ParseResponse(getDelete.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("delete lambda executed, error: %v", deleteResponse.IsError))
}
