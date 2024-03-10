#!/bin/env bash

# Ensure the script is run as root
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root or by sudo user" 1>&2
   exit 1
fi

# System vars
ARCHITECTURE=$(uname -m)

# Docker vars
BASE_IMAGE="python:3.12.2-slim"
DOCKER_VER="7.0.0"
DOCKER_COMPOSE_VER="v2.24.6"

# Ansible vars
ANSIBLE_VER="9.2.0"
ANSIBLE_TAG="ansible-runner"
ANSIBLE_DIR="/root/ansible-runner"
ANSIBLE_DOCKERFILE="$ANSIBLE_DIR/Dockerfile"
ANSIBLE_ENTRYPOINT="$ANSIBLE_DIR/entrypoint.sh"

# Function to update system
update_system() {
    echo "Updating system..."
    apt-get update || { echo "Failed to update system. Exiting..."; exit 1; }
}

# Function to install prerequisites
install_prerequisites() {
    echo "Installing prerequisites..."
    apt-get install -y apt-transport-https ca-certificates wget curl software-properties-common jq || { echo "Failed to install prerequisites. Exiting..."; exit 1; }
}

# Function to install Docker
install_docker() {
    echo "Installing Docker..."	
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null || { echo "Failed to add Docker repository to apt sources."; exit 1; }
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin || { echo "Failed to install Docker."; exit 1; }
}

create_dockerfile_ansible_runner() {
    cat > "$ANSIBLE_DOCKERFILE" <<EOF
# Use build argument for the base image
ARG BASE_IMAGE=$BASE_IMAGE

# Use an official Python runtime as a parent image
FROM \$BASE_IMAGE

# Environment variables to be added when building
ENV DOCKER_VER=$DOCKER_VER
ENV ANSIBLE_VER=$ANSIBLE_VER

# Install dependencies required for ansible and ssh connections
RUN apt-get update && \\
    apt-get install -y --no-install-recommends \\
    ssh \\
    openssh-client \\
    rsyslog \\
    systemd \\
    && apt-get clean \\
    && rm -rf /var/lib/apt/lists/* \\
    && python -m pip install --upgrade pip cffi \\
    && pip install ansible==\${ANSIBLE_VER} \\
    && pip install docker==\${DOCKER_VER}

# Set the working directory in the container to /src
WORKDIR /src

# Copy the entrypoint script and src directory into the container
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["ansible-playbook", "--version"]
EOF

    echo "Dockerfile created."
}

create_entrypoint_ansible_runner() {
    cat > "$ANSIBLE_ENTRYPOINT" <<EOF
#!/bin/bash

# Default to showing the Ansible version if no command is specified
if [ \$# -eq 0 ]; then
    exec ansible-playbook --version
else
    exec "\$@"
fi
EOF

    chmod +x "$ANSIBLE_ENTRYPOINT"
    echo "Entrypoint script created."
}

create_ansible_runner() {
    echo "Creating ansible-runner container..."
    mkdir -p "$ANSIBLE_DIR"
    create_dockerfile_ansible_runner || { echo "Failed to create Dockerfile."; exit 1; }
    create_entrypoint_ansible_runner || { echo "Failed to create entrypoint script."; exit 1; }
    cd "$ANSIBLE_DIR" || exit
    docker build -t "$ANSIBLE_TAG:$ANSIBLE_VER" . || { echo "Failed to build ansible-runner container."; exit 1; }
    echo "Ansible runner container created."
}

install_docker_compose() {
    echo "Installing dcoker compose..."
    DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
    mkdir -p $DOCKER_CONFIG/cli-plugins
    curl -SL "https://github.com/docker/compose/releases/download/$DOCKER_COMPOSE_VER/docker-compose-linux-$ARCHITECTURE" -o $DOCKER_CONFIG/cli-plugins/docker-compose
    mv $DOCKER_CONFIG/cli-plugins/docker-compose /usr/local/bin
    chmod 755 /usr/local/bin/docker-compose
    echo "Installed $(docker-compose version)"
}

main() {
    update_system
    install_prerequisites
    install_docker
    install_docker_compose
    create_ansible_runner
}

# Call the main function
main

