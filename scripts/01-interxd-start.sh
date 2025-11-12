curl -X POST http://localhost:8081/api/execute \
-H "Content-Type: application/json" \
-d '{
    "command": "start",
    "args": {
        "home": "/interx"
    }
}'