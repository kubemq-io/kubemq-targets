package sns

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Attributes struct {
	Name        string `json:"name"`
	DataType    string `json:"data_type"`
	StringValue string `json:"string_value"`
}

func (c *Client) createSNSMessage(meta metadata, data []byte) (*sns.PublishInput, error) {
	var s []Attributes
	if data != nil {
		err := json.Unmarshal(data, &s)
		if err != nil {
			return nil, err
		}
	}
	a := make(map[string]*sns.MessageAttributeValue)
	for _, i := range s {
		if _, ok := a[i.Name]; ok {
			return nil, fmt.Errorf("duplicate name value in attibutes")
		}
		a[i.Name] = &sns.MessageAttributeValue{
			DataType:    aws.String(i.DataType),
			StringValue: aws.String(i.StringValue),
		}
	}
	i := &sns.PublishInput{
		TopicArn:          aws.String(meta.topic),
		Message:           aws.String(meta.message),
		PhoneNumber:       aws.String(meta.phoneNumber),
		Subject:           aws.String(meta.subject),
		TargetArn:         aws.String(meta.targetArn),
		MessageAttributes: a,
	}
	return i, nil
}
