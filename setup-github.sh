#!/bin/bash
# =====================================================
# HustleX - GitHub Setup Script
# =====================================================
# Run this script after extracting the project zip
# Usage: chmod +x setup-github.sh && ./setup-github.sh
# =====================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "╔═══════════════════════════════════════════════════════╗"
echo "║           HustleX - GitHub Setup Script               ║"
echo "╚═══════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: git is not installed. Please install git first.${NC}"
    exit 1
fi

# Initialize git repository
echo -e "${YELLOW}Initializing Git repository...${NC}"
git init

# Add all files
echo -e "${YELLOW}Adding files to staging...${NC}"
git add .

# Create initial commit
echo -e "${YELLOW}Creating initial commit...${NC}"
git commit -m "🚀 Initial commit: HustleX Nigerian super app

Complete project including:

📱 Mobile App (Flutter)
  - Wallet management with deposits, withdrawals, transfers
  - Gig marketplace for freelancers
  - Ajo/Esusu savings circles
  - Alternative credit scoring & micro-loans
  - Full offline support with Hive caching
  - Riverpod state management
  - GoRouter navigation

🖥️ Backend (Go)
  - RESTful API with Chi router
  - PostgreSQL database with GORM
  - Redis caching & job queues
  - Asynq background workers
  - JWT authentication with OTP
  - Paystack payment integration

🏗️ Infrastructure
  - Docker & Docker Compose setup
  - Kubernetes manifests
  - CI/CD ready

Tech Stack: Go 1.21+, Flutter 3.16+, PostgreSQL, Redis, Riverpod, GoRouter"

echo ""
echo -e "${GREEN}✅ Git repository initialized and committed!${NC}"
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}Next Steps:${NC}"
echo ""
echo "Option A - Using GitHub CLI (recommended if installed):"
echo -e "${GREEN}  gh repo create hustlex --private --source=. --push${NC}"
echo ""
echo "Option B - Manual setup:"
echo "  1. Go to https://github.com/new"
echo "  2. Create a new repository named 'hustlex'"
echo "  3. Do NOT initialize with README, .gitignore, or license"
echo "  4. Run these commands:"
echo ""
echo -e "${GREEN}  git remote add origin https://github.com/YOUR_USERNAME/hustlex.git${NC}"
echo -e "${GREEN}  git branch -M main${NC}"
echo -e "${GREEN}  git push -u origin main${NC}"
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}After pushing, run these commands to set up:${NC}"
echo ""
echo "Backend:"
echo "  cd backend"
echo "  cp .env.example .env  # Edit with your values"
echo "  docker-compose up -d"
echo "  go run cmd/api/main.go"
echo ""
echo "Mobile:"
echo "  cd mobile"
echo "  cp .env.example .env  # Edit with your values"
echo "  flutter pub get"
echo "  flutter pub run build_runner build --delete-conflicting-outputs"
echo "  flutter run"
echo ""
echo -e "${GREEN}Happy coding! 🎉${NC}"
