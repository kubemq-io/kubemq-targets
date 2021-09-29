package firebase

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type testStructure struct {
	projectID string
	dbName    string
	cred      string
	token     string
	uid       string
	uid2      string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/firebaseDBName.txt")
	if err != nil {
		return nil, err
	}
	t.dbName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/projectID.txt")
	if err != nil {
		return nil, err
	}
	t.projectID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	if err != nil {
		return nil, err
	}
	t.cred = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/token.txt")
	if err != nil {
		return nil, err
	}
	t.token = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/uid.txt")
	if err != nil {
		return nil, err
	}
	t.uid = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/uid2.txt")
	if err != nil {
		return nil, err
	}
	t.uid2 = fmt.Sprintf("%s", dat)
	return t, nil
}

func TestClient_Init(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init - db client",
			cfg: config.Spec{
				Name: "gcp-firebase",
				Kind: "gcp.firebase",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
					"db_url":      dat.dbName,
					"db_client":   "true",
				},
			},
			wantErr: false,
		},
		{
			name: "init - auth client",
			cfg: config.Spec{
				Name: "gcp-firebase",
				Kind: "gcp.firebase",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
					"auth_client": "true",
				},
			},
			wantErr: false,
		},
		{
			name: "init - multiclient",
			cfg: config.Spec{
				Name: "gcp-firebase",
				Kind: "gcp.firebase",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"credentials": dat.cred,
					"db_url":      dat.dbName,
					"auth_client": "true",
					"db_client":   "true",
				},
			},
			wantErr: false,
		},
		{
			name: "init-firebase-project-id - missing project_id",
			cfg: config.Spec{
				Name: "gcp-firebase",
				Kind: "gcp.firebase",
				Properties: map[string]string{
					"credentials": dat.cred,
					"auth_client": "true",
					"db_client":   "true",
				},
			},
			wantErr: true,
		},
		{
			name: "init-firebase-instance missing credentials",
			cfg: config.Spec{
				Name: "gcp-firebase",
				Kind: "gcp.firebase",
				Properties: map[string]string{
					"project_id":  dat.projectID,
					"auth_client": "true",
					"db_client":   "true",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			err := c.Init(ctx, tt.cfg, nil)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)

		})
	}
}
