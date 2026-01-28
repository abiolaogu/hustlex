# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records (ADRs) for the HustleX platform. ADRs capture important architectural decisions made during the development of the platform.

## What is an ADR?

An Architecture Decision Record is a document that captures an important architectural decision made along with its context and consequences.

## ADR Index

| ADR | Title | Status | Date |
|-----|-------|--------|------|
| [ADR-001](001-go-fiber-backend.md) | Go with Fiber Framework for Backend | Accepted | 2024-01 |
| [ADR-002](002-flutter-mobile-framework.md) | Flutter for Cross-Platform Mobile Development | Accepted | 2024-01 |
| [ADR-003](003-postgresql-database.md) | PostgreSQL as Primary Database | Accepted | 2024-01 |
| [ADR-004](004-redis-caching-queues.md) | Redis for Caching and Job Queues | Accepted | 2024-01 |
| [ADR-005](005-riverpod-state-management.md) | Riverpod for State Management | Accepted | 2024-01 |
| [ADR-006](006-asynq-background-jobs.md) | Asynq for Background Job Processing | Accepted | 2024-01 |
| [ADR-007](007-jwt-authentication.md) | JWT-Based Authentication Strategy | Accepted | 2024-01 |
| [ADR-008](008-paystack-payment-integration.md) | Paystack as Primary Payment Gateway | Accepted | 2024-01 |
| [ADR-009](009-clean-architecture-mobile.md) | Clean Architecture for Mobile App | Accepted | 2024-01 |
| [ADR-010](010-kubernetes-deployment.md) | Kubernetes for Container Orchestration | Accepted | 2024-01 |
| [ADR-011](011-uuid-primary-keys.md) | UUID Primary Keys for Database Models | Accepted | 2024-01 |
| [ADR-012](012-otp-authentication.md) | OTP-Based User Authentication | Accepted | 2024-01 |

## ADR Status Definitions

- **Proposed**: The ADR is under discussion
- **Accepted**: The ADR has been accepted and is being implemented
- **Deprecated**: The ADR was accepted but is no longer relevant
- **Superseded**: The ADR was replaced by another ADR

## Template

New ADRs should follow the template in [TEMPLATE.md](TEMPLATE.md).
