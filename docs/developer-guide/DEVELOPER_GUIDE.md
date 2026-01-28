# HustleX Developer Guide

> Complete Guide for HustleX Developers

---

## Table of Contents

1. [Development Environment Setup](#1-development-environment-setup)
2. [Project Structure](#2-project-structure)
3. [Backend Development](#3-backend-development)
4. [Mobile Development](#4-mobile-development)
5. [Database Operations](#5-database-operations)
6. [Testing](#6-testing)
7. [Code Standards](#7-code-standards)
8. [CI/CD Pipeline](#8-cicd-pipeline)
9. [Debugging & Troubleshooting](#9-debugging--troubleshooting)
10. [Contributing](#10-contributing)

---

## 1. Development Environment Setup

### 1.1 Prerequisites

**Required Software:**

| Software | Version | Purpose |
|----------|---------|---------|
| Go | 1.21+ | Backend development |
| Flutter | 3.16+ | Mobile development |
| Docker | 24.x | Local services |
| PostgreSQL | 16 | Database |
| Redis | 7 | Cache/Queue |
| Git | 2.40+ | Version control |

**Recommended Tools:**
- VS Code with Go and Flutter extensions
- GoLand / Android Studio
- TablePlus / DBeaver (database GUI)
- Postman / Insomnia (API testing)
- Docker Desktop

### 1.2 Clone Repository

```bash
git clone https://github.com/billyronks/hustlex.git
cd hustlex
```

### 1.3 Backend Setup

```bash
# Navigate to backend
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your local settings
# Especially: DB credentials, JWT secret

# Start dependencies
docker-compose up -d

# Install Go dependencies
go mod download

# Run the API
go run cmd/api/main.go
```

**Verify Backend:**
```bash
curl http://localhost:8080/api/v1/health
# Should return: {"status":"healthy"}
```

### 1.4 Mobile Setup

```bash
# Navigate to mobile
cd mobile

# Copy environment file
cp .env.example .env

# Get Flutter dependencies
flutter pub get

# Generate code (Freezed, JSON, Riverpod)
flutter pub run build_runner build --delete-conflicting-outputs

# Run on device/simulator
flutter run
```

### 1.5 VS Code Configuration

**Recommended Extensions:**
```json
{
  "recommendations": [
    "golang.go",
    "dart-code.flutter",
    "dart-code.dart-code",
    "ms-azuretools.vscode-docker",
    "esbenp.prettier-vscode",
    "streetsidesoftware.code-spell-checker"
  ]
}
```

**Workspace Settings (.vscode/settings.json):**
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "[dart]": {
    "editor.defaultFormatter": "Dart-Code.dart-code",
    "editor.selectionHighlight": false,
    "editor.suggestSelection": "first",
    "editor.tabCompletion": "onlySnippets",
    "editor.wordBasedSuggestions": "off"
  }
}
```

---

## 2. Project Structure

### 2.1 Repository Layout

```
hustlex/
├── backend/                 # Go API server
│   ├── cmd/
│   │   ├── api/            # Main API entry point
│   │   └── worker/         # Background job worker
│   ├── internal/
│   │   ├── config/         # Configuration loading
│   │   ├── database/       # Database initialization
│   │   ├── handlers/       # HTTP handlers
│   │   ├── jobs/           # Background job definitions
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Database models
│   │   └── services/       # Business logic
│   ├── .env.example
│   ├── Dockerfile
│   └── go.mod
│
├── mobile/                  # Flutter mobile app
│   ├── lib/
│   │   ├── core/           # Shared infrastructure
│   │   │   ├── api/        # HTTP client
│   │   │   ├── config/     # App configuration
│   │   │   ├── di/         # Dependency injection
│   │   │   ├── router/     # Navigation
│   │   │   ├── services/   # Platform services
│   │   │   ├── storage/    # Local persistence
│   │   │   └── widgets/    # Shared widgets
│   │   ├── features/       # Feature modules
│   │   │   ├── auth/
│   │   │   ├── wallet/
│   │   │   ├── gigs/
│   │   │   ├── savings/
│   │   │   ├── credit/
│   │   │   └── profile/
│   │   └── main.dart
│   ├── android/
│   ├── ios/
│   └── pubspec.yaml
│
├── docs/                    # Documentation
├── k8s/                     # Kubernetes manifests
├── scripts/                 # Utility scripts
├── .github/                 # GitHub Actions
└── docker-compose.yml
```

### 2.2 Backend Structure Detail

```
backend/internal/
├── config/
│   └── config.go           # Environment variable loading
│
├── database/
│   └── database.go         # PostgreSQL connection & migrations
│
├── handlers/
│   ├── auth_handler.go     # Authentication endpoints
│   ├── wallet_handler.go   # Wallet operations
│   ├── gig_handler.go      # Gig marketplace
│   ├── savings_handler.go  # Savings circles
│   ├── credit_handler.go   # Credit & loans
│   └── common_handler.go   # Shared utilities
│
├── middleware/
│   ├── auth.go             # JWT validation
│   └── ratelimit.go        # Rate limiting
│
├── models/
│   └── models.go           # GORM model definitions
│
├── services/
│   ├── auth_service.go     # Authentication logic
│   ├── wallet_service.go   # Wallet business logic
│   ├── gig_service.go      # Gig operations
│   ├── savings_service.go  # Savings logic
│   ├── credit_service.go   # Credit scoring
│   └── payment_service.go  # Payment gateway
│
└── jobs/
    └── jobs.go             # Asynq job handlers
```

### 2.3 Mobile Structure Detail

```
mobile/lib/
├── core/
│   ├── api/
│   │   └── api_client.dart       # Dio HTTP client
│   ├── di/
│   │   └── providers.dart        # Riverpod providers
│   ├── router/
│   │   └── app_router.dart       # GoRouter configuration
│   ├── storage/
│   │   ├── local_storage.dart    # Hive database
│   │   └── secure_storage.dart   # Encrypted storage
│   └── widgets/
│       ├── buttons.dart
│       ├── inputs.dart
│       └── cards.dart
│
└── features/
    └── {feature}/
        ├── data/
        │   ├── models/           # Freezed data classes
        │   │   └── {name}_model.dart
        │   ├── repositories/
        │   │   └── {name}_repository.dart
        │   └── services/
        │       └── {name}_service.dart
        └── presentation/
            ├── providers/
            │   └── {name}_provider.dart
            ├── screens/
            │   └── {name}_screen.dart
            └── widgets/
                └── {name}_widget.dart
```

---

## 3. Backend Development

### 3.1 Creating a New Handler

**Step 1: Define Handler**

```go
// internal/handlers/example_handler.go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/billyronks/hustlex/internal/services"
)

type ExampleHandler struct {
    service *services.ExampleService
}

func NewExampleHandler(service *services.ExampleService) *ExampleHandler {
    return &ExampleHandler{service: service}
}

// Request/Response structs
type CreateExampleRequest struct {
    Name        string `json:"name" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"max=500"`
}

type ExampleResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

// Handler method
func (h *ExampleHandler) Create(c *fiber.Ctx) error {
    // Get user from context (set by auth middleware)
    userID := c.Locals("userID").(string)

    // Parse and validate request
    var req CreateExampleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "error":   "Invalid request body",
        })
    }

    // Call service
    result, err := h.service.Create(userID, req.Name, req.Description)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "error":   err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data": ExampleResponse{
            ID:          result.ID.String(),
            Name:        result.Name,
            Description: result.Description,
        },
    })
}
```

**Step 2: Register Routes**

```go
// In main.go or routes setup
exampleHandler := handlers.NewExampleHandler(exampleService)

api := app.Group("/api/v1")
api.Post("/examples", authMiddleware, exampleHandler.Create)
api.Get("/examples/:id", authMiddleware, exampleHandler.Get)
```

### 3.2 Creating a New Service

```go
// internal/services/example_service.go
package services

import (
    "github.com/google/uuid"
    "github.com/billyronks/hustlex/internal/models"
    "gorm.io/gorm"
)

type ExampleService struct {
    db *gorm.DB
}

func NewExampleService(db *gorm.DB) *ExampleService {
    return &ExampleService{db: db}
}

func (s *ExampleService) Create(userID, name, description string) (*models.Example, error) {
    example := &models.Example{
        ID:          uuid.New(),
        UserID:      uuid.MustParse(userID),
        Name:        name,
        Description: description,
    }

    if err := s.db.Create(example).Error; err != nil {
        return nil, err
    }

    return example, nil
}

func (s *ExampleService) GetByID(id string) (*models.Example, error) {
    var example models.Example
    if err := s.db.Where("id = ?", id).First(&example).Error; err != nil {
        return nil, err
    }
    return &example, nil
}
```

### 3.3 Creating a New Model

```go
// internal/models/models.go
type Example struct {
    ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    UserID      uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
    Name        string         `gorm:"size:100;not null" json:"name"`
    Description string         `gorm:"size:500" json:"description"`
    Status      string         `gorm:"size:20;default:'active'" json:"status"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

    // Relations
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

### 3.4 Creating Background Jobs

```go
// internal/jobs/jobs.go

// Define task type constant
const TaskExampleProcess = "example:process"

// Define payload
type ExampleProcessPayload struct {
    ExampleID string `json:"example_id"`
    Action    string `json:"action"`
}

// Register handler in worker
func (p *JobProcessor) registerHandlers() {
    p.mux.HandleFunc(TaskExampleProcess, p.HandleExampleProcess)
}

// Implement handler
func (p *JobProcessor) HandleExampleProcess(ctx context.Context, t *asynq.Task) error {
    var payload ExampleProcessPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }

    // Process the job
    log.Printf("Processing example: %s, action: %s", payload.ExampleID, payload.Action)

    // Do work...

    return nil
}

// Enqueue job from service
func (s *ExampleService) QueueProcess(exampleID string) error {
    payload, _ := json.Marshal(ExampleProcessPayload{
        ExampleID: exampleID,
        Action:    "process",
    })

    task := asynq.NewTask(TaskExampleProcess, payload,
        asynq.Queue("default"),
        asynq.MaxRetry(3),
        asynq.Timeout(30*time.Second),
    )

    _, err := s.asynqClient.Enqueue(task)
    return err
}
```

---

## 4. Mobile Development

### 4.1 Creating a New Feature

**Step 1: Create Feature Directory Structure**

```bash
mkdir -p lib/features/example/data/{models,repositories,services}
mkdir -p lib/features/example/presentation/{providers,screens,widgets}
```

**Step 2: Create Data Model**

```dart
// lib/features/example/data/models/example_model.dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'example_model.freezed.dart';
part 'example_model.g.dart';

@freezed
class ExampleModel with _$ExampleModel {
  const factory ExampleModel({
    required String id,
    required String name,
    String? description,
    required DateTime createdAt,
  }) = _ExampleModel;

  factory ExampleModel.fromJson(Map<String, dynamic> json) =>
      _$ExampleModelFromJson(json);
}
```

**Step 3: Create Repository**

```dart
// lib/features/example/data/repositories/example_repository.dart
import 'package:hustlex/core/api/api_client.dart';
import '../models/example_model.dart';

class ExampleRepository {
  final ApiClient _apiClient;

  ExampleRepository(this._apiClient);

  Future<List<ExampleModel>> getAll() async {
    final response = await _apiClient.get('/examples');
    final List<dynamic> data = response.data['data']['examples'];
    return data.map((e) => ExampleModel.fromJson(e)).toList();
  }

  Future<ExampleModel> getById(String id) async {
    final response = await _apiClient.get('/examples/$id');
    return ExampleModel.fromJson(response.data['data']['example']);
  }

  Future<ExampleModel> create({
    required String name,
    String? description,
  }) async {
    final response = await _apiClient.post('/examples', data: {
      'name': name,
      'description': description,
    });
    return ExampleModel.fromJson(response.data['data']['example']);
  }
}
```

**Step 4: Create Provider**

```dart
// lib/features/example/presentation/providers/example_provider.dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:hustlex/core/di/providers.dart';
import '../../data/models/example_model.dart';
import '../../data/repositories/example_repository.dart';

part 'example_provider.g.dart';

// Repository provider
@riverpod
ExampleRepository exampleRepository(ExampleRepositoryRef ref) {
  return ExampleRepository(ref.watch(apiClientProvider));
}

// State provider
@riverpod
class ExampleNotifier extends _$ExampleNotifier {
  @override
  FutureOr<List<ExampleModel>> build() async {
    return ref.watch(exampleRepositoryProvider).getAll();
  }

  Future<void> refresh() async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() =>
      ref.read(exampleRepositoryProvider).getAll()
    );
  }

  Future<void> create(String name, String? description) async {
    final example = await ref.read(exampleRepositoryProvider).create(
      name: name,
      description: description,
    );

    state = state.whenData((examples) => [...examples, example]);
  }
}
```

**Step 5: Create Screen**

```dart
// lib/features/example/presentation/screens/example_screen.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/example_provider.dart';

