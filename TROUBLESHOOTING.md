# PipeCD Troubleshooting Guide

This guide helps you diagnose and resolve common issues with PipeCD.

## Table of Contents

- [General Troubleshooting](#general-troubleshooting)
- [Control Plane Issues](#control-plane-issues)
- [Piped Agent Issues](#piped-agent-issues)
- [Deployment Issues](#deployment-issues)
- [Web UI Issues](#web-ui-issues)
- [Performance Issues](#performance-issues)
- [Network & Connectivity](#network--connectivity)
- [Database Issues](#database-issues)
- [Getting Help](#getting-help)

## General Troubleshooting

### Check Component Versions

Ensure all components are running compatible versions:

```bash
# Check PipeCD control plane version
kubectl exec -n pipecd deployment/pipecd -- /pipecd version

# Check Piped version
kubectl exec -n pipecd deployment/piped -- /piped version

# Check pipectl version
pipectl version
```

### Enable Debug Logging

Enable debug logging for more detailed information:

**Control Plane:**
```yaml
# In pipecd configuration
spec:
  logLevel: debug
```

**Piped:**
```yaml
# In piped configuration
spec:
  logLevel: debug
```

### View Logs

```bash
# Control Plane logs
kubectl logs -n pipecd deployment/pipecd -f

# Piped logs
kubectl logs -n pipecd deployment/piped -f

# View logs for specific container
kubectl logs -n pipecd deployment/pipecd -c pipecd -f

# Previous container logs (if crashed)
kubectl logs -n pipecd deployment/piped --previous
```

## Control Plane Issues

### Issue: Control Plane Fails to Start

**Symptoms**: Pod in CrashLoopBackOff state

**Common Causes & Solutions**:

1. **Database Connection Failure**
   ```bash
   # Check logs for database errors
   kubectl logs -n pipecd deployment/pipecd | grep -i "database\|connection"
   
   # Verify database credentials
   kubectl get secret -n pipecd pipecd-secrets -o yaml
   
   # Test database connectivity
   kubectl run -n pipecd mysql-client --rm -it --image=mysql:8.0 -- \
     mysql -h <db-host> -u <username> -p
   ```

2. **Missing Configuration**
   ```bash
   # Verify ConfigMap exists
   kubectl get configmap -n pipecd pipecd-config
   
   # Check configuration content
   kubectl get configmap -n pipecd pipecd-config -o yaml
   ```

3. **Insufficient Resources**
   ```bash
   # Check resource limits
   kubectl describe pod -n pipecd -l app=pipecd
   
   # Adjust resources in deployment
   kubectl edit deployment -n pipecd pipecd
   ```

### Issue: Unable to Login to Web UI

**Symptoms**: Authentication fails or redirects loop

**Solutions**:

1. **Check Static Admin Configuration**
   ```bash
   # Verify admin account exists in config
   kubectl get configmap -n pipecd pipecd-config -o yaml | grep -A 5 "staticAdmin"
   ```

2. **SSO Configuration Issues**
   ```bash
   # Check SSO settings
   kubectl get configmap -n pipecd pipecd-config -o yaml | grep -A 10 "sso"
   
   # Verify callback URLs match
   # GitHub: Settings -> Developer settings -> OAuth Apps
   # Google: Cloud Console -> APIs & Services -> Credentials
   ```

3. **Cookie/Session Issues**
   - Clear browser cache and cookies
   - Try incognito/private browsing mode
   - Check browser console for errors (F12)

### Issue: API Requests Failing

**Symptoms**: gRPC errors, timeout errors

**Solutions**:

```bash
# Check service endpoints
kubectl get svc -n pipecd

# Test gRPC connectivity
grpcurl -plaintext localhost:9080 list

# Verify TLS certificates (if using)
openssl s_client -connect <domain>:443 -servername <domain>

# Check for network policies
kubectl get networkpolicies -n pipecd
```

## Piped Agent Issues

### Issue: Piped Unable to Connect to Control Plane

**Symptoms**: "connection refused", "context deadline exceeded"

**Solutions**:

1. **Verify API Address**
   ```yaml
   # In piped config, check apiAddress
   spec:
     apiAddress: <control-plane-host>:443
   ```

2. **Check Network Connectivity**
   ```bash
   # From piped pod
   kubectl exec -n pipecd deployment/piped -- nc -zv <control-plane-host> 443
   
   # Check DNS resolution
   kubectl exec -n pipecd deployment/piped -- nslookup <control-plane-host>
   ```

3. **Verify Piped Credentials**
   ```bash
   # Check piped ID and key
   kubectl get secret -n pipecd piped-secret -o yaml
   
   # Regenerate piped key if needed (via Web UI)
   ```

4. **Firewall/Network Policies**
   - Ensure firewall allows outbound HTTPS (443)
   - Check cloud provider security groups
   - Verify Kubernetes NetworkPolicies

### Issue: Piped Not Syncing Git Repository

**Symptoms**: No deployments triggered despite Git changes

**Solutions**:

1. **Verify Git Configuration**
   ```yaml
   # Check piped config
   spec:
     repositories:
       - repoId: example
         remote: git@github.com:org/repo.git
         branch: main
   ```

2. **Check SSH Keys**
   ```bash
   # For SSH repositories
   kubectl exec -n pipecd deployment/piped -- ssh -T git@github.com
   
   # Verify SSH key secret
   kubectl get secret -n pipecd piped-git-ssh-key
   ```

3. **Check Git Sync Interval**
   ```yaml
   spec:
     syncInterval: 1m  # Increase if too frequent
   ```

4. **Review Git Sync Logs**
   ```bash
   kubectl logs -n pipecd deployment/piped | grep -i "git\|sync"
   ```

### Issue: Piped Using Old Version

**Symptoms**: Features not working, incompatibility errors

**Solutions**:

```bash
# Check current version
kubectl exec -n pipecd deployment/piped -- /piped version

# Update image in deployment
kubectl set image deployment/piped -n pipecd piped=ghcr.io/pipe-cd/piped:v0.x.x

# Or use launcher for auto-updates
# Configure remote upgrade in piped config
```

## Deployment Issues

### Issue: Deployment Stuck in Pending

**Symptoms**: Deployment doesn't progress

**Solutions**:

1. **Check Piped Status**
   ```bash
   # Ensure piped is running
   kubectl get pods -n pipecd -l app=piped
   
   # Check piped logs
   kubectl logs -n pipecd deployment/piped -f
   ```

2. **Verify Pipeline Configuration**
   ```yaml
   # Check .pipe.yaml in repository
   apiVersion: pipecd.dev/v1beta1
   kind: KubernetesApp
   spec:
     pipeline:
       stages:
         - name: K8S_SYNC
   ```

3. **Check Resource Quota**
   ```bash
   # Verify namespace quotas
   kubectl describe resourcequota -n <target-namespace>
   ```

### Issue: Deployment Fails with Analysis Error

**Symptoms**: Analysis stage fails, deployment rolls back

**Solutions**:

1. **Check Analysis Configuration**
   ```yaml
   # Verify metrics provider is configured
   spec:
     analysisProviders:
       - name: prometheus
         type: PROMETHEUS
         config:
           address: http://prometheus:9090
   ```

2. **Verify Metrics Query**
   ```bash
   # Test Prometheus query manually
   curl "http://prometheus:9090/api/v1/query?query=<your-query>"
   ```

3. **Review Analysis Logs**
   ```bash
   kubectl logs -n pipecd deployment/piped | grep -i "analysis"
   ```

### Issue: Kubernetes Manifest Apply Fails

**Symptoms**: "error applying manifest", "invalid resource"

**Solutions**:

1. **Validate Kubernetes Manifests**
   ```bash
   # Dry-run apply
   kubectl apply --dry-run=server -f manifest.yaml
   
   # Validate YAML
   kubectl apply --validate=true -f manifest.yaml
   ```

2. **Check RBAC Permissions**
   ```bash
   # Verify piped service account permissions
   kubectl auth can-i create deployments --as=system:serviceaccount:pipecd:piped -n <namespace>
   ```

3. **Review Piped Logs for Details**
   ```bash
   kubectl logs -n pipecd deployment/piped | grep -B5 -A5 "error"
   ```

### Issue: Terraform Deployment Fails

**Symptoms**: Terraform plan/apply errors

**Solutions**:

1. **Check Terraform Version**
   ```bash
   # Verify terraform binary version
   kubectl exec -n pipecd deployment/piped -- terraform version
   ```

2. **Verify Cloud Credentials**
   ```bash
   # Check environment variables or mounted secrets
   kubectl describe pod -n pipecd -l app=piped
   ```

3. **Enable Terraform Debug Logs**
   ```yaml
   # In pipeline config
   spec:
     pipeline:
       stages:
         - name: TERRAFORM_APPLY
           with:
             args:
               - -debug
   ```

## Web UI Issues

### Issue: Web UI Not Loading

**Symptoms**: Blank page, 404 errors

**Solutions**:

1. **Check Service and Ingress**
   ```bash
   # Verify service
   kubectl get svc -n pipecd pipecd
   
   # Check ingress
   kubectl get ingress -n pipecd
   kubectl describe ingress -n pipecd pipecd
   ```

2. **Browser Console Errors**
   - Open browser DevTools (F12)
   - Check Console tab for JavaScript errors
   - Check Network tab for failed requests

3. **Clear Browser Cache**
   ```
   Chrome: Ctrl+Shift+Delete
   Firefox: Ctrl+Shift+Delete
   Safari: Cmd+Option+E
   ```

### Issue: Deployment List Not Showing

**Symptoms**: Empty deployment list despite having deployments

**Solutions**:

1. **Check Project Filter**
   - Verify correct project selected in UI
   - Check project ID in URL parameters

2. **Verify Database Connection**
   ```bash
   # Check control plane logs
   kubectl logs -n pipecd deployment/pipecd | grep -i "database"
   ```

3. **Check API Communication**
   ```
   # Browser DevTools -> Network tab
   # Look for failed API calls to /grpc.*
   ```

## Performance Issues

### Issue: Slow Deployment Processing

**Symptoms**: Deployments take longer than expected

**Solutions**:

1. **Check Resource Usage**
   ```bash
   # Control Plane
   kubectl top pod -n pipecd -l app=pipecd
   
   # Piped
   kubectl top pod -n pipecd -l app=piped
   ```

2. **Optimize Database Queries**
   ```bash
   # Check for slow queries
   # MySQL: Enable slow query log
   # PostgreSQL: pg_stat_statements
   ```

3. **Reduce Git Sync Frequency**
   ```yaml
   spec:
     syncInterval: 5m  # Increase from 1m if needed
   ```

### Issue: High Memory Usage

**Symptoms**: OOMKilled errors, frequent restarts

**Solutions**:

```bash
# Increase memory limits
kubectl patch deployment -n pipecd pipecd -p '{"spec":{"template":{"spec":{"containers":[{"name":"pipecd","resources":{"limits":{"memory":"4Gi"}}}]}}}}'

# Enable garbage collection tuning
# Set GOGC environment variable
kubectl set env deployment/pipecd -n pipecd GOGC=80
```

## Network & Connectivity

### Issue: Cannot Reach External Services

**Symptoms**: Timeout connecting to webhooks, metrics providers

**Solutions**:

1. **Check Network Policies**
   ```bash
   kubectl get networkpolicies -n pipecd
   kubectl describe networkpolicy -n pipecd
   ```

2. **Verify DNS Resolution**
   ```bash
   kubectl exec -n pipecd deployment/piped -- nslookup google.com
   ```

3. **Test Connectivity**
   ```bash
   kubectl exec -n pipecd deployment/piped -- curl -v https://api.github.com
   ```

### Issue: Webhook Notifications Not Working

**Symptoms**: No Slack/Discord messages sent

**Solutions**:

1. **Verify Webhook URL**
   ```yaml
   # Check notification config
   spec:
     notifications:
       - name: slack
         type: SLACK
         config:
           webhookUrl: https://hooks.slack.com/services/...
   ```

2. **Test Webhook Manually**
   ```bash
   curl -X POST -H 'Content-type: application/json' \
     --data '{"text":"Test message"}' \
     <webhook-url>
   ```

3. **Check Firewall Rules**
   - Ensure outbound HTTPS allowed
   - Verify no proxy blocking requests

## Database Issues

### Issue: Database Migration Fails

**Symptoms**: Control plane won't start, migration errors

**Solutions**:

```bash
# Check migration status
kubectl logs -n pipecd deployment/pipecd | grep -i "migration"

# Manually run migrations (if needed)
kubectl exec -n pipecd deployment/pipecd -- /pipecd migrate

# Rollback migration (last resort)
# Backup database first!
```

### Issue: Database Connection Pool Exhausted

**Symptoms**: "too many connections", slow queries

**Solutions**:

```yaml
# Adjust connection pool settings
spec:
  datastore:
    maxOpenConns: 50
    maxIdleConns: 10
    connMaxLifetime: 300s
```

## Getting Help

If you've tried the above solutions and still have issues:

### 1. Gather Diagnostic Information

```bash
# Collect logs
kubectl logs -n pipecd deployment/pipecd > pipecd.log
kubectl logs -n pipecd deployment/piped > piped.log

# Get pod status
kubectl describe pod -n pipecd > pod-status.txt

# Export configurations (remove sensitive data!)
kubectl get configmap -n pipecd -o yaml > configs.yaml
```

### 2. Search Existing Issues

- Check [GitHub Issues](https://github.com/pipe-cd/pipecd/issues)
- Search [Discussions](https://github.com/pipe-cd/pipecd/discussions)

### 3. Ask the Community

- Join [CNCF Slack #pipecd](https://cloud-native.slack.com/archives/C01B27F9T0X)
- Attend [community meetings](https://bit.ly/pipecd-mtg-notes)
- Post in [GitHub Discussions](https://github.com/pipe-cd/pipecd/discussions)

### 4. Report a Bug

If you believe you've found a bug:

1. Create a [bug report](https://github.com/pipe-cd/pipecd/issues/new?template=bug-report.md)
2. Include:
   - PipeCD version
   - Component affected (control plane, piped, web UI)
   - Steps to reproduce
   - Expected vs actual behavior
   - Relevant logs (sanitized)
   - Environment details (Kubernetes version, cloud provider, etc.)

## Debugging Tools

### Useful Commands

```bash
# Port forward to control plane
kubectl port-forward -n pipecd svc/pipecd 9080:9080

# Port forward to piped (admin port)
kubectl port-forward -n pipecd svc/piped 9085:9085

# Access pprof for profiling
curl http://localhost:9085/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Check gRPC health
grpcurl -plaintext localhost:9080 grpc.health.v1.Health/Check

# Inspect database
kubectl exec -it -n pipecd deployment/pipecd -- sh
# Then use mysql or psql client
```

### Enable Profiling

```yaml
# In component config
spec:
  profiling:
    enabled: true
    port: 6060
```

## Common Error Messages

### "context deadline exceeded"
- Network connectivity issue
- Timeout too short
- Service not responding

### "permission denied"
- RBAC permissions insufficient
- Service account misconfigured
- File permissions issue

### "resource not found"
- Wrong namespace
- Resource deleted
- Wrong cluster context

### "invalid configuration"
- YAML syntax error
- Missing required fields
- Incompatible values

### "authentication failed"
- Invalid credentials
- Expired token
- Wrong API key

---

**Note**: This guide is continuously updated. For the latest troubleshooting tips, check the [official documentation](https://pipecd.dev/docs).
