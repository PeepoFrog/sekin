curl -X POST http://localhost:8081/api/execute -H "Content-Type: application/json" -d '{
  "command": "init", "args":{
  "faucet_mnemonic": "arch ride cement wink scale flight vault shrimp rigid scrap parade voice author cloth cigar crew ivory recall argue vicious index solve swing hill",
  "grpc": "dns:///sekin-sekaid_rpc-1:9090",
  "home": "/interx",
  "node_type": "validator",
  "port": "11000",
  "rpc": "http://sekin-sekaid_rpc-1:26657",
  "signing_mnemonic": "arch ride cement wink scale flight vault shrimp rigid scrap parade voice author cloth cigar crew ivory recall argue vicious index solve swing hill",
  "validator_node_id": "7ec83d60aba744504c8bd1fcdddb98b8f50652aa"
}}'