class ExampleScreen extends ConsumerWidget {
  const ExampleScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final examplesAsync = ref.watch(exampleNotifierProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Examples'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.read(exampleNotifierProvider.notifier).refresh(),
          ),
        ],
      ),
      body: examplesAsync.when(
        data: (examples) => ListView.builder(
          itemCount: examples.length,
          itemBuilder: (context, index) {
            final example = examples[index];
            return ListTile(
              title: Text(example.name),
              subtitle: Text(example.description ?? ''),
              onTap: () {
                // Navigate to detail
              },
            );
          },
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('Error: $error'),
              ElevatedButton(
                onPressed: () => ref.invalidate(exampleNotifierProvider),
                child: const Text('Retry'),
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showCreateDialog(context, ref),
        child: const Icon(Icons.add),
      ),
    );
  }

  void _showCreateDialog(BuildContext context, WidgetRef ref) {
    // Show dialog to create new example
  }
}
```

**Step 6: Add Route**

```dart
// lib/core/router/app_router.dart
GoRoute(
  path: '/examples',
  builder: (context, state) => const ExampleScreen(),
),
GoRoute(
  path: '/examples/:id',
  builder: (context, state) {
    final id = state.pathParameters['id']!;
    return ExampleDetailScreen(id: id);
  },
),
```

**Step 7: Generate Code**

```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

