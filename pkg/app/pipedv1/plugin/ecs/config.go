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
	"fmt"
)

// ECSDeployStageOptions contains configurable values for an ECS_DEPLOY stage.
type ECSDeployStageOptions struct {
	// Region is the AWS region where the ECS cluster is located.
	Region string `json:"region"`
	// Cluster is the name or ARN of the ECS cluster.
	Cluster string `json:"cluster"`
	// Service is the name of the ECS service to update.
	Service string `json:"service"`
	// TaskDefArn is the ARN of the task definition to deploy.
	TaskDefArn string `json:"taskDefArn"`
}

// validate validates the ECSDeployStageOptions.
func (o ECSDeployStageOptions) validate() error {
	if o.Region == "" {
		return fmt.Errorf("region is required")
	}
	if o.Cluster == "" {
		return fmt.Errorf("cluster is required")
	}
	if o.Service == "" {
		return fmt.Errorf("service is required")
	}
	if o.TaskDefArn == "" {
		return fmt.Errorf("taskDefArn is required")
	}
	return nil
}

// decodeECSDeployOptions decodes the raw JSON data and validates it.
func decodeECSDeployOptions(data json.RawMessage) (ECSDeployStageOptions, error) {
	var opts ECSDeployStageOptions
	if err := json.Unmarshal(data, &opts); err != nil {
		return ECSDeployStageOptions{}, fmt.Errorf("failed to unmarshal the stage config: %w", err)
	}
	if err := opts.validate(); err != nil {
		return ECSDeployStageOptions{}, fmt.Errorf("failed to validate the stage config: %w", err)
	}
	return opts, nil
}
