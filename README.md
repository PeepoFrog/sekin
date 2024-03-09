# Sekin 
## Overview
This repository, Sekin, serves as a CI/CD hub, specifically automating the Docker image creation for the "Sekai" and "Interx" repositories. It leverages GitHub Actions to build and push Docker images to a Docker registry whenever changes are pushed to the relevant project directories.

## Sekai REST

### Version

```bash
 curl -X POST http://localhost:8080/api/execute -H "Content-Type: application/json" -d '{"command":"version"}'
```
### Keys add

```bash
curl -X POST http://localhost:8080/api/execute -H "Content-Type: application/json" \
-d '{
    "command": "keys-add",
    "args": {
        "address": "validator",       
        "keyring-backend": "test",
        "home": "/sekai"
    }
}'
```



## Interx REST

### Version

```bash 
curl -X POST http://localhost:8081/api/execute -H "Content-Type: application/json" -d '{"command":"version"}'
```
### Init 

Full list of args

```bash
curl -X POST http://localhost:8081/api/execute -H "Content-Type: application/json" -d '{
  "command": "init", "args":{
  "addrbook": "/path/to/addrbook.json",
  "cache_dir": "/path/to/cache_dir",
  "caching_duration": 60,
  "download_file_size_limitation": "10MB",
  "faucet_amounts": "100000stake,100000ukex",
  "faucet_minimum_amounts": "1000stake,1000ukex",
  "faucet_mnemonic": "arch ride cement wink scale flight vault shrimp rigid scrap parade voice author cloth cigar crew ivory recall argue vicious index solve swing hill",
  "faucet_time_limit": 30,
  "fee_amounts": "stake 1000ukex",
  "grpc": "dns:///sekai:9090",
  "halted_avg_block_times": 10,
  "home": "/home/d/tmp/tmpInterx",
  "max_cache_size": "2GB",
  "node_discovery_interx_port": "11000",
  "node_discovery_tendermint_port": "26657",
  "node_discovery_timeout": "3s",
  "node_discovery_use_https": true,
  "node_key": "node_key.json",
  "node_type": "seed",
  "port": "8080",
  "rpc": "http://sekai:26657",
  "seed_node_id": "your_seed_node_id",
  "sentry_node_id": "your_sentry_node_id",
  "serve_https": false,
  "signing_mnemonic": "arch ride cement wink scale flight vault shrimp rigid scrap parade voice author cloth cigar crew ivory recall argue vicious index solve swing hill",
  "snapshot_interval": 100,
  "snapshot_node_id": "your_snapshot_node_id",
  "status_sync": 5,
  "tx_modes": "sync,async,block",
  "validator_node_id": "your_validator_node_id"
}}'
```

Realistic one 

```bash
curl -X POST http://localhost:8081/api/execute -H "Content-Type: application/json" -d '{
  "command": "init", "args":{
  "faucet_mnemonic": "arch ride cement wink scale flight vault shrimp rigid scrap parade voice author cloth cigar crew ivory recall argue vicious index solve swing hill",
  "grpc": "dns:///sekai:9090",
  "home": "/home/d/tmp/tmpInterx",
  "node_type": "validator",
  "port": "11000",
  "rpc": "http://sekai:26657",
  "signing_mnemonic": "arch ride cement wink scale flight vault shrimp rigid scrap parade voice author cloth cigar crew ivory recall argue vicious index solve swing hill",
  "validator_node_id": "your_validator_node_id",
}}'
```