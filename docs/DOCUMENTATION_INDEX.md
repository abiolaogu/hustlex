# HustleX Documentation Index

> Complete Index of All HustleX Documentation

---

## Quick Navigation

| Category | For | Location |
|----------|-----|----------|
| **Getting Started** | Everyone | [README.md](../README.md) |
| **User Manual** | End Users | [docs/user-manuals/](user-manuals/) |
| **Admin Manual** | Administrators | [docs/admin-manual/](admin-manual/) |
| **Developer Guide** | Developers | [docs/developer-guide/](developer-guide/) |
| **API Reference** | Developers | [docs/api/](api/) |
| **Training** | All Users | [docs/training/](training/) |

---

## 1. Project Documentation

### Core Documents

| Document | Description | Location |
|----------|-------------|----------|
| README | Project overview and quick start | [README.md](../README.md) |
| Development Status | Current progress and roadmap | [DEVELOPMENT_STATUS.md](../DEVELOPMENT_STATUS.md) |
| System Architecture | Technical architecture overview | [SYSTEM_ARCHITECTURE.md](SYSTEM_ARCHITECTURE.md) |
| Product Requirements | Feature specifications | [PRODUCT_REQUIREMENTS.md](PRODUCT_REQUIREMENTS.md) |

---

## 2. Architecture Decision Records (ADRs)

All architectural decisions are documented in the [docs/adr/](adr/) directory.

| ADR | Title | Status |
|-----|-------|--------|
| [ADR-001](adr/001-go-fiber-backend.md) | Go with Fiber Framework for Backend | Accepted |
| [ADR-002](adr/002-flutter-mobile-framework.md) | Flutter for Cross-Platform Mobile | Accepted |
| [ADR-003](adr/003-postgresql-database.md) | PostgreSQL as Primary Database | Accepted |
| [ADR-004](adr/004-redis-caching-queues.md) | Redis for Caching and Job Queues | Accepted |
| [ADR-005](adr/005-riverpod-state-management.md) | Riverpod for State Management | Accepted |
| [ADR-006](adr/006-asynq-background-jobs.md) | Asynq for Background Job Processing | Accepted |
| [ADR-007](adr/007-jwt-authentication.md) | JWT-Based Authentication Strategy | Accepted |
| [ADR-008](adr/008-paystack-payment-integration.md) | Paystack as Primary Payment Gateway | Accepted |
| [ADR-009](adr/009-clean-architecture-mobile.md) | Clean Architecture for Mobile App | Accepted |
| [ADR-010](adr/010-kubernetes-deployment.md) | Kubernetes for Container Orchestration | Accepted |
| [ADR-011](adr/011-uuid-primary-keys.md) | UUID Primary Keys for Database | Accepted |
| [ADR-012](adr/012-otp-authentication.md) | OTP-Based User Authentication | Accepted |

---

## 3. Technical Documentation

### Technical Specifications

| Document | Description | Location |
|----------|-------------|----------|
| Technical Specifications | Complete technical specs | [docs/technical/TECHNICAL_SPECIFICATIONS.md](technical/TECHNICAL_SPECIFICATIONS.md) |

### API Documentation

| Document | Description | Location |
|----------|-------------|----------|
| API Reference | Complete API endpoint reference | [docs/api/API_REFERENCE.md](api/API_REFERENCE.md) |

### Security Documentation

| Document | Description | Location |
|----------|-------------|----------|
| Security Documentation | Security policies and procedures | [docs/security/SECURITY_DOCUMENTATION.md](security/SECURITY_DOCUMENTATION.md) |

### Operations Documentation

| Document | Description | Location |
|----------|-------------|----------|
| Deployment & Operations | Deployment and operations guide | [docs/operations/DEPLOYMENT_OPERATIONS_GUIDE.md](operations/DEPLOYMENT_OPERATIONS_GUIDE.md) |

---

## 4. User Documentation

### End User Manuals

| Document | Audience | Location |
|----------|----------|----------|
| End User Manual | App users | [docs/user-manuals/END_USER_MANUAL.md](user-manuals/END_USER_MANUAL.md) |

### Admin Manuals

| Document | Audience | Location |
|----------|----------|----------|
| Admin Manual | Platform administrators | [docs/admin-manual/ADMIN_MANUAL.md](admin-manual/ADMIN_MANUAL.md) |

