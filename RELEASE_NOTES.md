# Release Notes

## Version [TBD] - January 2026

### 🔒 Security Enhancements

#### Path Traversal Protection
- **Added comprehensive path validation** to prevent directory traversal attacks
  - New `security.ValidatePath()` function validates all file paths
  - Rejects paths containing `..` sequences that could escape intended directories
  - Applied to configuration files, plan files, external variables, transaction includes, and response files
  - Affects: `config/config.go`, `handler/tests.go`, `structs/plan/plan.go`

#### TLS Security Improvements
- **SlackPost TLS configuration now configurable** (previously hardcoded to skip verification)
  - Added `insecureSkipVerify` parameter to control TLS certificate validation
  - Default behavior should be `false` in production environments
  - Only set to `true` in trusted internal networks
  - See [Breaking Changes](#breaking-changes) below

### ⚠️ Breaking Changes

#### SlackPost API Change
The `SlackPost` function signature has changed to include TLS verification control:

```go
// Previous signature:
func SlackPost(payload []byte, url string) error

// New signature:
func SlackPost(payload []byte, url string, insecureSkipVerify bool) error
```

**Migration:** Update all calls to `SlackPost()` to include the third parameter. For production use, pass `false` to ensure TLS certificate validation. Only use `true` in trusted internal networks.

```go
// Example migration:
// Before:
err := core.SlackPost(payload, url)

// After (production):
err := core.SlackPost(payload, url, false)

// After (trusted internal network):
err := core.SlackPost(payload, url, true)
```

### 🐛 Bug Fixes

- **Fixed SlackPost error handling** - Now properly returns an error when Slack API responds with non-2xx status codes (previously returned nil on failures)

### 📦 Dependency Updates

#### Go Version
- Upgraded from Go 1.13 to Go 1.24.0

#### Major Dependencies
- `gin-gonic/gin`: 1.5.0 → 1.11.0
- `prometheus/client_golang`: 1.3.0 → 1.23.2
- `stretchr/testify`: 1.4.0 → 1.11.1
- `gofrs/uuid`: 3.2.0+incompatible → 4.4.0+incompatible
- `juju/loggo`: 0.0.0-20190526231331 → 1.0.0
- `mitchellh/mapstructure`: 1.2.2 → 1.5.0

See `go.mod` for complete dependency list.

### 🔧 Internal Improvements

#### Deprecated API Replacements
Replaced deprecated `ioutil` package with modern `os` and `io` packages (Go 1.16+ best practices):
- `ioutil.ReadFile` → `os.ReadFile`
- `ioutil.ReadAll` → `io.ReadAll`
- `ioutil.TempDir` → `os.MkdirTemp`

Files affected: `config/config.go`, `core/alerts.go`, `handler/tests.go`, `structs/plan/plan.go`

#### Test Coverage Improvements
- Added comprehensive test coverage for security validations
- Added tests for path traversal protection
- Expanded test coverage for handlers, queue operations, metrics, and routing
- Added integration tests for end-to-end workflows

### 📝 Notes

- All changes maintain backward compatibility except for the `SlackPost` API modification
- Path validation changes may reject previously-accepted (but unsafe) file paths containing `..`
- Recommended to review all `SlackPost` calls before upgrading to ensure proper TLS verification settings

### 🔄 Upgrade Checklist

1. ✅ Update Go to version 1.24.0 or later
2. ✅ Review and update all `SlackPost()` function calls to include the `insecureSkipVerify` parameter
3. ✅ Ensure configuration files, plan files, and transaction includes use safe paths (no `..` sequences)
4. ✅ Run test suite to verify compatibility
5. ✅ Review TLS settings for Slack integrations

---

For questions or issues, please open an issue on the repository.