### 4.2 State Management Patterns

**Simple State (Single Value):**
```dart
@riverpod
class Counter extends _$Counter {
  @override
  int build() => 0;

  void increment() => state++;
  void decrement() => state--;
}
```

**Async State (API Data):**
```dart
@riverpod
class UserProfile extends _$UserProfile {
  @override
  FutureOr<UserModel> build() async {
    return ref.watch(userRepositoryProvider).getProfile();
  }

  Future<void> update(String name, String? email) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() =>
      ref.read(userRepositoryProvider).updateProfile(name: name, email: email)
    );
  }
}
```

**Dependent State:**
```dart
@riverpod
FutureOr<List<GigModel>> userGigs(UserGigsRef ref) async {
  // This automatically refreshes when auth state changes
  final user = await ref.watch(authNotifierProvider.future);
  if (user == null) return [];

  return ref.watch(gigsRepositoryProvider).getUserGigs(user.id);
}
```

---

## 5. Database Operations

### 5.1 Migrations

**Auto-migration (Development):**
```go
// database.go
db.AutoMigrate(
    &models.User{},
    &models.Wallet{},
    &models.Transaction{},
    // Add new models here
)
```

**Manual Migration (Production):**
```sql
-- migrations/001_create_examples.sql
CREATE TABLE examples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_examples_user_id ON examples(user_id);
CREATE INDEX idx_examples_status ON examples(status);
```

