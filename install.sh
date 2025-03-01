#!/bin/bash

REPO_URL="https://github.com/1704mori/registry.ui.git"
REPO_DIR="/tmp/registry.ui"

if [ ! -d "$REPO_DIR" ]; then
    echo "Cloning repository from $REPO_URL..."
    git clone "$REPO_URL" "$REPO_DIR" || { echo "Failed to clone repository."; exit 1; }
else
    echo "Repository already exists. Pulling latest changes..."
    cd "$REPO_DIR" && git pull && cd ..
fi

cd "$REPO_DIR" || exit 1

GREEN='\033[0;32m'
NC='\033[0m' # No Color

if [[ "$1" == "docker" ]]; then
    MODE="docker"
else
    MODE="service"
fi

read -p "Enter Docker Registry URL: " REGISTRY_URL
read -p "Enter Docker Registry Username: " REGISTRY_USERNAME
read -sp "Enter Docker Registry Password: " REGISTRY_PASSWORD
echo

ENV_FILE=$(pwd)/.env.registry_ui

echo "REGISTRY_URL=$REGISTRY_URL
REGISTRY_USERNAME=$REGISTRY_USERNAME
REGISTRY_PASSWORD=$REGISTRY_PASSWORD" > "$ENV_FILE"

echo -e "${GREEN}Environment file created at: $ENV_FILE${NC}"

if [ "$MODE" == "docker" ]; then
    CONTAINER_ID=$(docker ps -aq -f name=^/registry_ui$)
    if [ -n "$CONTAINER_ID" ]; then
        echo "Container 'registry_ui' already exists. Stopping and removing it..."
        docker stop registry_ui
        docker rm registry_ui
    fi

    docker build -t registry_ui .
    docker run --name registry_ui -d -p 8080:8080 --env-file "$ENV_FILE" registry_ui
    echo -e "${GREEN}Docker container built and started.${NC}"
else
    go build -o /usr/local/bin/registry-ui cmd/server/main.go
    echo -e "${GREEN}Binary built at: /usr/local/bin/registry-ui${NC}"

    # Determine the original user and group when running with sudo.
    if [ -n "$SUDO_USER" ]; then
        CURRENT_USER=$SUDO_USER
        CURRENT_GROUP=$(id -gn "$SUDO_USER")
    else
        CURRENT_USER=$(whoami)
        CURRENT_GROUP=$(id -gn)
    fi

    SERVICE_FILE=/etc/systemd/system/registry-ui.service

    echo "[Unit]
Description=Registry.UI
After=network.target

[Service]
User=${CURRENT_USER}
Group=${CURRENT_GROUP}
ExecStart=/usr/local/bin/registry-ui
WorkingDirectory=$(pwd)
Restart=always
EnvironmentFile=$ENV_FILE

[Install]
WantedBy=multi-user.target" > "$SERVICE_FILE"

    systemctl daemon-reload
    systemctl enable registry-ui.service
    systemctl start registry-ui.service

    echo -e "${GREEN}Registry.UI has been set up and started as user ${CURRENT_USER}.${NC}"
fi
