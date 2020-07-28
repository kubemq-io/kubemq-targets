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
	m := make(map[string]interface{})
	m["some_key"] = "some_value"
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	// get value
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get_db").
		SetMetadataKeyValue("ref_path", "test")
	queryGetResponse, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("query.gcp.firebase").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(queryGetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get request for ref path : est, is_error: %v", getResponse.IsError))

	// set value
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set_db").
		SetMetadataKeyValue("ref_path", "test").
		SetData(b)

	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.gcp.firebase").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	setResponse, err := types.ParseResponse(querySetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("set request for ref path :test is_error: %v", setResponse.IsError))
}
