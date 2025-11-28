#!/bin/bash
# Development helper script for running services locally

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if docker-compose is running
check_infra() {
    echo -e "${YELLOW}Checking infrastructure services...${NC}"
    
    if ! docker compose ps | grep -q "postgres.*running"; then
        echo -e "${RED}Postgres is not running. Start it with: make dev-infra${NC}"
        return 1
    fi
    
    if ! docker compose ps | grep -q "nats.*running"; then
        echo -e "${RED}NATS is not running. Start it with: make dev-infra${NC}"
        return 1
    fi
    
    if ! docker compose ps | grep -q "redis.*running"; then
        echo -e "${RED}Redis is not running. Start it with: make dev-infra${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✓ All infrastructure services are running${NC}"
}

# Check environment variables
check_env() {
    echo -e "${YELLOW}Checking environment variables...${NC}"
    
    missing=()
    
    [ -z "$DATABASE_URL" ] && missing+=("DATABASE_URL")
    [ -z "$JWT_SECRET" ] && missing+=("JWT_SECRET")
    
    if [ ${#missing[@]} -gt 0 ]; then
        echo -e "${RED}Missing environment variables: ${missing[*]}${NC}"
        echo ""
        echo "Set them with:"
        echo "export DATABASE_URL='postgres://postgres:password@localhost:5432/tickets?sslmode=disable'"
        echo "export JWT_SECRET='your-secret-key'"
        echo "export NATS_URL='nats://localhost:4222'"
        echo "export REDIS_HOST='localhost:6379'"
        return 1
    fi
    
    echo -e "${GREEN}✓ Environment variables are set${NC}"
}

# Run migrations
run_migrations() {
    echo -e "${YELLOW}Running database migrations...${NC}"
    
    for file in migrations/*.sql; do
        echo "Running $file..."
        PGPASSWORD=password psql -h localhost -U postgres -d tickets -f "$file" 2>/dev/null || true
    done
    
    echo -e "${GREEN}✓ Migrations completed${NC}"
}

# Start service
start_service() {
    local service=$1
    echo -e "${GREEN}Starting ${service} service...${NC}"
    go run "./cmd/${service}"
}

# Main menu
show_menu() {
    echo ""
    echo "==================================="
    echo "   Ticket Booking Dev Helper"
    echo "==================================="
    echo "1. Check infrastructure"
    echo "2. Start infrastructure (docker-compose)"
    echo "3. Stop infrastructure"
    echo "4. Run migrations"
    echo "5. Start Auth service"
    echo "6. Start Tickets service"
    echo "7. Start Orders service"
    echo "8. Start Payments service"
    echo "9. Start Expiration worker"
    echo "10. Run all tests"
    echo "11. Run tests with coverage"
    echo "0. Exit"
    echo "==================================="
    echo -n "Choose an option: "
}

while true; do
    show_menu
    read -r choice
    
    case $choice in
        1)
            check_infra
            ;;
        2)
            echo -e "${YELLOW}Starting infrastructure services...${NC}"
            docker compose up -d
            sleep 3
            check_infra
            ;;
        3)
            echo -e "${YELLOW}Stopping infrastructure services...${NC}"
            docker compose down
            ;;
        4)
            check_infra && run_migrations
            ;;
        5)
            check_infra && check_env && start_service "auth"
            ;;
        6)
            check_infra && check_env && start_service "tickets"
            ;;
        7)
            check_infra && check_env && start_service "orders"
            ;;
        8)
            check_infra && check_env && start_service "payments"
            ;;
        9)
            check_infra && start_service "expiration"
            ;;
        10)
            echo -e "${YELLOW}Running tests...${NC}"
            go test -v -race ./...
            ;;
        11)
            echo -e "${YELLOW}Running tests with coverage...${NC}"
            go test -v -race -coverprofile=coverage.out ./...
            go tool cover -html=coverage.out -o coverage.html
            echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
            ;;
        0)
            echo -e "${GREEN}Goodbye!${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}Invalid option${NC}"
            ;;
    esac
    
    echo ""
    echo -e "${YELLOW}Press Enter to continue...${NC}"
    read -r
done
