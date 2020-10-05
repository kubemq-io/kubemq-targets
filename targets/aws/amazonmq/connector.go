package amazonmq

import (
	"github.com/kubemq-hub/builder/common"
)

// TODO
func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.aws.amazonmq").
		SetDescription("AWS AmazonMQ Target")
}
