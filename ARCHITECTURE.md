# PipeCD Architecture

This document provides a comprehensive overview of PipeCD's architecture, components, and design principles.

## Table of Contents

- [Overview](#overview)
- [Core Components](#core-components)
- [Architecture Diagram](#architecture-diagram)
- [Component Interaction](#component-interaction)
- [Data Flow](#data-flow)
- [Security Model](#security-model)
- [Scalability & Performance](#scalability--performance)
- [Technology Stack](#technology-stack)

## Overview

PipeCD is a GitOps-based continuous delivery platform designed for managing deployments across multiple application types and cloud platforms. The architecture follows a control plane/agent model with strong security principles.

### Design Principles

1. **GitOps-First**: All deployment configurations and application manifests are stored in Git
2. **Multi-Cloud Native**: Support for Kubernetes, Terraform, Cloud Run, Lambda, ECS, and more
3. **Security by Design**: No credentials required outside the application cluster
4. **Scalability**: Designed to handle thousands of applications across multiple environments
5. **Observability**: Built-in metrics, logging, and deployment insights

## Core Components

### 1. Control Plane (PipeCD Server)

The control plane is the centralized management component that:

- **Responsibilities**:
  - Stores deployment metadata and history
  - Provides gRPC API for piped agents
  - Manages authentication and authorization
  - Serves the web UI
  - Aggregates metrics and insights

- **Technology**: Go, gRPC, Protocol Buffers
- **Storage**: Supports multiple datastores (MySQL, PostgreSQL, Firestore, etc.)
- **Deployment**: Can run as a Kubernetes deployment or standalone binary

### 2. Piped (Agent)

Piped is the agent component that runs in your infrastructure:

- **Responsibilities**:
  - Watches Git repositories for changes
  - Executes deployments according to pipeline definitions
  - Reports deployment status to control plane
  - Performs health checks and analysis
  - Manages platform-specific operations (kubectl, terraform, etc.)

- **Technology**: Go, plugin architecture
- **Deployment**: Runs as a pod in Kubernetes or as a standalone process
- **Security**: Holds all deployment credentials locally

### 3. Launcher

A helper component for remote upgrades:

- **Responsibilities**:
  - Manages piped agent lifecycle
  - Enables remote upgrade capabilities
  - Ensures piped availability

### 4. Pipectl

Command-line tool for PipeCD:

- **Responsibilities**:
  - Register/manage pipeds
  - Encrypt secrets
  - Trigger deployments
  - Query deployment status

### 5. Web UI

Modern web interface for PipeCD:

- **Responsibilities**:
  - Visualize deployments and pipelines
  - Manage applications and environments
  - View insights and metrics
  - Configure settings

- **Technology**: TypeScript, React, Material-UI

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                          Users/Developers                        │
└───────────────┬─────────────────────────────┬───────────────────┘
                │                             │
                │ Web UI                      │ pipectl CLI
                │                             │
                ▼                             ▼
┌───────────────────────────────────────────────────────────────────┐
│                                                                   │
│                        Control Plane (PipeCD)                     │
│                                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐    │
│  │   Web Server │  │  gRPC Server │  │  Auth Service      │    │
│  └──────────────┘  └──────────────┘  └────────────────────┘    │
│                                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐    │
│  │  API Service │  │  Datastore   │  │  Insight Aggregator│    │
│  └──────────────┘  └──────────────┘  └────────────────────┘    │
│                                                                   │
└─────────────────────────────┬─────────────────────────────────────┘
                              │
                              │ gRPC Communication
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│               │     │               │     │               │
│  Piped Agent  │     │  Piped Agent  │     │  Piped Agent  │
│  (Cluster A)  │     │  (Cluster B)  │     │  (Cloud Env)  │
│               │     │               │     │               │
│  ┌─────────┐  │     │  ┌─────────┐  │     │  ┌─────────┐  │
│  │Git Sync │  │     │  │Git Sync │  │     │  │Git Sync │  │
│  └─────────┘  │     │  └─────────┘  │     │  └─────────┘  │
│  ┌─────────┐  │     │  ┌─────────┐  │     │  ┌─────────┐  │
│  │Deployer │  │     │  │Deployer │  │     │  │Deployer │  │
│  └─────────┘  │     │  └─────────┘  │     │  └─────────┘  │
│  ┌─────────┐  │     │  ┌─────────┐  │     │  ┌─────────┐  │
│  │Analyzer │  │     │  │Analyzer │  │     │  │Analyzer │  │
│  └─────────┘  │     │  └─────────┘  │     │  └─────────┘  │
│               │     │               │     │               │
└───────┬───────┘     └───────┬───────┘     └───────┬───────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│  Kubernetes   │     │  Kubernetes   │     │  Cloud Run /  │
│   Cluster     │     │   Cluster     │     │  Lambda / ECS │
└───────────────┘     └───────────────┘     └───────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────────────────────────────────────────────────┐
│                     Git Repositories                       │
│            (Application Manifests & Pipeline Configs)      │
└───────────────────────────────────────────────────────────┘
```

## Component Interaction

### Deployment Flow

1. **Change Detection**:
   - Piped polls Git repository at regular intervals (configurable)
   - Detects changes in application manifests or pipeline definitions

2. **Deployment Creation**:
   - Piped creates deployment request
   - Sends deployment plan to Control Plane
   - Control Plane stores deployment metadata

3. **Deployment Execution**:
   - Piped executes deployment pipeline stages
   - Reports progress to Control Plane
   - Performs automated analysis (if configured)

4. **Status Updates**:
   - Real-time updates sent to Control Plane
   - Web UI reflects current deployment state
   - Notifications sent (if configured)

### Communication Patterns

- **Control Plane ↔ Piped**: Bidirectional gRPC streams
- **Control Plane ↔ Web UI**: HTTP/WebSocket
- **Control Plane ↔ Datastore**: Database protocol
- **Piped ↔ Git**: Git protocol (SSH/HTTPS)
- **Piped ↔ Platform**: Platform-specific APIs (kubectl, terraform, etc.)

## Data Flow

### Deployment Data

```
Git Repo → Piped (sync) → Piped (plan) → Control Plane (store)
                ↓
         Platform (apply)
                ↓
         Control Plane (update status)
                ↓
           Web UI (display)
```

### Metrics & Insights

```
Piped (collect) → Control Plane (aggregate) → Web UI (visualize)
```

## Security Model

### Key Security Features

1. **Credential Isolation**:
   - Deployment credentials never leave the piped environment
   - Control Plane has no access to cluster credentials

2. **Authentication**:
   - Static admin accounts (for quickstart)
   - SSO integration (GitHub, Google, etc.)
   - RBAC for fine-grained access control

3. **Encrypted Communication**:
   - TLS for all gRPC connections
   - mTLS support for enhanced security

4. **Secret Management**:
   - Secrets encrypted at rest
   - Support for external secret managers (coming soon)

5. **Audit Logging**:
   - All API calls logged
   - Deployment history maintained

## Scalability & Performance

### Horizontal Scaling

- **Control Plane**: Can run multiple replicas behind a load balancer
- **Piped Agents**: One per environment/cluster (scales with infrastructure)

### Performance Optimizations

- **Caching**: Redis cache for frequently accessed data
- **Database Indexing**: Optimized queries for deployment history
- **Streaming**: gRPC streaming for real-time updates
- **Connection Pooling**: Efficient resource utilization

### Resource Requirements

**Control Plane**:
- CPU: 2-4 cores (production)
- Memory: 4-8 GB
- Storage: Depends on deployment history retention

**Piped Agent**:
- CPU: 1-2 cores
- Memory: 1-2 GB
- Storage: Minimal (logs and temporary files)

## Technology Stack

### Backend

- **Language**: Go 1.21+
- **Framework**: gRPC, Protocol Buffers
- **Database**: MySQL, PostgreSQL, Firestore, DynamoDB
- **Cache**: Redis (optional)

### Frontend

- **Language**: TypeScript
- **Framework**: React 18
- **UI Library**: Material-UI
- **Build Tool**: Webpack, Yarn

### DevOps

- **Container**: Docker
- **Orchestration**: Kubernetes
- **Packaging**: Helm
- **CI/CD**: GitHub Actions

### Observability

- **Metrics**: Prometheus-compatible endpoints
- **Logging**: Structured logging (JSON)
- **Tracing**: OpenTelemetry support

## Plugin Architecture

PipeCD supports a plugin system for platform providers:

```
┌─────────────────────────────────────┐
│           Piped Core                │
├─────────────────────────────────────┤
│      Plugin Interface (gRPC)        │
├─────────────────────────────────────┤
│  ┌────────┐  ┌────────┐  ┌────────┐│
│  │K8s     │  │Terraform│ │Lambda  ││
│  │Plugin  │  │Plugin   │ │Plugin  ││
│  └────────┘  └────────┘  └────────┘│
└─────────────────────────────────────┘
```

### Supported Platforms

- Kubernetes
- Terraform
- Cloud Run
- AWS Lambda
- AWS ECS
- Azure (coming soon)

## References

- [PipeCD Documentation](https://pipecd.dev/docs)
- [API Reference](https://pipecd.dev/docs/api-reference)
- [Design Proposals (RFCs)](./docs/rfcs)

## Contributing

To contribute to PipeCD architecture:

1. Review existing [RFCs](./docs/rfcs)
2. Propose new features via GitHub Discussions
3. Submit design proposals for major changes
4. Follow the [Contributing Guide](./CONTRIBUTING.md)

---

For questions about the architecture, please:
- Join our [Slack channel](https://cloud-native.slack.com/archives/C01B27F9T0X)
- Attend our [community meetings](https://bit.ly/pipecd-mtg-notes)
- Open a [GitHub Discussion](https://github.com/pipe-cd/pipecd/discussions)
