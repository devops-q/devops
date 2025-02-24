#!/usr/bin/env bash

# Kill any existing server process on port 8181
PID=$(lsof -t -i:8181)
if [ -n "$PID" ]; then
    kill "$PID"
fi

# Remove test database if exists
rm -f database.test.sqlite

DB_PATH=./database.test.sqlite go run cmd/create_api_user/main.go -username=test -password=test

# Start the server in background
PORT=8181 DB_PATH=./database.test.sqlite go run cmd/server/main.go &
#PORT=8181 DB_PATH=./database.test.sqlite go run cmd/server/main.go > /dev/null 2>&1 &

# Wait for server to start
sleep 2

# Run tests
go test ./tests
TEST_EXIT_CODE=$?


# Kill any existing server process on port 8181
kill "$(lsof -t -i:8181)"

# Clean up test database
rm -f database.test.sqlite

exit $TEST_EXIT_CODE