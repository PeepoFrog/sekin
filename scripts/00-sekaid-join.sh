#!/bin/env bash
curl -X POST "http://localhost:11000/api/execute" \
     -H "Content-Type: application/json" \
     -d '{
            "command": "join",
            "args": {
                "ip": "'"$1"'",
                "interx_port": 80,
                "rpc_port": 26657,
                "p2p_port": 26656,
                "sekaidAddress": "sekai.local",
                "interxAddress": "interx.local",
                "mnemonic": "YOUR_MNEMONIC_PHRASE_HERE",
                "local": false
            }
         }'
