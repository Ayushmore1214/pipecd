# API Documentation Guide

This guide explains how to use and contribute to PipeCD's API documentation.

## Table of Contents

- [API Overview](#api-overview)
- [gRPC API](#grpc-api)
- [REST API](#rest-api)
- [API Authentication](#api-authentication)
- [Using the API](#using-the-api)
- [API Client Libraries](#api-client-libraries)
- [Contributing to API Docs](#contributing-to-api-docs)

## API Overview

PipeCD exposes two types of APIs:

1. **gRPC API**: Primary API for piped agents and programmatic access
2. **REST API**: HTTP endpoints for web UI and simple integrations

### API Versioning

PipeCD follows semantic versioning for its APIs:
- **Major version**: Breaking changes
- **Minor version**: New features, backward compatible
- **Patch version**: Bug fixes

Current API version: **v1beta1**

## gRPC API

### Service Definitions

PipeCD's gRPC services are defined in Protocol Buffer files:

```
pkg/app/api/
├── api.proto              # Common types
├── pipedservice.proto     # Piped agent API
├── webservice.proto       # Web UI API
└── service.proto          # Control plane API
```

### Generating API Clients

```bash
# Generate Go clients
make gen/api

# Generate TypeScript clients for web
make gen/web-api
```

### Available Services

#### 1. PipedService

Used by piped agents to communicate with control plane.

**Key methods**:
- `RegisterPiped`: Register a new piped agent
- `ReportDeployment`: Report deployment status
- `GetDeployment`: Retrieve deployment details
- `ListDeployments`: List deployments

**Example**:

```go
import (
    "context"
    "google.golang.org/grpc"
    pipedservice "github.com/pipe-cd/pipecd/pkg/app/api/pipedservice"
)

func main() {
    conn, err := grpc.Dial("control-plane:443", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pipedservice.NewPipedServiceClient(conn)
    
    resp, err := client.GetDeployment(context.Background(), &pipedservice.GetDeploymentRequest{
        DeploymentId: "deployment-123",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Deployment: %v\n", resp.Deployment)
}
```

#### 2. WebService

Used by the web UI and external integrations.

**Key methods**:
- `GetProject`: Get project details
- `ListApplications`: List applications
- `GetApplication`: Get application details
- `SyncApplication`: Trigger application sync

**Example**:

```go
import (
    webservice "github.com/pipe-cd/pipecd/pkg/app/api/webservice"
)

client := webservice.NewWebServiceClient(conn)

resp, err := client.ListApplications(ctx, &webservice.ListApplicationsRequest{
    ProjectId: "my-project",
})
```

#### 3. APIService

General API service for administrative tasks.

**Key methods**:
- `GetInsight`: Get deployment insights
- `GetDeploymentStatistics`: Get deployment statistics
- `ListDeploymentConfigTemplates`: List config templates

### API Request/Response Examples

#### Register Piped

**Request**:
```protobuf
message RegisterPipedRequest {
    string name = 1;
    string desc = 2;
    repeated string envIds = 3;
}
```

**Response**:
```protobuf
message RegisterPipedResponse {
    string id = 1;
    string key = 2;
}
```

#### Get Deployment

**Request**:
```protobuf
message GetDeploymentRequest {
    string deployment_id = 1;
}
```

**Response**:
```protobuf
message GetDeploymentResponse {
    Deployment deployment = 1;
}
```

## REST API

### Base URL

```
https://your-pipecd-instance.com/api/v1
```

### Authentication

All REST API requests require authentication via API key or session token:

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
     https://pipecd.example.com/api/v1/projects
```

### Endpoints

#### Projects

```bash
# List projects
GET /api/v1/projects

# Get project
GET /api/v1/projects/{project_id}
```

#### Applications

```bash
# List applications
GET /api/v1/projects/{project_id}/applications

# Get application
GET /api/v1/applications/{app_id}

# Sync application
POST /api/v1/applications/{app_id}/sync
```

#### Deployments

```bash
# List deployments
GET /api/v1/deployments?project_id={project_id}

# Get deployment
GET /api/v1/deployments/{deployment_id}

# Cancel deployment
POST /api/v1/deployments/{deployment_id}/cancel
```

#### Pipeds

```bash
# List pipeds
GET /api/v1/pipeds?project_id={project_id}

# Get piped
GET /api/v1/pipeds/{piped_id}

# Recreate piped key
POST /api/v1/pipeds/{piped_id}/recreate-key
```

## API Authentication

### Using API Keys

1. **Create API Key** (via Web UI):
   - Navigate to Settings → API Keys
   - Click "Add API Key"
   - Set permissions and expiration
   - Copy the generated key

2. **Use API Key**:
   ```bash
   curl -H "Authorization: Bearer YOUR_API_KEY" \
        https://pipecd.example.com/api/v1/applications
   ```

### Using Service Account

For programmatic access:

```yaml
# Service account configuration
apiVersion: pipecd.dev/v1beta1
kind: ServiceAccount
metadata:
  name: ci-pipeline
spec:
  role: DEPLOYER
  projectId: my-project
```

Generate token:
```bash
pipectl service-account create \
  --name ci-pipeline \
  --role DEPLOYER \
  --project my-project
```

## Using the API

### With curl

```bash
# Get deployment status
curl -X GET \
  -H "Authorization: Bearer YOUR_API_KEY" \
  https://pipecd.example.com/api/v1/deployments/deploy-123

# Trigger sync
curl -X POST \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"commandId": "cmd-123"}' \
  https://pipecd.example.com/api/v1/applications/app-123/sync
```

### With grpcurl

```bash
# List gRPC services
grpcurl -plaintext localhost:9080 list

# List methods for a service
grpcurl -plaintext localhost:9080 list pipecd.service.webservice.WebService

# Call a method
grpcurl -plaintext \
  -d '{"deployment_id": "deploy-123"}' \
  localhost:9080 \
  pipecd.service.webservice.WebService/GetDeployment
```

### With Go

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    
    webservice "github.com/pipe-cd/pipecd/pkg/app/api/webservice"
)

func main() {
    creds := credentials.NewTLS(&tls.Config{})
    conn, err := grpc.Dial(
        "pipecd.example.com:443",
        grpc.WithTransportCredentials(creds),
        grpc.WithPerRPCCredentials(newAPIKeyAuth("YOUR_API_KEY")),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    client := webservice.NewWebServiceClient(conn)
    
    // List applications
    resp, err := client.ListApplications(context.Background(), &webservice.ListApplicationsRequest{
        ProjectId: "my-project",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for _, app := range resp.Applications {
        log.Printf("App: %s (%s)\n", app.Name, app.Id)
    }
}

// API Key authentication
type apiKeyAuth struct {
    apiKey string
}

func newAPIKeyAuth(key string) credentials.PerRPCCredentials {
    return &apiKeyAuth{apiKey: key}
}

func (a *apiKeyAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
    return map[string]string{
        "authorization": "Bearer " + a.apiKey,
    }, nil
}

func (a *apiKeyAuth) RequireTransportSecurity() bool {
    return true
}
```

### With Python

```python
import grpc
from pipecd.api import webservice_pb2, webservice_pb2_grpc

# Create channel
channel = grpc.secure_channel(
    'pipecd.example.com:443',
    grpc.ssl_channel_credentials()
)

# Create client
client = webservice_pb2_grpc.WebServiceStub(channel)

# Add metadata (API key)
metadata = [('authorization', 'Bearer YOUR_API_KEY')]

# Make request
request = webservice_pb2.ListApplicationsRequest(
    project_id='my-project'
)

response = client.ListApplications(request, metadata=metadata)

for app in response.applications:
    print(f"App: {app.name} ({app.id})")
```

## API Client Libraries

### Official Clients

- **Go**: Built-in (generated from proto files)
- **TypeScript**: Used by web UI

### Community Clients

> Note: Community clients are maintained by third parties and may not be up-to-date.

- Python: [pipecd-python-client](https://github.com/example/pipecd-python-client)
- Ruby: [pipecd-ruby](https://github.com/example/pipecd-ruby)

## Contributing to API Docs

### Adding New API Methods

1. **Define in proto file**:
   ```protobuf
   service WebService {
     rpc GetNewResource(GetNewResourceRequest) returns (GetNewResourceResponse) {}
   }
   
   message GetNewResourceRequest {
     string resource_id = 1;
   }
   
   message GetNewResourceResponse {
     Resource resource = 1;
   }
   ```

2. **Regenerate code**:
   ```bash
   make gen/api
   ```

3. **Implement server**:
   ```go
   func (s *webService) GetNewResource(ctx context.Context, req *webservice.GetNewResourceRequest) (*webservice.GetNewResourceResponse, error) {
       // Implementation
   }
   ```

4. **Document in this file**:
   - Add method description
   - Provide request/response examples
   - Include usage examples

### API Documentation Standards

- **Clear descriptions**: Explain what the API does
- **Request/response examples**: Show actual usage
- **Error codes**: Document possible errors
- **Deprecation notices**: Mark deprecated APIs
- **Version information**: Note when API was added/changed

### Testing API Changes

```bash
# Run API tests
make test/api

# Test with grpcurl
grpcurl -plaintext localhost:9080 list

# Test with curl (for REST endpoints)
curl -v http://localhost:8080/api/v1/health
```

## API Rate Limiting

PipeCD implements rate limiting to prevent abuse:

- **Default limits**: 100 requests/minute per API key
- **Burst limit**: 200 requests
- **Headers**:
  - `X-RateLimit-Limit`: Max requests per window
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Time when limit resets

## Error Handling

### gRPC Status Codes

| Code | Description | Example |
|------|-------------|---------|
| OK | Success | Request completed successfully |
| INVALID_ARGUMENT | Invalid request | Missing required field |
| NOT_FOUND | Resource not found | Deployment ID doesn't exist |
| PERMISSION_DENIED | Insufficient permissions | API key lacks permission |
| UNAUTHENTICATED | Authentication failed | Invalid API key |
| INTERNAL | Server error | Database connection failed |

### Error Response Example

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Deployment not found: deploy-123",
    "details": []
  }
}
```

## Resources

- [Protocol Buffers](https://protobuf.dev/)
- [gRPC Documentation](https://grpc.io/docs/)
- [PipeCD API Reference](https://pipecd.dev/docs/api-reference/)

## Support

- [Slack Channel](https://cloud-native.slack.com/archives/C01B27F9T0X)
- [GitHub Discussions](https://github.com/pipe-cd/pipecd/discussions)
- [API Issues](https://github.com/pipe-cd/pipecd/issues?q=is%3Aissue+label%3Aapi)
