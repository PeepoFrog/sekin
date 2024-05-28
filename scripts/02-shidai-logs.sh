
curl -X POST "http://127.0.0.1:8282/api/execute" \
     -H "Content-Type: application/json" \
     -d '{
            "command": "logs",
            "args": {}
         }'
