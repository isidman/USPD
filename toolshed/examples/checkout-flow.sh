#!/usr/bin/env bash
# Scripted walkthrough of the full lending flow against a running server.
# Start the server first: `go run ./cmd/server` (from the toolshed/ root).
set -euo pipefail

BASE_URL="${TOOLSHED_URL:-http://localhost:8080}"

echo "== Creating a tool, a piece of hardware, and a software resource =="
DRILL=$(curl -sS -X POST "$BASE_URL/resources" -d '{"kind":"tool","name":"Cordless Drill"}')
echo "$DRILL"
DRILL_ID=$(echo "$DRILL" | python3 -c "import json,sys;print(json.load(sys.stdin)['id'])")

curl -sS -X POST "$BASE_URL/resources" -d '{"kind":"hardware","name":"Ender 3 3D Printer","metadata":{"serial":"E3-0042"}}'
echo
curl -sS -X POST "$BASE_URL/resources" -d '{"kind":"software","name":"FarmOS staging access","metadata":{"url":"https://farmos.example.org"}}'
echo

echo "== Listing resources (all available) =="
curl -sS "$BASE_URL/resources"
echo

echo "== Checking out the drill =="
LOAN=$(curl -sS -X POST "$BASE_URL/resources/$DRILL_ID/checkout" -d '{"borrower_id":"alice","duration_hours":48}')
echo "$LOAN"
LOAN_ID=$(echo "$LOAN" | python3 -c "import json,sys;print(json.load(sys.stdin)['id'])")

echo "== Trying to check it out again should fail with 409 =="
curl -sS -o /dev/null -w "HTTP %{http_code}\n" -X POST "$BASE_URL/resources/$DRILL_ID/checkout" -d '{"borrower_id":"bob"}'

echo "== Returning it =="
curl -sS -X POST "$BASE_URL/loans/$LOAN_ID/return"
echo

echo "== It's available again =="
curl -sS "$BASE_URL/resources"
echo
