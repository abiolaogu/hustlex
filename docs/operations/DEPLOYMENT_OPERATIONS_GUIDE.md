# HustleX Deployment & Operations Guide

> Complete Guide for Deploying and Operating HustleX

---

## Table of Contents

1. [Deployment Overview](#1-deployment-overview)
2. [Environment Setup](#2-environment-setup)
3. [Deployment Procedures](#3-deployment-procedures)
4. [Monitoring & Alerting](#4-monitoring--alerting)
5. [Scaling](#5-scaling)
6. [Backup & Recovery](#6-backup--recovery)
7. [Maintenance Procedures](#7-maintenance-procedures)
8. [Troubleshooting](#8-troubleshooting)
9. [Runbooks](#9-runbooks)

---

## 1. Deployment Overview

### 1.1 Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Production Environment                    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Kubernetes Cluster                    │    │
│  │                                                          │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐              │    │
│  │  │ API Pod  │  │ API Pod  │  │ API Pod  │   (3 replicas)│    │
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘              │    │
│  │       │             │             │                      │    │
│  │       └─────────────┼─────────────┘                      │    │
│  │                     │                                    │    │
│  │  ┌──────────────────▼──────────────────┐                │    │
│  │  │            Internal Service          │                │    │
│  │  └──────────────────┬──────────────────┘                │    │
│  │                     │                                    │    │
│  │     ┌───────────────┼───────────────┐                   │    │
│  │     ▼               ▼               ▼                   │    │
│  │  ┌──────┐       ┌──────┐       ┌──────┐                │    │
│  │  │Postgre│      │Redis │       │Worker│                │    │
│  │  │  SQL  │      │      │       │Pods  │                │    │
│  │  └──────┘       └──────┘       └──────┘                │    │
│  │                                                          │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Environment Matrix

| Environment | Purpose | URL | Auto-deploy |
|-------------|---------|-----|-------------|
| Development | Local dev | localhost:8080 | N/A |
| Staging | Pre-prod testing | staging-api.hustlex.app | develop branch |
| Production | Live service | api.hustlex.app | main branch (manual) |

### 1.3 Infrastructure Requirements

**Production (Minimum):**
| Component | Specification | Quantity |
|-----------|--------------|----------|
| API Server | 2 vCPU, 4GB RAM | 3 |
| Worker | 2 vCPU, 4GB RAM | 2 |
| PostgreSQL | 4 vCPU, 16GB RAM, 100GB SSD | 1 Primary + 1 Replica |
| Redis | 2 vCPU, 4GB RAM | 1 |
| Load Balancer | Managed | 1 |

---

## 2. Environment Setup

### 2.1 Kubernetes Cluster Setup

**Create cluster (example with eksctl):**
```bash
eksctl create cluster \
  --name hustlex-prod \
  --region af-south-1 \
  --nodegroup-name standard-workers \
  --node-type t3.medium \
  --nodes 3 \
  --nodes-min 3 \
  --nodes-max 10
```

### 2.2 Namespace Setup

```bash
# Create namespace
kubectl create namespace hustlex

# Set as default
kubectl config set-context --current --namespace=hustlex
```

### 2.3 Secrets Configuration

```bash
# Create secrets from .env file
kubectl create secret generic hustlex-secrets \
  --from-literal=JWT_SECRET='your-secret' \
  --from-literal=DB_PASSWORD='your-db-password' \
  --from-literal=REDIS_PASSWORD='your-redis-password' \
  --from-literal=PAYSTACK_SECRET_KEY='sk_live_xxx' \
  --from-literal=TERMII_API_KEY='your-termii-key'
```

### 2.4 ConfigMap Setup

```yaml
# k8s/base/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: hustlex-config
  namespace: hustlex
data:
  ENV: "production"
  PORT: "8080"
  DB_HOST: "postgres.hustlex.svc.cluster.local"
  DB_PORT: "5432"
  DB_NAME: "hustlex"
  REDIS_HOST: "redis.hustlex.svc.cluster.local"
  REDIS_PORT: "6379"
  LOG_LEVEL: "info"
```

### 2.5 Database Setup

**PostgreSQL StatefulSet:**
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: hustlex
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: hustlex
        - name: POSTGRES_USER
          value: hustlex
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: hustlex-secrets
              key: DB_PASSWORD
        volumeMounts:
        - name: postgres-data
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "4Gi"
            cpu: "1"
          limits:
            memory: "8Gi"
            cpu: "2"
  volumeClaimTemplates:
  - metadata:
      name: postgres-data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 100Gi
```

---

## 3. Deployment Procedures

### 3.1 CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build Docker image
        run: |
          docker build -t hustlex/api:${{ github.sha }} ./backend
          docker push hustlex/api:${{ github.sha }}

  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    environment: staging
    steps:
      - name: Deploy to staging
        run: |
          kubectl set image deployment/hustlex-api \
            api=hustlex/api:${{ github.sha }} \
            --namespace=hustlex-staging

  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production
    steps:
      - name: Deploy to production
        run: |
          kubectl set image deployment/hustlex-api \
            api=hustlex/api:${{ github.sha }} \
            --namespace=hustlex
```

### 3.2 Manual Deployment

**Step 1: Build and push image**
```bash
# Build
cd backend
docker build -t hustlex/api:v1.2.3 .

# Push to registry
docker push hustlex/api:v1.2.3
```

**Step 2: Update deployment**
```bash
# Update image
kubectl set image deployment/hustlex-api \
  api=hustlex/api:v1.2.3 \
  --namespace=hustlex

# Watch rollout
kubectl rollout status deployment/hustlex-api -n hustlex
```

**Step 3: Verify deployment**
```bash
# Check pods
kubectl get pods -n hustlex

# Check logs
kubectl logs -f deployment/hustlex-api -n hustlex

# Test health endpoint
curl https://api.hustlex.app/api/v1/health
```

### 3.3 Rollback Procedures

**Automatic rollback on failure:**
```yaml
spec:
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
    type: RollingUpdate
  minReadySeconds: 30
  progressDeadlineSeconds: 600
```

**Manual rollback:**
```bash
# View rollout history
kubectl rollout history deployment/hustlex-api -n hustlex

# Rollback to previous
kubectl rollout undo deployment/hustlex-api -n hustlex

# Rollback to specific revision
kubectl rollout undo deployment/hustlex-api --to-revision=2 -n hustlex
```

### 3.4 Database Migrations

**Run migrations:**
```bash
# Connect to a pod
kubectl exec -it deployment/hustlex-api -n hustlex -- /bin/sh

# Run migration
./api migrate up
```

**Rollback migration:**
```bash
./api migrate down 1  # Rollback last migration
```

### 3.5 Deployment Checklist

**Pre-deployment:**
- [ ] All tests passing
- [ ] Code reviewed and approved
- [ ] Database migrations tested
- [ ] Secrets updated (if needed)
- [ ] Notify team of deployment

**During deployment:**
- [ ] Monitor rollout status
- [ ] Check error rates
- [ ] Verify health checks
- [ ] Test critical paths

**Post-deployment:**
- [ ] Verify all pods running
- [ ] Check logs for errors
- [ ] Test user-facing features
- [ ] Update deployment log

---

## 4. Monitoring & Alerting

### 4.1 Monitoring Stack

```
┌─────────────────────────────────────────────────────────────────┐
│                     Monitoring Architecture                      │
│                                                                  │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐                  │
│  │   API    │───>│Prometheus│───>│ Grafana  │                  │
│  │  /metrics│    │          │    │Dashboard │                  │
│  └──────────┘    └──────────┘    └──────────┘                  │
│                                                                  │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐                  │
│  │   Logs   │───>│  Loki    │───>│ Grafana  │                  │
│  │          │    │          │    │  Logs    │                  │
│  └──────────┘    └──────────┘    └──────────┘                  │
│                                                                  │
│  ┌──────────┐    ┌──────────┐                                   │
│  │  Alerts  │───>│PagerDuty/│                                   │
│  │          │    │  Slack   │                                   │
│  └──────────┘    └──────────┘                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Key Metrics

**Application Metrics:**
| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `http_requests_total` | Request count | N/A |
| `http_request_duration_seconds` | Latency | p95 > 500ms |
| `http_requests_failed_total` | Error count | > 1% of requests |
| `active_connections` | Current connections | > 1000 |

**Business Metrics:**
| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `transactions_total` | Transaction count | Drops > 50% |
| `transaction_volume_ngn` | Transaction value | Drops > 50% |
| `registrations_total` | New signups | Drops > 80% |
| `active_users` | DAU | Drops > 30% |

**Infrastructure Metrics:**
| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `cpu_usage` | CPU utilization | > 80% |
| `memory_usage` | Memory utilization | > 85% |
| `disk_usage` | Disk utilization | > 80% |
| `postgres_connections` | DB connections | > 80% of max |

### 4.3 Alerting Rules

```yaml
# prometheus/alerts.yml
groups:
- name: hustlex
  rules:
  - alert: HighErrorRate
    expr: |
      sum(rate(http_requests_failed_total[5m])) /
      sum(rate(http_requests_total[5m])) > 0.01
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value | humanizePercentage }}"

  - alert: HighLatency
    expr: |
      histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 0.5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High latency detected"
      description: "p95 latency is {{ $value }}s"

  - alert: DatabaseDown
    expr: pg_up == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "PostgreSQL is down"

  - alert: RedisDown
    expr: redis_up == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Redis is down"
```

### 4.4 Dashboards

**API Overview Dashboard:**
- Request rate (RPS)
- Error rate
- Latency percentiles (p50, p95, p99)
- Active connections
- Top endpoints by request count
- Top endpoints by error rate

**Business Dashboard:**
- Transaction volume (hourly, daily)
- Registration count
- Active users
- Gigs created/completed
- Savings contributions
- Loan disbursements

---

## 5. Scaling

### 5.1 Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: hustlex-api-hpa
  namespace: hustlex
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: hustlex-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300  # 5 minutes
    scaleUp:
      stabilizationWindowSeconds: 60   # 1 minute
```

### 5.2 Manual Scaling

```bash
# Scale API pods
kubectl scale deployment hustlex-api --replicas=5 -n hustlex

# Scale worker pods
kubectl scale deployment hustlex-worker --replicas=3 -n hustlex
```

### 5.3 Database Scaling

**Read Replicas:**
```yaml
# Add read replica
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres-replica
spec:
  replicas: 2
  # ... replica configuration
```

**Connection Pooling (PgBouncer):**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pgbouncer
spec:
  replicas: 2
  template:
    spec:
      containers:
      - name: pgbouncer
        image: edoburu/pgbouncer:1.21.0
        env:
        - name: DATABASE_URL
          value: "postgres://..."
        - name: POOL_MODE
          value: "transaction"
        - name: MAX_CLIENT_CONN
          value: "1000"
        - name: DEFAULT_POOL_SIZE
          value: "20"
```

---

## 6. Backup & Recovery

### 6.1 Backup Strategy

| Data | Method | Frequency | Retention |
|------|--------|-----------|-----------|
| PostgreSQL | pg_dump | Hourly | 30 days |
| PostgreSQL | WAL archiving | Continuous | 7 days |
| Redis | RDB snapshot | Every 15 min | 7 days |
| Configs | Git | On change | Indefinite |

### 6.2 Database Backup

**Automated backup script:**
```bash
#!/bin/bash
# backup-postgres.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="hustlex_backup_${DATE}.sql.gz"
S3_BUCKET="hustlex-backups"

# Create backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | gzip > /tmp/$BACKUP_FILE

# Upload to S3
aws s3 cp /tmp/$BACKUP_FILE s3://$S3_BUCKET/postgres/$BACKUP_FILE

# Cleanup
rm /tmp/$BACKUP_FILE

# Keep only last 30 days
aws s3 ls s3://$S3_BUCKET/postgres/ | \
  awk '{print $4}' | \
  sort -r | \
  tail -n +721 | \  # Keep 30 days * 24 hours
  xargs -I {} aws s3 rm s3://$S3_BUCKET/postgres/{}
```

**CronJob for backups:**
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
  namespace: hustlex
spec:
  schedule: "0 * * * *"  # Every hour
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: hustlex/backup:latest
            command: ["/backup-postgres.sh"]
            envFrom:
            - secretRef:
                name: hustlex-secrets
          restartPolicy: OnFailure
```

### 6.3 Recovery Procedures

**Database Recovery:**
```bash
# Download backup
aws s3 cp s3://hustlex-backups/postgres/hustlex_backup_20240115.sql.gz /tmp/

# Decompress
gunzip /tmp/hustlex_backup_20240115.sql.gz

# Restore
psql -h $DB_HOST -U $DB_USER -d hustlex < /tmp/hustlex_backup_20240115.sql
```

**Point-in-Time Recovery:**
```bash
# Stop PostgreSQL
kubectl scale statefulset postgres --replicas=0 -n hustlex

# Restore base backup + replay WAL
pg_restore --target-time="2024-01-15 10:30:00" ...

# Restart PostgreSQL
kubectl scale statefulset postgres --replicas=1 -n hustlex
```

### 6.4 Disaster Recovery

**RTO/RPO Targets:**
| Scenario | RTO | RPO |
|----------|-----|-----|
| Pod failure | < 1 min | 0 |
| Node failure | < 5 min | 0 |
| AZ failure | < 15 min | < 1 min |
| Region failure | < 4 hours | < 1 hour |

---

## 7. Maintenance Procedures

### 7.1 Routine Maintenance

**Daily:**
- [ ] Review alerts and resolve issues
- [ ] Check error logs
- [ ] Verify backups completed

**Weekly:**
- [ ] Review performance metrics
- [ ] Check disk usage
- [ ] Review security alerts
- [ ] Apply minor patches

**Monthly:**
- [ ] Database VACUUM ANALYZE
- [ ] Review and rotate logs
- [ ] Security patches
- [ ] Capacity planning review

### 7.2 Database Maintenance

**VACUUM and ANALYZE:**
```sql
-- Run vacuum analyze on all tables
VACUUM (VERBOSE, ANALYZE);

-- Reindex if needed
REINDEX DATABASE hustlex;
```

**Index maintenance:**
```sql
-- Check index usage
SELECT
    indexrelname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Remove unused indexes
DROP INDEX CONCURRENTLY idx_unused;
```

### 7.3 Certificate Renewal

**TLS certificates (Let's Encrypt):**
```bash
# Auto-renewal handled by cert-manager
kubectl get certificates -n hustlex

# Manual renewal if needed
kubectl delete certificate hustlex-tls -n hustlex
kubectl apply -f certificate.yaml
```

### 7.4 Secret Rotation

**JWT Secret Rotation:**
1. Generate new secret
2. Update Kubernetes secret
3. Restart API pods (rolling)
4. Old tokens will expire naturally

```bash
# Generate new secret
NEW_SECRET=$(openssl rand -base64 32)

# Update secret
kubectl patch secret hustlex-secrets -n hustlex \
  -p '{"data":{"JWT_SECRET":"'$(echo -n $NEW_SECRET | base64)'"}}'

# Restart pods
kubectl rollout restart deployment/hustlex-api -n hustlex
```

---

## 8. Troubleshooting

### 8.1 Common Issues

**Pods not starting:**
```bash
# Check pod status
kubectl describe pod <pod-name> -n hustlex

# Check events
kubectl get events -n hustlex --sort-by='.lastTimestamp'

# Common causes:
# - Image pull error: Check registry credentials
# - CrashLoopBackOff: Check logs
# - Pending: Check resource requests
```

**High latency:**
```bash
# Check database connections
kubectl exec -it postgres-0 -n hustlex -- psql -c \
  "SELECT count(*) FROM pg_stat_activity;"

# Check Redis
kubectl exec -it redis-0 -n hustlex -- redis-cli info clients

# Profile slow queries
kubectl exec -it postgres-0 -n hustlex -- psql -c \
  "SELECT * FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
```

**Out of memory:**
```bash
# Check memory usage
kubectl top pods -n hustlex

# Check for memory leaks in logs
kubectl logs deployment/hustlex-api -n hustlex | grep -i "memory"

# Increase limits if needed
kubectl patch deployment hustlex-api -n hustlex -p \
  '{"spec":{"template":{"spec":{"containers":[{"name":"api","resources":{"limits":{"memory":"1Gi"}}}]}}}}'
```

### 8.2 Debug Commands

```bash
# Shell into pod
kubectl exec -it deployment/hustlex-api -n hustlex -- /bin/sh

# Port forward for local debugging
kubectl port-forward service/hustlex-api 8080:80 -n hustlex

# Copy files from pod
kubectl cp hustlex/pod-name:/path/to/file ./local-file

# View pod logs
kubectl logs -f deployment/hustlex-api -n hustlex --tail=100

# View previous container logs (if restarted)
kubectl logs deployment/hustlex-api -n hustlex --previous
```

---

## 9. Runbooks

### 9.1 Runbook: API Not Responding

**Symptoms:** Users cannot access the app, health checks failing

**Steps:**
1. Check pod status
   ```bash
   kubectl get pods -n hustlex
   ```

2. If pods are down, check events
   ```bash
   kubectl describe pods -n hustlex
   ```

3. Check logs for errors
   ```bash
   kubectl logs deployment/hustlex-api -n hustlex --tail=200
   ```

4. Verify database connectivity
   ```bash
   kubectl exec -it deployment/hustlex-api -n hustlex -- \
     nc -zv postgres 5432
   ```

5. Verify Redis connectivity
   ```bash
   kubectl exec -it deployment/hustlex-api -n hustlex -- \
     nc -zv redis 6379
   ```

6. If infrastructure is fine, restart deployment
   ```bash
   kubectl rollout restart deployment/hustlex-api -n hustlex
   ```

7. If still failing, rollback to previous version
   ```bash
   kubectl rollout undo deployment/hustlex-api -n hustlex
   ```

### 9.2 Runbook: Database Connection Issues

**Symptoms:** API errors mentioning database, slow responses

**Steps:**
1. Check PostgreSQL pod
   ```bash
   kubectl get pods -l app=postgres -n hustlex
   ```

2. Check connection count
   ```sql
   SELECT count(*) FROM pg_stat_activity;
   SELECT max_connections FROM pg_settings;
   ```

3. Kill idle connections if needed
   ```sql
   SELECT pg_terminate_backend(pid)
   FROM pg_stat_activity
   WHERE state = 'idle'
   AND query_start < now() - interval '10 minutes';
   ```

4. Restart PgBouncer if using
   ```bash
   kubectl rollout restart deployment/pgbouncer -n hustlex
   ```

### 9.3 Runbook: High Error Rate

**Symptoms:** Alerts firing, users reporting errors

**Steps:**
1. Check error rate in Grafana
2. Identify error types in logs
   ```bash
   kubectl logs deployment/hustlex-api -n hustlex | grep -i "error" | tail -50
   ```

3. Check recent deployments
   ```bash
   kubectl rollout history deployment/hustlex-api -n hustlex
   ```

4. If related to recent deploy, rollback
   ```bash
   kubectl rollout undo deployment/hustlex-api -n hustlex
   ```

5. If external service issue (Paystack, etc.), enable fallback or maintenance mode

### 9.4 Runbook: Disk Space Critical

**Symptoms:** Disk usage alerts, database errors

**Steps:**
1. Check disk usage
   ```bash
   kubectl exec -it postgres-0 -n hustlex -- df -h
   ```

2. Identify large tables
   ```sql
   SELECT relname, pg_size_pretty(pg_total_relation_size(relid))
   FROM pg_stat_user_tables
   ORDER BY pg_total_relation_size(relid) DESC
   LIMIT 10;
   ```

3. Clean up if possible
   ```sql
   VACUUM FULL verbose;  -- Warning: locks tables
   ```

4. Delete old data if applicable
   ```sql
   DELETE FROM logs WHERE created_at < now() - interval '90 days';
   ```

5. Expand volume if needed (PVC resize)

---

## Contact Information

**On-Call Rotation:**
- Primary: engineering@hustlex.app
- Escalation: cto@hustlex.app

**External Support:**
| Service | Contact |
|---------|---------|
| AWS Support | AWS Console |
| Paystack | support@paystack.com |
| Domain/DNS | registrar support |

---

*This guide is for internal use only.*

**Version 1.0 | Last Updated: January 2024**
