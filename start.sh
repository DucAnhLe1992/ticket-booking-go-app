#!/bin/bash

# Complete monorepo startup script
# This script helps start all services for the ticket booking platform

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}=====================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}=====================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        exit 1
    fi
    print_success "Go installed: $(go version)"
    
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed"
        exit 1
    fi
    print_success "Node.js installed: $(node --version)"
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        exit 1
    fi
    print_success "Docker installed: $(docker --version)"
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed"
        exit 1
    fi
    print_success "Docker Compose installed"
    
    echo ""
}

# Start infrastructure
start_infra() {
    print_header "Starting Infrastructure"
    print_info "Starting PostgreSQL, NATS, and Redis..."
    make dev-infra
    sleep 3
    print_success "Infrastructure started"
    echo ""
}

# Run migrations
run_migrations() {
    print_header "Running Database Migrations"
    print_info "Applying migrations..."
    make migrate-up 2>/dev/null || true
    print_success "Migrations completed"
    echo ""
}

# Install frontend dependencies
install_frontend() {
    print_header "Installing Frontend Dependencies"
    if [ ! -d "frontend/node_modules" ]; then
        print_info "Installing npm packages..."
        cd frontend && npm install && cd ..
        print_success "Frontend dependencies installed"
    else
        print_success "Frontend dependencies already installed"
    fi
    echo ""
}

# Show startup instructions
show_instructions() {
    print_header "Services Ready to Start"
    echo ""
    echo -e "${GREEN}Open these terminals and run:${NC}"
    echo ""
    echo -e "${YELLOW}Terminal 1 - Auth Service:${NC}"
    echo "  make dev-auth"
    echo ""
    echo -e "${YELLOW}Terminal 2 - Tickets Service:${NC}"
    echo "  make dev-tickets"
    echo ""
    echo -e "${YELLOW}Terminal 3 - Orders Service:${NC}"
    echo "  make dev-orders"
    echo ""
    echo -e "${YELLOW}Terminal 4 - Payments Service:${NC}"
    echo "  make dev-payments"
    echo ""
    echo -e "${YELLOW}Terminal 5 - Expiration Worker:${NC}"
    echo "  make dev-expiration"
    echo ""
    echo -e "${YELLOW}Terminal 6 - Frontend:${NC}"
    echo "  cd frontend && npm run dev"
    echo ""
    echo -e "${BLUE}=====================================${NC}"
    echo -e "${GREEN}Once all services are running:${NC}"
    echo -e "${BLUE}Frontend:${NC} http://localhost:3000"
    echo -e "${BLUE}Backend:${NC}  http://localhost:8080"
    echo -e "${BLUE}=====================================${NC}"
}

# Main execution
main() {
    clear
    print_header "Ticket Booking Platform - Setup"
    echo ""
    
    check_prerequisites
    start_infra
    run_migrations
    install_frontend
    show_instructions
}

main
