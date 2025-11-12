#!/bin/bash

# Ensure the script is run as root
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root or by sudo user" 1>&2
   exit 1
fi

set -e

# Configuration
KM_USER="km"
KM_HOME="/home/km"
SEKIN_REPO="https://github.com/KiraCore/sekin.git"
SEKIN_DIR="$KM_HOME/sekin"
COMPOSE_FILE="$SEKIN_DIR/compose.yml"

echo "======================================"
echo "SEKIN Bootstrap Script"
echo "======================================"
echo ""

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to update system
update_system() {
    echo "[1/5] Updating system packages..."
    apt-get update -qq || { echo "Failed to update system. Exiting..."; exit 1; }
    echo "✓ System updated"
}

# Function to install prerequisites
install_prerequisites() {
    echo "[2/5] Installing prerequisites..."
    apt-get install -y -qq \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg \
        lsb-release \
        git \
        jq \
        software-properties-common || { echo "Failed to install prerequisites. Exiting..."; exit 1; }
    echo "✓ Prerequisites installed"
}

# Function to install Docker
install_docker() {
    if command_exists docker; then
        echo "[3/5] Docker already installed ($(docker --version))"
        echo "✓ Skipping Docker installation"
    else
        echo "[3/5] Installing Docker..."

        # Add Docker's official GPG key
        install -m 0755 -d /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
        chmod a+r /etc/apt/keyrings/docker.asc

        # Add Docker repository
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
            tee /etc/apt/sources.list.d/docker.list > /dev/null

        # Install Docker
        apt-get update -qq
        apt-get install -y -qq \
            docker-ce \
            docker-ce-cli \
            containerd.io \
            docker-buildx-plugin \
            docker-compose-plugin || { echo "Failed to install Docker. Exiting..."; exit 1; }

        # Start and enable Docker service
        systemctl enable docker >/dev/null 2>&1
        systemctl start docker >/dev/null 2>&1

        echo "✓ Docker installed successfully ($(docker --version))"
    fi
}

# Function to create and configure km user
setup_km_user() {
    echo "[4/5] Setting up user 'km'..."

    # Create km user if it doesn't exist
    if id "$KM_USER" >/dev/null 2>&1; then
        echo "  - User 'km' already exists"
    else
        useradd -m -s /bin/bash "$KM_USER"
        echo "  - User 'km' created"
    fi

    # Add km to docker group
    if groups "$KM_USER" | grep -q docker; then
        echo "  - User 'km' already in docker group"
    else
        usermod -aG docker "$KM_USER"
        echo "  - User 'km' added to docker group"
    fi

    echo "✓ User 'km' configured"
}

# Function to clone repository
clone_repository() {
    echo "[5/5] Setting up SEKIN repository..."

    if [ -d "$SEKIN_DIR" ]; then
        echo "  - Repository already exists at $SEKIN_DIR"
        echo "  - Updating repository..."
        sudo -u "$KM_USER" git -C "$SEKIN_DIR" fetch --all || echo "    Warning: Could not fetch updates"
        sudo -u "$KM_USER" git -C "$SEKIN_DIR" pull || echo "    Warning: Could not pull updates"
    else
        echo "  - Cloning repository to $SEKIN_DIR..."
        sudo -u "$KM_USER" git clone "$SEKIN_REPO" "$SEKIN_DIR" || { echo "Failed to clone repository. Exiting..."; exit 1; }
    fi

    echo "✓ Repository ready at $SEKIN_DIR"
}

# Function to start services
start_services() {
    echo ""
    echo "======================================"
    echo "Starting SEKIN services..."
    echo "======================================"

    if [ ! -f "$COMPOSE_FILE" ]; then
        echo "Error: compose.yml not found at $COMPOSE_FILE"
        exit 1
    fi

    cd "$SEKIN_DIR"

    echo "Starting docker compose in detached mode..."
    sudo -u "$KM_USER" docker compose -f "$COMPOSE_FILE" up -d || { echo "Failed to start services. Exiting..."; exit 1; }

    echo ""
    echo "✓ Services started successfully!"
    echo ""
    echo "======================================"
    echo "Setup Complete!"
    echo "======================================"
    echo ""
    echo "Services running:"
    sudo -u "$KM_USER" docker compose -f "$COMPOSE_FILE" ps
    echo ""
    echo "To view logs: docker compose -f $COMPOSE_FILE logs -f"
    echo "To stop services: docker compose -f $COMPOSE_FILE down"
    echo ""
}

# Main execution
main() {
    update_system
    install_prerequisites
    install_docker
    setup_km_user
    clone_repository
    start_services
}

main
