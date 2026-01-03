# Security Policy

## Supported Versions

We release security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |
| < 0.x   | :x:                |

> **Note**: Replace with actual version numbers. Generally, we support the latest minor version and provide security patches.

## Reporting a Vulnerability

We take the security of PipeCD seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Where to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **Email**: Contact the maintainers directly (see table below)
2. **Private vulnerability disclosure**: Use [GitHub Security Advisories](https://github.com/pipe-cd/pipecd/security/advisories/new)

### Security Contacts

| Name                    | GitHub ID                                        | Email                                        |
|-------------------------|--------------------------------------------------|----------------------------------------------|
| Tran Cong Khanh         | [@khanhtc1202](https://github.com/khanhtc1202)   | khanhtc1202@gmail.com                        |
| Yoshiki Fujikane        | [@ffjlabo](https://github.com/ffjlabo)           | ffjlabo@gmail.com                            |
| Shinnosuke Sawada-Dazai | [@Warashi](https://github.com/Warashi)           | shin@warashi.dev                             |
| Tetsuya Kikuchi         | [@t-kikuc](https://github.com/t-kikuc)           | tkikuchi07f@gmail.com                        |

### What to Include

Please include the following information in your report:

- **Type of vulnerability** (e.g., XSS, SQL injection, authentication bypass)
- **Full path of affected source file(s)**
- **Location of the affected source code** (tag/branch/commit or direct URL)
- **Step-by-step instructions to reproduce the issue**
- **Proof-of-concept or exploit code** (if possible)
- **Impact of the issue**, including how an attacker might exploit it

### Response Timeline

- **Initial Response**: Within 48 hours of receiving your report
- **Status Update**: Within 7 days with an assessment of the report
- **Fix Timeline**: Depends on severity and complexity
  - Critical: Within 7 days
  - High: Within 30 days
  - Medium: Within 90 days
  - Low: Best effort

### Disclosure Policy

We follow a coordinated disclosure approach:

1. **Private disclosure**: Report is received and acknowledged
2. **Investigation**: We investigate and develop a fix
3. **Fix development**: Fix is developed and tested
4. **Release**: Security patch is released
5. **Public disclosure**: Details are made public after users have had time to update (typically 7-14 days after release)

### Security Updates

Security updates will be announced via:

- GitHub Security Advisories
- Release notes
- CNCF Slack #pipecd channel
- Twitter/X [@pipecd_dev](https://twitter.com/pipecd_dev)

## Security Best Practices

### For Users

1. **Keep PipeCD Updated**: Always use the latest stable version
2. **Use TLS**: Enable TLS for all communications
3. **Rotate Secrets**: Regularly rotate API keys and credentials
4. **Principle of Least Privilege**: Grant minimal necessary permissions
5. **Network Security**: Use network policies to restrict access
6. **Audit Logs**: Regularly review audit logs for suspicious activity

### For Contributors

1. **Input Validation**: Always validate and sanitize user inputs
2. **Authentication**: Use strong authentication mechanisms
3. **Authorization**: Implement proper RBAC checks
4. **Secrets Management**: Never commit secrets to the repository
5. **Dependencies**: Keep dependencies up-to-date and scan for vulnerabilities
6. **Code Review**: All code must be reviewed before merging
7. **Security Testing**: Run security scans (CodeQL, dependency checks)

## Security Features

PipeCD includes the following security features:

- **mTLS Support**: Mutual TLS for secure communication
- **RBAC**: Role-based access control
- **Audit Logging**: Comprehensive audit trail
- **Secret Encryption**: Secrets encrypted at rest
- **API Authentication**: API key and OAuth-based authentication
- **Network Isolation**: No credentials leave the deployment environment

## Known Security Limitations

- Secrets in Git repositories should be encrypted using external tools
- API keys have configurable but not enforced expiration
- Rate limiting is per-instance, not global

## Security Tools

We use the following tools to maintain security:

- **CodeQL**: Automated code scanning
- **Dependabot**: Dependency vulnerability scanning
- **gosec**: Go security checker
- **trivy**: Container image scanning
- **npm audit**: JavaScript dependency scanning

## Vulnerability Management

### CVE Assignment

For confirmed vulnerabilities, we will:
1. Request a CVE number
2. Create a GitHub Security Advisory
3. Publish details after the fix is released

### Hall of Fame

We recognize security researchers who responsibly disclose vulnerabilities:
- [Security Researchers Hall of Fame](https://pipecd.dev/security-hall-of-fame) (coming soon)

## Questions?

If you have questions about this policy, please:
- Join [#pipecd on CNCF Slack](https://cloud-native.slack.com/archives/C01B27F9T0X)
- Email the maintainers
- Open a [GitHub Discussion](https://github.com/pipe-cd/pipecd/discussions)

---

**Note**: This security policy may be updated from time to time. Please check back regularly for updates.
