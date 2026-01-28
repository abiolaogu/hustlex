#!/usr/bin/env bash
#
# generate-openapi.sh
#
# Generates OpenAPI specification and related documentation for HustleX API.
# This script validates, bundles, and generates client SDKs from the OpenAPI spec.
#
# Usage:
#   ./scripts/generate-openapi.sh [command]
#
# Commands:
#   validate    - Validate the OpenAPI specification
#   bundle      - Bundle multi-file spec into single file
#   generate    - Generate client SDKs (dart, typescript)
#   docs        - Generate HTML documentation
#   all         - Run all commands
#   help        - Show this help message
#
# Requirements:
#   - Node.js 18+
#   - npm packages: @redocly/cli, openapi-generator-cli
#
# Environment:
#   OPENAPI_SPEC_PATH - Path to OpenAPI spec (default: docs/api/openapi.yaml)
#   OUTPUT_DIR        - Output directory (default: docs/api/generated)
#

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
DOCS_DIR="${ROOT_DIR}/../docs/api"
SPEC_PATH="${OPENAPI_SPEC_PATH:-${DOCS_DIR}/openapi.yaml}"
OUTPUT_DIR="${OUTPUT_DIR:-${DOCS_DIR}/generated}"
VERSION=$(grep -E "^  version:" "${SPEC_PATH}" | head -1 | awk '{print $2}' | tr -d '"' || echo "1.0.0")

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."

    local missing_deps=()

    if ! command -v node &> /dev/null; then
        missing_deps+=("node")
    fi

    if ! command -v npx &> /dev/null; then
        missing_deps+=("npx")
    fi

    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing dependencies: ${missing_deps[*]}"
        log_info "Please install Node.js 18+ from https://nodejs.org/"
        exit 1
    fi

    # Check npm packages
    if ! npx @redocly/cli --version &> /dev/null 2>&1; then
        log_warn "@redocly/cli not found, installing..."
        npm install -g @redocly/cli
    fi

    log_success "All dependencies available"
}

# Validate OpenAPI specification
validate_spec() {
    log_info "Validating OpenAPI specification: ${SPEC_PATH}"

    if [ ! -f "${SPEC_PATH}" ]; then
        log_error "OpenAPI spec not found: ${SPEC_PATH}"
        exit 1
    fi

    # Validate with Redocly
    if npx @redocly/cli lint "${SPEC_PATH}" --skip-rule=no-server-example.com; then
        log_success "OpenAPI specification is valid"
    else
        log_error "OpenAPI specification validation failed"
        exit 1
    fi
}

# Bundle multi-file spec into single file
bundle_spec() {
    log_info "Bundling OpenAPI specification..."

    mkdir -p "${OUTPUT_DIR}"

    local bundled_yaml="${OUTPUT_DIR}/openapi-bundled.yaml"
    local bundled_json="${OUTPUT_DIR}/openapi-bundled.json"

    # Bundle to YAML
    npx @redocly/cli bundle "${SPEC_PATH}" -o "${bundled_yaml}"
    log_success "Bundled YAML: ${bundled_yaml}"

    # Convert to JSON
    npx @redocly/cli bundle "${SPEC_PATH}" -o "${bundled_json}"
    log_success "Bundled JSON: ${bundled_json}"

    # Add version info
    echo "# Generated: $(date -u +"%Y-%m-%dT%H:%M:%SZ")" >> "${bundled_yaml}.meta"
    echo "# Version: ${VERSION}" >> "${bundled_yaml}.meta"
}

