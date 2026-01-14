# Security

This document describes the security measures and best practices implemented in Spotigo.

## Overview

Spotigo takes security seriously and implements multiple layers of protection to safeguard user data, particularly OAuth tokens and personal Spotify information.

## Security Measures

### 1. OAuth Token Protection

#### Encryption at Rest
- All OAuth tokens are encrypted using **AES-256-GCM** before being written to disk
- Encryption keys are derived from machine-specific entropy sources
- Token files use restrictive permissions (`0600` - owner read/write only)

#### Key Derivation
The encryption key is derived from multiple machine-specific sources:
- Username (from `USER` or `USERNAME` environment variables)
- Home directory path
- OS and architecture information
- User config directory
- Application-specific salt

This ensures that:
- Tokens cannot be easily moved between machines
- Each installation has a unique encryption key
- Attackers cannot decrypt tokens without the specific machine context

#### Legacy Token Support
- Spotigo can read older plaintext tokens for backward compatibility
- On next save, plaintext tokens are automatically upgraded to encrypted format
- A warning is displayed if encryption fails and plaintext storage is used

### 2. Path Traversal Prevention

All file operations use `filepath.Clean()` to sanitize paths and prevent directory traversal attacks:

#### Protected Operations
- **Config loading** (`internal/config/models.go`)
- **Token storage** (`internal/spotify/client.go`, `internal/crypto/crypto.go`)
- **Backup files** (`internal/storage/store.go`, `internal/cmd/backup.go`)
- **JSON data files** (`internal/jsonutil/jsonutil.go`, `internal/jsonquery/query.go`)

#### Implementation
```go
// Example: Clean path before file operations
cleanPath := filepath.Clean(userProvidedPath)
data, err := os.ReadFile(cleanPath) // #nosec G304 - path is sanitized
```

### 3. Secure File Permissions

Spotigo creates files and directories with restrictive permissions:

| Type | Permission | Description |
|------|-----------|-------------|
| Token files | `0600` | Owner read/write only |
| Encryption directories | `0700` | Owner full access only |
| Backup directories | `0750` | Owner full, group read/execute |
| Data directories | `0750` | Owner full, group read/execute |

### 4. HTTP Server Security

The OAuth callback server implements security best practices:

#### Timeouts
- **Read header timeout**: 5 seconds
- **Overall timeout**: 5 minutes for the entire auth flow
- Prevents slowloris and similar attacks

#### State Parameter Validation
- Random 32-byte state parameter for CSRF protection
- State is verified on callback before token exchange
- Invalid state immediately rejects the request

#### Graceful Shutdown
- Server properly shuts down after successful authentication
- 2-second grace period for final response delivery
- Prevents resource leaks

### 5. Subprocess Execution

Browser opening uses platform-specific, safe methods:

- **Windows**: `rundll32 url.dll,FileProtocolHandler` - prevents command injection
- **macOS**: `open` command with URL argument
- **Linux**: `xdg-open` with URL argument

All subprocess executions are marked with `#nosec G204` and include justification comments.

### 6. Code Scanning

#### gosec Integration
Spotigo is regularly scanned with gosec (Go Security Checker):
- **G101**: Hardcoded credentials detection
- **G104**: Unchecked errors
- **G204**: Subprocess command injection
- **G304**: File path traversal

False positives are explicitly marked with `#nosec` directives and justification comments.

#### CodeQL Analysis
GitHub CodeQL scanning runs on:
- Every push to main
- Pull requests
- Weekly scheduled scans

### 7. Dependency Security

#### Secure Dependencies
- `golang.org/x/oauth2` - Official OAuth2 library
- `github.com/zmb3/spotify/v2` - Well-maintained Spotify SDK
- `golang.org/x/crypto` - Official cryptography libraries (future use)

#### Regular Updates
Dependencies are regularly reviewed and updated for security patches.

## Threat Model

### In Scope
1. **Local file access**: Protection against unauthorized file reading
2. **Token theft**: Encryption prevents casual token extraction
3. **Path traversal**: Sanitization prevents directory escape
4. **CSRF**: State parameter prevents cross-site attacks
5. **Command injection**: Safe subprocess execution

### Out of Scope
1. **Memory dumps**: Running process memory is not protected
2. **Privileged attackers**: Root/admin users can bypass file permissions
3. **Physical access**: Full disk encryption is user's responsibility
4. **Network attacks**: TLS is provided by Spotify's infrastructure

## Best Practices for Users

### Recommended Setup
1. **Use strong OS passwords**: Protects the machine-specific encryption key
2. **Enable full disk encryption**: Adds an additional layer of protection
3. **Keep software updated**: Run `spotigo` updates regularly
4. **Secure your environment**: Don't share your `.spotify_token` file

### Configuration Security
1. **Client secrets**: Store in environment variables, not in config files
   ```bash
   export SPOTIFY_CLIENT_SECRET="your_secret_here"
   ```
2. **Config files**: Keep `config.yaml` in a secure location with `0600` permissions
3. **Backup files**: Store backups in a secure location, consider additional encryption

### Token Management
1. **Regular rotation**: Re-authenticate periodically (tokens auto-refresh but rotation helps)
2. **Logout when needed**: Use `spotigo auth logout` to remove tokens
3. **Check permissions**: Verify only necessary OAuth scopes are granted

## Security Reporting

### Reporting Vulnerabilities

If you discover a security vulnerability in Spotigo, please report it privately:

1. **Do NOT** open a public GitHub issue
2. Email the maintainer with details (see README for contact)
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### Response Timeline
- **Initial response**: Within 48 hours
- **Status update**: Within 7 days
- **Fix timeline**: Varies by severity
  - Critical: 24-48 hours
  - High: 7 days
  - Medium: 30 days
  - Low: Next release

## Security Checklist for Developers

When adding new features:

- [ ] Validate all user input
- [ ] Sanitize all file paths with `filepath.Clean()`
- [ ] Use restrictive file permissions
- [ ] Avoid hardcoding secrets (use env vars or secure storage)
- [ ] Implement proper error handling
- [ ] Add `#nosec` directives with justification for false positives
- [ ] Run `gosec ./...` before committing
- [ ] Review CodeQL findings
- [ ] Update tests to cover security-critical paths
- [ ] Document security implications in code comments

## Security Testing

### Running Security Scans

```bash
# Install gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run security scan
gosec -exclude-generated ./...

# Run with specific checks
gosec -include=G101,G104,G204,G304 ./...

# Generate report
gosec -fmt=json -out=results.json ./...
```

### Running Tests

```bash
# Run all tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Check coverage
go tool cover -html=coverage.out

# Run security-focused tests
go test -v -run=TestSecurity ./...
```

## Compliance

### Data Protection
- **No external transmission**: All data stays local except Spotify API calls
- **No analytics**: Spotigo doesn't collect usage data
- **No telemetry**: No phone-home functionality

### OAuth Best Practices
- Uses Authorization Code flow (most secure)
- Implements state parameter for CSRF protection
- Stores tokens encrypted at rest
- Auto-refreshes tokens securely
- Supports token revocation via logout

## Additional Resources

- [OWASP Go Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Go_Security_Cheat_Sheet.html)
- [OAuth 2.0 Security Best Current Practice](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics)
- [Spotify OAuth Documentation](https://developer.spotify.com/documentation/web-api/concepts/authorization)
- [gosec Documentation](https://github.com/securego/gosec)

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2024-01-XX | Initial security documentation |

---

**Last Updated**: 2024
**Security Contact**: See README for maintainer contact information