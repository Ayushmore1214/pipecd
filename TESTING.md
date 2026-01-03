# Testing Guide for PipeCD

Comprehensive guide for writing and running tests in PipeCD.

## Table of Contents

- [Testing Philosophy](#testing-philosophy)
- [Test Types](#test-types)
- [Running Tests](#running-tests)
- [Writing Unit Tests](#writing-unit-tests)
- [Writing Integration Tests](#writing-integration-tests)
- [Testing Best Practices](#testing-best-practices)
- [Mocking](#mocking)
- [Test Coverage](#test-coverage)
- [CI/CD Testing](#cicd-testing)

## Testing Philosophy

PipeCD follows these testing principles:

1. **Test Pyramid**: More unit tests, fewer integration tests, even fewer E2E tests
2. **Fast Feedback**: Tests should run quickly for rapid development
3. **Reliability**: Tests should be deterministic and not flaky
4. **Maintainability**: Tests should be easy to understand and update
5. **Coverage**: Critical paths should have high test coverage

## Test Types

### Unit Tests

**Purpose**: Test individual functions/methods in isolation

**Characteristics**:
- Fast execution (< 1 second per test)
- No external dependencies
- Use mocks for dependencies
- High code coverage

**Location**: `*_test.go` files alongside source code

### Integration Tests

**Purpose**: Test interaction between components

**Characteristics**:
- Test real integrations (database, APIs)
- Slower than unit tests
- May require setup/teardown
- Test realistic scenarios

**Location**: `test/integration/` directory

### End-to-End Tests

**Purpose**: Test complete user workflows

**Characteristics**:
- Test full system
- Slowest tests
- Most realistic
- Fewer in number

**Location**: `test/e2e/` directory

## Running Tests

### Go Tests

```bash
# Run all unit tests
make test/go

# Run tests with coverage
make test/go COVERAGE=true

# Run specific package
go test ./pkg/app/piped/deployer/...

# Run specific test
go test -run TestDeploymentController ./pkg/app/piped/deployer/

# Run with verbose output
go test -v ./pkg/...

# Run with race detector
go test -race ./pkg/...

# Run integration tests
make test/integration

# Run with timeout
go test -timeout 30s ./pkg/...
```

### Web Tests

```bash
# Run all web tests
make test/web

# Run specific test file
yarn --cwd web test DeploymentList.test.tsx

# Run in watch mode
yarn --cwd web test --watch

# Run with coverage
yarn --cwd web test --coverage

# Update snapshots
yarn --cwd web test -u
```

### Quick Test Commands

```bash
# Test everything
make test

# Run only fast tests
make test/go SKIP_INTEGRATION=true

# Run specific module tests
make test/go MODULES=./pkg/app/piped
```

## Writing Unit Tests

### Go Unit Tests

**Structure**: Follow table-driven test pattern

```go
package deployer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployApplication(t *testing.T) {
	t.Parallel() // Run tests in parallel when possible

	tests := []struct {
		name        string
		app         *Application
		expectError bool
		errContains string
	}{
		{
			name: "successful deployment",
			app: &Application{
				ID:   "app-1",
				Name: "test-app",
			},
			expectError: false,
		},
		{
			name: "missing app ID",
			app: &Application{
				Name: "test-app",
			},
			expectError: true,
			errContains: "app ID is required",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run sub-tests in parallel

			// Arrange
			deployer := NewDeployer()

			// Act
			err := deployer.Deploy(context.Background(), tt.app)

			// Assert
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

**Using Assertions**:

```go
import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	// Use 'require' for critical checks (stops test on failure)
	require.NotNil(t, obj)
	require.NoError(t, err)

	// Use 'assert' for non-critical checks (continues test on failure)
	assert.Equal(t, expected, actual)
	assert.True(t, condition)
	assert.Contains(t, slice, element)
}
```

### TypeScript/React Tests

**Component Testing**:

```typescript
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { DeploymentList } from './DeploymentList';

describe('DeploymentList', () => {
  it('renders deployment list', async () => {
    // Arrange
    const deployments = [
      { id: '1', name: 'app-1', status: 'SUCCESS' },
      { id: '2', name: 'app-2', status: 'RUNNING' },
    ];

    // Act
    render(<DeploymentList deployments={deployments} />);

    // Assert
    expect(screen.getByText('app-1')).toBeInTheDocument();
    expect(screen.getByText('app-2')).toBeInTheDocument();
  });

  it('filters deployments by status', async () => {
    // Arrange
    const deployments = [...];
    render(<DeploymentList deployments={deployments} />);

    // Act
    fireEvent.click(screen.getByRole('button', { name: /filter/i }));
    fireEvent.click(screen.getByText('SUCCESS'));

    // Assert
    await waitFor(() => {
      expect(screen.queryByText('app-2')).not.toBeInTheDocument();
    });
  });
});
```

**Hook Testing**:

```typescript
import { renderHook, act } from '@testing-library/react-hooks';
import { useDeployments } from './useDeployments';

describe('useDeployments', () => {
  it('fetches deployments', async () => {
    const { result, waitForNextUpdate } = renderHook(() =>
      useDeployments('project-1')
    );

    expect(result.current.loading).toBe(true);

    await waitForNextUpdate();

    expect(result.current.loading).toBe(false);
    expect(result.current.deployments).toHaveLength(2);
  });
});
```

## Writing Integration Tests

### Database Integration Tests

```go
// +build integration

package datastore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeploymentStore(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	store := NewDeploymentStore(db)

	t.Run("create and get deployment", func(t *testing.T) {
		ctx := context.Background()

		// Create deployment
		deployment := &Deployment{
			ID:   "deploy-1",
			Name: "test-deployment",
		}
		err := store.Create(ctx, deployment)
		require.NoError(t, err)

		// Retrieve deployment
		retrieved, err := store.Get(ctx, deployment.ID)
		require.NoError(t, err)
		require.Equal(t, deployment.Name, retrieved.Name)
	})
}

func setupTestDB(t *testing.T) *sql.DB {
	// Setup test database (e.g., SQLite in-memory)
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Run migrations
	err = runMigrations(db)
	require.NoError(t, err)

	return db
}
```

### API Integration Tests

```go
// +build integration

func TestPipedAPI(t *testing.T) {
	// Start test server
	server := startTestServer(t)
	defer server.Close()

	client := NewPipedClient(server.URL)

	t.Run("register piped", func(t *testing.T) {
		req := &RegisterPipedRequest{
			Name: "test-piped",
		}

		resp, err := client.RegisterPiped(context.Background(), req)
		require.NoError(t, err)
		require.NotEmpty(t, resp.PipedID)
	})
}
```

## Testing Best Practices

### 1. Test Naming

```go
// Good: Descriptive names
func TestDeploymentController_HandleSuccessfulDeployment(t *testing.T) {}
func TestKubernetesProvider_ApplyManifestWithInvalidYAML(t *testing.T) {}

// Bad: Vague names
func TestDeploy(t *testing.T) {}
func TestHandle(t *testing.T) {}
```

### 2. Arrange-Act-Assert Pattern

```go
func TestExample(t *testing.T) {
	// Arrange - Set up test data and dependencies
	deployer := NewDeployer()
	app := &Application{ID: "app-1"}

	// Act - Execute the code under test
	err := deployer.Deploy(context.Background(), app)

	// Assert - Verify the results
	require.NoError(t, err)
}
```

### 3. Test Independence

```go
// Good: Each test is independent
func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)
	// test logic
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)
	// test logic
}

// Bad: Tests depend on each other
var globalDB *sql.DB

func TestCreate(t *testing.T) {
	globalDB = setupTestDB(t) // Affects other tests
	// test logic
}
```

### 4. Use Test Fixtures

```go
// fixtures.go
func NewTestDeployment(opts ...func(*Deployment)) *Deployment {
	d := &Deployment{
		ID:     "deploy-1",
		Name:   "test-deployment",
		Status: DeploymentStatus_RUNNING,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Usage in tests
func TestDeployment(t *testing.T) {
	deployment := NewTestDeployment(
		func(d *Deployment) {
			d.Status = DeploymentStatus_SUCCESS
		},
	)
	// test logic
}
```

### 5. Test Error Conditions

```go
func TestDeployWithError(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() error
		expectError string
	}{
		{
			name:        "database error",
			setup:       func() error { return errors.New("db error") },
			expectError: "failed to save deployment",
		},
		{
			name:        "network timeout",
			setup:       func() error { return context.DeadlineExceeded },
			expectError: "deployment timeout",
		},
	}
	// test logic
}
```

## Mocking

### Using mockgen

```go
//go:generate mockgen -source=deployer.go -destination=mock_deployer.go -package=mocks

// Interface to mock
type Deployer interface {
	Deploy(ctx context.Context, app *Application) error
}

// Using the mock in tests
import "github.com/golang/mock/gomock"

func TestWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeployer := mocks.NewMockDeployer(ctrl)
	mockDeployer.EXPECT().
		Deploy(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	// Use mock in test
}
```

### Manual Mocks

```go
// Simple mock implementation
type MockDatastore struct {
	GetFunc func(ctx context.Context, id string) (*Deployment, error)
}

func (m *MockDatastore) Get(ctx context.Context, id string) (*Deployment, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

// Usage
func TestWithManualMock(t *testing.T) {
	mock := &MockDatastore{
		GetFunc: func(ctx context.Context, id string) (*Deployment, error) {
			return &Deployment{ID: id}, nil
		},
	}
	// test logic
}
```

### TypeScript Mocking with Jest

```typescript
// Mock API client
jest.mock('../api/deployment', () => ({
  getDeployment: jest.fn(),
  listDeployments: jest.fn(),
}));

import { getDeployment } from '../api/deployment';

describe('DeploymentComponent', () => {
  it('fetches deployment data', async () => {
    // Setup mock
    (getDeployment as jest.Mock).mockResolvedValue({
      id: '1',
      name: 'test-app',
    });

    // Test component
    render(<DeploymentDetail id="1" />);

    // Verify mock was called
    expect(getDeployment).toHaveBeenCalledWith('1');
  });
});
```

## Test Coverage

### Measuring Coverage

```bash
# Go coverage
make test/go COVERAGE=true

# View coverage report
go tool cover -html=coverage.out

# Web coverage
yarn --cwd web test --coverage

# View web coverage
open web/coverage/lcov-report/index.html
```

### Coverage Goals

- **Critical code**: 90%+ coverage
- **Business logic**: 80%+ coverage
- **UI components**: 70%+ coverage
- **Overall project**: 75%+ coverage

### What to Cover

**High Priority**:
- Business logic
- Error handling
- Security-critical code
- Data transformations

**Lower Priority**:
- Simple getters/setters
- Type definitions
- Generated code

## CI/CD Testing

### GitHub Actions Workflow

Tests run automatically on:
- Pull requests
- Pushes to master
- Release branches

### Test Matrix

```yaml
strategy:
  matrix:
    go-version: [1.21, 1.22]
    os: [ubuntu-latest, macos-latest]
```

### Debugging CI Failures

```bash
# Run tests with same flags as CI
go test -race -coverprofile=coverage.out ./...

# Check for flaky tests
go test -count=100 ./pkg/...

# Run with verbose output
go test -v ./...
```

## Common Testing Patterns

### Table-Driven Tests

```go
tests := []struct {
	name    string
	input   string
	want    string
	wantErr bool
}{
	{"valid input", "test", "TEST", false},
	{"empty input", "", "", true},
}

for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		got, err := Transform(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("unexpected error: %v", err)
		}
		if got != tt.want {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	})
}
```

### Subtests

```go
func TestDeployment(t *testing.T) {
	t.Run("kubernetes", func(t *testing.T) {
		// k8s specific tests
	})

	t.Run("terraform", func(t *testing.T) {
		// terraform specific tests
	})
}
```

### Test Helpers

```go
// Helper function
func assertDeploymentStatus(t *testing.T, d *Deployment, expected DeploymentStatus) {
	t.Helper() // Mark as helper to get correct line numbers
	if d.Status != expected {
		t.Errorf("got status %v, want %v", d.Status, expected)
	}
}

// Usage
func TestDeployment(t *testing.T) {
	deployment := createTestDeployment()
	assertDeploymentStatus(t, deployment, DeploymentStatus_SUCCESS)
}
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- [Jest Documentation](https://jestjs.io/docs/getting-started)

## Questions?

- Ask in [#pipecd Slack](https://cloud-native.slack.com/archives/C01B27F9T0X)
- Open a [GitHub Discussion](https://github.com/pipe-cd/pipecd/discussions)
