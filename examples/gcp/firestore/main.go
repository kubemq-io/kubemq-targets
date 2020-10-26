package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/topicID.txt")
	if err != nil {
		panic(err)
	}
	topicID := string(dat)
	user := map[string]interface{}{
		"first": "test-kubemq",
		"last":  "test-kubemq-last",
		"id":    123,
	}
	bUser, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	// add file
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "add").
		SetMetadataKeyValue("collection", "my_collection").
		SetData(bUser)
	querySetResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.gcp.firestore").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	setResponse, err := types.ParseResponse(querySetResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("add file request for topic : %s executed, is_error: %v", topicID, setResponse.IsError))

	dat, err = ioutil.ReadFile("./credentials/deleteKey.txt")
	if err != nil {
		panic(err)
	}
	deleteKey := string(dat)
	// delete file
	deleteRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete_document_key").
		SetMetadataKeyValue("item", deleteKey).
		SetMetadataKeyValue("collection", "my_collection")

	queryDeleteResponse, err := client.SetQuery(deleteRequest.ToQuery()).
		SetChannel("query.gcp.firestore").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	deleteResponse, err := types.ParseResponse(queryDeleteResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("delete file request for topic : %s executed, with id %s , is_error: %v", topicID, deleteKey, deleteResponse.IsError))
}
