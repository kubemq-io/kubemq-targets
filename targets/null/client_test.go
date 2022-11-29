package null

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
)

func TestClient_Do(t *testing.T) {
	type fields struct {
		Delay         time.Duration
		DoError       error
		ResponseError error
	}
	tests := []struct {
		name    string
		fields  fields
		req     *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "do",
			fields: fields{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			req:     types.NewRequest().SetData([]byte("data")),
			want:    types.NewResponse().SetData([]byte("data")),
			wantErr: false,
		},
		{
			name: "do with DoError",
			fields: fields{
				Delay:         0,
				DoError:       fmt.Errorf("do-error"),
				ResponseError: nil,
			},
			req:     types.NewRequest().SetData([]byte("data")),
			want:    nil,
			wantErr: true,
		},
		{
			name: "do with response error",
			fields: fields{
				Delay:         0,
				DoError:       nil,
				ResponseError: fmt.Errorf("response-error"),
			},
			req:     types.NewRequest().SetData([]byte("data")),
			want:    nil,
			wantErr: true,
		},
		{
			name: "do cancel ctx",
			fields: fields{
				Delay:         5 * time.Second,
				DoError:       nil,
				ResponseError: nil,
			},
			req:     types.NewRequest().SetData([]byte("data")),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			c := &Client{
				Delay:         tt.fields.Delay,
				DoError:       tt.fields.DoError,
				ResponseError: tt.fields.ResponseError,
			}
			_ = c.Init(ctx, config.Spec{
				Name:       "null",
				Kind:       "",
				Properties: nil,
			}, nil)
			got, err := c.Do(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Do() got = %v, want %v", got, tt.want)
			}
		})
	}
}
