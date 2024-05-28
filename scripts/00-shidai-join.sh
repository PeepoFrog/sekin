
curl -X POST "http://localhost:8282/api/execute" \
     -H "Content-Type: application/json" \
     -d '{
            "command": "join",
            "args": {
                "ip": "10.43.239.82",
                "interxPort": 11000,
                "rpcPort": 26657,
                "p2pPort": 26656,
                "sekaidAddress": "sekai.local",
                "interxAddress": "interx.local",
                "mnemonic": "bargain erosion electric skill extend aunt unfold cricket spice sudden insane shock purpose trumpet holiday tornado fiction check pony acoustic strike side gold resemble"
            }
         }'
