# PipeCD Roadmap

This document outlines the future direction and planned features for PipeCD. Items are subject to change based on community feedback and priorities.

## Vision

PipeCD aims to be the most versatile, secure, and developer-friendly GitOps continuous delivery platform for multi-cloud environments.

## Current Status (as of 2026)

- ✅ Stable v0.x releases
- ✅ Support for Kubernetes, Terraform, Cloud Run, Lambda, ECS
- ✅ Multi-cloud deployment capabilities
- ✅ Built-in deployment analysis
- ✅ CNCF Sandbox Project

## Short-term Goals (Q1-Q2 2026)

### Enhanced Security
- [ ] External secret management integration (Vault, AWS Secrets Manager)
- [ ] Enhanced RBAC with fine-grained permissions
- [ ] Audit log improvements with filtering and export
- [ ] mTLS enforcement options
- [ ] Security compliance reports (SOC2, ISO27001)

### Developer Experience
- [ ] Improved CLI with interactive mode
- [ ] VS Code extension for pipeline editing
- [ ] Local development environment improvements
- [ ] Better error messages and debugging tools
- [ ] Pipeline dry-run capabilities

### Platform Support
- [ ] Azure DevOps integration
- [ ] Google Cloud Deploy integration
- [ ] Enhanced GitLab CI/CD integration
- [ ] Bitbucket Pipelines support

### Observability
- [ ] Enhanced metrics and dashboards
- [ ] Distributed tracing support
- [ ] Better log aggregation
- [ ] Real-time deployment notifications
- [ ] SLO/SLI tracking

## Mid-term Goals (Q3-Q4 2026)

### Advanced Deployment Strategies
- [ ] Flagger integration for progressive delivery
- [ ] Multi-region deployment orchestration
- [ ] Canary deployment with automated rollback
- [ ] Blue-green deployment improvements
- [ ] Shadow deployments

### Scalability & Performance
- [ ] Horizontal scaling improvements
- [ ] Database sharding support
- [ ] Caching layer enhancements
- [ ] Rate limiting per tenant
- [ ] Resource usage optimization

### Plugin Ecosystem
- [ ] Plugin marketplace
- [ ] Community plugin contributions
- [ ] Plugin SDK improvements
- [ ] Plugin versioning and compatibility
- [ ] Plugin documentation generator

### AI/ML Integration
- [ ] AI-powered deployment recommendations
- [ ] Anomaly detection in deployments
- [ ] Predictive rollback suggestions
- [ ] Intelligent resource optimization
- [ ] ChatOps with AI assistance

## Long-term Goals (2027+)

### Multi-Cluster Management
- [ ] Federation support
- [ ] Cross-cluster deployments
- [ ] Unified dashboard for multiple clusters
- [ ] Disaster recovery automation
- [ ] Global deployment policies

### Advanced GitOps
- [ ] Git-native secret management
- [ ] Multi-repository support improvements
- [ ] Monorepo optimization
- [ ] GitOps workflows for infrastructure
- [ ] Policy-as-code enforcement

### Enterprise Features
- [ ] Multi-tenancy improvements
- [ ] SSO integrations (LDAP, SAML, OIDC)
- [ ] Cost tracking and optimization
- [ ] Compliance automation
- [ ] SLA management

### Platform Expansion
- [ ] Azure Kubernetes Service native support
- [ ] OpenShift dedicated support
- [ ] Nomad deployment support
- [ ] Edge computing platforms
- [ ] Serverless framework support

## Community & Ecosystem

### Documentation
- [x] Architecture documentation
- [x] Troubleshooting guide
- [x] Development guide
- [x] Testing guide
- [x] API documentation
- [ ] Video tutorials
- [ ] Interactive tutorials
- [ ] Localization (i18n)

### Community Growth
- [ ] More community meetings (regional)
- [ ] PipeCD certification program
- [ ] Ambassador program
- [ ] Conference talks and workshops
- [ ] University partnerships

### Governance
- [ ] Formal governance model
- [ ] Technical steering committee
- [ ] Special interest groups (SIGs)
- [ ] Working groups for major features
- [ ] CNCF Incubating project status

## How to Contribute to the Roadmap

We welcome community input on the roadmap!

### Suggest New Features

1. **Search existing proposals**: Check [GitHub Discussions](https://github.com/pipe-cd/pipecd/discussions)
2. **Create a proposal**: Open a new discussion with:
   - Problem statement
   - Proposed solution
   - Expected impact
   - Implementation considerations
3. **Gather feedback**: Engage with the community
4. **Submit RFC**: For major features, submit a formal RFC

### Vote on Features

- 👍 React to discussions you're interested in
- 💬 Comment with use cases
- 🎯 Share your priorities in community meetings

### Implement Features

1. **Claim an item**: Comment on the roadmap item
2. **Create design doc**: For significant features
3. **Submit PR**: Follow [contributing guidelines](./CONTRIBUTING.md)
4. **Iterate**: Work with maintainers on feedback

## Roadmap Status Legend

- 🎯 **Planned**: On the roadmap, not started
- 🚧 **In Progress**: Actively being worked on
- ✅ **Completed**: Delivered in a release
- 🔄 **Under Review**: Design or implementation review
- ⏸️ **Paused**: Temporarily on hold
- ❌ **Cancelled**: Removed from roadmap

## Versioning Strategy

PipeCD follows semantic versioning:

- **Major (v1.0, v2.0)**: Breaking changes, major features
- **Minor (v0.1, v0.2)**: New features, backward compatible
- **Patch (v0.1.1, v0.1.2)**: Bug fixes, security patches

### Path to v1.0

Requirements for v1.0 release:
- [ ] API stability guarantee
- [ ] Comprehensive documentation
- [ ] Security audit completed
- [ ] Performance benchmarks met
- [ ] Production deployments at scale (>1000 apps)
- [ ] 6+ months of stable releases

**Estimated timeline**: Q4 2026

## Release Cadence

- **Minor releases**: Every 2-3 months
- **Patch releases**: As needed for critical fixes
- **Security releases**: Immediate for critical vulnerabilities

## Deprecation Policy

- **Announcement**: Features marked deprecated in release notes
- **Deprecation period**: Minimum 2 minor releases (6 months)
- **Removal**: After deprecation period expires
- **Migration guide**: Provided for all deprecated features

## Recent Completions

### Q4 2025
- ✅ Improved documentation and guides
- ✅ Enhanced CI/CD workflows
- ✅ Dependabot integration
- ✅ Pre-commit hooks
- ✅ Comprehensive testing guide

### Q3 2025
- ✅ Plugin architecture improvements
- ✅ Web UI performance enhancements
- ✅ Multi-platform build support
- ✅ ARM64 support

## Feedback

We value your feedback on this roadmap!

- **Community Meetings**: Join bi-weekly meetings
- **Slack**: [#pipecd on CNCF Slack](https://cloud-native.slack.com/archives/C01B27F9T0X)
- **Discussions**: [GitHub Discussions](https://github.com/pipe-cd/pipecd/discussions)
- **Issues**: [Feature Requests](https://github.com/pipe-cd/pipecd/issues/new?template=new-feature.md)

---

**Last Updated**: January 2026  
**Next Review**: March 2026

*This roadmap is a living document and will be updated regularly based on community feedback and project evolution.*
