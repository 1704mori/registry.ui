#!/bin/bash

GREEN='\033[0;32m'
NC='\033[0m' # No Color

if [[ "$1" == "docker" ]]; then
    MODE="docker"
else
    MODE="service"
fi

read -p "Enter Docker Registry URL: " DOCKER_REGISTRY_URL
read -p "Enter Docker Registry Username: " DOCKER_REGISTRY_USERNAME
read -sp "Enter Docker Registry Password: " DOCKER_REGISTRY_PASSWORD
echo

ENV_FILE=$(pwd)/.env.registry_ui

echo "DOCKER_REGISTRY_URL=$DOCKER_REGISTRY_URL
DOCKER_REGISTRY_USERNAME=$DOCKER_REGISTRY_USERNAME
DOCKER_REGISTRY_PASSWORD=$DOCKER_REGISTRY_PASSWORD" > "$ENV_FILE"

echo -e "${GREEN}Environment file created at: $ENV_FILE${NC}"

if [ "$MODE" == "docker" ]; then
    docker build -t registry_ui .
    docker run -d -p 8080:8080 --env-file "$ENV_FILE" registry_ui
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
