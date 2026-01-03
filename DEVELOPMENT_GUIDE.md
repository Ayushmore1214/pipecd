# PipeCD Development Guide

A comprehensive guide for developers contributing to PipeCD.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Debugging Tips](#debugging-tips)
- [Common Tasks](#common-tasks)
- [Performance Optimization](#performance-optimization)
- [Best Practices](#best-practices)

## Getting Started

### Prerequisites

Ensure you have the following tools installed:

```bash
# Required
- Go 1.21+ (check go.mod for exact version)
- Node.js 18.12+ (for web development)
- Yarn (for web dependencies)
- Docker (for container builds)
- kubectl (for Kubernetes operations)
- Make (for build automation)

# Recommended
- Git 2.30+
- A good IDE (VS Code, GoLand, or similar)
- pre-commit (for git hooks)
- grpcurl (for testing gRPC endpoints)
```

### Quick Setup

```bash
# Clone the repository
git clone https://github.com/pipe-cd/pipecd.git
cd pipecd

# Install pre-commit hooks (recommended)
pip install pre-commit
pre-commit install

# Update dependencies
make update/go-deps
make update/web-deps

# Verify everything works
make check
```

## Development Environment Setup

### VS Code Setup

Recommended extensions:

```json
{
  "recommendations": [
    "golang.go",
    "dbaeumer.vscode-eslint",
    "esbenp.prettier-vscode",
    "ms-kubernetes-tools.vscode-kubernetes-tools",
    "redhat.vscode-yaml",
    "streetsidesoftware.code-spell-checker"
  ]
}
```

Save this as `.vscode/extensions.json`.

Workspace settings (`.vscode/settings.json`):

```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--config=.golangci.yml"],
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.insertSpaces": false
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[yaml]": {
    "editor.defaultFormatter": "redhat.vscode-yaml"
  }
}
```

### GoLand/IntelliJ Setup

1. **Import Project**: Open the project directory
2. **Go SDK**: Set Go SDK to version specified in go.mod
3. **Code Style**: 
   - Go: Use tabs, tab size 4
   - TypeScript: Use spaces, indent size 2
4. **Enable golangci-lint**: 
   - Settings → Tools → File Watchers → Add golangci-lint
   - Arguments: `run --config .golangci.yml`

## Project Structure

```
pipecd/
├── cmd/                    # Main applications
│   ├── pipecd/            # Control plane server
│   ├── piped/             # Agent (piped)
│   ├── pipectl/           # CLI tool
│   └── launcher/          # Launcher for remote upgrades
├── pkg/                   # Shared libraries
│   ├── app/               # Application logic
│   ├── config/            # Configuration handling
│   ├── datastore/         # Database interfaces
│   ├── model/             # Data models
│   └── version/           # Version information
├── web/                   # Frontend React application
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── modules/       # Feature modules
│   │   └── api/           # API clients
│   └── package.json
├── docs/                  # Documentation site
├── manifests/             # Kubernetes manifests & Helm charts
├── examples/              # Example configurations
├── test/                  # Integration tests
├── hack/                  # Scripts for development
└── tool/                  # Development tools
```

### Key Directories

- **cmd/**: Each subdirectory contains a `main.go` for a binary
- **pkg/**: Reusable packages, organized by functionality
- **web/**: TypeScript/React frontend
- **manifests/**: Helm charts for deployment
- **docs/**: Hugo-based documentation site

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

### 2. Make Changes

Follow the [coding standards](#coding-standards) and [testing guidelines](#testing-guidelines).

### 3. Run Checks Locally

```bash
# Format code
make fmt/go
make fmt/web

# Run linters
make lint/go
make lint/web

# Run tests
make test/go
make test/web

# Or run all checks at once
make check
```

### 4. Commit Changes

```bash
# Sign your commits (required)
git commit -s -m "Add feature X"

# Follow commit message conventions
# Format: <type>(<scope>): <subject>
# Examples:
# feat(piped): add support for Azure deployments
# fix(web): resolve deployment list filtering issue
# docs: update contributing guide
```

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
# Then create a PR on GitHub
```

## Coding Standards

### Go

**Follow these conventions:**

1. **Use gofmt and goimports**:
   ```bash
   gofmt -s -w .
   goimports -w .
   ```

2. **Error Handling**:
   ```go
   // Good: Handle errors explicitly
   if err != nil {
       return fmt.Errorf("failed to do X: %w", err)
   }
   
   // Bad: Ignoring errors
   _ = doSomething()
   ```

3. **Context Propagation**:
   ```go
   // Always accept context as first parameter
   func DoSomething(ctx context.Context, param string) error {
       // ...
   }
   ```

4. **Naming Conventions**:
   ```go
   // Interfaces: end with -er
   type Deployer interface { ... }
   
   // Constructors: NewXxx
   func NewDeployer() *Deployer { ... }
   
   // Acronyms: keep consistent case
   var apiURL string  // Good
   var apiUrl string  // Bad
   ```

5. **Comments**:
   ```go
   // Package comment: describe package purpose
   package deployer
   
   // Public function: describe what it does
   // DeployApplication deploys the specified application to the target environment.
   func DeployApplication(ctx context.Context, app *Application) error {
       // ...
   }
   ```

6. **Testing**:
   ```go
   // Test files: _test.go suffix
   // Test functions: TestXxx
   func TestDeployApplication(t *testing.T) {
       // Use table-driven tests
       tests := []struct {
           name string
           app  *Application
           want error
       }{
           // test cases
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // test logic
           })
       }
   }
   ```

### TypeScript/React

1. **Use TypeScript Strictly**:
   ```typescript
   // Good: explicit types
   interface DeploymentProps {
     id: string;
     status: DeploymentStatus;
   }
   
   // Bad: using any
   const data: any = ...
   ```

2. **Component Structure**:
   ```typescript
   // Functional components with hooks
   export const DeploymentList: FC<DeploymentListProps> = ({ projectId }) => {
     const [deployments, setDeployments] = useState<Deployment[]>([]);
     
     useEffect(() => {
       // fetch deployments
     }, [projectId]);
     
     return <div>...</div>;
   };
   ```

3. **File Organization**:
   ```
   components/
     DeploymentList/
       index.tsx          # Main component
       index.test.tsx     # Tests
       styles.ts          # Styles (if needed)
   ```

4. **Styling**:
   ```typescript
   // Use Material-UI's styling solution
   import { makeStyles } from '@material-ui/core/styles';
   
   const useStyles = makeStyles((theme) => ({
     root: {
       padding: theme.spacing(2),
     },
   }));
   ```

### YAML

1. **Indentation**: Use 2 spaces
2. **Quotes**: Use double quotes for strings with special characters
3. **Comments**: Explain non-obvious configurations

```yaml
# Good
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  name: my-app
  labels:
    team: platform
  pipeline:
    stages:
      - name: K8S_SYNC
        with:
          # Wait for rollout to complete before marking as success
          waitForRollout: true
```

## Testing Guidelines

### Unit Tests

**Go:**

```bash
# Run all tests
make test/go

# Run specific package tests
go test ./pkg/app/piped/...

# Run with coverage
make test/go COVERAGE=true

# Run specific test
go test -run TestDeployApplication ./pkg/app/piped/deployer
```

**TypeScript:**

```bash
# Run all tests
make test/web

# Run specific test file
yarn --cwd web test DeploymentList.test.tsx

# Run in watch mode
yarn --cwd web test --watch
```

### Integration Tests

```bash
# Run integration tests
make test/integration

# Run specific integration test
go test -tags=integration ./test/integration/...
```

### Writing Good Tests

1. **Table-Driven Tests** (Go):
   ```go
   func TestCalculate(t *testing.T) {
       tests := []struct {
           name     string
           input    int
           expected int
       }{
           {"zero", 0, 0},
           {"positive", 5, 10},
           {"negative", -5, -10},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               got := Calculate(tt.input)
               if got != tt.expected {
                   t.Errorf("got %d, want %d", got, tt.expected)
               }
           })
       }
   }
   ```

2. **Use Mocks Appropriately**:
   ```go
   // Generate mocks with mockgen
   //go:generate mockgen -source=interface.go -destination=mock.go -package=mocks
   ```

3. **Test Coverage**:
   - Aim for >80% coverage on new code
   - Focus on critical paths
   - Don't test trivial getters/setters

## Debugging Tips

### Debugging Go Code

1. **Use Delve**:
   ```bash
   # Install delve
   go install github.com/go-delve/delve/cmd/dlv@latest
   
   # Debug a test
   dlv test ./pkg/app/piped -- -test.run TestDeployApplication
   
   # Debug a binary
   dlv exec .artifacts/piped -- --config-file=/path/to/config.yaml
   ```

2. **Print Debugging**:
   ```go
   // Use structured logging
   import "go.uber.org/zap"
   
   logger.Debug("deployment started", zap.String("id", deploymentID))
   ```

3. **Remote Debugging**:
   ```bash
   # Forward debug port
   kubectl port-forward -n pipecd pod/piped-xxx 2345:2345
   
   # Connect with IDE debugger
   ```

### Debugging TypeScript/React

1. **Browser DevTools**:
   - Use React DevTools extension
   - Check Console for errors
   - Use Network tab for API calls
   - Use Profiler for performance

2. **VS Code Debugging**:
   ```json
   {
     "type": "chrome",
     "request": "launch",
     "name": "Debug Web",
     "url": "http://localhost:9090",
     "webRoot": "${workspaceFolder}/web"
   }
   ```

## Common Tasks

### Adding a New gRPC API

1. **Define in proto file**:
   ```protobuf
   service PipedService {
     rpc GetDeployment(GetDeploymentRequest) returns (GetDeploymentResponse) {}
   }
   ```

2. **Generate code**:
   ```bash
   make gen/api
   ```

3. **Implement server**:
   ```go
   func (s *PipedService) GetDeployment(ctx context.Context, req *pipedservice.GetDeploymentRequest) (*pipedservice.GetDeploymentResponse, error) {
       // implementation
   }
   ```

4. **Add client call** (web):
   ```typescript
   const deployment = await pipedClient.getDeployment({ id: deploymentId });
   ```

### Adding a New Platform Provider

1. Create plugin directory:
   ```
   pkg/app/pipedv1/plugin/yourplatform/
   ```

2. Implement provider interface
3. Register provider
4. Add tests
5. Update documentation

### Updating Dependencies

```bash
# Go dependencies
make update/go-deps

# Web dependencies
make update/web-deps

# Check for outdated dependencies
go list -u -m all
yarn --cwd web outdated
```

## Performance Optimization

### Profiling

**CPU Profiling**:

```bash
# Enable profiling in config
# Access at http://localhost:6060/debug/pprof/

# Capture CPU profile
curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof
```

**Memory Profiling**:

```bash
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

### Database Optimization

1. **Add indexes for frequent queries**
2. **Use connection pooling**
3. **Batch operations when possible**
4. **Monitor slow queries**

### Frontend Optimization

1. **Code splitting**: Use React.lazy()
2. **Memoization**: Use React.memo, useMemo, useCallback
3. **Virtualization**: For long lists
4. **Bundle analysis**: `yarn --cwd web analyze`

## Best Practices

### Security

1. **Never commit secrets**
2. **Use context for cancellation**
3. **Validate all inputs**
4. **Sanitize user-provided data**
5. **Use prepared statements for SQL**

### Code Organization

1. **Keep functions small** (< 50 lines)
2. **Single responsibility** per function/component
3. **DRY** (Don't Repeat Yourself)
4. **SOLID principles**

### Git Workflow

1. **Commit early, commit often**
2. **Write meaningful commit messages**
3. **Keep commits atomic**
4. **Rebase before merging** (if needed)
5. **Sign all commits** (`git commit -s`)

### Documentation

1. **Update docs with code changes**
2. **Add comments for complex logic**
3. **Keep README files updated**
4. **Document breaking changes**

## Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
- [React Best Practices](https://react.dev/learn)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Material-UI Documentation](https://material-ui.com/)

## Getting Help

- **Slack**: [#pipecd on CNCF Slack](https://cloud-native.slack.com/archives/C01B27F9T0X)
- **Meetings**: [Community meetings](https://bit.ly/pipecd-mtg-notes)
- **Discussions**: [GitHub Discussions](https://github.com/pipe-cd/pipecd/discussions)

---

Happy coding! 🚀
