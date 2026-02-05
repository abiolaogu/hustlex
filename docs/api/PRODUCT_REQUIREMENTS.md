# HustleX - Product Requirements Document

## Executive Summary

**HustleX** is a mobile-first super app designed for Nigerian millennials and Gen Z (ages 13-44, representing 50.1% of Nigeria's 238 million population). The platform addresses the critical intersection of unemployment, entrepreneurship aspirations, and financial inclusion by combining a gig economy marketplace, social savings (Ajo/Esusu), skill development, and community features into one seamless experience.

### Vision Statement
*"Turn your skills into steady income, build your hustle credit, and save with your squad."*

### Market Opportunity
- **Target Population**: 119+ million Nigerians (Millennials + Gen Z)
- **Problem**: Real youth unemployment/underemployment exceeds 50%; 93% work in informal sector
- **Opportunity**: 60% of Nigerian youth want to start their own businesses; $80B in informal savings (Ajo/Esusu)
- **Market Gap**: No unified platform combining skill monetization, social savings, and credit building

---

## Target Users

### Primary Personas

#### 1. "The Side Hustler" - Adaeze (24, Female)
- **Background**: University graduate, currently underemployed
- **Skills**: Graphic design, social media management
- **Goals**: Find consistent freelance work, build client base, save for laptop upgrade
- **Pain Points**: Hard to find clients, no portfolio platform, income is inconsistent
- **Tech Behavior**: Heavy Instagram/TikTok user, WhatsApp for business communication

#### 2. "The Aspiring Entrepreneur" - Chukwuemeka (28, Male)
- **Background**: Works in retail, side hustle selling phone accessories
- **Skills**: Sales, basic tech support, customer service
- **Goals**: Build capital for his own shop, learn e-commerce skills
- **Pain Points**: No access to business loans, can't save consistently, wants to learn new skills
- **Tech Behavior**: Uses Opay, active on Twitter, learns from YouTube

#### 3. "The Skill Seeker" - Fatima (20, Female)
- **Background**: NYSC corps member, looking to build employable skills
- **Skills**: Currently learning digital marketing
- **Goals**: Land her first paid gig, build experience portfolio, network with professionals
- **Pain Points**: No experience = no jobs, no jobs = no experience
- **Tech Behavior**: Mobile-first, TikTok for learning, participates in WhatsApp communities

#### 4. "The Community Saver" - Tunde (32, Male)
- **Background**: Shop owner, participates in traditional Ajo
- **Skills**: Retail management, customer relations
- **Goals**: Digitize his Ajo group, access small loans for inventory, grow his business
- **Pain Points**: Manual Ajo tracking is error-prone, trust issues with some members
- **Tech Behavior**: Moniepoint user, prefers USSD backup, WhatsApp power user

---

## Core Features

### 1. Skills Marketplace ("Hustle Hub")

#### 1.1 Gig Categories
**Tier 1 - Digital Services**
- Graphic Design (logos, flyers, social media graphics)
- Content Writing (blogs, copywriting, academic)
- Digital Marketing (social media management, ads)
- Video Editing (YouTube, TikTok, promotional)
- Web/App Development
- Virtual Assistance
- Data Entry

**Tier 2 - Physical Services**
- Photography/Videography
- Event Planning/Coordination
- Tutoring (in-person and virtual)
- Beauty Services (makeup, hair)
- Fashion/Tailoring
- Delivery/Errands
- Home Services (cleaning, repairs)

**Tier 3 - Professional Services**
- Accounting/Bookkeeping
- Legal Document Preparation
- Business Consulting
- Translation Services

#### 1.2 Gig Workflow
```
Client Posts Gig → Matched Hustlers Apply → Client Selects → 
Escrow Payment → Work Delivered → Client Approves → Payment Released
```

#### 1.3 Pricing Model
- **Commission**: 10% platform fee (5% from seller, 5% from buyer)
- **Featured Listings**: ₦500-2,000/week for priority visibility
- **Instant Withdrawal**: 2% fee (vs free 24-hour withdrawal)

### 2. Social Savings ("Squad Save")

#### 2.1 Ajo/Esusu Digital Circles
- Create or join savings circles (5-30 members)
- Automated contribution reminders via push/SMS
- Transparent payout schedule
- Member reputation tracking
- Integration with wallet for seamless contributions

#### 2.2 Circle Types
- **Rotational (Classic Ajo)**: Each member receives pool in turns
- **Fixed Target**: Group saves toward shared goal, split at end
- **Emergency Fund**: Pooled fund for member emergencies (voting mechanism)

#### 2.3 Trust Features
- Contribution history visible to all members
- Late payment tracking and warnings
- Member referral chains (accountability)
- Optional collateral locking from wallet

### 3. Hustle Credit System

#### 3.1 Credit Score Components
| Factor | Weight | Description |
|--------|--------|-------------|
| Gig Completion Rate | 25% | % of accepted gigs completed |
| Client Ratings | 20% | Average rating from clients |
| Ajo Contribution Record | 20% | On-time contribution % |
| Account Age & Activity | 15% | Tenure and engagement |
| Skill Verification | 10% | Verified skills/certifications |
| Community Standing | 10% | Reviews, reports, referrals |

#### 3.2 Credit Benefits
- **Bronze (0-300)**: Basic features, 10% fee, no credit access
- **Silver (301-500)**: 8% fee, micro-loans up to ₦50,000
- **Gold (501-700)**: 6% fee, loans up to ₦200,000, featured profile
- **Platinum (701-850)**: 5% fee, loans up to ₦500,000, priority support, badge

### 4. Skill Development ("Level Up")

#### 4.1 Learning Tracks
- Curated micro-courses (15-30 min modules)
- Local expert-created content
- Certificates upon completion
- Direct pathway to marketplace gigs

#### 4.2 Learn-to-Earn
- Complete courses → Earn XP → Unlock badges
- Top learners get marketplace visibility boost
- Skill assessments unlock premium gig tiers

### 5. Wallet & Payments

#### 5.1 Wallet Features
- Instant deposits (card, bank transfer, USSD)
- Hold/escrow for gig payments
- Split payments across savings goals
- Bill payments (airtime, data, utilities)

#### 5.2 Withdrawal Options
- Bank transfer (free 24hr, 2% instant)
- Mobile money (OPay, PalmPay, etc.)
- Cash pickup (agent network - future)

### 6. Community Features

#### 6.1 Tribes
- Industry-specific communities (Designers, Writers, etc.)
- Local communities (Lagos Hustlers, Abuja Creatives)
- Mentorship matching

#### 6.2 Feed
- Success stories and tips
- Job/opportunity sharing
- Skill showcases
- Community challenges

---

## Technical Requirements

### Performance Metrics
| Metric | Target |
|--------|--------|
| App Load Time | < 3 seconds on 3G |
| API Response Time | < 200ms (p95) |
| Uptime | 99.9% |
| Concurrent Users | 100,000+ |
| Transaction TPS | 1,000+ |

### Offline Capability
- Cached user profile and recent data
- Queue actions for sync when online
- USSD fallback for critical transactions

### Security Requirements
- End-to-end encryption for messages
- PCI DSS compliance for payments
- Biometric authentication option
- 2FA for sensitive operations
- Data encryption at rest and transit

### Localization
- English and Pidgin English
- Nigerian phone number verification
- Naira (₦) as primary currency
- Local bank integration (all major banks)
- NIN/BVN verification integration

---

## Monetization Strategy

### Revenue Streams
1. **Transaction Fees**: 10% commission on gig transactions
2. **Premium Subscriptions**: ₦2,000/month for enhanced features
3. **Featured Listings**: ₦500-5,000 for visibility boosts
4. **Micro-Lending Interest**: 3-5% monthly on credit products
5. **Partner Services**: Affiliate commissions on bill payments
6. **Data Insights**: Anonymized market intelligence (B2B)

### Projected Revenue (Year 1)
- Target: 500,000 active users
- Average monthly transactions: ₦50,000/user
- Revenue: ~₦300M/month (transaction fees only)

---

## Success Metrics (KPIs)

### Growth Metrics
- Monthly Active Users (MAU)
- Daily Active Users (DAU)
- New User Registrations
- User Retention (D1, D7, D30)

### Engagement Metrics
- Gigs Posted per Day
- Gigs Completed per Day
- Savings Circles Created
- Average Savings per Circle
- Messages Sent

### Financial Metrics
- Gross Transaction Volume (GTV)
- Revenue per User (ARPU)
- Customer Acquisition Cost (CAC)
- Lifetime Value (LTV)

### Trust Metrics
- Average Hustle Credit Score
- Dispute Rate
- Ajo Default Rate
- NPS Score

---

## Go-to-Market Strategy

### Phase 1: Lagos Launch (Months 1-3)
- Focus on university campuses and tech hubs
- Partner with skills training centers
- Influencer marketing (micro-influencers)
- WhatsApp community building

### Phase 2: Southwest Expansion (Months 4-6)
- Ibadan, Abeokuta, Osogbo
- Agent network for onboarding
- Radio advertising
- Strategic partnerships with cooperatives

### Phase 3: National Rollout (Months 7-12)
- Abuja, Port Harcourt, Kano
- TV advertising
- Government partnership (3MTT alignment)
- Corporate enterprise features

---

## Competitive Analysis

| Feature | HustleX | Fiverr | PiggyVest | Bumpa |
|---------|---------|--------|-----------|-------|
| Gig Marketplace | ✅ | ✅ | ❌ | ❌ |
| Social Savings | ✅ | ❌ | ✅ | ❌ |
| Credit Building | ✅ | ❌ | ❌ | ❌ |
| Skill Learning | ✅ | ✅ | ❌ | ❌ |
| Local Focus | ✅ | ❌ | ✅ | ✅ |
| Offline Support | ✅ | ❌ | ❌ | ❌ |
| Community | ✅ | ❌ | ❌ | ❌ |

### Competitive Advantage
HustleX is the only platform combining gig economy, social savings, and credit building - creating a unique flywheel where earning leads to saving, saving builds credit, and credit enables more earning.

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Regulatory changes | Medium | High | Proactive CBN engagement, compliance team |
| Payment fraud | High | High | ML fraud detection, escrow system, KYC |
| Ajo defaults | Medium | Medium | Credit scoring, collateral, group liability |
| Competition | High | Medium | First-mover advantage, network effects |
| Scalability issues | Medium | High | Cloud-native architecture, load testing |

---

## Timeline

### Q1 2026
- MVP development complete
- Beta launch with 1,000 users
- Core features: Gig marketplace, basic wallet, savings circles

### Q2 2026
- Public launch in Lagos
- Target: 50,000 users
- Add: Hustle Credit, skill learning basics

### Q3 2026
- Southwest expansion
- Target: 200,000 users
- Add: Premium features, micro-lending pilot

### Q4 2026
- National rollout begins
- Target: 500,000 users
- Add: Full credit products, enterprise features

---

## Implementation Status

### Current Progress (Q1 2026)

As of **February 5, 2026**, we are in the MVP development phase. The following has been accomplished:

**Completed:**
- ✅ Repository structure and monorepo setup
- ✅ Backend architecture (Go with Clean Architecture)
- ✅ Domain layer structure for all core modules (gig, wallet, savings, identity, notification, credit, diaspora)
- ✅ Infrastructure setup (Docker, PostgreSQL, Hasura, DragonflyDB, RabbitMQ)
- ✅ Mobile app scaffolding (Flutter consumer-app, provider-app)
- ✅ Admin dashboard scaffolding (React/Refine)

**In Progress:**
- 🚧 MVP feature implementation (see detailed tracker)
- 🚧 Mobile app development (Flutter)
- 🚧 API endpoint implementation
- 🚧 Database schema and migrations

**Next Steps:**
1. **Conduct comprehensive code audit** to map existing implementation to PRD requirements
2. **Implement core gig workflow** (post → apply → select → escrow → deliver → release)
3. **Build MVP mobile screens** (login, marketplace, wallet, savings circles)
4. **Setup testing infrastructure** (unit tests, integration tests, CI/CD)

### Tracking Documentation

For detailed implementation status, task breakdowns, and progress tracking, see:
- **[MVP Status Tracker](../MVP_STATUS_TRACKER.md)** - Comprehensive feature completion matrix
- **[Platform Concepts Summary](../business/PLATFORM_CONCEPTS_SUMMARY.md)** - Strategic platform evolution

### Updated Timeline

**Target Beta Launch:** April 15, 2026 (10 weeks from now)

| Phase | Timeline | Status |
|-------|----------|--------|
| Code Audit & Gap Analysis | Weeks 1-2 | 🔄 Current Phase |
| Core Feature Implementation | Weeks 3-4 | ⏳ Pending |
| Mobile App Development | Weeks 5-6 | ⏳ Pending |
| Integration Testing | Week 7 | ⏳ Pending |
| Security Audit & Bug Fixes | Week 8 | ⏳ Pending |
| Beta Preparation | Week 9 | ⏳ Pending |
| Beta Launch (1,000 users) | Week 10 | ⏳ Target: Apr 15 |

---

*Document Version: 1.1*
*Last Updated: February 5, 2026*
*Next Review: February 12, 2026*
