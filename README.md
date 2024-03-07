# Sekin 
## Overview
This repository, Sekin, serves as a CI/CD hub, specifically automating the Docker image creation for the "Sekai" and "Interx" repositories. It leverages GitHub Actions to build and push Docker images to a Docker registry whenever changes are pushed to the relevant project directories.

## Sekai

To pull the image, execute:
```bash
docker pull ghcr.io/kiracore/sekin/sekai:v0.3.42
```

To initialize Sekai before starting the server, you can run:
```bash
docker run --rm -it -v $(pwd)/sekai:/sekai ghcr.io/kiracore/sekin/sekai:v0.3.42 init --overwrite --chain-id=testnet-1 "KIRA TEST LOCAL VALIDATOR NODE" --home=/sekai
```

To add a genesis account:
```bash
docker run -it --rm -v $(pwd)/sekai:/sekai $DOCKER_IMAGE keys add genesis --keyring-backend=test --home=/sekai --output=json
```

To add a validator account:
```bash
docker run -it --rm -v $(pwd)/sekai:/sekai $DOCKER_IMAGE add-genesis-account validator 300000ukex --keyring-backend=test --home=/sekai
```

To add a genesis account to the genesis file:
```bash
docker run -it --rm -v $(pwd)/sekai:/sekai $DOCKER_IMAGE add-genesis-account genesis 300000000000000ukex --keyring-backend=test --home=/sekai
```

To add a validator account to the genesis file:
```bash
docker run -it --rm -v $(pwd)/sekai:/sekai $DOCKER_IMAGE add-genesis-account validator 300000ukex --keyring-backend=test --home=/sekai
```

To initialize the genesis validator:
```bash
docker run -it --rm -v $(pwd)/sekai:/sekai $DOCKER_IMAGE gentx-claim genesis --keyring-backend=test --moniker="GENESIS VALIDATOR" --home=/sekai
```

To start the sekai node:
```bash
docker run -d -v $(pwd)/sekai:/sekai -p 26657:26657 -p 26656:26656 -p 26660:26660 --name sekai --restart always $DOCKER_IMAGE start --home=/sekai --rpc.laddr "tcp://0.0.0.0:26657"
```

To execute commands on a running Sekai node:
```bash
docker exec sekai /sekaid
```