---

## 5. Developer Documentation

| Document | Description | Location |
|----------|-------------|----------|
| Developer Guide | Complete dev setup and guide | [docs/developer-guide/DEVELOPER_GUIDE.md](developer-guide/DEVELOPER_GUIDE.md) |
| Mobile Setup | Flutter app setup | [mobile/SETUP.md](../mobile/SETUP.md) |
| Mobile README | Mobile app overview | [mobile/README.md](../mobile/README.md) |

---

## 6. Training Materials

| Manual | Target Audience | Location |
|--------|-----------------|----------|
| Gig Workers Training | Freelancers | [docs/training/TRAINING_MANUAL_GIG_WORKERS.md](training/TRAINING_MANUAL_GIG_WORKERS.md) |
| Savings Members Training | Savings circle participants | [docs/training/TRAINING_MANUAL_SAVINGS_MEMBERS.md](training/TRAINING_MANUAL_SAVINGS_MEMBERS.md) |

---

## 7. Feature Documentation

| Feature | Description | Location |
|---------|-------------|----------|
| Gig Marketplace | Freelance marketplace feature | [docs/features/GIG_MARKETPLACE.md](features/GIG_MARKETPLACE.md) |
| Savings Circles | Ajo/Esusu savings feature | [docs/features/SAVINGS_CIRCLES.md](features/SAVINGS_CIRCLES.md) |

---

## 8. Infrastructure

### Kubernetes

| File | Description | Location |
|------|-------------|----------|
| Base Configuration | K8s base manifests | [k8s/base/](../k8s/base/) |
| Production Overlay | Production-specific config | [k8s/overlays/production/](../k8s/overlays/production/) |

### Docker

| File | Description | Location |
|------|-------------|----------|
| Docker Compose | Local development setup | [docker-compose.yml](../docker-compose.yml) |
| Backend Dockerfile | API container | [backend/Dockerfile](../backend/Dockerfile) |

### CI/CD

| File | Description | Location |
|------|-------------|----------|
| CI/CD Pipeline | GitHub Actions workflow | [.github/workflows/ci-cd.yml](../.github/workflows/ci-cd.yml) |

---

## Documentation Standards

### Writing Guidelines

1. **Use clear headings** - H1 for title, H2 for sections, H3 for subsections
2. **Include examples** - Code samples, screenshots, diagrams
3. **Keep it updated** - Update docs with code changes
4. **Version documents** - Include version and last updated date

### File Naming

- Use `UPPERCASE_WITH_UNDERSCORES.md` for major documents
- Use `lowercase-with-dashes.md` for supporting files
- Include document type in name (e.g., `MANUAL`, `GUIDE`, `SPEC`)

### Review Process

1. Technical accuracy review
2. Grammar and clarity review
3. Format consistency check
4. Link verification

---

## Contributing to Documentation

1. Fork the repository
2. Create a branch for your changes
3. Follow the documentation standards
4. Submit a pull request
5. Request review from documentation maintainers

---

## Document Maintenance

| Task | Frequency | Responsible |
|------|-----------|-------------|
| Review accuracy | Monthly | Tech Lead |
| Update screenshots | On UI changes | Product |
| API documentation | On API changes | Backend Dev |
| Security docs | Quarterly | Security |

---

## Quick Links

**For Users:**
- [Getting Started](../README.md#getting-started)
- [User Manual](user-manuals/END_USER_MANUAL.md)
- [FAQs](user-manuals/END_USER_MANUAL.md#8-faqs)

**For Developers:**
- [Development Setup](developer-guide/DEVELOPER_GUIDE.md#1-development-environment-setup)
- [API Reference](api/API_REFERENCE.md)
- [Code Standards](developer-guide/DEVELOPER_GUIDE.md#7-code-standards)

**For Operations:**
- [Deployment Guide](operations/DEPLOYMENT_OPERATIONS_GUIDE.md)
- [Monitoring](operations/DEPLOYMENT_OPERATIONS_GUIDE.md#4-monitoring--alerting)
- [Runbooks](operations/DEPLOYMENT_OPERATIONS_GUIDE.md#9-runbooks)

---

*Documentation Version 1.0 | Last Updated: January 2024*

*© 2024 BillyRonks Global Limited. All rights reserved.*
