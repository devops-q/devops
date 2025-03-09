#!/usr/bin/env bash

# Kill any existing server process on port 8181
PID=$(lsof -t -i:8181)
if [ -n "$PID" ]; then
    echo "Killing existing server process on port 8181"
    kill "$PID"
fi

# Remove test database if exists
rm -f database.test.sqlite

# Create the test user and database
export DB_HOST=localhost
export DB_NAME=minitwit_db_test
export DB_USER=minitwit_test_user
export DB_PASSWORD=minitwit_test_password
export DB_PORT=5432

docker run --name minitwit-postgres-test -e POSTGRES_USER=$DB_USER -e POSTGRES_PASSWORD=$DB_PASSWORD -e POSTGRES_DB=$DB_NAME -p $DB_PORT:5432 -d postgres:alpine
CONTAINER_ID=$(docker ps -q -f name=minitwit-postgres-test)

# Wait for postgres to start
for i in {1..10}; do
  if docker exec $CONTAINER_ID pg_isready -U $DB_USER -d $DB_NAME -h localhost -p $DB_PORT > /dev/null; then
    echo "Postgres is up and running."
    break
  else
    echo "Waiting for postgres to start... Attempt $i"
    sleep 2
  fi
done

go run cmd/create_api_user/main.go -username=test -password=test

# Start the server in the background
PORT=8181 go run cmd/server/main.go &
SERVER_PID=$!

# Wait for server to start by checking for a response on port 8181
for i in {1..10}; do
  if curl -s http://127.0.0.1:8181 > /dev/null; then
    echo "Server is up and running."
    break
  else
    echo "Waiting for server to start... Attempt $i"
    sleep 2
  fi
done

# Run tests
go test ./tests
TEST_EXIT_CODE=$?

# Kill the server process
echo "Killing the server process"
kill "$SERVER_PID"

# Stop the postgres container
docker stop $CONTAINER_ID

# Exit with the test exit code
exit $TEST_EXIT_CODE