### 5.2 Common Queries

**With Pagination:**
```go
func (s *ExampleService) List(userID string, page, limit int) ([]models.Example, int64, error) {
    var examples []models.Example
    var total int64

    offset := (page - 1) * limit

    // Get total count
    s.db.Model(&models.Example{}).
        Where("user_id = ?", userID).
        Count(&total)

    // Get paginated results
    err := s.db.Where("user_id = ?", userID).
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&examples).Error

    return examples, total, err
}
```

**With Joins:**
```go
func (s *GigService) GetWithProposals(gigID string) (*models.Gig, error) {
    var gig models.Gig
    err := s.db.Preload("Proposals").
        Preload("Proposals.Freelancer").
        Where("id = ?", gigID).
        First(&gig).Error
    return &gig, err
}
```

**With Transactions:**
```go
func (s *WalletService) Transfer(fromID, toID string, amount float64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Debit sender
        if err := tx.Model(&models.Wallet{}).
            Where("user_id = ? AND balance >= ?", fromID, amount).
            Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
            return err
        }

        // Credit receiver
        if err := tx.Model(&models.Wallet{}).
            Where("user_id = ?", toID).
            Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
            return err
        }

        // Create transactions
        // ...

        return nil
    })
}
```

---

## 6. Testing

### 6.1 Backend Testing

