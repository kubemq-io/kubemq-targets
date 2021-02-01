package main

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
	"github.com/kubemq-hub/kubemq-targets/types"
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
	dat, err := ioutil.ReadFile("./credentials/aws/athena/db.txt")
	if err != nil {
		log.Fatal(err)
	}
	db := string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/athena/catalog.txt")
	if err != nil {
		log.Fatal(err)
	}
	catalog := string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/athena/query.txt")
	if err != nil {
		log.Fatal(err)
	}
	query := string(dat)

	dat, err = ioutil.ReadFile("./credentials/aws/athena/outputLocation.txt")
	if err != nil {
		log.Fatal(err)
	}
	outputLocation := string(dat)
	// query Request
	queryRequest := types.NewRequest().
		SetMetadataKeyValue("method", "query").
		SetMetadataKeyValue("db", db).
		SetMetadataKeyValue("output_location", outputLocation).
		SetMetadataKeyValue("catalog", catalog).
		SetMetadataKeyValue("query", query)
	queryResponse, err := client.SetQuery(queryRequest.ToQuery()).
		SetChannel("query.aws.athena").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	qryResponse, err := types.ParseResponse(queryResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	executionCode := qryResponse.Metadata["execution_id"]
	log.Println(fmt.Sprintf("qry executed, executionCode: %s", executionCode))
	//Give query time to end
	time.Sleep(2 * time.Second)
	// get query result
	getResultRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get_query_result").
		SetMetadataKeyValue("execution_id", executionCode)

	getResult, err := client.SetQuery(getResultRequest.ToQuery()).
		SetChannel("query.aws.athena").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(getResult.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get resoponse %v, error: %v", getResponse.Data, getResponse.IsError))
}
