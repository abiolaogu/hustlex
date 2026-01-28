# ADR-010: Kubernetes for Container Orchestration

## Status

Accepted

## Date

2024-01-15

## Context

HustleX requires a deployment platform that:
- Scales API instances based on load
- Handles failover automatically
- Supports zero-downtime deployments
- Manages secrets and configuration
- Works across cloud providers (avoid vendor lock-in)
- Supports staging and production environments

## Decision

We chose **Kubernetes (K8s)** with **Kustomize** for configuration management.

### Key Reasons:

1. **Industry Standard**: Kubernetes is the de facto standard for container orchestration.

2. **Auto-Scaling**: Horizontal Pod Autoscaler (HPA) scales based on CPU/memory/custom metrics.

3. **Self-Healing**: Automatic restart of failed containers, rescheduling on node failures.

4. **Rolling Updates**: Zero-downtime deployments with configurable rollout strategies.

5. **Cloud Agnostic**: Works on AWS EKS, Google GKE, Azure AKS, and bare metal.

6. **Ecosystem**: Rich ecosystem of tools (Prometheus, Grafana, Istio, ArgoCD).

## Consequences

### Positive

- **Scalability**: Handle traffic spikes automatically
- **Reliability**: Self-healing, multi-replica deployments
- **Portability**: Same manifests work across cloud providers
- **Observability**: Native support for metrics, logging, tracing
- **Security**: RBAC, network policies, secret management
- **GitOps Ready**: Declarative manifests enable GitOps workflows

### Negative

- **Complexity**: Steep learning curve for operations team
- **Resource Overhead**: Control plane requires resources (~3 nodes minimum)
- **Cost**: Managed Kubernetes (EKS, GKE) adds cost
- **Debugging**: Distributed system debugging is harder

### Neutral

- Requires container registry for images
- Need for monitoring/alerting setup
- Ingress controller selection required

## Kubernetes Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                       │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                    hustlex namespace                 │    │
│  │  ┌─────────────────────────────────────────────┐    │    │
│  │  │              Ingress Controller              │    │    │
│  │  │         (nginx-ingress / traefik)           │    │    │
│  │  └──────────────────┬──────────────────────────┘    │    │
│  │                     │                                │    │
│  │     ┌───────────────┼───────────────┐               │    │
│  │     ▼               ▼               ▼               │    │
│  │  ┌──────┐       ┌──────┐       ┌──────┐            │    │
│  │  │ API  │       │ API  │       │ API  │            │    │
│  │  │ Pod  │       │ Pod  │       │ Pod  │            │    │
│  │  │  1   │       │  2   │       │  3   │            │    │
│  │  └──┬───┘       └──┬───┘       └──┬───┘            │    │
│  │     │              │              │                 │    │
│  │     └──────────────┼──────────────┘                 │    │
│  │                    │                                │    │
│  │     ┌──────────────┴──────────────┐                │    │
│  │     ▼                             ▼                │    │
│  │  ┌────────────┐            ┌────────────┐         │    │
│  │  │ PostgreSQL │            │   Redis    │         │    │
│  │  │StatefulSet │            │ Deployment │         │    │
│  │  └────────────┘            └────────────┘         │    │
│  │                                                    │    │
│  │  ┌─────────────────────────────────────────────┐  │    │
│  │  │              Worker Deployment               │  │    │
│  │  │         (Background job processor)          │  │    │
│  │  └─────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Manifest Structure

```
k8s/
├── base/
│   ├── kustomization.yaml    # Base resource list
│   ├── namespace.yaml        # Namespace definition
│   ├── configmap.yaml        # Non-sensitive config
│   ├── secrets.yaml          # Sensitive config (base64)
│   ├── deployment.yaml       # API deployment
│   ├── service.yaml          # LoadBalancer service
│   ├── postgres.yaml         # PostgreSQL StatefulSet
│   ├── redis.yaml            # Redis deployment
│   └── worker.yaml           # Background worker
│
└── overlays/
    ├── staging/
    │   ├── kustomization.yaml
    │   └── patches/
    │       ├── replicas.yaml    # 2 replicas
    │       └── resources.yaml   # Lower resources
    │
    └── production/
        ├── kustomization.yaml
        └── patches/
            ├── replicas.yaml    # 5 replicas
            ├── resources.yaml   # Higher resources
            └── hpa.yaml         # Auto-scaling
```

## Key Manifests

### Deployment

```yaml
# k8s/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hustlex-api
  namespace: hustlex
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hustlex-api
  template:
    metadata:
      labels:
        app: hustlex-api
    spec:
      containers:
      - name: api
        image: hustlex/api:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: hustlex-config
        - secretRef:
            name: hustlex-secrets
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Horizontal Pod Autoscaler

```yaml
# k8s/overlays/production/patches/hpa.yaml
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
```

### PostgreSQL StatefulSet

```yaml
# k8s/base/postgres.yaml
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
        envFrom:
        - secretRef:
            name: postgres-secrets
        volumeMounts:
        - name: postgres-data
          mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
  - metadata:
      name: postgres-data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 50Gi
```

## Deployment Strategy

### Rolling Update (Default)

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # 1 extra pod during update
      maxUnavailable: 0  # Zero downtime
```

### Blue-Green (For Major Releases)

```bash
# Deploy new version as separate deployment
kubectl apply -f deployment-v2.yaml

# Test new version
curl https://v2.api.hustlex.app/health

# Switch traffic
kubectl patch service hustlex-api -p '{"spec":{"selector":{"version":"v2"}}}'

# Cleanup old version
kubectl delete deployment hustlex-api-v1
```

## Alternatives Considered

### Alternative 1: Docker Swarm

**Pros**: Simpler setup, built into Docker, good for small teams
**Cons**: Limited ecosystem, fewer features, declining community

**Rejected because**: Kubernetes ecosystem is more mature for production workloads.

### Alternative 2: AWS ECS/Fargate

**Pros**: Managed service, AWS integration, simpler than K8s
**Cons**: AWS lock-in, less portable, limited customization

**Rejected because**: Prefer cloud-agnostic solution for flexibility.

### Alternative 3: Nomad

**Pros**: Simpler than K8s, good for mixed workloads
**Cons**: Smaller ecosystem, less tooling, smaller community

**Rejected because**: Kubernetes has better tooling and community support.

### Alternative 4: Bare Metal / VM Deployment

**Pros**: Full control, no abstraction overhead
**Cons**: Manual scaling, no self-healing, operational burden

**Rejected because**: Doesn't meet scalability and reliability requirements.

## References

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Kustomize Documentation](https://kustomize.io/)
- [Kubernetes Best Practices](https://cloud.google.com/blog/products/containers-kubernetes/your-guide-kubernetes-best-practices)
- [HPA Configuration](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