**Unit Test Example:**
```go
// internal/services/wallet_service_test.go
package services_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/billyronks/hustlex/internal/services"
)

func TestWalletService_Credit(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    service := services.NewWalletService(db)

    // Create test user and wallet
    user := createTestUser(t, db)

    // Test
    err := service.Credit(user.ID.String(), 1000.00, "test-ref", "Test credit")

    // Assert
    assert.NoError(t, err)

    wallet, _ := service.GetByUserID(user.ID.String())
    assert.Equal(t, 1000.00, wallet.Balance)
}
```

**Run Backend Tests:**
```bash
cd backend
go test ./... -v
go test ./... -cover  # With coverage
```

### 6.2 Mobile Testing

**Unit Test Example:**
```dart
// test/features/wallet/wallet_repository_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  late WalletRepository repository;
  late MockApiClient mockApiClient;

  setUp(() {
    mockApiClient = MockApiClient();
    repository = WalletRepository(mockApiClient);
  });

  test('getWallet returns wallet model', () async {
    // Arrange
    when(() => mockApiClient.get('/wallet')).thenAnswer(
      (_) async => Response(data: {
        'data': {
          'id': '123',
          'balance': 50000.0,
          'currency': 'NGN',
        }
      }),
    );

    // Act
    final wallet = await repository.getWallet();

    // Assert
    expect(wallet.balance, 50000.0);
    verify(() => mockApiClient.get('/wallet')).called(1);
  });
}
```

**Widget Test Example:**
```dart
// test/features/wallet/wallet_screen_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

void main() {
  testWidgets('WalletScreen displays balance', (tester) async {
    // Build widget
    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          walletNotifierProvider.overrideWith(
            () => FakeWalletNotifier(WalletModel(
              id: '123',
              balance: 50000,
              currency: 'NGN',
            )),
          ),
        ],
        child: MaterialApp(home: WalletScreen()),
      ),
    );

    // Wait for async loading
    await tester.pumpAndSettle();

    // Assert
    expect(find.text('₦50,000.00'), findsOneWidget);
  });
}
```

**Run Mobile Tests:**
```bash
cd mobile
flutter test
flutter test --coverage  # With coverage
```

---

## 7. Code Standards

### 7.1 Go Code Style

**Naming:**
- Exported functions: `PascalCase`
- Private functions: `camelCase`
- Constants: `PascalCase` or `SCREAMING_SNAKE_CASE`
- Acronyms: `userID` not `userId`, `httpClient` not `HTTPClient`

**Error Handling:**
```go
// Good
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Bad
if err != nil {
    return err  // No context
}
```

**Comments:**
```go
// CreateUser creates a new user with the given phone number.
// It returns an error if the phone number is already registered.
func (s *AuthService) CreateUser(phone string) (*models.User, error) {
    // ...
}
```

