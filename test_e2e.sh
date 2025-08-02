#!/bin/bash
set -e

echo "Building je..."
go build -o je ./cmd/je

echo "Running end-to-end tests..."

# Test 1: Basic string assignment
echo '{}' > test1.json
./je test1.json name=john city="New York"
result=$(cat test1.json)
if ! echo "$result" | grep -q '"name":"john"' || ! echo "$result" | grep -q '"city":"New York"'; then
    echo "Test 1 failed: got $result"
    exit 1
fi
echo "✓ Test 1: Basic string assignment"

# Test 2: JSON value assignment
echo '{}' > test2.json
./je test2.json age:=30 active:=true balance:=99.50
result=$(cat test2.json)
if ! echo "$result" | grep -q '"age":30' || ! echo "$result" | grep -q '"active":true' || ! echo "$result" | grep -q '"balance":99.5'; then
    echo "Test 2 failed: got $result"
    exit 1
fi
echo "✓ Test 2: JSON value assignment"

# Test 3: Nested paths
echo '{}' > test3.json
./je test3.json user.name=john user.age:=30 config.port:=8080
result=$(cat test3.json)
if ! echo "$result" | grep -q '"user":{' || ! echo "$result" | grep -q '"name":"john"' || ! echo "$result" | grep -q '"age":30' || ! echo "$result" | grep -q '"config":{"port":8080}'; then
    echo "Test 3 failed: got $result"
    exit 1
fi
echo "✓ Test 3: Nested paths"

# Test 4: Array operations
echo '{"tags":["old"]}' > test4.json
./je test4.json tags[]=new tags[]=another
result=$(cat test4.json)
if ! echo "$result" | grep -q '"tags":\["old","new","another"\]'; then
    echo "Test 4 failed: got $result"
    exit 1
fi
echo "✓ Test 4: Array append"

# Test 5: Delete operation
echo '{"name":"john","age":30}' > test5.json
./je test5.json age:=
result=$(cat test5.json)
if echo "$result" | grep -q 'age'; then
    echo "Test 5 failed: age key not deleted, got $result"
    exit 1
fi
echo "✓ Test 5: Delete operation"

# Test 6: Pretty print
echo '{}' > test6.json
./je test6.json name=john --pretty
result=$(cat test6.json)
if ! echo "$result" | grep -q '  "name": "john"'; then
    echo "Test 6 failed: not pretty printed, got $result"
    exit 1
fi
echo "✓ Test 6: Pretty print"

# Test 7: stdin/stdout
result=$(echo '{}' | ./je - name=test)
expected='{"name":"test"}'
if [ "$result" != "$expected" ]; then
    echo "Test 7 failed: got $result, expected $expected"
    exit 1
fi
echo "✓ Test 7: stdin/stdout"

# Clean up
rm -f test*.json je

echo "All tests passed!"