# Generate client SDKs
generate_sdks() {
    log_info "Generating client SDKs..."

    mkdir -p "${OUTPUT_DIR}/sdks"

    local bundled_spec="${OUTPUT_DIR}/openapi-bundled.yaml"

    if [ ! -f "${bundled_spec}" ]; then
        log_warn "Bundled spec not found, running bundle first..."
        bundle_spec
    fi

    # Generate Dart/Flutter client
    log_info "Generating Dart client..."
    if command -v openapi-generator-cli &> /dev/null || npx @openapitools/openapi-generator-cli version &> /dev/null 2>&1; then
        npx @openapitools/openapi-generator-cli generate \
            -i "${bundled_spec}" \
            -g dart-dio \
            -o "${OUTPUT_DIR}/sdks/dart" \
            --additional-properties=pubName=hustlex_api,pubAuthor=HustleX,pubVersion="${VERSION}" \
            --skip-validate-spec \
            2>/dev/null || log_warn "Dart SDK generation skipped (generator not available)"
    else
        log_warn "OpenAPI Generator not found, skipping Dart SDK generation"
        log_info "Install with: npm install -g @openapitools/openapi-generator-cli"
    fi

    # Generate TypeScript client
    log_info "Generating TypeScript client..."
    if npx @openapitools/openapi-generator-cli version &> /dev/null 2>&1; then
        npx @openapitools/openapi-generator-cli generate \
            -i "${bundled_spec}" \
            -g typescript-fetch \
            -o "${OUTPUT_DIR}/sdks/typescript" \
            --additional-properties=npmName=@hustlex/api-client,npmVersion="${VERSION}",supportsES6=true \
            --skip-validate-spec \
            2>/dev/null || log_warn "TypeScript SDK generation skipped"
    else
        log_warn "Skipping TypeScript SDK generation"
    fi

    log_success "SDK generation complete"
}

# Generate HTML documentation
generate_docs() {
    log_info "Generating HTML documentation..."

    mkdir -p "${OUTPUT_DIR}/html"

    local bundled_spec="${OUTPUT_DIR}/openapi-bundled.yaml"

    if [ ! -f "${bundled_spec}" ]; then
        bundle_spec
    fi

    # Generate Redoc HTML
    npx @redocly/cli build-docs "${bundled_spec}" \
        -o "${OUTPUT_DIR}/html/index.html" \
        --title "HustleX API Documentation" \
        --disableGoogleFont

    log_success "HTML documentation generated: ${OUTPUT_DIR}/html/index.html"

    # Generate Swagger UI bundle
    cat > "${OUTPUT_DIR}/html/swagger.html" << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HustleX API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "./openapi-bundled.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>
EOF

    # Copy bundled spec to HTML directory for Swagger UI
    cp "${OUTPUT_DIR}/openapi-bundled.json" "${OUTPUT_DIR}/html/"

    log_success "Swagger UI generated: ${OUTPUT_DIR}/html/swagger.html"
}

# Generate API diff between versions
diff_specs() {
    local old_spec="${1:-}"
    local new_spec="${SPEC_PATH}"

    if [ -z "${old_spec}" ]; then
        log_error "Usage: $0 diff <old-spec-path>"
        exit 1
    fi

    log_info "Comparing API specifications..."

    # Use oasdiff for comparison
    if ! command -v oasdiff &> /dev/null; then
        log_warn "oasdiff not found, installing..."
        go install github.com/tufin/oasdiff@latest 2>/dev/null || {
            log_error "Failed to install oasdiff. Please install Go first."
            exit 1
        }
    fi

    oasdiff diff "${old_spec}" "${new_spec}" --format text
}

# Print help
show_help() {
    cat << EOF
HustleX OpenAPI Generator

Usage: $0 [command]

Commands:
  validate    Validate the OpenAPI specification
  bundle      Bundle multi-file spec into single file
  generate    Generate client SDKs (dart, typescript)
  docs        Generate HTML documentation
  diff <old>  Compare with older spec version
  all         Run all commands (validate, bundle, generate, docs)
  help        Show this help message

Environment Variables:
  OPENAPI_SPEC_PATH  Path to OpenAPI spec (default: docs/api/openapi.yaml)
  OUTPUT_DIR         Output directory (default: docs/api/generated)

Examples:
  $0 validate                    # Validate spec
  $0 all                         # Run complete pipeline
  $0 diff v0.9.0/openapi.yaml   # Compare with old version

EOF
}

# Main entry point
main() {
    local command="${1:-help}"

    case "${command}" in
        validate)
            check_dependencies
            validate_spec
            ;;
        bundle)
            check_dependencies
            validate_spec
            bundle_spec
            ;;
        generate)
            check_dependencies
            validate_spec
            bundle_spec
            generate_sdks
            ;;
        docs)
            check_dependencies
            validate_spec
            bundle_spec
            generate_docs
            ;;
        diff)
            shift
            diff_specs "$@"
            ;;
        all)
            check_dependencies
            validate_spec
            bundle_spec
            generate_sdks
            generate_docs
            log_success "All tasks completed successfully!"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown command: ${command}"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
