# Testing Guide for Podgrab

This document provides an overview of the testing infrastructure and best practices for the Podgrab project.

## Table of Contents

- [Overview](#overview)
- [Running Tests](#running-tests)
- [Test Structure](#test-structure)
- [Test Coverage](#test-coverage)
- [Writing New Tests](#writing-new-tests)
- [Testing Best Practices](#testing-best-practices)
- [CI/CD Integration](#cicd-integration)

## Overview

Podgrab uses Go's built-in testing framework along with the [testify](https://github.com/stretchr/testify) assertion library for writing clean, readable tests. The test suite covers:

- **Pure functions** (sanitization, time formatting, validation)
- **Database operations** (using in-memory SQLite)
- **Business logic** (service layer functions)
- **HTTP handlers** (controller tests)

## Running Tests

### Run all tests
```bash
go test ./...
```

### Run tests with verbose output
```bash
go test -v ./...
```

### Run tests with coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Run tests for a specific package
```bash
go test ./internal/sanitize/... -v
go test ./db/... -v
go test ./service/... -v
go test ./model/... -v
```

### Run tests with race detection
```bash
go test -race ./...
```

### Run a specific test
```bash
go test -run TestFunctionName ./package/...
```

## Test Structure

### Directory Structure
```
podgrab/
├── internal/sanitize/
│   ├── sanitize.go
│   └── sanitize_test.go          # Tests for sanitization functions
├── service/
│   ├── naturaltime.go
│   ├── naturaltime_test.go       # Tests for time formatting
│   ├── podcastService.go
│   └── podcastService_test.go    # Tests for podcast business logic
├── model/
│   ├── queryModels.go
│   └── queryModels_test.go       # Tests for validation logic
├── db/
│   ├── dbfunctions.go
│   ├── dbfunctions_test.go       # Tests for database operations
│   └── test_helpers.go           # Database test utilities
└── controllers/
    ├── podcast.go
    └── podcast_test.go            # Tests for HTTP handlers
```

### Test File Naming
- Test files are named `*_test.go`
- Test functions are named `Test<FunctionName>`
- Helper functions are in `test_helpers.go`

## Test Coverage

Current test coverage by package:

| Package | Coverage | Description |
|---------|----------|-------------|
| internal/sanitize | 92.5% | HTML/string sanitization functions |
| model | 88.9% | Query models and validation |
| db | 21.8% | Database operations |
| service | 8.0% | Business logic (naturaltime tested) |
| controllers | 0.0% | HTTP handlers (future work) |

### View Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Writing New Tests

### 1. Pure Function Tests

For functions with no external dependencies (like `internal/sanitize`):

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "test case 1",
            input:    "input value",
            expected: "expected output",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionName(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 2. Database Tests

For testing database operations:

```go
func TestDatabaseFunction(t *testing.T) {
    // Setup test database
    db, err := SetupTestDB()
    require.NoError(t, err)
    defer TeardownTestDB(db)

    // Create test data
    podcast, err := CreateTestPodcast(db, "Test Podcast")
    require.NoError(t, err)

    // Run your test
    var result Podcast
    err = GetPodcastByURL(podcast.URL, &result)

    // Assert results
    assert.NoError(t, err)
    assert.Equal(t, podcast.Title, result.Title)
}
```

### 3. Service Layer Tests (with Mocks)

For testing business logic with external dependencies:

```go
// TODO: Add examples when service tests are implemented
```

### 4. HTTP Handler Tests

For testing controllers:

```go
// TODO: Add examples when controller tests are implemented
```

## Testing Best Practices

### 1. Table-Driven Tests
Use table-driven tests for testing multiple scenarios:
```go
tests := []struct {
    name     string
    input    interface{}
    expected interface{}
}{
    // Test cases...
}
```

### 2. Use testify Assertions
Prefer `assert` and `require` from testify:
- Use `assert` when test can continue after failure
- Use `require` when test should stop after failure

```go
require.NoError(t, err)  // Stop if error
assert.Equal(t, expected, actual)  // Continue if not equal
```

### 3. Test Naming
- Test function: `TestFunctionName`
- Test case: Descriptive name explaining what's being tested
- Use underscores in test case names for readability

```go
t.Run("returns_error_when_input_is_empty", func(t *testing.T) {
    // ...
})
```

### 4. Setup and Teardown
Use `defer` for cleanup:
```go
db, err := SetupTestDB()
require.NoError(t, err)
defer TeardownTestDB(db)
```

### 5. Test Independence
Each test should be independent and not rely on state from other tests.

### 6. Use Test Helpers
Create helper functions for common test setup:
```go
// Located in db/test_helpers.go
func CreateTestPodcast(db *gorm.DB, title string) (*Podcast, error)
func CreateTestPodcastItem(db *gorm.DB, podcast *Podcast, title string, status DownloadStatus) (*PodcastItem, error)
```

### 7. Mock External Dependencies
For HTTP clients, file systems, and other external resources:
- Create interfaces for external dependencies
- Use mock implementations in tests
- Consider using [testify/mock](https://pkg.go.dev/github.com/stretchr/testify/mock)

## CI/CD Integration

### GitHub Actions
Tests automatically run on:
- Every push to `main`, `master`, or `claude/**` branches
- Every pull request to `main` or `master`

See `.github/workflows/test.yml` for configuration.

### Coverage Reports
Test coverage is automatically:
- Generated for each test run
- Uploaded to Codecov (if configured)
- Displayed in GitHub Actions summary

## Database Testing

### Important: Parallel Test Limitations
⚠️ **Note on Parallel Testing**: The current database layer uses a global `DB` variable (`db.DB`), which creates potential race conditions when running tests in parallel. For this reason:

- Database tests should NOT use `t.Parallel()`
- Run tests sequentially to avoid race conditions
- Future improvement: refactor to use dependency injection instead of global DB

### In-Memory SQLite
Tests use in-memory SQLite databases for fast, isolated testing:

```go
db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
```

Benefits:
- No external dependencies
- Fast test execution
- Clean state for each test
- Automatic cleanup

### Test Helpers
Use provided helpers in `db/test_helpers.go`:
- `SetupTestDB()` - Create and migrate test database
- `TeardownTestDB(db)` - Close test database
- `CreateTestPodcast()` - Create test podcast
- `CreateTestPodcastItem()` - Create test episode
- `CreateTestTag()` - Create test tag

## Future Improvements

### High Priority
- [ ] Add tests for `service/podcastService.go` (core business logic)
- [ ] Add tests for `service/fileService.go` (file operations)
- [ ] Add tests for `controllers/podcast.go` (HTTP handlers)

### Medium Priority
- [ ] Add integration tests for end-to-end workflows
- [ ] Add tests for RSS parsing logic
- [ ] Add tests for download concurrency

### Low Priority
- [ ] Add benchmark tests for performance-critical functions
- [ ] Add fuzzing tests for parser functions
- [ ] Set up mutation testing

## Troubleshooting

### Tests Fail Locally But Pass in CI
- Ensure you're using the correct Go version (check `go.mod`)
- Run `go mod tidy` to sync dependencies
- Check for environment-specific issues

### Database Tests Fail
- Verify SQLite is available
- Check that migrations run successfully
- Ensure test database is properly cleaned up

### Coverage Seems Low
- Run with `-coverprofile` to generate detailed report
- Use `go tool cover -html` to visualize untested code
- Focus on testing critical business logic first

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [GORM Testing Guide](https://gorm.io/docs/testing.html)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

## Contributing

When adding new features:
1. Write tests first (TDD approach recommended)
2. Ensure all tests pass: `go test ./...`
3. Check coverage: `go test -cover ./...`
4. Aim for >80% coverage on new code
5. Document any new test helpers or patterns

## Questions?

If you have questions about testing or need help writing tests for a specific feature, please open an issue or reach out to the maintainers.
