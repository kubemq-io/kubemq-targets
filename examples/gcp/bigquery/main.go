package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/query.txt")
	if err != nil {
		panic(err)
	}
	query := string(dat)
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// query
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "query").
		SetMetadataKeyValue("query", query)
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.gcp.bigquery").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	setResponse, err := types.ParseResponse(querySetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get query request for query: %s executed, response: %s", query, setResponse.Data))
	// get data sets

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get_data_sets")

	getDataSets, err := client.SetQuery(delRequest.ToQuery()).
		SetChannel("query.gcp.bigquery").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getSetsResponse, err := types.ParseResponse(getDataSets.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get data sets executed, response: %s", getSetsResponse.Data))
}
