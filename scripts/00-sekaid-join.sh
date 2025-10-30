#!/bin/env bash
curl -X POST "http://localhost:8282/api/execute" \
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
                "mnemonic": "burden size near tragic fitness couch search suffer fluid output expire swap poem utility sing genuine replace dune tenant monkey sauce soccer twin sentence",
                "local": false
            }
         }'
