# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 2.x.x   | :white_check_mark: |
| < 2.0   | :x:                |

## Reporting a Vulnerability

We take the security of Spotigo seriously. If you believe you have found a security vulnerability, please report it to us responsibly.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:

**Email:** [baalateja.k@gmail.com](mailto:baalateja.k@gmail.com)

Please include the following information in your report:

- **Type of vulnerability** (e.g., token exposure, injection, authentication bypass)
- **Location** of the affected source code (file, line number if known)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept** or exploit code (if available)
- **Impact** of the vulnerability and potential exploitation scenarios
- **Suggested fix** (if you have one)

### What to Expect

1. **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours.

2. **Assessment**: We will investigate and validate the reported vulnerability within 7 days.

3. **Resolution**: We aim to resolve critical vulnerabilities within 30 days. You will be kept informed of our progress.

4. **Disclosure**: Once the vulnerability is fixed, we will:
   - Release a security update
   - Publish a security advisory on GitHub
   - Credit you in the advisory (unless you prefer to remain anonymous)

### Safe Harbor

We consider security research conducted in accordance with this policy to be:

- Authorized in accordance with the Computer Fraud and Abuse Act (CFAA)
- Exempt from DMCA restrictions on circumvention of technology controls
- Lawful and welcome

We will not pursue legal action against researchers who:

- Follow this responsible disclosure policy
- Make a good faith effort to avoid privacy violations and data destruction
- Do not exploit vulnerabilities beyond proof-of-concept

## Security Measures in Spotigo

### Token Encryption

Spotify OAuth tokens are encrypted using **AES-256-GCM** with machine-specific key derivation:

- Keys are derived using PBKDF2 with SHA-256
- Salt includes machine-specific identifiers (hostname, username)
- Tokens are automatically encrypted when saved
- Backward compatible with plaintext tokens (auto-migrates)

### Local Processing

All AI processing happens **locally** via Ollama:

- No data is sent to external AI services
- Your music library data never leaves your machine
- Embeddings are generated and stored locally

### Data Privacy

- **No telemetry**: Spotigo does not collect or transmit usage data
- **No analytics**: No third-party analytics or tracking
- **Local storage**: All data is stored in your local filesystem
- **No cloud sync**: Data is not synced to any cloud service

### Minimal Permissions

Spotigo requests only the necessary Spotify API scopes:

- `user-library-read` - Read saved tracks and albums
- `playlist-read-private` - Read private playlists
- `user-follow-read` - Read followed artists
- `user-read-recently-played` - Read recently played tracks
- `user-top-read` - Read top artists and tracks

### Secure Defaults

- OAuth callback runs on localhost only
- Token files are created with restrictive permissions (0600)
- Configuration files should not be committed to version control

## Security Best Practices for Users

### Configuration Security

1. **Never commit credentials** to version control:
   ```bash
   # .gitignore should include:
   spotigo.yaml
   .spotify_token
   *.enc
   ```

2. **Use environment variables** for sensitive data:
   ```bash
   export SPOTIFY_CLIENT_ID=your_client_id
   export SPOTIFY_CLIENT_SECRET=your_client_secret
   ```

3. **Restrict file permissions**:
   ```bash
   chmod 600 spotigo.yaml
   chmod 600 .spotify_token
   ```

### Token Management

1. **Rotate tokens periodically** by logging out and re-authenticating:
   ```bash
   spotigo auth logout
   spotigo auth
   ```

2. **Revoke access** if you suspect compromise:
   - Visit [Spotify Account Settings](https://www.spotify.com/account/apps/)
   - Remove Spotigo from authorized apps
   - Delete local token file
   - Re-authenticate

### Network Security

1. **Use local Ollama** - Do not expose Ollama to the internet
2. **Firewall rules** - Ensure port 11434 is not publicly accessible
3. **VPN considerations** - OAuth callback requires localhost access

## Known Security Considerations

### Current Limitations

1. **Token file location** - Token file path is configurable but not encrypted at rest on disk (the token content is encrypted)

2. **Config file** - `spotigo.yaml` may contain sensitive data; ensure proper file permissions

3. **Backup data** - Backup files contain your music library data in plaintext JSON; store securely

### Planned Improvements

- [ ] Keyring/keychain integration for token storage
- [ ] Encrypted backup files
- [ ] Config file encryption option

## Vulnerability Disclosure History

| Date | Description | Severity | Fixed In |
|------|-------------|----------|----------|
| - | No vulnerabilities disclosed yet | - | - |

---

Thank you for helping keep Spotigo and its users safe!
