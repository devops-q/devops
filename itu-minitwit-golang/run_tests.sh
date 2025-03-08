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
DB_PATH=./database.test.sqlite go run cmd/create_api_user/main.go -username=test -password=test

# Start the server in the background
PORT=8181 DB_PATH=./database.test.sqlite go run cmd/server/main.go &
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

# Clean up test database
rm -f database.test.sqlite

# Exit with the test exit code
exit $TEST_EXIT_CODE
