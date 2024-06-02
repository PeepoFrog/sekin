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

# HOST
KM_DIR="/home/km"
SEKIN_GIT="https://github.com/KiraCore/sekin.git"
SEKIN_DIR="$KM_DIR/sekin"
COMPOSE_URL="https://raw.githubusercontent.com/KiraCore/sekin/main/compose.yml"
COMPOSE_PATH="/home/km/sekin/compose.yml"

# Function to update system
update_system() {
    echo "Updating system..."
    sudo apt-get update || { echo "Failed to update system. Exiting..."; exit 1; }
}

# Function to install prerequisites
install_prerequisites() {
    echo "Installing prerequisites..."
    sudo apt-get install -y apt-transport-https ca-certificates wget curl software-properties-common jq || { echo "Failed to install prerequisites. Exiting..."; exit 1; }
}

# Function to install Docker
install_docker() {
    echo "Installing Docker..."	
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null || { echo "Failed to add Docker repository to apt sources."; exit 1; }
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin || { echo "Failed to install Docker."; exit 1; }
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
    echo "Installing Docker Compose..."

    # Setting up the Docker config directory
    DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
    mkdir -p "$DOCKER_CONFIG/cli-plugins" || { echo "Failed to create directory: $DOCKER_CONFIG/cli-plugins"; exit 1; }

    # Downloading Docker Compose
    echo "Downloading Docker Compose version $DOCKER_COMPOSE_VER for $ARCHITECTURE..."
    curl -SL "https://github.com/docker/compose/releases/download/$DOCKER_COMPOSE_VER/docker-compose-linux-$ARCHITECTURE" -o "$DOCKER_CONFIG/cli-plugins/docker-compose" || { echo "Download failed"; exit 1; }

    # Moving Docker Compose to a system-wide directory (requires root privileges)
    echo "Moving Docker Compose to /usr/local/bin..."
    mv "$DOCKER_CONFIG/cli-plugins/docker-compose" /usr/local/bin || { echo "Move failed"; exit 1; }

    # Making Docker Compose executable
    echo "Setting execute permissions for Docker Compose..."
    sudo chmod 755 /usr/local/bin/docker-compose || { echo "Failed to set permissions"; exit 1; }

    # Verifying installation
    if docker-compose version; then
        echo "Docker Compose installed successfully."
    else
        echo "Failed to verify Docker Compose installation."
        exit 1
    fi
}

add_user_km() {
        echo "Creating user 'km' with a home directory and bash as the default shell..."
        sudo useradd -m -s /bin/bash km || :
}

add_km_to_docker_group() {
    echo "Adding user 'km' to the docker group..."
    sudo usermod -aG docker km || :
}


download_compose_and_change_owner(){
    # Attempt to download the file up to 5 times
    local max_attempts=5
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        echo "Attempt $attempt: Downloading compose.yml from ${COMPOSE_URL}..."
        curl -o "${COMPOSE_PATH}" "${COMPOSE_URL}"

        # Check if the download was successful
        if [ $? -eq 0 ]; then
            echo "Download successful. Setting ownership..."
            # Set km as the owner of the downloaded file
            sudo chown km:km "${COMPOSE_PATH}"

            # Check if the chown command was successful
            if [ $? -eq 0 ]; then
                echo "Ownership set to 'km'."
                return 0
            else
                echo "Failed to set ownership. Check permissions."
                return 1
            fi
        else
            echo "Download failed. Retrying..."
            # Random delay between 1 and 3 seconds before retrying
            sleep $((RANDOM % 3 + 1))
            ((attempt++))
        fi
    done

    echo "Failed to download the file after $max_attempts attempts."
    return 1
}

# Function to run docker-compose up -d as user km
run_docker_compose_as_km() {
    echo "Starting docker-compose up as user km..."

    # Execute docker-compose up -d as user km
    if ! sudo -u km docker-compose -f "$COMPOSE_PATH" up -d; then
        echo "Failed to start docker-compose. Exiting."
        exit 1
    fi

    echo "docker-compose started successfully."
}

clone_repo_as_km() {
    # Specify the directory where you want to clone the repo
    echo "Cloning the repository as user 'km'..."

    # Clone the repository using sudo -u to run as user 'km'
    if sudo -u km git clone "$SEKIN_GIT" "$SEKIN_DIR"; then
        echo "Repository successfully cloned into $SEKIN_DIR."
    else
        echo "Failed to clone the repository. Please check permissions and repository URL."
        exit 1
    fi
}

main() {
    update_system
    install_prerequisites
    install_docker
    install_docker_compose
    add_user_km
    add_km_to_docker_group
    # download_compose_and_change_owner
    clone_repo_as_km
    run_docker_compose_as_km
}

main

