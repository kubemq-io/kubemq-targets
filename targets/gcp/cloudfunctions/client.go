package cloudfunctions

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubemq-hub/kubemq-targets/config"
	gf "github.com/kubemq-hub/kubemq-targets/targets/gcp/cloudfunctions/functions/apiv1"
	"github.com/kubemq-hub/kubemq-targets/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	functionspb "google.golang.org/genproto/googleapis/cloud/functions/v1"
)

type Client struct {
	name           string
	opts           options
	client         *gf.CloudFunctionsClient
	parrantProject string
	list           []string
	//nameFunctions  map[string]string
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	var err error

	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	b := []byte(c.opts.credentials)

	client, err := gf.NewCloudFunctionsClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	c.client = client
	c.parrantProject = c.opts.parrentProject

	if c.opts.locationMatch {
		it := client.ListFunctions(ctx, &functionspb.ListFunctionsRequest{
			Parent: fmt.Sprintf("projects/%s/locations/-", c.opts.parrentProject),
		})

		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				// TODO: Handle error.
			}
			if resp != nil {
				c.list = append(c.list, resp.GetName())
				/*
					if strings.Contains(fmt.Sprintf("%v", resp.GetTrigger()), "event_type") {
						c.nameFunctions[resp.GetName()] = "event_type"
					} else {
						c.nameFunctions[resp.GetName()] = "url"
					}
				*/
			}
		}
	}
	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	m, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}

	if m.project == "" {
		m.project = c.parrantProject
	}

	name := fmt.Sprintf("projects/%s/locations/%s/functions/%s", m.project, m.location, m.name)
	if m.location == "" {
		for _, n := range c.list {
			if strings.Contains(n, m.name) && strings.Contains(n, m.project) {
				m.location = "added from match"
				name = n
				break
			}
		}
	}
	if m.location == "" {
		return nil, fmt.Errorf("no location found for function")
	}

	// if c.nameFunctions[name] == "event_type" {
	// 	data = b64.StdEncoding.EncodeToString(request.Data)
	// }

	cfo := &functionspb.CallFunctionRequest{
		Name: name,
		Data: string(request.Data),
	}

	res, err := c.client.CallFunction(ctx, cfo)
	if err != nil {
		return nil, err
	}
	if res.Error != "" {
		return nil, fmt.Errorf(res.Error)
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", res.Result).
		SetMetadataKeyValue("executionid", res.ExecutionId).
		SetData([]byte(res.Result)), nil

}
