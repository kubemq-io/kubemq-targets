package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
)

func getUserStruct() map[string]interface{} {
	m := make(map[string]interface{})
	m["disabled"] = false
	m["display_name"] = "test"
	m["email"] = "fakeuserkube123@test.com"
	m["email_verified"] = true
	m["password"] = "testPassword"
	m["phone_number"] = "+12343678123"
	m["photo_url"] = "https://kubemq.io/wp-content/uploads/2018/11/24350KubeMQ_clean.png"

	return m
}

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}

	u := getUserStruct()
	b, err := json.Marshal(&u)
	if err != nil {
		log.Fatal(err)
	}
	// create user
	createRequest := types.NewRequest().
		SetMetadataKeyValue("method", "create_user").
		SetData(b)
	queryCreateResponse, err := client.SetQuery(createRequest.ToQuery()).
		SetChannel("query.gcp.firebase").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	createResponse, err := types.ParseResponse(queryCreateResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("create request done, is_error: %v", createResponse.IsError))

	// retrieve by email value
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "retrieve_user").
		SetMetadataKeyValue("email", fmt.Sprintf("%s", u["email"])).
		SetMetadataKeyValue("retrieve_by", "by_email")

	queryRetrieveResponse, err := client.SetQuery(setRequest.ToQuery()).
		SetChannel("query.gcp.firebase").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	retrieveResponse, err := types.ParseResponse(queryRetrieveResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("retrieve user for ref path :test is_error: %v , user :%s", retrieveResponse.IsError, retrieveResponse.Data))
}
