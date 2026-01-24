// Copyright 2025 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestECSDeployStageOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    ECSDeployStageOptions
		wantErr string
	}{
		{
			name: "valid options",
			opts: ECSDeployStageOptions{
				Region:     "ap-south-1",
				Cluster:    "my-cluster",
				Service:    "my-service",
				TaskDefArn: "arn:aws:ecs:ap-south-1:123456789012:task-definition/my-task:1",
			},
			wantErr: "",
		},
		{
			name: "missing region",
			opts: ECSDeployStageOptions{
				Cluster:    "my-cluster",
				Service:    "my-service",
				TaskDefArn: "arn:aws:ecs:ap-south-1:123456789012:task-definition/my-task:1",
			},
			wantErr: "region is required",
		},
		{
			name: "missing cluster",
			opts: ECSDeployStageOptions{
				Region:     "ap-south-1",
				Service:    "my-service",
				TaskDefArn: "arn:aws:ecs:ap-south-1:123456789012:task-definition/my-task:1",
			},
			wantErr: "cluster is required",
		},
		{
			name: "missing service",
			opts: ECSDeployStageOptions{
				Region:     "ap-south-1",
				Cluster:    "my-cluster",
				TaskDefArn: "arn:aws:ecs:ap-south-1:123456789012:task-definition/my-task:1",
			},
			wantErr: "service is required",
		},
		{
			name: "missing taskDefArn",
			opts: ECSDeployStageOptions{
				Region:  "ap-south-1",
				Cluster: "my-cluster",
				Service: "my-service",
			},
			wantErr: "taskDefArn is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

func TestDecodeECSDeployOptions(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    ECSDeployStageOptions
		wantErr string
	}{
		{
			name: "valid JSON",
			json: `{
				"region": "ap-south-1",
				"cluster": "my-cluster",
				"service": "my-service",
				"taskDefArn": "arn:aws:ecs:ap-south-1:123456789012:task-definition/my-task:1"
			}`,
			want: ECSDeployStageOptions{
				Region:     "ap-south-1",
				Cluster:    "my-cluster",
				Service:    "my-service",
				TaskDefArn: "arn:aws:ecs:ap-south-1:123456789012:task-definition/my-task:1",
			},
			wantErr: "",
		},
		{
			name:    "invalid JSON",
			json:    `{invalid}`,
			want:    ECSDeployStageOptions{},
			wantErr: "failed to unmarshal the stage config",
		},
		{
			name: "missing required field",
			json: `{
				"region": "ap-south-1",
				"cluster": "my-cluster"
			}`,
			want:    ECSDeployStageOptions{},
			wantErr: "service is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := decodeECSDeployOptions(json.RawMessage(tt.json))
			if tt.wantErr == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.want, opts)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}
