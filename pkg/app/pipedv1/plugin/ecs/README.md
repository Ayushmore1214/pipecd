# PipeCD ECS Deploy Plugin

This plugin adds support for deploying to Amazon ECS by implementing the `ECS_DEPLOY` stage type for PipeCD pipelines.

## Stage Type

- `ECS_DEPLOY` - Triggers a new deployment to Amazon ECS by calling UpdateService with a specified task definition.

## Configuration

The `ECS_DEPLOY` stage accepts the following configuration:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | Yes | AWS region where the ECS cluster is located |
| `cluster` | string | Yes | Name or ARN of the ECS cluster |
| `service` | string | Yes | Name of the ECS service to update |
| `taskDefArn` | string | Yes | ARN of the task definition to deploy |

## Example Pipeline Configuration

```yaml
stages:
  - name: deploy
    type: ECS_DEPLOY
    config:
      region: ap-south-1
      cluster: my-cluster
      service: my-service
      taskDefArn: arn:aws:ecs:ap-south-1:123456789012:task-definition/my-task:1
```

## AWS Credentials

The plugin uses the AWS SDK v2 default credential chain, which supports:

- Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
- IAM Roles for Service Accounts (IRSA) when running on EKS
- EC2 Instance Profiles
- Shared credentials file (`~/.aws/credentials`)
- ECS container credentials

## Behavior

When executed, the plugin will:

1. Parse and validate the stage configuration
2. Create an AWS ECS client using the default credential chain
3. Call `UpdateService` with:
   - The specified cluster
   - The specified service
   - The specified task definition ARN
   - `ForceNewDeployment = true`
4. Return `SUCCESS` if ECS accepts the deployment, `FAILURE` otherwise
5. Log meaningful messages to the PipeCD UI for deployment visibility

## Limitations

This is a minimal implementation focused on triggering deployments. The following features are not yet implemented:

- Rollout status tracking
- Health check verification
- Drift detection
- Rollback stages
- Canary/progressive deployments
- Traffic routing

## Building

```bash
cd pkg/app/pipedv1/plugin/ecs
go build -o ecs-plugin .
```

## Testing

```bash
cd pkg/app/pipedv1/plugin/ecs
go test -v ./...
```
