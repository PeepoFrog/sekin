# Sekin 
## Overview
This repository, Sekin, serves as a CI/CD hub, specifically automating the Docker image creation for the "Sekai" and "Interx" repositories. It leverages GitHub Actions to build and push Docker images to a Docker registry whenever changes are pushed to the relevant project directories.


## Sekai REST

### Version

```bash
 curl -X POST http://localhost:8080/api/execute -H "Content-Type: application/json" -d '{"command":"version", "args":{}}'
```
