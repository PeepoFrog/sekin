# Sekin 
## Overview
This repository, Sekin, serves as a CI/CD hub, specifically automating the Docker image creation for the "Sekai" and "Interx" repositories. It leverages GitHub Actions to build and push Docker images to a Docker registry whenever changes are pushed to the relevant project directories.

### Building

1. Prepare your host with `./scripts/bootstrap.sh` script:

```bash
chmod +x ./scripts/*
```

2. Run the script. It will install all the dependencies:
```bash
./scripts/bootstrap.sh
```

3. To run both Sekai and Interx we can use compose file. Run the command:

```bash
docker compose build
docker compose up
```

Apart from using compose we can build Sekai and Interx independently using Docker. 

1. Make scripts executable running command from above.

2. To build the Sekai image run:

```bash
./scripts/docker-sekaid-build.sh v0.0.1
```

3. To build the Interx image run:

```bash
./scripts/docker-interxd-build.sh v0.0.1
```

4. To run the Sekai container:

```bash
./scripts/docker-sekaid-run.sh v0.0.1
```
5. To run the Interx container:

```bash
./scripts/docker-interxd-run.sh v0.0.1
```

