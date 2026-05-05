#!/bin/bash

echo "🎮 Starting Tic-Tac-Toe API Server..."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running"
    echo "Please start Docker and try again"
    exit 1
fi

# Start PostgreSQL
echo "📦 Starting PostgreSQL database..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for database to be ready..."
sleep 3

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  Warning: .env file not found, copying from .env.example"
    cp .env.example .env
fi

# Start the backend server
echo "🚀 Starting backend server..."
echo ""
cd backend
go run cmd/server/main.go
