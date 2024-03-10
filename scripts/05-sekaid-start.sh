#!/bin/env bash
curl -X POST http://localhost:8080/api/execute \
-H "Content-Type: application/json" \
-d '{
    "command": "start",
    "args": {
        "home": "/sekai"
    }
}'
