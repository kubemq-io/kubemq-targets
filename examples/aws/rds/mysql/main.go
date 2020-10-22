package main

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"log"
	"time"
)

var (
	transactionString = `DROP TABLE IF EXISTS post;
	       CREATE TABLE post (
	         ID bigint,
	         TITLE varchar(40),
	         CONTENT varchar(255),
			 BIGNUMBER bigint,
			 BOOLVALUE boolean,
	         CONSTRAINT pk_post PRIMARY KEY(ID)
	       );
	       INSERT INTO post(ID,TITLE,CONTENT,BIGNUMBER,BOOLVALUE) VALUES
	                       (0,NULL,'Content One',1231241241231231123,true),
	                       (1,'Title Two','Content Two',123125241231231123,false);`
	queryString = `SELECT id,title,content,bignumber,boolvalue FROM post;`
)

func main() {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("kubemq-cluster", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}

	transactionRequest := types.NewRequest().
		SetMetadataKeyValue("method", "transaction").
		SetData([]byte(transactionString))
	queryTransactionResponse, err := client.SetQuery(transactionRequest.ToQuery()).
		SetChannel("query.aws.rds.mysql").
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
		SetChannel("query.aws.rds.mysql").
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
