#!/bin/bash
# Install Goose if not present
if ! command -v goose &> /dev/null; then
    echo "Goose not found, installing..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi



echo "Waiting for PostgreSQL to start..."

# Wait for PostgreSQL to be ready (adjust this if necessary)
until pg_isready -h db -p 5432 -U postgres; do
  echo "Waiting for database to be ready..."
  sleep 2
done

# Run Goose migrations
echo "Running Goose migrations..."
goose -dir /migrations postgres "user=postgres password=postgres dbname=db sslmode=disable" up

echo "Goose migrations completed."
