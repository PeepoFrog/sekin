#!/bin/env bash
curl -X POST http://localhost:8080/api/execute -H "Content-Type: application/json" \
-d '{
    "command": "keys-add",
    "args": {
        "address": "genesis",       
        "keyring-backend": "test",
        "home": "/sekai"
    }
}'

