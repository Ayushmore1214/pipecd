# Code Style Guide for PipeCD

This document outlines the coding standards and style conventions for contributing to PipeCD.

## Table of Contents

- [General Principles](#general-principles)
- [Go Style Guide](#go-style-guide)
- [TypeScript/React Style Guide](#typescriptreact-style-guide)
- [YAML Style Guide](#yaml-style-guide)
- [Shell Script Style Guide](#shell-script-style-guide)
- [Commit Message Style](#commit-message-style)
- [Documentation Style](#documentation-style)

## General Principles

1. **Readability First**: Code is read more often than written
2. **Consistency**: Follow existing patterns in the codebase
3. **Simplicity**: Prefer simple, straightforward solutions
4. **Self-Documenting**: Write code that explains itself
5. **DRY**: Don't Repeat Yourself - extract common logic
6. **SOLID**: Follow SOLID principles for object-oriented code

## Go Style Guide

### Formatting

**Use gofmt and goimports**:
```bash
# Format all Go files
gofmt -s -w .

# Organize imports
goimports -w .
```

### Naming Conventions

**Packages**:
```go
// Good: lowercase, singular, concise
package deployer
package kubernetes

// Bad: mixed case, plural, verbose
package deployerPackage
package kubernetesResources
```

**Variables**:
```go
// Good: camelCase, descriptive
var deploymentID string
var maxRetryCount int

// Bad: snake_case, abbreviations
var deployment_id string
var max_retry_cnt int
```

**Constants**:
```go
// Good: camelCase (not SCREAMING_SNAKE_CASE in Go)
const defaultTimeout = 30 * time.Second
const maxRetries = 3

// Exception: Exported constants can be PascalCase
const DefaultNamespace = "default"
```

**Functions**:
```go
// Good: PascalCase for exported, camelCase for private
func DeployApplication(ctx context.Context, app *Application) error {}
func validateConfig(cfg *Config) error {}

// Bad: unclear or abbreviated
func Deploy(a *App) error {}  // Too generic
func valCfg(c *Config) error {}  // Too abbreviated
```

**Interfaces**:
```go
// Good: noun + -er suffix
type Deployer interface {}
type ConfigValidator interface {}

// Bad: prefix or unclear
type IDeployer interface {}  // Don't use "I" prefix
type DeployInterface interface {}  // Don't use "Interface" suffix
```

### Code Organization

**File Structure**:
```go
// 1. Package declaration
package deployer

// 2. Imports (grouped: stdlib, external, internal)
import (
    "context"
    "fmt"
    
    "github.com/pkg/errors"
    "go.uber.org/zap"
    
    "github.com/pipe-cd/pipecd/pkg/model"
)

// 3. Constants
const (
    defaultRetries = 3
)

// 4. Types
type Deployer struct {
    logger *zap.Logger
}

// 5. Constructor
func NewDeployer(logger *zap.Logger) *Deployer {
    return &Deployer{logger: logger}
}

// 6. Public methods
func (d *Deployer) Deploy(ctx context.Context, app *Application) error {
    // implementation
}

// 7. Private methods
func (d *Deployer) validateApp(app *Application) error {
    // implementation
}
```

### Error Handling

**Wrap errors with context**:
```go
// Good: provides context
if err := validateConfig(cfg); err != nil {
    return fmt.Errorf("failed to validate config: %w", err)
}

// Bad: loses context
if err := validateConfig(cfg); err != nil {
    return err
}
```

**Check errors explicitly**:
```go
// Good: handle all errors
result, err := doSomething()
if err != nil {
    return err
}

// Bad: ignore errors
result, _ := doSomething()
```

### Comments

**Package comments**:
```go
// Package deployer provides deployment orchestration functionality
// for managing application deployments across multiple platforms.
package deployer
```

**Function comments**:
```go
// DeployApplication initiates a deployment for the specified application.
// It validates the configuration, creates a deployment plan, and executes
// the deployment pipeline stages.
//
// Returns an error if validation fails or deployment cannot be started.
func DeployApplication(ctx context.Context, app *Application) error {
    // implementation
}
```

**Inline comments**:
```go
// Good: explain WHY, not WHAT
// Wait for rollout to complete to avoid race conditions
time.Sleep(1 * time.Second)

// Bad: explain WHAT (obvious from code)
// Sleep for 1 second
time.Sleep(1 * time.Second)
```

### Context Usage

```go
// Good: context as first parameter
func ProcessDeployment(ctx context.Context, id string) error {
    // Use context for cancellation and timeouts
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-processChan:
        return handleResult(result)
    }
}

// Bad: no context support
func ProcessDeployment(id string) error {
    // Cannot be cancelled
}
```

### Testing

**Table-driven tests**:
```go
func TestValidateConfig(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name        string
        config      *Config
        expectError bool
        errContains string
    }{
        {
            name: "valid config",
            config: &Config{
                Name: "test",
                URL:  "https://example.com",
            },
            expectError: false,
        },
        {
            name: "missing name",
            config: &Config{
                URL: "https://example.com",
            },
            expectError: true,
            errContains: "name is required",
        },
    }
    
    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            err := validateConfig(tt.config)
            
            if tt.expectError {
                require.Error(t, err)
                if tt.errContains != "" {
                    assert.Contains(t, err.Error(), tt.errContains)
                }
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

## TypeScript/React Style Guide

### Formatting

Use Prettier with project configuration:
```bash
yarn --cwd web format
```

### Naming Conventions

**Files**:
```
// Components: PascalCase
DeploymentList.tsx
ApplicationDetail.tsx

// Hooks: camelCase with 'use' prefix
useDeployments.ts
useApplications.ts

// Utilities: camelCase
formatDate.ts
apiClient.ts

// Types: PascalCase
types.ts
models.ts
```

**Variables and Functions**:
```typescript
// Good: camelCase
const deploymentId = "deploy-123";
const fetchDeployments = async () => {};

// Bad: PascalCase or snake_case
const DeploymentId = "deploy-123";
const fetch_deployments = async () => {};
```

**Types and Interfaces**:
```typescript
// Good: PascalCase, descriptive
interface DeploymentProps {
    id: string;
    status: DeploymentStatus;
}

type ApplicationKind = "KUBERNETES" | "TERRAFORM" | "CLOUDRUN";

// Bad: prefix or unclear
interface IDeploymentProps {}  // No "I" prefix
type AppKind = string;  // Not specific enough
```

### React Components

**Functional components with TypeScript**:
```typescript
import { FC } from 'react';

interface DeploymentListProps {
    projectId: string;
    onSelect?: (id: string) => void;
}

export const DeploymentList: FC<DeploymentListProps> = ({ 
    projectId, 
    onSelect 
}) => {
    const [deployments, setDeployments] = useState<Deployment[]>([]);
    
    useEffect(() => {
        fetchDeployments(projectId).then(setDeployments);
    }, [projectId]);
    
    return (
        <div>
            {deployments.map(deployment => (
                <DeploymentCard 
                    key={deployment.id}
                    deployment={deployment}
                    onClick={() => onSelect?.(deployment.id)}
                />
            ))}
        </div>
    );
};
```

### Hooks

**Custom hooks**:
```typescript
// Good: starts with 'use', typed return
function useDeployments(projectId: string) {
    const [deployments, setDeployments] = useState<Deployment[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);
    
    useEffect(() => {
        setLoading(true);
        fetchDeployments(projectId)
            .then(setDeployments)
            .catch(setError)
            .finally(() => setLoading(false));
    }, [projectId]);
    
    return { deployments, loading, error };
}
```

### TypeScript Best Practices

**Avoid `any`**:
```typescript
// Good: specific types
function formatDeployment(deployment: Deployment): string {
    return `${deployment.name} - ${deployment.status}`;
}

// Bad: using any
function formatDeployment(deployment: any): string {
    return `${deployment.name} - ${deployment.status}`;
}
```

**Use union types**:
```typescript
// Good: explicit union types
type DeploymentStatus = 
    | "PENDING"
    | "RUNNING"
    | "SUCCESS"
    | "FAILURE"
    | "CANCELLED";

// Bad: string without constraints
type DeploymentStatus = string;
```

## YAML Style Guide

### Indentation

```yaml
# Good: 2 spaces
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: my-app
  pipeline:
    stages:
      - name: K8S_SYNC

# Bad: 4 spaces or tabs
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
    name: my-app
```

### Quotes

```yaml
# Good: use quotes for strings with special characters
name: "my-app"
url: "https://example.com"
description: "This app does X, Y, and Z"

# Unquoted for simple strings
environment: production
replicas: 3
```

### Comments

```yaml
# Good: explain non-obvious configurations
spec:
  pipeline:
    stages:
      - name: K8S_SYNC
        with:
          # Wait for rollout to prevent premature success reporting
          waitForRollout: true
```

## Shell Script Style Guide

### Shebang and Options

```bash
#!/bin/bash
set -euo pipefail  # Exit on error, undefined vars, pipe failures

# Good: strict error handling
do_something || {
    echo "Error: something failed"
    exit 1
}
```

### Variables

```bash
# Good: uppercase for constants, lowercase for local
readonly MAX_RETRIES=3
local deployment_id="deploy-123"

# Use braces for clarity
echo "Deployment: ${deployment_id}"
```

### Functions

```bash
# Good: descriptive names, local variables
deploy_application() {
    local app_name="$1"
    local environment="$2"
    
    echo "Deploying ${app_name} to ${environment}"
    # implementation
}
```

## Commit Message Style

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

### Examples

```bash
# Feature
feat(piped): add support for Azure deployments

# Bug fix
fix(web): resolve deployment list filtering issue

# Documentation
docs: update contributing guide with new process

# Breaking change
feat(api)!: change deployment status enum values

BREAKING CHANGE: DeploymentStatus enum values are now uppercase
```

## Documentation Style

### Markdown

**Headings**:
```markdown
# Title (H1) - Only one per document

## Section (H2)

### Subsection (H3)
```

**Lists**:
```markdown
- Use hyphens for unordered lists
- Keep consistent punctuation
- End with period if full sentences

1. Use numbers for ordered lists
2. Start each item with capital letter
3. Keep format consistent
```

**Code blocks**:
````markdown
```go
// Always specify language
func example() {}
```
````

**Links**:
```markdown
[Descriptive text](https://example.com)
[Reference-style link][ref]

[ref]: https://example.com
```

## Tools and Automation

### Pre-commit Hooks

Install and use pre-commit hooks:
```bash
pip install pre-commit
pre-commit install
pre-commit run --all-files
```

### Linters

Run before committing:
```bash
# Go
make lint/go

# TypeScript
make lint/web

# YAML
yamllint .

# All checks
make check
```

## Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google TypeScript Style Guide](https://google.github.io/styleguide/tsguide.html)
- [Conventional Commits](https://www.conventionalcommits.org/)

## Questions?

- [#pipecd on CNCF Slack](https://cloud-native.slack.com/archives/C01B27F9T0X)
- [GitHub Discussions](https://github.com/pipe-cd/pipecd/discussions)
