package main

import (
	"context"
	"encoding/json"
	"fmt"
	aero "github.com/aerospike/aerospike-client-go"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
	"github.com/kubemq-hub/kubemq-targets/targets/stores/aerospike"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"
)

func main() {
	k := aerospike.PutRequest{
		UserKey:   "user_key1",
		KeyName:   "some-key",
		Namespace: "test",
		BinMap: // define some bins with data
		aero.BinMap{
			"bin1": 42,
			"bin2": "An elephant is a mouse with an operating system",
			"bin3": []interface{}{"Go", 2009},
		},
	}
	req, err := json.Marshal(k)
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	randomKey := uuid.New().String()
	// set request
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set").
		SetData(req)
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.aerospike").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	setResponse, err := types.ParseResponse(querySetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("set request for key: %s executed, response: %s", randomKey, setResponse.Metadata.String()))

	// get request
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("namespace", "test").
		SetMetadataKeyValue("key", "some-key").
		SetMetadataKeyValue("user_key", "user_key1")

	queryGetResponse, err := client.SetQuery(getRequest.ToQuery()).
		SetChannel("query.aerospike").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	getResponse, err := types.ParseResponse(queryGetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("get request for key: %s executed, response: %s, data: %s", randomKey, getResponse.Metadata.String(), string(getResponse.Data)))
}
