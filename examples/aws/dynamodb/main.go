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
	// listRequest
	listRequest := types.NewRequest().
		SetMetadataKeyValue("method", "list_tables")
	queryListResponse, err := client.SetQuery(listRequest.ToQuery()).
		SetChannel("query.aws.dynamodb").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listResponse, err := types.ParseResponse(queryListResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("list tables executed, response: %s", listResponse.Data))
	dat, err := ioutil.ReadFile("./credentials/aws/dynamodb/tableName.txt")
	if err != nil {
		panic(err)
	}
	tableName := string(dat)
	//Create Table
	input := fmt.Sprintf(`{
					"AttributeDefinitions": [
						{
							"AttributeName": "Year",
							"AttributeType": "N"
						},
						{
							"AttributeName": "Title",
							"AttributeType": "S"
						}
					],
					"BillingMode": null,
					"GlobalSecondaryIndexes": null,
					"KeySchema": [
						{
							"AttributeName": "Year",
							"KeyType": "HASH"
						},
						{
							"AttributeName": "Title",
							"KeyType": "RANGE"
						}
					],
					"LocalSecondaryIndexes": null,
					"ProvisionedThroughput": {
						"ReadCapacityUnits": 10,
						"WriteCapacityUnits": 10
					},
					"SSESpecification": null,
					"StreamSpecification": null,
					"TableName": "%s",
					"Tags": null
				}`, tableName)
	createRequest := types.NewRequest().
		SetMetadataKeyValue("method", "create_table").
		SetData([]byte(input))

	getCreate, err := client.SetQuery(createRequest.ToQuery()).
		SetChannel("query.aws.dynamodb").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	createResponse, err := types.ParseResponse(getCreate.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("create table executed, error: %v", createResponse.IsError))
}
