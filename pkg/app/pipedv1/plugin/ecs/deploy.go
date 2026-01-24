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
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

// ecsClient is an interface for ECS operations to allow mocking in tests.
type ecsClient interface {
	UpdateService(ctx context.Context, params *ecs.UpdateServiceInput, optFns ...func(*ecs.Options)) (*ecs.UpdateServiceOutput, error)
}

// executeECSDeploy executes the ECS_DEPLOY stage.
func executeECSDeploy(ctx context.Context, request sdk.ExecuteStageRequest[struct{}], lp sdk.StageLogPersister) sdk.StageStatus {
	lp.Infof("Start executing the ECS_DEPLOY stage")

	// Decode and validate stage options
	opts, err := decodeECSDeployOptions(request.StageConfig)
	if err != nil {
		lp.Errorf("Failed to decode stage config: %v", err)
		return sdk.StageStatusFailure
	}

	lp.Infof("Deploying to ECS:")
	lp.Infof("  Region: %s", opts.Region)
	lp.Infof("  Cluster: %s", opts.Cluster)
	lp.Infof("  Service: %s", opts.Service)
	lp.Infof("  Task Definition: %s", opts.TaskDefArn)

	// Create ECS client
	client, err := createECSClient(ctx, opts.Region)
	if err != nil {
		lp.Errorf("Failed to create ECS client: %v", err)
		return sdk.StageStatusFailure
	}

	// Execute the deployment
	if err := updateECSService(ctx, client, opts, lp); err != nil {
		lp.Errorf("Failed to update ECS service: %v", err)
		return sdk.StageStatusFailure
	}

	lp.Infof("Successfully triggered ECS deployment")
	lp.Infof("The ECS service will now deploy the new task definition")

	return sdk.StageStatusSuccess
}

// createECSClient creates an ECS client using the default credential chain.
// This supports AWS credentials from environment variables, IRSA, or instance profiles.
func createECSClient(ctx context.Context, region string) (*ecs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return ecs.NewFromConfig(cfg), nil
}

// updateECSService calls UpdateService on ECS to trigger a new deployment.
func updateECSService(ctx context.Context, client ecsClient, opts ECSDeployStageOptions, lp sdk.StageLogPersister) error {
	lp.Infof("Calling ECS UpdateService with ForceNewDeployment=true...")

	input := &ecs.UpdateServiceInput{
		Cluster:            aws.String(opts.Cluster),
		Service:            aws.String(opts.Service),
		TaskDefinition:     aws.String(opts.TaskDefArn),
		ForceNewDeployment: true,
	}

	output, err := client.UpdateService(ctx, input)
	if err != nil {
		return fmt.Errorf("ECS UpdateService failed: %w", err)
	}

	if output.Service != nil {
		lp.Infof("ECS service update accepted:")
		if output.Service.ServiceArn != nil {
			lp.Infof("  Service ARN: %s", *output.Service.ServiceArn)
		}
		if output.Service.TaskDefinition != nil {
			lp.Infof("  New Task Definition: %s", *output.Service.TaskDefinition)
		}
		lp.Infof("  Desired Count: %d", output.Service.DesiredCount)
		lp.Infof("  Running Count: %d", output.Service.RunningCount)
	}

	return nil
}
