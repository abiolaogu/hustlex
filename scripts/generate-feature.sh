#!/bin/bash

# HustleX Pro Feature Generator
# Usage: ./scripts/generate-feature.sh <platform> <feature-name>
# Example: ./scripts/generate-feature.sh backend wallet
#          ./scripts/generate-feature.sh flutter remittance
#          ./scripts/generate-feature.sh android savings

set -e

PLATFORM=$1
FEATURE_NAME=$2
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

function log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

function log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

if [ -z "$PLATFORM" ] || [ -z "$FEATURE_NAME" ]; then
    echo "Usage: $0 <platform> <feature-name>"
    echo ""
    echo "Platforms:"
    echo "  backend  - Go backend feature"
    echo "  flutter  - Flutter consumer app feature"
    echo "  admin    - React admin dashboard page"
    echo "  android  - Android native feature"
    echo "  ios      - iOS native feature"
    echo ""
    echo "Example:"
    echo "  $0 backend notifications"
    echo "  $0 flutter savings"
    exit 1
fi

# Convert feature name to proper case formats
FEATURE_LOWER=$(echo "$FEATURE_NAME" | tr '[:upper:]' '[:lower:]')
FEATURE_UPPER=$(echo "$FEATURE_NAME" | tr '[:lower:]' '[:upper:]')
FEATURE_PASCAL=$(echo "$FEATURE_NAME" | sed -r 's/(^|_)([a-z])/\U\2/g')
FEATURE_CAMEL=$(echo "$FEATURE_PASCAL" | sed 's/^./\L&/')

case $PLATFORM in
    backend)
        log_info "Generating Go backend feature: $FEATURE_NAME"

        DOMAIN_DIR="$PROJECT_ROOT/apps/api/internal/domain/$FEATURE_LOWER"
        HANDLER_DIR="$PROJECT_ROOT/apps/api/internal/handlers"

        mkdir -p "$DOMAIN_DIR/entity" "$DOMAIN_DIR/repository" "$DOMAIN_DIR/service"

        # Entity template
        cat > "$DOMAIN_DIR/entity/${FEATURE_LOWER}.go" << EOF
package entity

import (
	"time"

	"github.com/google/uuid"
)