### 7.2 Dart Code Style

**Naming:**
- Classes: `PascalCase`
- Variables/functions: `camelCase`
- Constants: `camelCase` or `SCREAMING_SNAKE_CASE`
- Private: `_prefixWithUnderscore`

**Widget Structure:**
```dart
class MyWidget extends StatelessWidget {
  // 1. Final fields
  final String title;

  // 2. Constructor
  const MyWidget({required this.title, super.key});

  // 3. Build method
  @override
  Widget build(BuildContext context) {
    return Container();
  }

  // 4. Private methods
  void _handleTap() {}
}
```

**Imports Order:**
```dart
// 1. Dart imports
import 'dart:async';

// 2. Flutter imports
import 'package:flutter/material.dart';

// 3. Package imports
import 'package:riverpod/riverpod.dart';

// 4. Local imports
import 'package:hustlex/core/api/api_client.dart';
```

### 7.3 Git Commit Messages

**Format:**
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance

**Examples:**
```
feat(wallet): add bank withdrawal feature

- Implement withdrawal to saved bank accounts
- Add withdrawal fee calculation
- Add confirmation dialog

Closes #123
```

---

## 8. CI/CD Pipeline

### 8.1 GitHub Actions Workflow

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  backend-test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: hustlex_test
        ports:
          - 5432:5432
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        working-directory: ./backend
        run: go test ./... -v -cover

  mobile-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/flutter-action@v2
        with:
          flutter-version: '3.16.0'

      - name: Get dependencies
        working-directory: ./mobile
        run: flutter pub get

      - name: Run tests
        working-directory: ./mobile
        run: flutter test
```

### 8.2 Deployment

**Staging:**
```bash
# Automatic on merge to develop
git checkout develop
git merge feature/my-feature
git push origin develop
# CI/CD deploys to staging
```

**Production:**
```bash
# Manual approval required
git checkout main
git merge develop
git push origin main
# CI/CD deploys to production after approval
```

---

## 9. Debugging & Troubleshooting

### 9.1 Backend Debugging

**Enable Debug Logging:**
```bash
export LOG_LEVEL=debug
go run cmd/api/main.go
```

**Database Query Logging:**
```go
// In development
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

**Redis Monitoring:**
```bash
redis-cli monitor
```

### 9.2 Mobile Debugging

**Flutter DevTools:**
```bash
flutter run --debug
# Press 'r' for hot reload
# Press 'R' for hot restart
# Press 'v' for DevTools
```

**Network Debugging:**
```dart
// Add interceptor for logging
dio.interceptors.add(LogInterceptor(
  requestBody: true,
  responseBody: true,
));
```

### 9.3 Common Issues

**Issue: "Connection refused" to database**
- Check if Docker containers are running: `docker-compose ps`
- Verify database credentials in `.env`
- Try: `docker-compose down && docker-compose up -d`

**Issue: "Token is expired"**
- JWT tokens expire after 15 minutes
- Implement token refresh in mobile app
- Check system time sync

**Issue: Build runner fails**
- Delete generated files: `find . -name "*.g.dart" -delete`
- Run clean build: `flutter clean && flutter pub get`
- Regenerate: `flutter pub run build_runner build --delete-conflicting-outputs`

---

## 10. Contributing

### 10.1 Development Workflow

1. Create feature branch from `develop`
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/my-feature
   ```

2. Make changes and commit
   ```bash
   git add .
   git commit -m "feat(scope): description"
   ```

3. Push and create PR
   ```bash
   git push origin feature/my-feature
   # Create PR on GitHub
   ```

4. Address review feedback

5. Merge after approval

### 10.2 Code Review Checklist

- [ ] Code follows style guidelines
- [ ] Tests added for new features
- [ ] Documentation updated
- [ ] No sensitive data exposed
- [ ] Error handling is proper
- [ ] No console.log/print statements
- [ ] Migrations are reversible

---

*Happy coding! For questions, reach out to the engineering team.*

**Version 1.0 | Last Updated: January 2024**
