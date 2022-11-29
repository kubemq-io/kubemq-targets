package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/tableName.txt")
	if err != nil {
		panic(err)
	}
	tableName := string(dat)
	dat, err = ioutil.ReadFile("./credentials/querySpanner.txt")
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
	// set request
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "query").
		SetMetadataKeyValue("query", query)
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.gcp.spanner").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	queryResponse, err := types.ParseResponse(querySetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get query request for query: %s executed, response: %s", query, queryResponse.Data))
	// read

	columnNames := []string{"id", "name"}
	b, err := json.Marshal(columnNames)
	if err != nil {
		log.Fatal(err)
	}
	readRequest := types.NewRequest().
		SetMetadataKeyValue("method", "read").
		SetMetadataKeyValue("table_name", tableName).
		SetData(b)

	read, err := client.SetQuery(readRequest.ToQuery()).
		SetChannel("query.gcp.spanner").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getReadResponse, err := types.ParseResponse(read.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("read executed, response: %s", getReadResponse.Data))
}