// ${FEATURE_PASCAL} represents a ${FEATURE_LOWER} entity
type ${FEATURE_PASCAL} struct {
	ID        uuid.UUID \`json:"id" db:"id"\`
	UserID    uuid.UUID \`json:"user_id" db:"user_id"\`
	// TODO: Add fields specific to ${FEATURE_LOWER}
	CreatedAt time.Time \`json:"created_at" db:"created_at"\`
	UpdatedAt time.Time \`json:"updated_at" db:"updated_at"\`
}

// New${FEATURE_PASCAL} creates a new ${FEATURE_LOWER}
func New${FEATURE_PASCAL}(userID uuid.UUID) *${FEATURE_PASCAL} {
	now := time.Now()
	return &${FEATURE_PASCAL}{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
EOF

        # Repository template
        cat > "$DOMAIN_DIR/repository/${FEATURE_LOWER}_repository.go" << EOF
package repository

import (
	"context"

	"github.com/google/uuid"
	"hustlex/internal/domain/${FEATURE_LOWER}/entity"
)

// ${FEATURE_PASCAL}Repository defines the interface for ${FEATURE_LOWER} persistence
type ${FEATURE_PASCAL}Repository interface {
	Create(ctx context.Context, ${FEATURE_CAMEL} *entity.${FEATURE_PASCAL}) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.${FEATURE_PASCAL}, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.${FEATURE_PASCAL}, error)
	Update(ctx context.Context, ${FEATURE_CAMEL} *entity.${FEATURE_PASCAL}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
EOF

        # Service template
        cat > "$DOMAIN_DIR/service/${FEATURE_LOWER}_service.go" << EOF
package service

import (
	"context"

	"github.com/google/uuid"
	"hustlex/internal/domain/${FEATURE_LOWER}/entity"
	"hustlex/internal/domain/${FEATURE_LOWER}/repository"
)

// ${FEATURE_PASCAL}Service provides ${FEATURE_LOWER} business logic
type ${FEATURE_PASCAL}Service struct {
	repo repository.${FEATURE_PASCAL}Repository
}

// New${FEATURE_PASCAL}Service creates a new ${FEATURE_LOWER} service
func New${FEATURE_PASCAL}Service(repo repository.${FEATURE_PASCAL}Repository) *${FEATURE_PASCAL}Service {
	return &${FEATURE_PASCAL}Service{repo: repo}
}

// Create creates a new ${FEATURE_LOWER}
func (s *${FEATURE_PASCAL}Service) Create(ctx context.Context, userID uuid.UUID) (*entity.${FEATURE_PASCAL}, error) {
	${FEATURE_CAMEL} := entity.New${FEATURE_PASCAL}(userID)
	if err := s.repo.Create(ctx, ${FEATURE_CAMEL}); err != nil {
		return nil, err
	}
	return ${FEATURE_CAMEL}, nil
}

// GetByID retrieves a ${FEATURE_LOWER} by ID
func (s *${FEATURE_PASCAL}Service) GetByID(ctx context.Context, id uuid.UUID) (*entity.${FEATURE_PASCAL}, error) {
	return s.repo.GetByID(ctx, id)
}
EOF

        # Test template
        cat > "$DOMAIN_DIR/entity/${FEATURE_LOWER}_test.go" << EOF
package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew${FEATURE_PASCAL}(t *testing.T) {
	userID := uuid.New()
	${FEATURE_CAMEL} := New${FEATURE_PASCAL}(userID)

	assert.NotEqual(t, uuid.Nil, ${FEATURE_CAMEL}.ID)
	assert.Equal(t, userID, ${FEATURE_CAMEL}.UserID)
	assert.False(t, ${FEATURE_CAMEL}.CreatedAt.IsZero())
	assert.False(t, ${FEATURE_CAMEL}.UpdatedAt.IsZero())
}
EOF

        log_info "Created Go backend files for $FEATURE_NAME"
        ;;

    flutter)
        log_info "Generating Flutter feature: $FEATURE_NAME"

        FEATURE_DIR="$PROJECT_ROOT/apps/consumer-app/flutter/lib/features/$FEATURE_LOWER"

        mkdir -p "$FEATURE_DIR/domain/entities" "$FEATURE_DIR/domain/repositories" "$FEATURE_DIR/domain/usecases"
        mkdir -p "$FEATURE_DIR/data/models" "$FEATURE_DIR/data/repositories"
        mkdir -p "$FEATURE_DIR/presentation/providers" "$FEATURE_DIR/presentation/screens" "$FEATURE_DIR/presentation/widgets"

        # Entity template
        cat > "$FEATURE_DIR/domain/entities/${FEATURE_LOWER}.dart" << EOF
import 'package:freezed_annotation/freezed_annotation.dart';

part '${FEATURE_LOWER}.freezed.dart';
part '${FEATURE_LOWER}.g.dart';

@freezed
class ${FEATURE_PASCAL} with _\$${FEATURE_PASCAL} {
  const factory ${FEATURE_PASCAL}({
    required String id,
    required String userId,
    // TODO: Add fields specific to ${FEATURE_LOWER}
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _${FEATURE_PASCAL};

  factory ${FEATURE_PASCAL}.fromJson(Map<String, dynamic> json) =>
      _\$${FEATURE_PASCAL}FromJson(json);
}
EOF

        # Repository interface
        cat > "$FEATURE_DIR/domain/repositories/${FEATURE_LOWER}_repository.dart" << EOF
import 'package:either_dart/either.dart';
import '../entities/${FEATURE_LOWER}.dart';

abstract class ${FEATURE_PASCAL}Repository {
  Future<Either<Exception, ${FEATURE_PASCAL}>> create(String userId);
  Future<Either<Exception, ${FEATURE_PASCAL}>> getById(String id);
  Future<Either<Exception, List<${FEATURE_PASCAL}>>> getByUserId(String userId);
  Future<Either<Exception, ${FEATURE_PASCAL}>> update(${FEATURE_PASCAL} ${FEATURE_CAMEL});
  Future<Either<Exception, void>> delete(String id);
}
EOF

        # Provider
        cat > "$FEATURE_DIR/presentation/providers/${FEATURE_LOWER}_provider.dart" << EOF
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../domain/entities/${FEATURE_LOWER}.dart';
import '../../domain/repositories/${FEATURE_LOWER}_repository.dart';

part '${FEATURE_LOWER}_provider.g.dart';

@riverpod
class ${FEATURE_PASCAL}Notifier extends _\$${FEATURE_PASCAL}Notifier {
  @override
  AsyncValue<List<${FEATURE_PASCAL}>> build() {
    return const AsyncValue.loading();
  }

  Future<void> load(String userId) async {
    state = const AsyncValue.loading();
    // TODO: Implement load logic
  }
}
EOF

        # Screen
        cat > "$FEATURE_DIR/presentation/screens/${FEATURE_LOWER}_screen.dart" << EOF
import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

class ${FEATURE_PASCAL}Screen extends HookConsumerWidget {
  const ${FEATURE_PASCAL}Screen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('${FEATURE_PASCAL}'),
      ),
      body: const Center(
        child: Text('${FEATURE_PASCAL} Screen'),
      ),
    );
  }
}
EOF

        log_info "Created Flutter feature files for $FEATURE_NAME"
        ;;

    admin)
        log_info "Generating Admin page: $FEATURE_NAME"

        PAGE_DIR="$PROJECT_ROOT/apps/admin-web/src/pages/$FEATURE_LOWER"
        mkdir -p "$PAGE_DIR"

        cat > "$PAGE_DIR/index.tsx" << EOF
import { List, useTable, ShowButton, Show, Edit, useForm } from "@refinedev/antd";
import { Table, Space, Card, Descriptions, Form, Input } from "antd";
import { useShow } from "@refinedev/core";

export const ${FEATURE_PASCAL}List: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="ID" />
        {/* TODO: Add columns */}
        <Table.Column
          title="Actions"
          dataIndex="actions"
          render={(_, record: any) => (
            <Space>
              <ShowButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};

export const ${FEATURE_PASCAL}Show: React.FC = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Card title="${FEATURE_PASCAL} Details">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="ID">{record?.id}</Descriptions.Item>
          {/* TODO: Add fields */}
        </Descriptions>
      </Card>
    </Show>
  );
};

export const ${FEATURE_PASCAL}Edit: React.FC = () => {
  const { formProps, saveButtonProps } = useForm();

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        {/* TODO: Add form fields */}
      </Form>
    </Edit>
  );
};
EOF

        log_info "Created Admin page files for $FEATURE_NAME"
        ;;

    android)
        log_info "Generating Android feature: $FEATURE_NAME"

        FEATURE_DIR="$PROJECT_ROOT/apps/android/app/src/main/java/ng/hustlex/features/$FEATURE_LOWER"
        mkdir -p "$FEATURE_DIR/presentation" "$FEATURE_DIR/domain" "$FEATURE_DIR/data"

        cat > "$FEATURE_DIR/presentation/${FEATURE_PASCAL}Screen.kt" << EOF
package ng.hustlex.features.${FEATURE_LOWER}.presentation

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ${FEATURE_PASCAL}Screen(
    onNavigateBack: () -> Unit
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("${FEATURE_PASCAL}") },
                navigationIcon = {
                    IconButton(onClick = onNavigateBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Back")
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(16.dp)
        ) {
            Text("${FEATURE_PASCAL} Screen")
        }
    }
}
EOF

        log_info "Created Android feature files for $FEATURE_NAME"
        ;;

    ios)
        log_info "Generating iOS feature: $FEATURE_NAME"

        FEATURE_DIR="$PROJECT_ROOT/apps/ios/HustleX/Features/${FEATURE_PASCAL}"
        mkdir -p "$FEATURE_DIR"

        cat > "$FEATURE_DIR/${FEATURE_PASCAL}View.swift" << EOF
import SwiftUI

struct ${FEATURE_PASCAL}View: View {
    var body: some View {
        NavigationView {
            VStack {
                Text("${FEATURE_PASCAL} Screen")
            }
            .navigationTitle("${FEATURE_PASCAL}")
        }
    }
}

struct ${FEATURE_PASCAL}View_Previews: PreviewProvider {
    static var previews: some View {
        ${FEATURE_PASCAL}View()
    }
}
EOF

        log_info "Created iOS feature files for $FEATURE_NAME"
        ;;

    *)
        log_error "Unknown platform: $PLATFORM"
        ;;
esac

log_info "Feature generation complete!"
echo ""
echo "Next steps:"
echo "  1. Implement the generated boilerplate"
echo "  2. Add tests"
echo "  3. Run: git add . && git commit -m 'feat($FEATURE_NAME): add $FEATURE_NAME feature'"
