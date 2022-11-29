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
	i := &sns.PublishInput{
		Message: aws.String(meta.message),
	}
	if meta.topic != "" {
		i.TopicArn = aws.String(meta.topic)
	} else {
		i.TargetArn = aws.String(meta.targetArn)
	}

	if meta.phoneNumber != "" {
		i.PhoneNumber = aws.String(meta.phoneNumber)
	}

	if meta.subject != "" {
		i.Subject = aws.String(meta.subject)
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
	if len(a) > 0 {
		i.MessageAttributes = a
	}

	return i, nil
}
