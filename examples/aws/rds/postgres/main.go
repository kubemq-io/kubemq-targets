package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
)

var (
	transactionString = `	DROP TABLE IF EXISTS post;
	       CREATE TABLE post (
	         ID serial,
	         TITLE varchar(40),
	         CONTENT varchar(255),
	         CONSTRAINT pk_post PRIMARY KEY(ID)
	       );
	       INSERT INTO post(ID,TITLE,CONTENT) VALUES
	                       (1,NULL,'Content One'),
	                       (2,'Title Two','Content Two');`
	queryString = `SELECT id,title,content FROM post;`
)

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster:50000", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}

	transactionRequest := types.NewRequest().
		SetMetadataKeyValue("method", "transaction").
		SetData([]byte(transactionString))
	queryTransactionResponse, err := client.SetQuery(transactionRequest.ToQuery()).
		SetChannel("query.aws.rds.postgres").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	transactionResponse, err := types.ParseResponse(queryTransactionResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("transaction request result: %s ", transactionResponse.Metadata.String()))

	queryRequest := types.NewRequest().
		SetMetadataKeyValue("method", "query").
		SetData([]byte(queryString))

	queryResponse, err := client.SetQuery(queryRequest.ToQuery()).
		SetChannel("query.aws.rds.postgres").
		SetTimeout(10 * time.Second).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	response, err := types.ParseResponse(queryResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("query request results: %s ", string(response.Data)))
}
