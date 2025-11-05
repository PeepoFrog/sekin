#!/bin/env bash
curl -X POST "http://localhost:8282/api/execute" \
     -H "Content-Type: application/json" \
     -d '{
            "command": "join",
            "args": {
                "ip": "VALIDATOR_NODE_IP_ADDRESS_HERE",
                "interx_port": 11000,
                "rpc_port": 26657,
                "p2p_port": 26656,
                "sekaidAddress": "sekai.local",
                "interxAddress": "proxy.local",
                "mnemonic": "YOUR_MNEMONIC_PHRASE_HERE",
                "local": false
            }
         }'
